package execution

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/trigger"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowUpdater 流程更新者
type WorkflowUpdater interface {
	SendWorkflowDefExternalEvent(namespace string, defID string, eventType event.ExternalEventType) error
	SendWorkflowInstExternalEvent(inst *entity.WorkflowInst, eventType event.ExternalEventType) error
	SendNodeExternalEvent(nodeInst *entity.NodeInst, eventType event.ExternalEventType) error
	SendWorkflowStartDriveEvent(namespace string, defID string, instID string, fromResumeInst bool) error
	SendNodeScheduleDriveEvent(inst *entity.WorkflowInst, nodesToBeScheduled []string) error
	SendNodeCompleteDriveEvent(nodeInst *entity.NodeInst, fromResumeInst bool) error
	SendDelayNodeExecuteDriveEvent(nodeInst *entity.NodeInst, deliverAfter time.Duration) error
	SendPresetNodeExecuteDriveEvent(nodeInst *entity.NodeInst, deliverAt time.Time) error
	SendNodeDriveEvent(nodeInst *entity.NodeInst, eventType event.DriveEventType) error
	SendDelayNodeDriveEvent(nodeInst *entity.NodeInst, deliverAfter time.Duration, eventType event.DriveEventType) error
	UpdateNodeInstWithStatus(nodeInst *entity.NodeInst) error
	UpdateWorkflowInstWithStatus(inst *entity.WorkflowInst) error
	CreateWorkflowInst(req *dto.StartWorkflowInstDTO, workflowDef *entity.WorkflowDef) (string, error)
	UpdateWorkflowInst(inst *entity.WorkflowInst) error
	RegisterTriggers(workflowDef *entity.WorkflowDef, instID string, operator string) error
	UnRegisterTriggers(inst *entity.WorkflowInst) error
	UpdateWorkflowInstFailed(req *dto.UpdateWorkflowInstFailedDTO) error
}

// DefaultWorkflowUpdater 流程更新
type DefaultWorkflowUpdater struct {
	workflowInstRepo ports.WorkflowInstRepository
	triggerRegistry  trigger.Registry
	nodeInstRepo     ports.NodeInstRepository
	eventBusRepo     ports.EventBusRepository
	exprEvaluator    expr.Evaluator
}

// NewDefaultWorkflowUpdater 新建更新者
func NewDefaultWorkflowUpdater(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet,
	triggerRegistry *trigger.DefaultRegistry) *DefaultWorkflowUpdater {
	return &DefaultWorkflowUpdater{
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:     repoProviderSet.NodeInstRepo(),
		eventBusRepo:     repoProviderSet.EventBusRepo(),
		exprEvaluator:    workflowProviderSet.ExprEvaluator(),
		triggerRegistry:  triggerRegistry,
	}
}

// SendWorkflowStartDriveEvent 发送流程启动驱动事件
func (w *DefaultWorkflowUpdater) SendWorkflowStartDriveEvent(namespace, defID, instID string, fromResumeInst bool) error {
	driveEvent := event.WorkflowStartDriveEvent{
		BasicEvent:     event.NewDriveBasicEvent(event.WorkflowStartDrive, namespace, defID, instID),
		InstID:         instID,
		DefID:          defID,
		FromResumeInst: fromResumeInst,
	}
	return w.eventBusRepo.SendDriveEvent(context.Background(), driveEvent)
}

// SendWorkflowDefExternalEvent 发送外部事件 现在定义的结构相同 为了后续的处理还是分开处理
func (w *DefaultWorkflowUpdater) SendWorkflowDefExternalEvent(namespace string,
	defID string, eventType event.ExternalEventType) error {
	var sendMsg interface{}
	basicEvent := event.NewExternalBasicEvent(eventType, namespace, defID, "", nil)
	switch eventType {
	case event.DefCreate:
		sendMsg = event.DefCreateEvent{
			DefID:      defID,
			BasicEvent: basicEvent,
		}
	case event.DefUpdate:
		sendMsg = event.DefUpdateEvent{
			DefID:      defID,
			BasicEvent: basicEvent,
		}
	case event.DefEnable:
		sendMsg = event.DefEnableEvent{
			DefID:      defID,
			BasicEvent: basicEvent,
		}
	case event.DefDisable:
		sendMsg = event.DefDisableEvent{
			DefID:      defID,
			BasicEvent: basicEvent,
		}
	default:
		return fmt.Errorf("illegal eventType=%s", eventType)
	}
	return w.eventBusRepo.SendExternalEvent(context.Background(), sendMsg)
}

