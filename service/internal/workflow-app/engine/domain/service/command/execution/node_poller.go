package execution

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/nodeexecutor"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	minPollFailedCount = 0
	maxPollFailedCount = 10
)

// NodePoller 节点轮询者
type NodePoller interface {
	Polling(ctx context.Context, req *dto.PollingNodeDTO) error
}

// DefaultNodePoller 节点轮询者
type DefaultNodePoller struct {
	workflowInstRepo     ports.WorkflowInstRepository
	nodeInstRepo         ports.NodeInstRepository
	eventBusRepo         ports.EventBusRepository
	workflowUpdater      WorkflowUpdater
	exprEvaluator        expr.Evaluator
	nodeExecutorRegistry nodeexecutor.Registry
}

// NewDefaultNodePoller 初始化
func NewDefaultNodePoller(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet,
	workflowUpdater *DefaultWorkflowUpdater) *DefaultNodePoller {
	r := &DefaultNodePoller{
		workflowInstRepo:     repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:         repoProviderSet.NodeInstRepo(),
		eventBusRepo:         repoProviderSet.EventBusRepo(),
		exprEvaluator:        workflowProviderSet.ExprEvaluator(),
		nodeExecutorRegistry: workflowProviderSet.NodeExecutorRegistry(),
		workflowUpdater:      workflowUpdater,
	}
	return r
}

// Polling 轮询节点
func (r *DefaultNodePoller) Polling(ctx context.Context, req *dto.PollingNodeDTO) error {
	nodeInst, err := r.nodeInstRepo.Get(&dto.GetNodeInstDTO{
		NodeInstID: req.NodeInstID,
		InstID:     req.InstID,
		DefID:      req.DefID,
		DefVersion: req.DefVersion,
	})

	// 如果查询实例出了问题, 直接返回, 这种情况也不会继续轮询也不用继续更新结果
	if err != nil {
		return err
	}

	// 节点如果已经结束, 就没必要再轮询了
	if nodeInst.Status.IsTerminal() {
		log.Warnf(nodeInstTerminalErrFormat+", skip polling",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
		return nil
	}

	if err := r.doPolling(ctx, nodeInst); err != nil {
		return r.handleErrInPolling(err, nodeInst)
	}

	return r.afterPolling(nodeInst)
}

func (r *DefaultNodePoller) handleErrInPolling(pollingErr error, nodeInst *entity.NodeInst) error {
	if needRetryPolling(pollingErr, nodeInst) {
		if err := r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst); err != nil {
			log.Errorf("Failed to update node inst, caused by %s", err)
			return fmt.Errorf("failed to update node inst, %w", err)
		}

		return r.sendPollDriveEvent(nodeInst)
	}

	if nodeInst.Status.IsTerminal() {
		return nil
	}

	nodeInst.Status = entity.NodeInstFailed
	nodeInst.Reason.FailedReason = pollingErr.Error()
	return r.updateNodeInstStatusAndSendDriveEvent(nodeInst)
}

func (r *DefaultNodePoller) updateNodeInstStatusAndSendDriveEvent(nodeInst *entity.NodeInst) error {
	if err := r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst); err != nil {
		log.Errorf("Failed to update node inst, caused by %s", err)
		return fmt.Errorf("failed to update node inst, %w", err)
	}

	return r.workflowUpdater.SendNodeDriveEvent(nodeInst, event.NodeCompleteDrive)
}

func needRetryPolling(err error, nodeInst *entity.NodeInst) bool {
	if nodeInst.PollFailedCount > minPollFailedCount && nodeInst.PollFailedCount < maxPollFailedCount {
		return true
	}

	return errno.NeedRetryErr(err)
}

func (r *DefaultNodePoller) needContinuePolling(inst *entity.WorkflowInst, nodeInst *entity.NodeInst) bool {
	match, err := r.isMatchPollingCondition(inst, nodeInst)
	if err != nil {
		return true
	}

	return match
}

func (r *DefaultNodePoller) afterPolling(nodeInst *entity.NodeInst) error {
	inst, err := r.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(nodeInst.InstID, nodeInst.DefID, ""))
	if err != nil {
		return err
	}

	if r.needContinuePolling(inst, nodeInst) {
		return r.sendPollDriveEvent(nodeInst)
	}

	nodeInst.Status, err = r.checkNodeInstStatus(inst, nodeInst)
	if err != nil {
		nodeInst.Status = entity.NodeInstFailed
		nodeInst.Reason.FailedReason = fmt.Sprintf("failed to evaluate node status, caused by %s", err)
	}

	return r.updateNodeInstStatusAndSendDriveEvent(nodeInst)
}

func (r *DefaultNodePoller) sendPollDriveEvent(nodeInst *entity.NodeInst) error {
	deliverAfter, err := r.evaluatePollingDurationTime(nodeInst)
	if err != nil {
		return err
	}
	if err := r.workflowUpdater.SendDelayNodeDriveEvent(nodeInst, deliverAfter, event.NodePollDrive); err != nil {
		return fmt.Errorf("failed to produce poll drive msg, caused by %s: %w", err, errno.Unavailable)
	}

	return nil
}

func (r *DefaultNodePoller) isMatchPollingCondition(inst *entity.WorkflowInst,
	nodeInst *entity.NodeInst) (bool, error) {
	basicArgs, err := entity.GetServiceNodeBasicArgs(nodeInst.NodeDef, entity.PollingArgs)
	if err != nil {
		return false, err
	}

	if basicArgs.PollCondition == "" {
		return false, fmt.Errorf("args `pollCondition` must not be empty")
	}

	match, err := matchConditionForCurNodeInst(r.exprEvaluator, inst, nodeInst, basicArgs.PollCondition)
	if err != nil {
		return false, err
	}

	return match, nil
}

