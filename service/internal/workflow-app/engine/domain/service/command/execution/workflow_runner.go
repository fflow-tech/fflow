package execution

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/validator"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowRunner 流程实例管理
type WorkflowRunner interface {
	Start(ctx context.Context, req *dto.StartWorkflowInstDTO) (string, error)    // 启动流程
	DriveNext(ctx context.Context, req *dto.DriveWorkflowInstNextNodesDTO) error // 驱动下一个节点
	Restart(ctx context.Context, req *dto.RestartWorkflowInstDTO) error          // 重启流程
	Pause(ctx context.Context, req *dto.PauseWorkflowInstDTO) error              // 暂停流程
	Resume(ctx context.Context, req *dto.ResumeWorkflowInstDTO) error            // 恢复流程
	Complete(ctx context.Context, req *dto.CompleteWorkflowInstDTO) error        // 标记流程结束
	Cancel(ctx context.Context, req *dto.CancelWorkflowInstDTO) error            // 取消流程
	SetTimeout(ctx context.Context, req *dto.SetWorkflowInstTimeoutDTO) error    // 标记流程超时
	UpdateCtx(ctx context.Context, req *dto.UpdateWorkflowInstCtxDTO) error      // 更新流程上下文
	Debug(ctx context.Context, req *dto.DebugWorkflowInstDTO) error              // 更新调试信息
}

// DefaultWorkflowRunner 流程实例管理
type DefaultWorkflowRunner struct {
	workflowDefRepo  ports.WorkflowDefRepository
	workflowInstRepo ports.WorkflowInstRepository
	nodeInstRepo     ports.NodeInstRepository
	cacheRepo        ports.CacheRepository
	nodeRunner       NodeRunner
	msgSender        common.MsgSender
	workflowExecutor WorkflowExecutor
	workflowUpdater  WorkflowUpdater
}

// NewDefaultWorkflowRunner 初始化
func NewDefaultWorkflowRunner(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet,
	workflowExecutor *DefaultWorkflowExecutor,
	nodeRunner *DefaultNodeRunner,
	workflowUpdater *DefaultWorkflowUpdater) *DefaultWorkflowRunner {
	r := &DefaultWorkflowRunner{
		workflowDefRepo:  repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:     repoProviderSet.NodeInstRepo(),
		cacheRepo:        repoProviderSet.CacheRepo(),
		msgSender:        workflowProviderSet.MsgSender(),
		workflowExecutor: workflowExecutor,
		workflowUpdater:  workflowUpdater,
	}
	instExprEvaluator := common.NewInstExprEvaluator(repoProviderSet.WorkflowInstRepo(),
		workflowProviderSet.ExprEvaluator())
	nodeRunner.nodeExecutorRegistry.Register(NewSubWorkflowExecutor(
		r, r.workflowDefRepo, repoProviderSet, instExprEvaluator))
	r.nodeRunner = nodeRunner
	return r
}

// Start 启动
func (e *DefaultWorkflowRunner) Start(ctx context.Context, req *dto.StartWorkflowInstDTO) (string, error) {
	getDefDTO := &dto.GetWorkflowDefDTO{DefID: req.DefID, ReadFromSlave: true}
	def, err := e.workflowDefRepo.GetLastVersion(getDefDTO)
	if err != nil {
		return "", fmt.Errorf("failed to get enabled workflow def: %w", err)
	}

	// 0. 校验请求
	if err := e.validateStartReq(req, def); err != nil {
		return "", err
	}

	// 1. 若用户未填写 input 值则将其设置为默认值
	fillWithDefaultValueIfNotPass(req, def)

	// 2. 在数据库表里创建实例
	instID, err := e.workflowUpdater.CreateWorkflowInst(req, def)
	if err != nil {
		log.Errorf("[%s]Failed to create workflow inst in db, caused by %s",
			logs.GetFlowTraceID(req.DefID, instID), err)
		return "", fmt.Errorf("failed to create workflow inst in db: %w", err)
	}

	// 3. 发送流程开始驱动事件
	if err := e.workflowUpdater.SendWorkflowStartDriveEvent(def.Namespace, def.DefID, instID,
		false); err != nil {
		log.Errorf("[%s]Failed to send WorkflowStartDriveEvent, caused by %s",
			logs.GetFlowTraceID(req.DefID, instID), err)
		return "", fmt.Errorf("failed to send WorkflowStartDriveEvent: %w", err)
	}

	return instID, nil
}