// SendWorkflowInstExternalEvent 发送流程外部事件
func (w *DefaultWorkflowUpdater) SendWorkflowInstExternalEvent(inst *entity.WorkflowInst,
	eventType event.ExternalEventType) error {
	externalEvent, err := convertor.EventConvertor.ConvertEntityToWorkflowExternalEvent(inst, "", eventType)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendExternalEvent(context.Background(), externalEvent)
}

// SendWorkflowStartEvent 发送流程外部事件
func (w *DefaultWorkflowUpdater) SendWorkflowStartEvent(inst *entity.WorkflowInst) error {
	externalEvent, err := convertor.EventConvertor.ConvertEntityToWorkflowStartEvent(inst)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendExternalEvent(context.Background(), externalEvent)
}

// SendWorkflowExternalEventWithNodeRefName 发送流程外部事件
func (w *DefaultWorkflowUpdater) SendWorkflowExternalEventWithNodeRefName(inst *entity.WorkflowInst, nodeRefName string,
	eventType event.ExternalEventType) error {
	externalEvent, err := convertor.EventConvertor.ConvertEntityToWorkflowExternalEvent(inst, nodeRefName, eventType)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendExternalEvent(context.Background(), externalEvent)
}

// SendNodeExternalEvent 发送节点外部事件
func (w *DefaultWorkflowUpdater) SendNodeExternalEvent(nodeInst *entity.NodeInst,
	eventType event.ExternalEventType) error {
	externalEvent, err := convertor.EventConvertor.ConvertEntityToNodeExternalEvent(nodeInst, eventType)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendExternalEvent(context.Background(), externalEvent)
}

// SendNodeScheduleDriveEvent 发送节点被调度驱动事件
func (w *DefaultWorkflowUpdater) SendNodeScheduleDriveEvent(inst *entity.WorkflowInst,
	nodesToBeScheduled []string) error {
	driveEvent, err := convertor.EventConvertor.ConvertEntityToNodeScheduleDriveEvent(inst, nodesToBeScheduled)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendDriveEvent(context.Background(), driveEvent)
}

// SendNodeCompleteDriveEvent 发送节点执行完成驱动事件
func (w *DefaultWorkflowUpdater) SendNodeCompleteDriveEvent(nodeInst *entity.NodeInst, fromResumeInst bool) error {
	driveEvent, err := convertor.EventConvertor.ConvertEntityToNodeCompleteDriveEvent(nodeInst, fromResumeInst)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendDriveEvent(context.Background(), driveEvent)
}

// SendNodeDriveEvent 发送节点驱动事件
func (w *DefaultWorkflowUpdater) SendNodeDriveEvent(nodeInst *entity.NodeInst, eventType event.DriveEventType) error {
	driveEvent, err := convertor.EventConvertor.ConvertEntityToNodeDriveEvent(nodeInst, eventType)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendDriveEvent(context.Background(), driveEvent)
}

// SendDelayNodeDriveEvent 发送延迟节点驱动事件
func (w *DefaultWorkflowUpdater) SendDelayNodeDriveEvent(nodeInst *entity.NodeInst,
	deliverAfter time.Duration, eventType event.DriveEventType) error {
	driveEvent, err := convertor.EventConvertor.ConvertEntityToNodeDriveEvent(nodeInst, eventType)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendDelayDriveEvent(context.Background(), deliverAfter, driveEvent)
}

