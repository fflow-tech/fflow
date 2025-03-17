package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// NodeInstCommandService 写服务
// 上层RPC/HTTP接口对接的是这个服务，不和底层的Executor对接
type NodeInstCommandService struct {
	workflowDefRepo  ports.WorkflowDefRepository
	workflowInstRepo ports.WorkflowInstRepository
	nodeInstRepo     ports.NodeInstRepository
	cacheRepo        ports.CacheRepository
	remoteRepo       ports.RemoteRepository
	exprEvaluator    expr.Evaluator
	nodeRunner       execution.NodeRunner
}

// NewNodeInstCommandService 新建服务
func NewNodeInstCommandService(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *execution.WorkflowProviderSet,
	nodeRunner *execution.DefaultNodeRunner) *NodeInstCommandService {
	r := &NodeInstCommandService{
		workflowDefRepo:  repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:     repoProviderSet.NodeInstRepo(),
		cacheRepo:        repoProviderSet.CacheRepo(),
		exprEvaluator:    workflowProviderSet.ExprEvaluator(),
		nodeRunner:       nodeRunner,
	}
	return r
}

// ConsumeNodeScheduleDriveEvent 消费节点的调度事件
func (m *NodeInstCommandService) ConsumeNodeScheduleDriveEvent(ctx context.Context, req *dto.DriveEventDTO) error {
	msg := req.Message
	driveEvent := event.NodeScheduleDriveEvent{}
	if err := json.Unmarshal(msg.Payload(), &driveEvent); err != nil {
		return err
	}

	inst, err := m.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(driveEvent.InstID, driveEvent.DefID, ""))
	if err != nil {
		return err
	}
	if eventNeedIgnore(inst, driveEvent.EventTime) {
		log.Warnf("[%s]The event is expired, skip it, eventTime:%s, event:%s",
			logs.GetFlowTraceID(driveEvent.DefID, driveEvent.InstID),
			driveEvent.EventTime, utils.StructToJsonStr(driveEvent))
		return nil
	}

	for _, nodeInstID := range driveEvent.NodeInstIDs {
		scheduleDTO := &dto.ScheduleNodeDTO{
			DefID:      driveEvent.DefID,
			DefVersion: driveEvent.DefVersion,
			InstID:     driveEvent.InstID,
			NodeInstID: nodeInstID,
		}
		if err := m.nodeRunner.Schedule(ctx, scheduleDTO); err != nil {
			log.Errorf("[%s]Failed to schedule one node inst, caused by %s",
				logs.GetFlowTraceID(driveEvent.DefID, driveEvent.InstID), err)
			return err
		}
	}

	return nil
}

// ConsumeNodeExecuteDriveEvent 消费节点执行事件
func (m *NodeInstCommandService) ConsumeNodeExecuteDriveEvent(ctx context.Context, req *dto.DriveEventDTO) error {
	msg := req.Message
	driveEvent := event.NodeExecuteDriveEvent{}
	if err := json.Unmarshal(msg.Payload(), &driveEvent); err != nil {
		return err
	}

	if driveEvent.InstID == "" {
		return fmt.Errorf("illegal workflow instID=0")
	}

	lock, err := execution.GetInstDistributeLock(m.cacheRepo, driveEvent.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	inst, err := m.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(driveEvent.InstID, driveEvent.DefID, ""))
	if err != nil {
		return err
	}

	if eventNeedIgnoreExecuteDriveEvent(inst, driveEvent.EventTime) {
		log.Warnf("[%s]The event is expired, skip it, eventTime:%s, event:%s",
			logs.GetFlowTraceID(driveEvent.DefID, driveEvent.InstID),
			driveEvent.EventTime, utils.StructToJsonStr(driveEvent))
		return nil
	}

	// 如果流程处于暂停状态, 需要把要执行的节点保存下来
	if inst.Status == entity.InstPaused {
		inst.WaitCompletedNodeInstIDsAfterPaused = append(inst.WaitCompletedNodeInstIDsAfterPaused,
			driveEvent.NodeInstID)
		updateDTO, err := convertor.InstConvertor.ConvertEntityToUpdateDTO(inst)
		if err != nil {
			log.Errorf("Failed to convert entity to update dto, caused by %s", err)
			return err
		}
		return m.workflowInstRepo.UpdateWithDefID(updateDTO)
	}

	runDTO := &dto.RunNodeDTO{
		DefID:      driveEvent.DefID,
		InstID:     driveEvent.InstID,
		NodeInstID: driveEvent.NodeInstID,
	}
	return m.nodeRunner.Run(ctx, runDTO)
}

