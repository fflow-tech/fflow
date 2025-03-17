package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowInstCommandService 写服务
// 上层RPC/HTTP接口对接的是这个服务，不和底层的workflowRunner对接
type WorkflowInstCommandService struct {
	workflowDefRepo  ports.WorkflowDefRepository
	workflowInstRepo ports.WorkflowInstRepository
	eventBusRepo     ports.EventBusRepository
	cacheRepo        ports.CacheRepository
	workflowRunner   execution.WorkflowRunner
	timeoutChecker   execution.TimeoutChecker
	historyArchiver  execution.HistoryArchiver
}

// NewWorkflowInstCommandService 新建服务
func NewWorkflowInstCommandService(repoProviderSet *ports.RepoProviderSet,
	timeoutChecker *execution.DefaultTimeoutChecker,
	historyArchiver *execution.DefaultHistoryArchiver,
	workflowRunner *execution.DefaultWorkflowRunner) *WorkflowInstCommandService {
	return &WorkflowInstCommandService{
		workflowDefRepo:  repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		eventBusRepo:     repoProviderSet.EventBusRepo(),
		cacheRepo:        repoProviderSet.CacheRepo(),
		workflowRunner:   workflowRunner,
		timeoutChecker:   timeoutChecker,
		historyArchiver:  historyArchiver,
	}
}

// StartWorkflowInst 创建工作流实例
func (m *WorkflowInstCommandService) StartWorkflowInst(ctx context.Context,
	req *dto.StartWorkflowInstDTO) (string, error) {
	return m.workflowRunner.Start(ctx, req)
}

// ConsumeWorkflowStartDriveEvent 消费驱动事件
func (m *WorkflowInstCommandService) ConsumeWorkflowStartDriveEvent(ctx context.Context, d *dto.DriveEventDTO) error {
	msg := d.Message
	driveEvent := event.WorkflowStartDriveEvent{}
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
	driveReq := &dto.DriveWorkflowInstNextNodesDTO{
		InstID:         driveEvent.InstID,
		DefID:          driveEvent.DefID,
		FromResumeInst: driveEvent.FromResumeInst,
	}
	return m.workflowRunner.DriveNext(ctx, driveReq)
}

// ConsumeNodeCompleteDriveEvent 消费驱动事件
func (m *WorkflowInstCommandService) ConsumeNodeCompleteDriveEvent(ctx context.Context, req *dto.DriveEventDTO) error {
	msg := req.Message
	driveEvent := event.NodeCompleteDriveEvent{}
	if err := json.Unmarshal(msg.Payload(), &driveEvent); err != nil {
		return err
	}
	inst, err := m.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(driveEvent.InstID, driveEvent.DefID, ""))
	if err != nil {
		return err
	}
	if driveEvent.EventTime.Before(inst.LastRestartAt) {
		log.Warnf("[%s]The event is expired, skip it, eventTime:%s, event:%s",
			logs.GetFlowTraceID(driveEvent.DefID, driveEvent.InstID),
			driveEvent.EventTime, utils.StructToJsonStr(driveEvent))
		return nil
	}

	driveReq := &dto.DriveWorkflowInstNextNodesDTO{
		InstID:         driveEvent.InstID,
		DefID:          driveEvent.DefID,
		CurNodeInstID:  driveEvent.NodeInstID,
		FromResumeInst: driveEvent.FromResumeInst,
	}

	return m.workflowRunner.DriveNext(ctx, driveReq)
}

func eventNeedIgnore(inst *entity.WorkflowInst, eventTime time.Time) bool {
	// 如果流程已经不在执行态了, 直接忽略掉消息
	if inst.Status != entity.InstRunning {
		return true
	}

	return eventTime.Before(inst.LastRestartAt)
}

// RestartWorkflowInst 重启
func (m *WorkflowInstCommandService) RestartWorkflowInst(ctx context.Context, req *dto.RestartWorkflowInstDTO) error {
	return m.workflowRunner.Restart(ctx, req)
}

// CompleteWorkflowInst 标记完成
func (m *WorkflowInstCommandService) CompleteWorkflowInst(ctx context.Context, req *dto.CompleteWorkflowInstDTO) error {
	return m.workflowRunner.Complete(ctx, req)
}

// UpdateWorkflowInstCtx 更新流程实例上下文
func (m *WorkflowInstCommandService) UpdateWorkflowInstCtx(ctx context.Context,
	req *dto.UpdateWorkflowInstCtxDTO) error {
	return m.workflowRunner.UpdateCtx(ctx, req)
}

// CancelWorkflowInst 取消
func (m *WorkflowInstCommandService) CancelWorkflowInst(ctx context.Context, req *dto.CancelWorkflowInstDTO) error {
	return m.workflowRunner.Cancel(ctx, req)
}

// PauseWorkflowInst 暂停
func (m *WorkflowInstCommandService) PauseWorkflowInst(ctx context.Context, req *dto.PauseWorkflowInstDTO) error {
	return m.workflowRunner.Pause(ctx, req)
}

// ResumeWorkflowInst 恢复
func (m *WorkflowInstCommandService) ResumeWorkflowInst(ctx context.Context, req *dto.ResumeWorkflowInstDTO) error {
	return m.workflowRunner.Resume(ctx, req)
}

// CheckTimeout 检查超时
func (m *WorkflowInstCommandService) CheckTimeout(ctx context.Context) error {
	return m.timeoutChecker.CheckAll()
}

// ArchiveHistoryWorkflowInsts 归档流程实例
func (m *WorkflowInstCommandService) ArchiveHistoryWorkflowInsts(ctx context.Context,
	req *dto.ArchiveHistoryWorkflowInstsDTO) error {
	return m.historyArchiver.Archive(req)
}

// DebugWorkflowInst 调试流程
func (m *WorkflowInstCommandService) DebugWorkflowInst(ctx context.Context,
	req *dto.DebugWorkflowInstDTO) error {
	return m.workflowRunner.Debug(ctx, req)
}
