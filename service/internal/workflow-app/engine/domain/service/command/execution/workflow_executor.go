package execution

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/pkg/errors"
)

// WorkflowExecutor 流程实例执行器
// 注意调这些接口的时候需要把原因塞上，表示因为什么做了这个操作
type WorkflowExecutor interface {
	DriveNext(ctx context.Context, req *dto.DriveWorkflowInstNextNodesDTO) error
}

// DefaultWorkflowExecutor 流程驱动器实现
type DefaultWorkflowExecutor struct {
	workflowDefRepo  ports.WorkflowDefRepository
	workflowInstRepo ports.WorkflowInstRepository
	nodeInstRepo     ports.NodeInstRepository
	cacheRepo        ports.CacheRepository
	workflowUpdater  WorkflowUpdater
	workflowDecider  WorkflowDecider
}

// NewDefaultWorkflowExecutor 初始化
func NewDefaultWorkflowExecutor(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet,
	workflowUpdater *DefaultWorkflowUpdater) *DefaultWorkflowExecutor {
	return &DefaultWorkflowExecutor{
		workflowDefRepo:  repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:     repoProviderSet.NodeInstRepo(),
		cacheRepo:        repoProviderSet.CacheRepo(),
		workflowDecider:  workflowProviderSet.WorkflowDecider(),
		workflowUpdater:  workflowUpdater,
	}
}

// DriveNext 驱动下一个节点
func (e *DefaultWorkflowExecutor) DriveNext(ctx context.Context, req *dto.DriveWorkflowInstNextNodesDTO) error {
	if err := e.doDriveNext(ctx, req); err != nil {
		if errno.NeedRetryErr(err) {
			return err
		}

		return e.workflowUpdater.UpdateWorkflowInstFailed(&dto.UpdateWorkflowInstFailedDTO{
			DefID:  req.DefID,
			InstID: req.InstID,
			Reason: err.Error(),
		})
	}

	return nil
}

func (e *DefaultWorkflowExecutor) doDriveNext(ctx context.Context, req *dto.DriveWorkflowInstNextNodesDTO) error {
	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, req.CurNodeInstID))
	if err != nil {
		return err
	}

	if inst.NodeInstsCount >= config.GetValidationRulesConfig().MaxNodeInstsForOneFlow {
		return fmt.Errorf("[%s]inst exceed %d node insts",
			logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID),
			config.GetValidationRulesConfig().MaxNodeInstsForOneFlow)
	}

	if inst.Status.IsTerminal() {
		return nil
	}

	if inst.Status == entity.InstPaused {
		return e.handleForInstAlreadyPaused(inst)
	}

	// 1. 决策器决策
	decideResult, err := e.workflowDecider.Decide(inst)
	if err != nil {
		return errors.Wrapf(err, "failed to decide")
	}

	// 2. 更新上下文
	if err := e.updateWorkflowInstByDecideResult(inst, decideResult); err != nil {
		return err
	}

	if decideResult.InstStatus.IsCompleted() || decideResult.InstStatus.IsTerminal() {
		return nil
	}

	// 3. 发送驱动消息
	return e.workflowUpdater.SendNodeScheduleDriveEvent(inst, e.getNodesToBeScheduled(decideResult))
}

func (e *DefaultWorkflowExecutor) handleForInstAlreadyPaused(inst *entity.WorkflowInst) error {
	if inst.CurNodeInst == nil {
		return nil
	}

	inst.RunCompletedNodeInstIDsAfterPaused = buildRunCompletedNodeInstIDsAfterPaused(inst)
	return e.workflowUpdater.UpdateWorkflowInst(inst)
}

func buildRunCompletedNodeInstIDsAfterPaused(inst *entity.WorkflowInst) []string {
	if len(inst.RunCompletedNodeInstIDsAfterPaused) == 0 {
		return []string{inst.CurNodeInst.NodeInstID}
	}

	return append(inst.RunCompletedNodeInstIDsAfterPaused, inst.CurNodeInst.NodeInstID)
}

func (e *DefaultWorkflowExecutor) updateWorkflowInstByDecideResult(inst *entity.WorkflowInst,
	decideResult *entity.DecideResult) error {
	if err := e.updateNodeInsts(decideResult); err != nil {
		log.Errorf("[%s]Failed to update node insts, caused by %s",
			logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID), err)
		return err
	}

	if err := e.creatNodeInsts(decideResult); err != nil {
		log.Errorf("Failed to create node insts, caused by %s", err)
		return err
	}

	appendExecutePath(inst, decideResult)

	inst.Status = decideResult.InstStatus
	if decideResult.InstStatus == entity.InstFailed {
		inst.Reason.FailedRootCause = *decideResult.InstFailedRootCause
	}

	return e.workflowUpdater.UpdateWorkflowInstWithStatus(inst)
}

func (e *DefaultWorkflowExecutor) creatNodeInsts(result *entity.DecideResult) error {
	for _, nodeInst := range result.NodesToBeScheduled {
		createNodeInstDTO, err := convertor.NodeInstConvertor.ConvertEntityToCreateDTO(nodeInst)
		if err != nil {
			return err
		}
		nodeInstID, err := e.nodeInstRepo.Create(createNodeInstDTO)
		if err != nil {
			return err
		}
		nodeInst.NodeInstID = nodeInstID
	}
	return nil
}

func (e *DefaultWorkflowExecutor) updateNodeInsts(result *entity.DecideResult) error {
	for _, nodeInst := range result.NodesToBeUpdated {
		if err := e.workflowUpdater.UpdateNodeInstWithStatus(nodeInst); err != nil {
			return err
		}
	}
	return nil
}

func (e *DefaultWorkflowExecutor) getNodesToBeScheduled(decideResult *entity.DecideResult) []string {
	var nodesToBeScheduled []string
	for _, nodeInst := range decideResult.NodesToBeScheduled {
		nodesToBeScheduled = append(nodesToBeScheduled, nodeInst.NodeInstID)
	}
	return nodesToBeScheduled
}