func eventNeedIgnoreExecuteDriveEvent(inst *entity.WorkflowInst, eventTime time.Time) bool {
	if eventTime.Before(inst.LastRestartAt) {
		return true
	}

	// 如果流程已经不在执行态且也不是暂停的情况, 直接忽略掉消息
	return inst.Status != entity.InstRunning && inst.Status != entity.InstPaused
}

// ConsumeNodePollDriveEvent 消费节点轮询事件
func (m *NodeInstCommandService) ConsumeNodePollDriveEvent(ctx context.Context, req *dto.DriveEventDTO) error {
	msg := req.Message
	driveEvent := event.NodePollDriveEvent{}
	if err := json.Unmarshal(msg.Payload(), &driveEvent); err != nil {
		return err
	}

	inst, err := m.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(driveEvent.InstID, driveEvent.DefID, ""))
	if err != nil {
		return err
	}
	if msg.EventTime().Before(inst.LastRestartAt) {
		log.Warnf("[%s]The event is expired, skip it, eventTime:%s, event:%s",
			logs.GetFlowTraceID(driveEvent.DefID, driveEvent.InstID),
			driveEvent.EventTime, utils.StructToJsonStr(driveEvent))
		return nil
	}

	pollingDTO := &dto.PollingNodeDTO{
		DefID:      driveEvent.DefID,
		DefVersion: driveEvent.DefVersion,
		InstID:     driveEvent.InstID,
		NodeInstID: driveEvent.NodeInstID,
	}
	return m.nodeRunner.Polling(ctx, pollingDTO)
}

// ConsumeNodeRetryDriveEvent 消费节点重跑事件
func (m *NodeInstCommandService) ConsumeNodeRetryDriveEvent(ctx context.Context, req *dto.DriveEventDTO) error {
	msg := req.Message
	driveEvent := event.NodeRetryDriveEvent{}
	if err := json.Unmarshal(msg.Payload(), &driveEvent); err != nil {
		return err
	}

	inst, err := m.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(driveEvent.InstID, driveEvent.DefID, ""))
	if err != nil {
		return err
	}
	if eventNeedIgnore(inst, driveEvent.EventTime) {
		log.Warnf("[%s]The event is expired, skip it, eventTime:%s, event:%s",
			logs.GetFlowTraceID(driveEvent.DefID, driveEvent.InstID),
			driveEvent.EventTime, utils.StructToJsonStr(driveEvent))
		return nil
	}

	runDTO := &dto.RunNodeDTO{
		DefID:      driveEvent.DefID,
		InstID:     driveEvent.InstID,
		NodeInstID: driveEvent.NodeInstID,
	}
	return m.nodeRunner.Run(ctx, runDTO)
}

// RerunNode 重跑
func (m *NodeInstCommandService) RerunNode(ctx context.Context, req *dto.RerunNodeDTO) error {
	return m.nodeRunner.Rerun(ctx, req)
}

// SkipNode 跳过
func (m *NodeInstCommandService) SkipNode(ctx context.Context, req *dto.SkipNodeDTO) error {
	return m.nodeRunner.Skip(ctx, req)
}

// CancelSkipNode 取消跳过
func (m *NodeInstCommandService) CancelSkipNode(ctx context.Context, req *dto.CancelSkipNodeDTO) error {
	return m.nodeRunner.CancelSkip(ctx, req)
}

// ResumeNode 恢复节点执行
func (m *NodeInstCommandService) ResumeNode(ctx context.Context, req *dto.ResumeNodeDTO) error {
	return m.nodeRunner.Resume(ctx, req)
}

// CancelNode 取消节点执行
func (m *NodeInstCommandService) CancelNode(ctx context.Context, req *dto.CancelNodeDTO) error {
	return m.nodeRunner.Cancel(ctx, req)
}

// CompleteNode 标记节点完成
func (m *NodeInstCommandService) CompleteNode(ctx context.Context, req *dto.CompleteNodeDTO) error {
	return m.nodeRunner.Complete(ctx, req)
}

// SetTimeout 标记节点超时
func (m *NodeInstCommandService) SetTimeout(ctx context.Context, req *dto.SetNodeTimeoutDTO) error {
	return m.nodeRunner.SetTimeout(ctx, req)
}

// SetNearTimeout 标记节点接近超时
func (m *NodeInstCommandService) SetNearTimeout(ctx context.Context, req *dto.SetNodeNearTimeoutDTO) error {
	return m.nodeRunner.SetNearTimeout(ctx, req)
}