func (r *DefaultNodePoller) checkNodeInstStatus(inst *entity.WorkflowInst,
	nodeInst *entity.NodeInst) (entity.NodeInstStatus, error) {
	serviceNodeBasicArgs, err := entity.GetServiceNodeBasicArgs(nodeInst.NodeDef, entity.PollingArgs)
	if err != nil {
		return entity.NodeInstStatus{}, err
	}

	// 如果取消条件没配, 默认直接检查是否成功
	if serviceNodeBasicArgs.CancelCondition == "" {
		return r.checkNodeInstSucceed(inst, nodeInst, serviceNodeBasicArgs)
	}

	return r.checkNodeInstCancelled(inst, nodeInst, serviceNodeBasicArgs)
}

func (r *DefaultNodePoller) checkNodeInstCancelled(inst *entity.WorkflowInst, nodeInst *entity.NodeInst,
	serviceNodeBasicArgs *entity.ServiceNodeBasicArgs) (entity.NodeInstStatus, error) {
	isCancelled, err := matchConditionForCurNodeInst(r.exprEvaluator, inst,
		nodeInst, serviceNodeBasicArgs.CancelCondition)
	if err != nil {
		return entity.NodeInstStatus{}, err
	}

	if isCancelled {
		return entity.NodeInstCancelled, nil
	}

	return r.checkNodeInstSucceed(inst, nodeInst, serviceNodeBasicArgs)
}

func (r *DefaultNodePoller) checkNodeInstSucceed(inst *entity.WorkflowInst, nodeInst *entity.NodeInst,
	serviceNodeBasicArgs *entity.ServiceNodeBasicArgs) (entity.NodeInstStatus, error) {
	// 如果条件没配, 默认为执行成功
	if serviceNodeBasicArgs.SuccessCondition == "" {
		return entity.NodeInstSucceed, nil
	}

	isSucceed, err := matchConditionForCurNodeInst(r.exprEvaluator, inst, nodeInst, serviceNodeBasicArgs.SuccessCondition)
	if err != nil {
		return entity.NodeInstStatus{}, err
	}
	if !isSucceed {
		nodeInst.Reason.FailedReason = fmt.Sprintf(successConditionNotMatchErrFormat,
			serviceNodeBasicArgs.SuccessCondition, utils.StructToJsonStr(nodeInst.PollOutput))
		return entity.NodeInstFailed, nil
	}

	return entity.NodeInstSucceed, nil
}

func (r *DefaultNodePoller) doPolling(ctx context.Context, nodeInst *entity.NodeInst) error {
	executor, exists := r.nodeExecutorRegistry.GetExecutor(nodeInst.BasicNodeDef.Type)
	if !exists {
		return fmt.Errorf("failed to get nodeType=%s executor", nodeInst.BasicNodeDef.Type)
	}
	nodeInst.PollCount += 1
	if err := executor.Polling(ctx, nodeInst); err != nil {
		return err
	}
	return r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst)
}

var (
	pollDefaultInitialDuration = 10 * time.Second
	pollMaxDuration            = 30 * time.Minute
	pollMinDuration            = 5 * time.Second
	pollDefaultMaxDuration     = 60 * time.Second
)

// evaluatePollingDurationTime 计算延迟时间
func (r *DefaultNodePoller) evaluatePollingDurationTime(nodeInst *entity.NodeInst) (time.Duration, error) {
	serviceNodeBasicArgs, err := entity.GetServiceNodeBasicArgs(nodeInst.NodeDef, entity.PollingArgs)
	if err != nil {
		return 0, err
	}

	initialDuration := getPollInitialDuration(serviceNodeBasicArgs.InitialDuration)
	if strings.EqualFold(serviceNodeBasicArgs.Policy.String(), entity.Fixed.String()) {
		return initialDuration, nil
	}

	curDuration := time.Duration(nodeInst.PollCount) * initialDuration
	maxDuration := getPollMaxDuration(serviceNodeBasicArgs.MaxDuration)
	if curDuration.Milliseconds() > maxDuration.Milliseconds() {
		return maxDuration, nil
	}

	if curDuration.Milliseconds() < pollMinDuration.Milliseconds() {
		return pollMinDuration, nil
	}

	return curDuration, nil
}

func getPollMaxDuration(durationStr string) time.Duration {
	if durationStr == "" {
		return pollDefaultMaxDuration
	}

	duration, err := expr.ParseDuration(durationStr)
	if err != nil {
		log.Warnf("Failed to ParseDuration, duration=[%s], caused by %s", durationStr, err)
		return pollDefaultMaxDuration
	}

	if duration.Milliseconds() > pollMaxDuration.Milliseconds() {
		return pollMaxDuration
	}

	return duration
}

func getPollInitialDuration(initialDurationStr string) time.Duration {
	if initialDurationStr == "" {
		return pollDefaultInitialDuration
	}

	duration, err := expr.ParseDuration(initialDurationStr)
	if err != nil {
		log.Warnf("Failed to ParseDuration, duration=[%s], caused by %s", initialDurationStr, err)
		return pollDefaultInitialDuration
	}

	if duration.Milliseconds() > pollMaxDuration.Milliseconds() {
		return pollMaxDuration
	}

	return duration
}