func (e *DefaultWorkflowRunner) validateStartReq(req *dto.StartWorkflowInstDTO,
	def *entity.WorkflowDef) error {
	if def.Status != entity.Enabled {
		return fmt.Errorf("def %d is not enabled", req.DefID)
	}

	count, err := e.workflowInstRepo.Count(&dto.GetWorkflowInstListDTO{DefID: req.DefID})
	if err != nil {
		return fmt.Errorf("failed to count inst for one def: %w", err)
	}

	if int(count) > config.GetValidationRulesConfig().MaxWorkflowInstsForOneDef {
		return fmt.Errorf("one def inst count %d must <= %d",
			count, config.GetValidationRulesConfig().MaxWorkflowInstsForOneDef)
	}

	return validator.ValidateStartInstInput(def, req.Input)
}

// fillWithDefaultValueIfNotPass 若用户未填写 input 值则将其设置为默认值
func fillWithDefaultValueIfNotPass(req *dto.StartWorkflowInstDTO, def *entity.WorkflowDef) {
	for _, inputKeyMap := range def.Input {
		for optionName, inputKeyDef := range inputKeyMap {
			// 该字段在流程定义中不存在，或者用户没有输入默认值
			inputOption, ok := req.Input[optionName]
			if !ok || inputOption == nil {
				req.Input[optionName] = inputKeyDef.Default
			}
		}
	}
}

// DriveNext 调度下一个节点
func (e *DefaultWorkflowRunner) DriveNext(ctx context.Context, req *dto.DriveWorkflowInstNextNodesDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.workflowExecutor.DriveNext(ctx, req)
}

// Restart 重启
func (e *DefaultWorkflowRunner) Restart(ctx context.Context, req *dto.RestartWorkflowInstDTO) error {
	if req.NodeRefName == "" {
		return fmt.Errorf("illegal RestartWorkflowInstDTO, `Node` must not be empty")
	}

	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doRestart(ctx, req)
}

func (e *DefaultWorkflowRunner) doRestart(ctx context.Context, req *dto.RestartWorkflowInstDTO) error {
	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	// 1. 如果填写了流程入参, 重新校验入参
	if len(req.Input) > 0 {
		if err := validator.ValidateStartInstInput(inst.WorkflowDef, req.Input); err != nil {
			return err
		}
	}

	// 2. 取消所有执行的节点
	if !inst.Status.IsTerminal() {
		cancelDTO := &dto.CancelWorkflowInstDTO{
			DefID:    inst.WorkflowDef.DefID,
			InstID:   inst.InstID,
			Operator: req.Operator,
			Reason:   req.Reason,
		}
		if err := e.Cancel(ctx, cancelDTO); err != nil {
			return fmt.Errorf("[%d]failed to cancel inst: %w", cancelDTO.InstID, err)
		}
	}

	// 3. 重新调度流程
	return e.rescheduleInst(ctx, req, inst)
}

// rescheduleInst 重新调度流程
func (e *DefaultWorkflowRunner) rescheduleInst(ctx context.Context,
	req *dto.RestartWorkflowInstDTO, inst *entity.WorkflowInst) error {
	inst.LastRestartNode = req.NodeRefName
	inst.BeforeLastRestartMaxNodeInstID = getMaxNodeInstID(inst)
	inst.LastRestartAt = time.Now()
	inst.Operator.RestartOperator = req.Operator
	inst.Reason.RestartReason = req.Reason
	inst.Status = entity.InstRunning
	// 如果流程被重启，之前记录的暂停相关信息需要被清空
	inst.RunCompletedNodeInstIDsAfterPaused = []string{}
	inst.WaitCompletedNodeInstIDsAfterPaused = []string{}
	if len(req.Input) > 0 {
		inst.Input = req.Input
	}
	if err := e.workflowUpdater.UpdateWorkflowInstWithStatus(inst); err != nil {
		return err
	}

	if err := e.workflowUpdater.RegisterTriggers(inst.WorkflowDef, inst.InstID,
		inst.Operator.RestartOperator); err != nil {
		return err
	}

	if err := e.workflowUpdater.SendWorkflowInstExternalEvent(inst, event.WorkflowStart); err != nil {
		return err
	}

	driveNextReq := &dto.DriveWorkflowInstNextNodesDTO{InstID: req.InstID, DefID: req.DefID}
	return e.DriveNext(ctx, driveNextReq)
}