// SendDelayNodeExecuteDriveEvent 发送节点延迟执行驱动事件
func (w *DefaultWorkflowUpdater) SendDelayNodeExecuteDriveEvent(nodeInst *entity.NodeInst,
	deliverAfter time.Duration) error {
	driveEvent, err := convertor.EventConvertor.ConvertEntityToNodeExecuteDriveEvent(nodeInst)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendDelayDriveEvent(context.Background(), deliverAfter, driveEvent)
}

// SendPresetNodeExecuteDriveEvent 发送节点定时执行驱动事件
func (w *DefaultWorkflowUpdater) SendPresetNodeExecuteDriveEvent(nodeInst *entity.NodeInst,
	deliverAt time.Time) error {
	driveEvent, err := convertor.EventConvertor.ConvertEntityToNodeExecuteDriveEvent(nodeInst)
	if err != nil {
		return err
	}
	return w.eventBusRepo.SendPresetDriveEvent(context.Background(), deliverAt, driveEvent)
}

// UpdateNodeInstWithStatus 更新节点实例的状态
func (w *DefaultWorkflowUpdater) UpdateNodeInstWithStatus(nodeInst *entity.NodeInst) error {
	switch nodeInst.Status {
	case entity.NodeInstSucceed:
		nodeInst.CompletedAt = time.Now()
		return w.updateNodeInstWithExternalEventType(nodeInst, event.NodeSuccess)
	case entity.NodeInstFailed:
		nodeInst.CompletedAt = time.Now()
		return w.updateNodeInstWithExternalEventType(nodeInst, event.NodeFail)
	case entity.NodeInstCancelled:
		nodeInst.CompletedAt = time.Now()
		return w.updateNodeInstWithExternalEventType(nodeInst, event.NodeCancel)
	case entity.NodeInstTimeout:
		nodeInst.CompletedAt = time.Now()
		return w.updateNodeInstWithExternalEventType(nodeInst, event.NodeTimeout)
	case entity.NodeInstRunning:
		return w.updateNodeInst(nodeInst)
	case entity.NodeInstWaiting:
		return w.updateNodeInstForWaitStatus(nodeInst)
	case entity.NodeInstScheduled:
		return w.updateNodeInst(nodeInst)
	default:
		return fmt.Errorf("illegal node inst status:%s", nodeInst.Status)
	}
}

// UpdateWorkflowInstWithStatus 更新流程实例状态
func (w *DefaultWorkflowUpdater) UpdateWorkflowInstWithStatus(inst *entity.WorkflowInst) error {
	switch inst.Status {
	case entity.InstSucceed:
		return w.updateWorkflowInstWithExternalEventType(inst, event.WorkflowSuccess)
	case entity.InstFailed:
		return w.updateWorkflowInstWithExternalEventType(inst, event.WorkflowFail)
	case entity.InstRunning:
		return w.UpdateWorkflowInst(inst)
	case entity.InstPaused:
		return w.updateWorkflowInstWithExternalEventType(inst, event.WorkflowPause)
	case entity.InstTimeout:
		return w.updateWorkflowInstWithExternalEventType(inst, event.WorkflowTimeout)
	case entity.InstCancelled:
		return w.updateWorkflowInstWithExternalEventType(inst, event.WorkflowCancel)
	default:
		return fmt.Errorf("illegal workflow inst status:%s", inst.Status)
	}
}

// CreateWorkflowInst 创建流程实例
func (w *DefaultWorkflowUpdater) CreateWorkflowInst(req *dto.StartWorkflowInstDTO,
	workflowDef *entity.WorkflowDef) (string, error) {
	inst := &entity.WorkflowInst{
		WorkflowDef:      workflowDef,
		ParentInstID:     req.ParentInstID,
		ParentNodeInstID: req.ParentNodeInstID,
		Name:             req.Name,
		Creator:          req.Creator,
		Status:           entity.InstRunning,
		StartAt:          time.Now(),
		Input:            req.Input,
		Variables:        workflowDef.Variables,
		Biz:              workflowDef.Biz,
		Owner:            &workflowDef.Owner,
		Reason:           &entity.InstReason{StartReason: req.Reason},
	}

	if req.DebugMode {
		inst.CurDebugMode = entity.SingleStepMode
	}

	repoDTO, err := convertor.InstConvertor.ConvertEntityToCreateRepoDTO(inst)
	if err != nil {
		return "", err
	}
	instID, err := w.workflowInstRepo.Create(repoDTO)
	if err != nil {
		return "", err
	}

	err = w.RegisterTriggers(workflowDef, instID, req.Creator)
	if err != nil {
		return "", err
	}

	inst.InstID = instID
	if err := w.SendWorkflowStartEvent(inst); err != nil {
		return "", err
	}

	return inst.InstID, nil
}