// Pause 暂停
func (e *DefaultWorkflowRunner) Pause(ctx context.Context, req *dto.PauseWorkflowInstDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doPause(ctx, req)
}

func (e *DefaultWorkflowRunner) doPause(ctx context.Context, req *dto.PauseWorkflowInstDTO) error {
	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	// 2. 如果已经是终态不允许暂停
	if inst.Status.IsTerminal() {
		return fmt.Errorf(instTerminalErrFormat, logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID))
	}
	// 3. 当前流程状态已经是暂停了
	if inst.Status == entity.InstPaused {
		return fmt.Errorf(instRepeatOperationErrFormat, logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID),
			inst.Status.String())
	}

	// 4. 更新流程的状态
	inst.Status = entity.InstPaused
	return e.workflowUpdater.UpdateWorkflowInstWithStatus(inst)
}

// Resume 恢复
func (e *DefaultWorkflowRunner) Resume(ctx context.Context, req *dto.ResumeWorkflowInstDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doResume(ctx, req)
}

func (e *DefaultWorkflowRunner) doResume(ctx context.Context, req *dto.ResumeWorkflowInstDTO) error {
	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	// 2. 如果不是暂停的状态不允许执行恢复
	if inst.Status != entity.InstPaused {
		return fmt.Errorf("[%d]Failed to resume inst, caused by cur status is not paused", inst.InstID)
	}

	runCompletedNodeInstIDsAfterPaused := inst.RunCompletedNodeInstIDsAfterPaused
	waitCompletedNodeInstIDsAfterPaused := inst.WaitCompletedNodeInstIDsAfterPaused

	// 3. 更新流程的状态
	inst.Status = entity.InstRunning
	inst.RunCompletedNodeInstIDsAfterPaused = []string{}
	inst.WaitCompletedNodeInstIDsAfterPaused = []string{}
	if err := e.workflowUpdater.UpdateWorkflowInstWithStatus(inst); err != nil {
		return err
	}

	// 4. 继续向后调度
	if err := e.continueDriveNextNode(runCompletedNodeInstIDsAfterPaused, inst); err != nil {
		return err
	}

	// 5. 如果有等待完成的节点会继续执行
	if err := e.continueRunWaitCompletedNode(ctx, waitCompletedNodeInstIDsAfterPaused, inst); err != nil {
		return err
	}

	// 6. 发送流程被恢复事件
	return e.workflowUpdater.SendWorkflowInstExternalEvent(inst, event.WorkflowResume)
}

func (e *DefaultWorkflowRunner) continueRunWaitCompletedNode(ctx context.Context,
	waitCompletedNodeInstIDsAfterPaused []string,
	inst *entity.WorkflowInst) error {
	for _, nodeInstID := range waitCompletedNodeInstIDsAfterPaused {
		runDTO := &dto.RunNodeDTO{
			DefID:      inst.WorkflowDef.DefID,
			InstID:     inst.InstID,
			NodeInstID: nodeInstID,
		}
		if err := e.nodeRunner.Run(ctx, runDTO); err != nil {
			return err
		}
	}

	return nil
}
func (e *DefaultWorkflowRunner) continueDriveNextNode(runCompletedNodeInstIDsAfterPaused []string,
	inst *entity.WorkflowInst) error {
	if len(runCompletedNodeInstIDsAfterPaused) == 0 && inst.InDebugMode() {
		// 当没有暂停的节点时 则意味着需要发送开始节点
		return e.workflowUpdater.SendWorkflowStartDriveEvent(inst.WorkflowDef.Namespace,
			inst.WorkflowDef.DefID, inst.InstID, true)
	}

	for _, nodeInstID := range runCompletedNodeInstIDsAfterPaused {
		nodeInst, err := e.nodeInstRepo.Get(&dto.GetNodeInstDTO{
			NodeInstID: nodeInstID,
			DefID:      inst.WorkflowDef.DefID,
			DefVersion: inst.WorkflowDef.Version,
			InstID:     inst.InstID,
		})
		if err != nil {
			return err
		}
		if err := e.workflowUpdater.SendNodeCompleteDriveEvent(nodeInst, true); err != nil {
			return err
		}
	}

	return nil
}

// Complete 标记结束
// 已经执行节点的让继续执行, 只是不会继续向后调度
func (e *DefaultWorkflowRunner) Complete(ctx context.Context, req *dto.CompleteWorkflowInstDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doComplete(ctx, req)
}

func (e *DefaultWorkflowRunner) doComplete(ctx context.Context, req *dto.CompleteWorkflowInstDTO) error {
	if req.Status != entity.InstSucceed && req.Status != entity.InstFailed {
		return fmt.Errorf("illegal complete status:%s", req.Status.String())
	}

	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	// 1. 如果已经是目标状态直接返回
	if inst.Status == req.Status {
		return nil
	}

	// 2. 如果已经是终态不允许变更为另一个终态
	if inst.Status.IsTerminal() {
		return fmt.Errorf(instTerminalErrFormat, logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID))
	}

	// 3. 更新流程的状态
	inst.Status = req.Status
	return e.workflowUpdater.UpdateWorkflowInstWithStatus(inst)
}

// Cancel 取消
func (e *DefaultWorkflowRunner) Cancel(ctx context.Context, req *dto.CancelWorkflowInstDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doCancel(ctx, req)
}

func (e *DefaultWorkflowRunner) doCancel(ctx context.Context, req *dto.CancelWorkflowInstDTO) error {
	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	// 如果当前状态和设置的相同直接返回
	if inst.Status == entity.InstCancelled {
		return nil
	}

	// 1. 取消前校验
	if inst.Status.IsTerminal() {
		return fmt.Errorf(instTerminalErrFormat, logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID))
	}

	// 2. 取消所有的节点
	if err := e.cancelAllRunningNodeInsts(ctx, inst, req.Operator, req.Operator); err != nil {
		return err
	}

	// 3. 更新流程的状态
	inst.Status = entity.InstCancelled
	inst.Reason.CancelledReason = req.Reason
	return e.workflowUpdater.UpdateWorkflowInstWithStatus(inst)
}

func (e *DefaultWorkflowRunner) cancelAllRunningNodeInsts(ctx context.Context,
	inst *entity.WorkflowInst, operator, reason string) error {
	for _, nodeInst := range inst.SchedNodeInsts {
		cancelDTO := &dto.CancelNodeDTO{
			DefID:       nodeInst.DefID,
			InstID:      nodeInst.InstID,
			NodeRefName: nodeInst.BasicNodeDef.RefName,
			Operator:    operator,
			Reason:      reason,
		}
		if err := e.nodeRunner.Cancel(ctx, cancelDTO); err != nil {
			log.Errorf("[%s]Failed to cancel node inst, caused by %s",
				logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID), err)
			continue
		}
	}

	return nil
}

// SetTimeout 设置为超时
func (e *DefaultWorkflowRunner) SetTimeout(ctx context.Context, req *dto.SetWorkflowInstTimeoutDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doSetTimeout(ctx, req)
}

// UpdateCtx 更新上下文
func (e *DefaultWorkflowRunner) UpdateCtx(ctx context.Context, req *dto.UpdateWorkflowInstCtxDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doUpdateCtx(ctx, req)
}

func (e *DefaultWorkflowRunner) doUpdateCtx(ctx context.Context, req *dto.UpdateWorkflowInstCtxDTO) error {
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	oldCtxMap, err := utils.StructToMap(inst)
	if err != nil {
		return err
	}
	newCtxMap, err := utils.MergeMap(oldCtxMap, req.Context)
	if err != nil {
		return err
	}
	updateDTO := &dto.UpdateWorkflowInstDTO{
		DefID:   inst.WorkflowDef.DefID,
		InstID:  inst.InstID,
		Context: utils.StructToJsonStr(newCtxMap),
	}
	return e.workflowInstRepo.UpdateWithDefID(updateDTO)
}

func (e *DefaultWorkflowRunner) doSetTimeout(ctx context.Context, req *dto.SetWorkflowInstTimeoutDTO) error {
	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	// 如果当前状态和设置的相同直接返回
	if inst.Status == entity.InstTimeout {
		return nil
	}

	// 1. 设置为超时前校验
	if inst.Status.IsTerminal() {
		return fmt.Errorf(instTerminalErrFormat, logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID))
	}

	return e.doTimeoutOperationByPolicy(ctx, req, inst)
}