// RegisterTriggers 注册触发器
func (w *DefaultWorkflowUpdater) RegisterTriggers(workflowDef *entity.WorkflowDef,
	instID string, operator string) error {
	registerTriggerDTO := &dto.RegisterTriggersDTO{
		DefID:      workflowDef.DefID,
		DefVersion: workflowDef.Version,
		InstID:     instID,
		Operator:   operator,
		Level:      entity.InstTrigger,
		Triggers:   workflowDef.Triggers,
	}

	return w.triggerRegistry.Register(registerTriggerDTO)
}

// UnRegisterTriggers 反注册触发器
func (w *DefaultWorkflowUpdater) UnRegisterTriggers(inst *entity.WorkflowInst) error {
	unRegisterTriggersDTO := &dto.UnRegisterTriggersDTO{
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: inst.WorkflowDef.Version,
		InstID:     inst.InstID,
		Level:      entity.InstTrigger,
	}
	return w.triggerRegistry.UnRegister(unRegisterTriggersDTO)
}

func (w *DefaultWorkflowUpdater) updateWorkflowInstWithExternalEventType(inst *entity.WorkflowInst,
	eventType event.ExternalEventType) error {
	if err := w.UpdateWorkflowInst(inst); err != nil {
		return err
	}
	if inst.ParentInstID != "" && inst.ParentNodeInstID != "" {
		// 当流程实例有对应的父流程实例和节点实例时，同时更新节点状态
		if err := w.updateNodeInstForSubWorkflow(inst); err != nil {
			return err
		}
	}
	if inst.Status.IsTerminal() {
		if err := w.UnRegisterTriggers(inst); err != nil {
			return err
		}
	}

	if inst.PreStatus == inst.Status {
		return nil
	}

	return w.SendWorkflowInstExternalEvent(inst, eventType)
}

// updateNodeInstForSubWorkflow 根据流程状态更新对应的父流程中节点状态
func (w *DefaultWorkflowUpdater) updateNodeInstForSubWorkflow(inst *entity.WorkflowInst) error {
	// 获取更新后的节点实例和对应的事件类型
	nodeInst, nodeEventType, err := w.getNodeInstAndEventTypeForUpdate(inst)
	if err != nil {
		return err
	}
	// 节点实例为空时直接返回不需要发送事件
	if nodeInst == nil {
		return nil
	}
	// 更新节点实例并发送外部事件
	if err := w.updateNodeInstWithExternalEventType(nodeInst, nodeEventType); err != nil {
		return err
	}
	// 节点状态为成功或失败时，发送节点完成驱动事件
	if inst.Status == entity.InstSucceed || inst.Status == entity.InstFailed {
		return w.SendNodeDriveEvent(nodeInst, event.NodeCompleteDrive)
	}
	return nil
}