func (e *DefaultWorkflowRunner) doTimeoutOperationByPolicy(ctx context.Context, req *dto.SetWorkflowInstTimeoutDTO,
	inst *entity.WorkflowInst) error {
	if inst.WorkflowDef.Timeout.Policy == entity.TimeoutWf {
		return e.terminalInstIfTimeout(ctx, req, inst)
	}

	return e.sendAlertIfTimeout(inst)
}

func (e *DefaultWorkflowRunner) sendAlertIfTimeout(inst *entity.WorkflowInst) error {
	alreadySend, err := e.isAlreadySendInstTimeoutAlert(inst, common.InstTimeoutAlert)
	if err != nil {
		return err
	}
	if alreadySend {
		return nil
	}

	msgInfo := map[string]interface{}{
		"InstID": inst.InstID,
	}
	if inst.Owner.ChatGroup != "" {
		return e.msgSender.SendChatGroupMsg(inst.Owner.ChatGroup, common.InstTimeoutAlert, msgInfo)
	} else if inst.Owner.Wechat != "" {
		return e.msgSender.SendWeChatMsg(inst.Owner.Wechat, common.InstTimeoutAlert, msgInfo)
	}

	return nil
}

func (e *DefaultWorkflowRunner) isAlreadySendInstTimeoutAlert(inst *entity.WorkflowInst,
	template common.MsgTemplate) (bool, error) {
	alreadySend, err := e.cacheRepo.SetNX(
		getInstTimeoutSendAlertKey(inst, template), "true", sendAlertKeyTtl.Milliseconds()/1000)
	if err != nil {
		return false, err
	}
	return alreadySend != "", nil
}

func (e *DefaultWorkflowRunner) terminalInstIfTimeout(ctx context.Context, req *dto.SetWorkflowInstTimeoutDTO,
	inst *entity.WorkflowInst) error {
	// 0. 取消当前实例所有节点运行
	if err := e.cancelAllRunningNodeInsts(ctx, inst, "timer", req.Reason); err != nil {
		return err
	}

	// 1. 修改流程实例状态为超时
	inst.Status = entity.InstTimeout
	return e.workflowUpdater.UpdateWorkflowInstWithStatus(inst)
}

// Debug 调试流程
func (e *DefaultWorkflowRunner) Debug(ctx context.Context, req *dto.DebugWorkflowInstDTO) error {
	lock, err := GetInstDistributeLock(e.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return e.doDebug(ctx, req)
}

func (e *DefaultWorkflowRunner) doDebug(ctx context.Context, req *dto.DebugWorkflowInstDTO) error {
	// 0. 获取当前流程实体
	inst, err := e.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}

	// 1. 更新断点
	inst.Breakpoints = utils.AddElementsToSliceIfNotExists(inst.Breakpoints, req.AddBreakpoints...)
	inst.Breakpoints = utils.DeleteElementsFromSlice(inst.Breakpoints, req.DeleteBreakpoints...)

	// 2. 更新需要 MOCK 的节点
	inst.DebugMockNodes = utils.AddElementsToSliceIfNotExists(inst.DebugMockNodes, req.AddMockNodes...)
	inst.DebugMockNodes = utils.DeleteElementsFromSlice(inst.DebugMockNodes, req.DeleteMockNodes...)

	// 3. 当前调试模式
	inst.CurDebugMode = req.DebugMode

	// 4. 更新流程实例信息
	if err := e.workflowUpdater.UpdateWorkflowInstWithStatus(inst); err != nil {
		return err
	}

	// 5. 恢复节点
	return e.resumeNodeForDebug(ctx, inst, req)
}

func (e *DefaultWorkflowRunner) resumeNodeForDebug(ctx context.Context,
	inst *entity.WorkflowInst, req *dto.DebugWorkflowInstDTO) error {
	if !inst.InDebugMode() {
		return nil
	}

	for _, nodeInst := range inst.SchedNodeInsts {
		if nodeInst.Status == entity.NodeInstWaiting && nodeInst.WaitForDebug {
			if err := e.nodeRunner.Resume(ctx, &dto.ResumeNodeDTO{
				Namespace:   nodeInst.Namespace,
				DefID:       nodeInst.DefID,
				InstID:      nodeInst.InstID,
				NodeRefName: nodeInst.BasicNodeDef.RefName,
				Input:       nodeInst.Input,
				Operator:    req.Operator,
				Reason:      req.Reason,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