// getNodeInstAndEventTypeForUpdate 获取节点实例和事件类型
func (w *DefaultWorkflowUpdater) getNodeInstAndEventTypeForUpdate(inst *entity.WorkflowInst) (*entity.NodeInst,
	event.ExternalEventType, error) {
	// 获取子流程对应的父流程中的节点实例
	nodeInst, err := w.nodeInstRepo.Get(&dto.GetNodeInstDTO{
		DefID:      inst.WorkflowDef.ParentDefID,
		InstID:     inst.ParentInstID,
		NodeInstID: inst.ParentNodeInstID,
	})
	if err != nil {
		return nil, "", err
	}
	// 将流程状态和节点状态做映射
	var nodeEventType event.ExternalEventType
	// 子流程的实例 ID 和定义 ID
	instID := inst.InstID
	defID := inst.WorkflowDef.DefID
	nodeInst.Input = inst.Input
	nodeInst.Output = inst.Output
	switch inst.Status {
	case entity.InstSucceed:
		nodeInst.CompletedAt = time.Now()
		nodeInst.Status = entity.NodeInstSucceed
		nodeInst.Reason.SucceedReason = fmt.Sprintf("subworkflow[instID: %s, DefID: %s] succeed", instID, defID)
		nodeEventType = event.NodeSuccess
	case entity.InstFailed:
		nodeInst.CompletedAt = time.Now()
		nodeInst.Status = entity.NodeInstFailed
		nodeInst.Reason.FailedReason = fmt.Sprintf("subworkflow[instID: %s, DefID: %s,] failed, caused by %s",
			instID, defID, inst.Reason.FailedRootCause.FailedReason)
		nodeEventType = event.NodeFail
	case entity.InstCancelled:
		nodeInst.CompletedAt = time.Now()
		nodeInst.Status = entity.NodeInstCancelled
		nodeInst.Reason.CancelledReason = fmt.Sprintf("subworkflow[instID: %s, DefID: %s,] cancelled", instID, defID)
		nodeEventType = event.NodeCancel
	case entity.InstTimeout:
		nodeInst.CompletedAt = time.Now()
		nodeInst.Status = entity.NodeInstTimeout
		nodeInst.Reason.TimeoutReason = fmt.Sprintf("subworkflow[instID: %s, DefID: %s,] timeout", instID, defID)
		nodeEventType = event.NodeTimeout
	default:
		// 对于其他的流程状态，不需要更新节点，直接忽略
		log.Infof("ignored inst status to node inst status:%s", inst.Status)
		return nil, "", nil
	}

	return nodeInst, nodeEventType, nil
}

// UpdateWorkflowInst 更新流程实例
func (w *DefaultWorkflowUpdater) UpdateWorkflowInst(inst *entity.WorkflowInst) error {
	updateDTO, err := convertor.InstConvertor.ConvertEntityToUpdateDTO(inst)
	if err != nil {
		log.Errorf("Failed to convert entity to update dto, caused by %s", err)
		return err
	}

	return w.workflowInstRepo.UpdateWithDefID(updateDTO)
}

func (w *DefaultWorkflowUpdater) updateNodeInstWithExternalEventType(nodeInst *entity.NodeInst,
	eventType event.ExternalEventType) error {
	if err := w.updateNodeInst(nodeInst); err != nil {
		return err
	}

	err := w.updateWorkflowInstOutputIfNodeHasReturn(nodeInst)
	if err != nil {
		return err
	}

	if nodeInst.PreStatus == nodeInst.Status {
		return nil
	}

	return w.SendNodeExternalEvent(nodeInst, eventType)
}

func (w *DefaultWorkflowUpdater) updateWorkflowInstOutputIfNodeHasReturn(nodeInst *entity.NodeInst) error {
	if !nodeInst.Status.IsTerminal() || len(nodeInst.BasicNodeDef.Return) <= 0 {
		return nil
	}
	inst, err := w.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(nodeInst.InstID, nodeInst.DefID, ""))
	if err != nil {
		return fmt.Errorf("failed to update workflow inst output: %w", err)
	}
	err = w.setWorkflowInstOutput(inst, nodeInst)
	if err != nil {
		inst.Reason.FailedRootCause.FailedReason = fmt.Sprintf("failed to evaluate inst output, caused by %s", err)
		inst.Status = entity.InstFailed
	}

	return w.UpdateWorkflowInst(inst)
}

func (w *DefaultWorkflowUpdater) setWorkflowInstOutput(inst *entity.WorkflowInst, nodeInst *entity.NodeInst) error {
	ctx, err := entity.ConvertToCtx(inst)
	if err != nil {
		return err
	}
	if err := entity.AppendNodeInfoToCtxKey(ctx, nodeInst, constants.ThisNode); err != nil {
		return err
	}
	inst.Output, err = w.exprEvaluator.EvaluateMap(ctx, nodeInst.BasicNodeDef.Return)
	return err
}

func (w *DefaultWorkflowUpdater) updateNodeInst(nodeInst *entity.NodeInst) error {
	// 每次更新前将biz赋值, 如果计算失败, 忽略掉错误
	if err := w.setNodeInstBiz(nodeInst); err != nil {
		log.Errorf("Failed to set node inst biz, skip this err, caused by %s", err)
	}

	updateNodeInstDTO, err := convertor.NodeInstConvertor.ConvertEntityToUpdateDTO(nodeInst)
	if err != nil {
		return err
	}

	return w.nodeInstRepo.UpdateWithDefID(updateNodeInstDTO)
}

// setNodeInstBiz 设置节点实例中的业务数据
func (w *DefaultWorkflowUpdater) setNodeInstBiz(nodeInst *entity.NodeInst) error {
	inst, err := w.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(nodeInst.InstID, nodeInst.DefID, ""))
	if err != nil {
		return err
	}
	ctx, err := entity.ConvertToCtx(inst)
	if err != nil {
		return err
	}
	nodeInst.Biz, err = w.exprEvaluator.EvaluateMap(ctx, nodeInst.BasicNodeDef.Biz)
	if err != nil {
		return err
	}
	// 增加节点的输入输出信息到业务字段中
	nodeInst.Biz["input"] = nodeInst.Input
	nodeInst.Biz["output"] = nodeInst.Output
	return nil
}

func (w *DefaultWorkflowUpdater) updateNodeInstForWaitStatus(nodeInst *entity.NodeInst) error {
	if err := w.updateNodeInst(nodeInst); err != nil {
		return err
	}

	if nodeInst.PreStatus == nodeInst.Status {
		return nil
	}

	// 如果是因为调试设置的等待, 必须等待用户使用 Resume 操作恢复
	if nodeInst.WaitForDebug {
		return nil
	}

	// 配置为 max 的时候, 必须等待用户使用 Resume 操作恢复
	if strings.EqualFold(nodeInst.BasicNodeDef.Wait.Duration, "max") {
		return nil
	}

	// 优先使用 duration 的配置
	if nodeInst.BasicNodeDef.Wait.Duration != "" {
		deliveryAfter, err := expr.ParseDuration(nodeInst.BasicNodeDef.Wait.Duration)
		if err != nil {
			return err
		}
		return w.SendDelayNodeExecuteDriveEvent(nodeInst, deliveryAfter)
	}

	deliveryAt, err := utils.GetNextTimeByExpr(nodeInst.BasicNodeDef.Wait.Expr, time.Now())
	if err != nil {
		return err
	}
	return w.SendPresetNodeExecuteDriveEvent(nodeInst, deliveryAt)
}

// UpdateWorkflowInstFailed 更新流程实例为失败
// 通过决策器决策出来的失败不走这个逻辑, 正常更新就好
// 不是通过决策器决策出来的失败可以走这个逻辑, 尝试只更新错误信息
func (w *DefaultWorkflowUpdater) UpdateWorkflowInstFailed(req *dto.UpdateWorkflowInstFailedDTO) error {
	inst, err := w.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		if err := w.workflowInstRepo.UpdateWorkflowInstFailed(req); err != nil {
			return err
		}

		externalEvent := event.WorkflowFailEvent{
			BasicEvent: event.NewExternalBasicEventWithReason(event.WorkflowFail, req.Namespace,
				req.DefID, req.InstID, "", req.Reason, map[string]interface{}{}),
			DefID:  req.DefID,
			InstID: req.InstID,
		}
		return w.eventBusRepo.SendExternalEvent(context.Background(), externalEvent)
	}

	if inst.Status.IsTerminal() {
		return fmt.Errorf(instTerminalErrFormat, logs.GetFlowTraceID(req.DefID, req.InstID))
	}

	inst.Status = entity.InstFailed
	inst.Reason.FailedRootCause.FailedNodeRefNames = []string{}
	inst.Reason.FailedRootCause.FailedReason = req.Reason
	return w.UpdateWorkflowInstWithStatus(inst)
}
