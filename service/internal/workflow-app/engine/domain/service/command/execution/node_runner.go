package execution

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/nodeexecutor"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/jinzhu/copier"
)

// NodeRunner 节点实例执行器
type NodeRunner interface {
	Run(ctx context.Context, req *dto.RunNodeDTO) error
	Rerun(ctx context.Context, req *dto.RerunNodeDTO) error
	Skip(ctx context.Context, req *dto.SkipNodeDTO) error
	CancelSkip(ctx context.Context, req *dto.CancelSkipNodeDTO) error
	Resume(ctx context.Context, req *dto.ResumeNodeDTO) error
	Cancel(ctx context.Context, req *dto.CancelNodeDTO) error
	Polling(ctx context.Context, req *dto.PollingNodeDTO) error
	Schedule(ctx context.Context, req *dto.ScheduleNodeDTO) error
	Complete(ctx context.Context, req *dto.CompleteNodeDTO) error
	SetTimeout(ctx context.Context, req *dto.SetNodeTimeoutDTO) error
	SetNearTimeout(ctx context.Context, req *dto.SetNodeNearTimeoutDTO) error
}

// DefaultNodeRunner 默认节点执行器
type DefaultNodeRunner struct {
	workflowDefRepo      ports.WorkflowDefRepository
	workflowInstRepo     ports.WorkflowInstRepository
	nodeInstRepo         ports.NodeInstRepository
	eventBusRepo         ports.EventBusRepository
	cacheRepo            ports.CacheRepository
	nodeExecutorRegistry nodeexecutor.Registry
	msgSender            common.MsgSender
	nodePoller           NodePoller
	workflowUpdater      WorkflowUpdater
}

// NewDefaultNodeRunner 初始化
func NewDefaultNodeRunner(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet,
	nodePoller *DefaultNodePoller,
	workflowUpdater *DefaultWorkflowUpdater) *DefaultNodeRunner {
	r := &DefaultNodeRunner{
		workflowDefRepo:      repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo:     repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:         repoProviderSet.NodeInstRepo(),
		eventBusRepo:         repoProviderSet.EventBusRepo(),
		cacheRepo:            repoProviderSet.CacheRepo(),
		nodeExecutorRegistry: workflowProviderSet.NodeExecutorRegistry(),
		msgSender:            workflowProviderSet.MsgSender(),
	}
	instExprEvaluator := common.NewInstExprEvaluator(repoProviderSet.WorkflowInstRepo(),
		workflowProviderSet.ExprEvaluator())
	r.nodeExecutorRegistry.Register(nodeexecutor.NewAssignNodeExecutor(
		repoProviderSet.WorkflowInstRepo(), workflowProviderSet.ExprEvaluator()))
	r.nodeExecutorRegistry.Register(NewServiceNodeExecutor(repoProviderSet, workflowProviderSet))
	r.nodeExecutorRegistry.Register(nodeexecutor.NewJoinNodeExecutor())
	r.nodeExecutorRegistry.Register(nodeexecutor.NewSwitchNodeExecutor())
	r.nodeExecutorRegistry.Register(nodeexecutor.NewForkNodeExecutor())
	r.nodeExecutorRegistry.Register(nodeexecutor.NewWaitNodeExecutor())
	r.nodeExecutorRegistry.Register(nodeexecutor.NewExclusiveJoinNodeExecutor(repoProviderSet.WorkflowInstRepo()))
	r.nodeExecutorRegistry.Register(nodeexecutor.NewTransformNodeExecutor(instExprEvaluator))
	r.nodeExecutorRegistry.Register(NewEventNodeExecutor(repoProviderSet, workflowProviderSet))
	r.workflowUpdater = workflowUpdater
	r.nodePoller = nodePoller
	return r
}

var (
	nodeInstTerminalErrFormat = "[%s]node inst [%s][%d] is already terminal"
	nodeInstRunningErrFormat  = "[%s]node inst [%s][%d] is running"
)

// Run 运行
func (r *DefaultNodeRunner) Run(ctx context.Context, req *dto.RunNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doRun(ctx, req)
}

func (r *DefaultNodeRunner) doRun(ctx context.Context, req *dto.RunNodeDTO) error {
	nodeInst, executor, err := r.getNodeInstAndExecutorForExecuteDriveEvent(&dto.GetNodeInstDTO{
		NodeInstID: req.NodeInstID,
		InstID:     req.InstID,
		DefID:      req.DefID,
	})
	if err != nil {
		return err
	}

	if err := r.beforeExecute(req, nodeInst); err != nil {
		return err
	}

	if err := r.doExecute(ctx, executor, nodeInst); err != nil {
		return err
	}

	return r.afterExecute(executor, nodeInst)
}

// Schedule 调度
func (r *DefaultNodeRunner) Schedule(ctx context.Context, req *dto.ScheduleNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doSchedule(req)
}

// doSchedule 执行真正的调度
func (r *DefaultNodeRunner) doSchedule(req *dto.ScheduleNodeDTO) error {
	nodeInst, err := r.nodeInstRepo.Get(&dto.GetNodeInstDTO{
		NodeInstID: req.NodeInstID,
		InstID:     req.InstID,
		DefID:      req.DefID,
		DefVersion: req.DefVersion,
	})
	if err != nil {
		log.Errorf("Failed to get node inst, caused by %s", err)
		return err
	}

	if nodeInst.Status != entity.NodeInstScheduled {
		log.Warnf("[%s]Node inst [%s][%d] already scheduled",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
		return nil
	}

	if r.needWaitForDebug(nodeInst) {
		nodeInst.WaitForDebug = true
		nodeInst.Status = entity.NodeInstWaiting
		nodeInst.WaitAt = time.Now()
		return r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst)
	}

	// 0. 如果需要等待, 则发送延迟消费的执行消息, 并修改为等待状态
	if needDelay(nodeInst) {
		nodeInst.Status = entity.NodeInstWaiting
		nodeInst.WaitAt = time.Now()
		return r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst)
	}

	// 1. 如果不需要等待, 则发送实时消费的执行消息, 不需要修改执行状态
	return r.workflowUpdater.SendNodeDriveEvent(nodeInst, event.NodeExecuteDrive)
}

// needWaitForDebug 因为处于调试模式需要等待
func (r *DefaultNodeRunner) needWaitForDebug(nodeInst *entity.NodeInst) bool {
	inst, err := r.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(nodeInst.InstID, nodeInst.DefID, ""))
	if err != nil {
		return false
	}

	return inst.CurDebugMode == entity.SingleStepMode ||
		(inst.CurDebugMode == entity.BreakpointMode &&
			utils.StrContains(inst.Breakpoints, nodeInst.BasicNodeDef.RefName))
}

// Rerun 重跑
func (r *DefaultNodeRunner) Rerun(ctx context.Context, req *dto.RerunNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doRerun(ctx, req)
}

// beforeRerun 根据节点实例状态判断可否重跑
func (r *DefaultNodeRunner) beforeRerun(ctx context.Context, req *dto.RerunNodeDTO) error {
	nodeInst, err := r.getNodeInstByInstIDAndRefName(req.InstID, req.NodeRefName)
	if err != nil {
		return err
	}
	// 只有节点状态为终态才能 rerun
	if !nodeInst.Status.IsTerminal() {
		return fmt.Errorf(nodeInstRunningErrFormat+", skip rerun",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
	}
	// 防止DefID未传时后续操作报错
	req.DefID = nodeInst.DefID
	return nil
}

// doRerun 重跑节点
func (r *DefaultNodeRunner) doRerun(ctx context.Context, req *dto.RerunNodeDTO) error {
	// 判断是否可以重跑
	if err := r.beforeRerun(ctx, req); err != nil {
		return err
	}
	// 创建新的节点实例
	newNodeInstID, err := r.createNodeInst(req)
	if err != nil {
		return err
	}
	// 重跑节点
	runReq := &dto.RunNodeDTO{}
	if err := copier.Copy(runReq, req); err != nil {
		return err
	}
	runReq.NodeInstID = newNodeInstID
	return r.Run(ctx, runReq)
}

// createNodeInst 创建新的节点实例
func (r *DefaultNodeRunner) createNodeInst(req *dto.RerunNodeDTO) (string, error) {
	// 获取流程实例
	inst, err := r.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return "", err
	}
	// 创建节点实例 entity
	nodeInst, err := entity.NewNodeInstByNodeRefName(*inst, req.NodeRefName)
	if err != nil {
		return "", err
	}
	nodeInst.FromRerun = true
	nodeInst.Input = req.Input
	nodeInst.Reason.RerunReason = req.Reason
	nodeInst.Operator.RerunOperator = req.Operator
	createNodeInstDTO, err := convertor.NodeInstConvertor.ConvertEntityToCreateDTO(nodeInst)
	if err != nil {
		return "", err
	}
	// 创建新的节点实例
	nodeInstID, err := r.nodeInstRepo.Create(createNodeInstDTO)
	if err != nil {
		return "", err
	}
	return nodeInstID, nil
}

// Skip 跳过节点执行
// 具体是指下次调度到跳过的节点会选择忽略执行，流程实例内有效
func (r *DefaultNodeRunner) Skip(ctx context.Context, req *dto.SkipNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doSkipNode(req)
}

func (r *DefaultNodeRunner) doSkipNode(req *dto.SkipNodeDTO) error {
	inst, err := r.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}
	if !r.instHasSkipNode(inst, req.NodeRefName) {
		// 找不到对应的节点直接返回
		return fmt.Errorf("failed to get skip node, caused by not find:%s", req.NodeRefName)
	}

	if r.instSkipNodeHasInputNode(inst, req.NodeRefName) {
		// 已经跳过的节点不再写入
		return nil
	}
	inst.SkipNodes = append(inst.SkipNodes, req.NodeRefName)

	return r.workflowUpdater.UpdateWorkflowInst(inst)
}

// instHasSkipNode 判断流程是否包含跳过的节点
func (r *DefaultNodeRunner) instHasSkipNode(inst *entity.WorkflowInst, skipNodeName string) bool {
	for _, node := range inst.WorkflowDef.Nodes {
		if _, ok := node[skipNodeName]; ok {
			return true
		}
	}
	for _, subWorkflow := range inst.WorkflowDef.Subworkflows {
		if _, ok := subWorkflow[skipNodeName]; ok {
			return true
		}
	}
	return false
}

// instSkipNodeHasInputNode 实例的跳过节点是否包含传入节点
func (r *DefaultNodeRunner) instSkipNodeHasInputNode(inst *entity.WorkflowInst, skipNodeName string) bool {
	for _, nodeName := range inst.SkipNodes {
		if nodeName == skipNodeName {
			return true
		}
	}
	return false
}

// CancelSkip 取消跳过
func (r *DefaultNodeRunner) CancelSkip(ctx context.Context, req *dto.CancelSkipNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doCancelSkipNode(req)
}

func (r *DefaultNodeRunner) doCancelSkipNode(req *dto.CancelSkipNodeDTO) error {
	inst, err := r.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(req.InstID, req.DefID, ""))
	if err != nil {
		return err
	}
	if !r.instHasSkipNode(inst, req.NodeRefName) {
		// 找不到对应的节点直接返回
		return fmt.Errorf("failed to get cancel skip node, caused by not find:%s", req.NodeRefName)
	}

	if !r.instSkipNodeHasInputNode(inst, req.NodeRefName) {
		// 已经跳过的节点未发现需要取消的节点则直接返回
		return nil
	}
	// 删除指定节点
	for index, skipNode := range inst.SkipNodes {
		if skipNode == req.NodeRefName {
			inst.SkipNodes = append(inst.SkipNodes[:index], inst.SkipNodes[index+1:]...)
			break
		}
	}

	return r.workflowUpdater.UpdateWorkflowInst(inst)
}

// Resume 恢复节点执行
// 恢复的前置条件为节点处于等待状态
func (r *DefaultNodeRunner) Resume(ctx context.Context, req *dto.ResumeNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	return r.doResume(ctx, req)
}

func (r *DefaultNodeRunner) doResume(ctx context.Context, req *dto.ResumeNodeDTO) error {
	nodeInst, err := r.getNodeInstByInstIDAndNodeInstID(req.InstID, req.NodeInstID, req.NodeRefName)
	if err != nil {
		return err
	}
	// 只有等待的节点可以恢复
	if nodeInst.Status != entity.NodeInstWaiting {
		return fmt.Errorf("[%s]node inst [%s][%d] is not waiting, skip resume",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
	}
	// 恢复节点
	runReq := &dto.RunNodeDTO{}
	if err := copier.Copy(runReq, req); err != nil {
		return err
	}
	runReq.DefID = nodeInst.DefID
	runReq.NodeInstID = nodeInst.NodeInstID
	runReq.Reason = fmt.Sprintf("[resume]%s", req.Reason)
	runReq.Operator = req.Operator
	return r.Run(ctx, runReq)
}

// Cancel 取消节点的执行
func (r *DefaultNodeRunner) Cancel(ctx context.Context, req *dto.CancelNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doCancel(ctx, req)
}

func (r *DefaultNodeRunner) doCancel(ctx context.Context, req *dto.CancelNodeDTO) error {
	nodeInst, err := r.getNodeInstByInstIDAndNodeInstID(req.InstID, req.NodeInstID, req.NodeRefName)
	if err != nil {
		return err
	}

	// 如果当前状态和设置的相同直接返回
	if nodeInst.Status == entity.NodeInstCancelled {
		return nil
	}

	if nodeInst.Status.IsTerminal() {
		return fmt.Errorf(nodeInstTerminalErrFormat+", skip cancel",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
	}

	executor, exists := r.nodeExecutorRegistry.GetExecutor(nodeInst.BasicNodeDef.Type)
	if !exists {
		return fmt.Errorf("failed to get executor, nodeType=%s", nodeInst.BasicNodeDef.Type)
	}

	nodeInst.Reason.CancelledReason = req.Reason
	nodeInst.Operator.CancelledOperator = req.Operator

	if !executor.AsyncComplete(nodeInst) {
		nodeInst.Status = entity.NodeInstCancelled
		return r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst)
	}

	if err := executor.Cancel(ctx, nodeInst); err != nil {
		return err
	}

	return r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst)
}

// Complete 标记完成
func (r *DefaultNodeRunner) Complete(ctx context.Context, req *dto.CompleteNodeDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doComplete(ctx, req)
}

func (r *DefaultNodeRunner) doComplete(ctx context.Context, req *dto.CompleteNodeDTO) error {
	nodeInst, err := r.getNodeInstByInstIDAndNodeInstID(req.InstID, req.NodeInstID, req.NodeRefName)
	if err != nil {
		return err
	}
	// 如果当前状态和设置的相同直接返回
	if nodeInst.Status == req.Status {
		return nil
	}
	// 如果已经是终态则不能设置为另外一种完成状态
	if nodeInst.Status.IsTerminal() {
		return fmt.Errorf(nodeInstTerminalErrFormat+", skip complete",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
	}

	// 更新inst
	nodeInst.Status = req.Status
	nodeInst.Output = req.Output
	if req.Status == entity.NodeInstFailed {
		nodeInst.Reason.FailedReason = req.Reason
		nodeInst.Operator.FailedOperator = req.Operator
	} else {
		nodeInst.Reason.SucceedReason = req.Reason
		nodeInst.Operator.SucceedOperator = req.Operator
	}

	// 更新节点状态并发送外部事件
	if err := r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst); err != nil {
		return err
	}
	// 发送 nodeDrive 事件
	return r.workflowUpdater.SendNodeDriveEvent(nodeInst, event.NodeCompleteDrive)
}

// Polling 轮询
func (r *DefaultNodeRunner) Polling(ctx context.Context, req *dto.PollingNodeDTO) error {
	return r.nodePoller.Polling(ctx, req)
}

// SetTimeout 标记为超时
func (r *DefaultNodeRunner) SetTimeout(ctx context.Context, req *dto.SetNodeTimeoutDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doSetTimeout(ctx, req)
}

func (r *DefaultNodeRunner) doSetTimeout(ctx context.Context, req *dto.SetNodeTimeoutDTO) error {
	// 0. 获取当前节点实例
	nodeInst, err := r.getNodeInstByInstIDAndRefName(req.InstID, req.NodeRefName)
	if err != nil {
		return err
	}
	// 如果当前状态和设置的相同直接返回
	if nodeInst.Status == entity.NodeInstTimeout {
		return nil
	}

	// 1. 设置节点状态为超时之前进行校验
	if nodeInst.Status.IsTerminal() {
		return fmt.Errorf(nodeInstTerminalErrFormat+", skip set timeout",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
	}

	nodeInst.Reason.TimeoutReason = req.Reason
	// 2. 根据超时策略进行操作
	return r.doTimeoutOperationByPolicy(ctx, nodeInst)
}

// SetNearTimeout 标记节点为接近超时
func (r *DefaultNodeRunner) SetNearTimeout(ctx context.Context, req *dto.SetNodeNearTimeoutDTO) error {
	lock, err := GetInstDistributeLock(r.cacheRepo, req.InstID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	return r.doSetNearTimeout(req)
}

func (r *DefaultNodeRunner) doSetNearTimeout(req *dto.SetNodeNearTimeoutDTO) error {
	// 0. 获取节点实例
	nodeInst, err := r.getNodeInstByInstIDAndRefName(req.InstID, req.NodeRefName)
	if err != nil {
		log.Errorf("Failed to get node inst, caused by %s", err)
		return err
	}

	// 1. 设置节点为接近超时之前进行校验
	if nodeInst.Status.IsTerminal() {
		return fmt.Errorf(nodeInstTerminalErrFormat+", skip set near timeout",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
	}

	// 2. 根据接近超时策略进行操作
	return r.doNearTimeoutOperationByPolicy(nodeInst)
}

func (r *DefaultNodeRunner) doTimeoutOperationByPolicy(ctx context.Context, nodeInst *entity.NodeInst) error {
	// 0. 判断超时策略, 当没有值时则为默认的超时策略 entity.TimeoutWf
	timeoutPolicy := nodeInst.BasicNodeDef.Timeout.Policy
	if timeoutPolicy == "" || timeoutPolicy == entity.TimeoutWf {
		return r.terminalNodeInstIfTimeout(ctx, nodeInst)
	}

	return r.sendAlertIfTimeout(nodeInst)
}

func (r *DefaultNodeRunner) doNearTimeoutOperationByPolicy(nodeInst *entity.NodeInst) error {
	// 0. 判断节点接近超时策略，目前只有 AlertOnly 一种，其他策略返回 nil
	nearTimeoutPolicy := nodeInst.BasicNodeDef.Timeout.NearTimeoutPolicy
	if nearTimeoutPolicy == "" || nearTimeoutPolicy == entity.AlertOnly {
		if err := r.sendAlertIfNearTimeout(nodeInst); err != nil {
			return err
		}
		// 发送节点接近超时事件
		return r.workflowUpdater.SendNodeExternalEvent(nodeInst, event.NodeNearTimeout)
	}

	return nil
}

// terminalNodeInstIfTimeout 函数处理节点超时策略为 TIME_OUT_WF 时的逻辑，节点类型为异步时终止节点执行，同步则直接返回
// 节点超时不会影响流程的状态
func (r *DefaultNodeRunner) terminalNodeInstIfTimeout(ctx context.Context, nodeInst *entity.NodeInst) error {
	// 0. 获取节点执行器
	executor, exists := r.nodeExecutorRegistry.GetExecutor(nodeInst.BasicNodeDef.Type)
	if !exists {
		log.Errorf("[%s]Failed to get executor, nodeType=%s",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID),
			nodeInst.BasicNodeDef.Type)
		return nil
	}
	// 1. 如果节点为异步节点，取消节点实例执行；同步节点取消没有任何操作
	if err := executor.Cancel(ctx, nodeInst); err != nil {
		return err
	}
	// 2. 设置节点实例状态为超时并发送外部事件
	nodeInst.Status = entity.NodeInstTimeout
	return r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst)
}

// sendAlertIfTimeout 节点超时时发送告警
func (r *DefaultNodeRunner) sendAlertIfTimeout(nodeInst *entity.NodeInst) error {
	alreadySend, err := r.isAlreadySendNodeInstTimeoutAlert(nodeInst, common.NodeInstTimeoutAlert)
	if err != nil {
		return err
	}
	if alreadySend {
		return nil
	}

	msgInfo := map[string]interface{}{
		"InstID":      nodeInst.InstID,
		"NodeRefName": nodeInst.BasicNodeDef.RefName,
	}
	if nodeInst.Owner.ChatGroup != "" {
		return r.msgSender.SendChatGroupMsg(nodeInst.Owner.ChatGroup, common.NodeInstTimeoutAlert, msgInfo)
	} else if nodeInst.Owner.Wechat != "" {
		return r.msgSender.SendWeChatMsg(nodeInst.Owner.Wechat, common.NodeInstTimeoutAlert, msgInfo)
	}

	return nil
}

// sendAlertIfNearTimeout 节点接近超时时发送告警
func (r *DefaultNodeRunner) sendAlertIfNearTimeout(nodeInst *entity.NodeInst) error {
	alreadySend, err := r.isAlreadySendNodeInstTimeoutAlert(nodeInst, common.NodeInstNearTimeoutAlert)
	if err != nil {
		return err
	}
	if alreadySend {
		return nil
	}

	msgInfo := map[string]interface{}{
		"InstID":      nodeInst.InstID,
		"NodeRefName": nodeInst.BasicNodeDef.RefName,
	}
	if nodeInst.Owner.ChatGroup != "" {
		return r.msgSender.SendChatGroupMsg(nodeInst.Owner.ChatGroup,
			common.NodeInstNearTimeoutAlert, msgInfo)
	} else if nodeInst.Owner.Wechat != "" {
		return r.msgSender.SendWeChatMsg(nodeInst.Owner.Wechat, common.NodeInstNearTimeoutAlert, msgInfo)
	}

	return nil
}

// isAlreadySendNodeInstTimeoutAlert 判断是否已经发送过告警
func (r *DefaultNodeRunner) isAlreadySendNodeInstTimeoutAlert(nodeInst *entity.NodeInst,
	template common.MsgTemplate) (bool, error) {
	alreadySend, err := r.cacheRepo.SetNX(
		getNodeTimeoutSendAlertKey(nodeInst, template), "true", sendAlertKeyTtl.Milliseconds()/1000)
	if err != nil {
		return false, err
	}
	return alreadySend != "", nil
}

func (r *DefaultNodeRunner) getNodeInstAndExecutorForExecuteDriveEvent(getNodeInstDTO *dto.GetNodeInstDTO) (
	*entity.NodeInst, nodeexecutor.NodeExecutor, error) {
	nodeInst, err := r.nodeInstRepo.Get(getNodeInstDTO)

	if err != nil {
		log.Errorf("Failed to get node inst, caused by %s", err)
		return nil, nil, err
	}
	executor, exists := r.nodeExecutorRegistry.GetExecutor(nodeInst.BasicNodeDef.Type)
	if !exists {
		return nil, nil, fmt.Errorf("failed to get nodeType=%s executor", nodeInst.BasicNodeDef.Type)
	}
	return nodeInst, executor, nil
}

func (r *DefaultNodeRunner) beforeExecute(req *dto.RunNodeDTO, nodeInst *entity.NodeInst) error {
	startTime := time.Now()
	defer func() {
		log.Infof("[%s]Before execute node [%d] costs %s",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), nodeInst.NodeInstID, time.Since(startTime))
	}()

	// 如果不是重试，但节点又不处于前置状态，说明是重复推送执行消息，忽略执行
	if !nodeInst.Retrying && nodeInst.Status != entity.NodeInstScheduled && nodeInst.Status != entity.NodeInstWaiting {
		return fmt.Errorf("[%s]failed to execute node [%s][%d], caused by node already running",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), nodeInst.BasicNodeDef.RefName, nodeInst.NodeInstID)
	}

	nodeInst.Input = req.Input
	nodeInst.Operator.RunOperator = req.Operator
	nodeInst.Reason.RunReason = req.Reason
	nodeInst.ExecuteAt = time.Now()
	// 重置等待重试状态
	nodeInst.Retrying = false
	// 当已经是执行状态时，又有执行命令 则表示为重试
	if nodeInst.Status == entity.NodeInstRunning {
		nodeInst.RetryCount += 1
	}
	err := r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst)
	if err != nil {
		return err
	}
	return r.workflowUpdater.SendNodeExternalEvent(nodeInst, event.NodeStart)
}

func (r *DefaultNodeRunner) afterExecute(executor nodeexecutor.NodeExecutor, nodeInst *entity.NodeInst) error {
	startTime := time.Now()
	defer func() {
		log.Infof("[%s]After execute node [%d] costs %s",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), nodeInst.NodeInstID, time.Since(startTime))
	}()

	if err := r.workflowUpdater.UpdateNodeInstWithStatus(nodeInst); err != nil {
		return err
	}

	if executor.AsyncComplete(nodeInst) && !nodeInst.Status.IsTerminal() && !nodeInst.Retrying {
		if err := r.workflowUpdater.SendNodeExternalEvent(nodeInst, event.NodeAsynWait); err != nil {
			log.Warnf("Failed to send NodeAsynWaitEvent, caused by %s", err)
		}
	}

	if executor.AsyncByPolling(nodeInst) && !nodeInst.Status.IsTerminal() && !nodeInst.Retrying {
		return r.workflowUpdater.SendNodeDriveEvent(nodeInst, event.NodePollDrive)
	}

	if nodeInst.Status.IsCompleted() {
		return r.workflowUpdater.SendNodeDriveEvent(nodeInst, event.NodeCompleteDrive)
	}

	return nil
}

func (r *DefaultNodeRunner) doExecute(ctx context.Context, executor nodeexecutor.NodeExecutor,
	nodeInst *entity.NodeInst) error {
	startTime := time.Now()
	defer func() {
		log.Infof("[%s]Execute node [%d] costs %s",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), nodeInst.NodeInstID, time.Since(startTime))
	}()
	// 实际执行失败的情况, 不用返回错误, 直接将错误放到错误原因里面
	if err := executor.Execute(ctx, nodeInst); err != nil {
		nodeInst.Status = entity.NodeInstFailed
		nodeInst.Reason.FailedReason = err.Error()
	}

	// 当节点失败时判断本身是否重试
	if needRetry(nodeInst) {
		return r.retryNode(nodeInst)
	}

	if nodeInst.Status.IsTerminal() {
		return nil
	}

	if executor.AsyncComplete(nodeInst) {
		nodeInst.AsynWaitResAt = time.Now()
		nodeInst.Status = entity.NodeInstRunning
		return nil
	}

	nodeInst.Status = entity.NodeInstSucceed
	return nil
}

// retryNode 节点重试
func (r *DefaultNodeRunner) retryNode(curNodeInst *entity.NodeInst) error {
	retryDelay, err := getRetryDelay(curNodeInst)
	if err != nil {
		return err
	}
	// 更新节点状态，这里是为了让节点保持运行状态
	curNodeInst.Status = entity.NodeInstRunning
	curNodeInst.Retrying = true
	if err := r.workflowUpdater.UpdateNodeInstWithStatus(curNodeInst); err != nil {
		return err
	}
	// 发送节点重试消息
	return r.workflowUpdater.SendDelayNodeExecuteDriveEvent(curNodeInst, retryDelay)
}

// getNodeInstByInstIDAndRefName 根据引用名称获取节点实例
func (r *DefaultNodeRunner) getNodeInstByInstIDAndRefName(instID string, refName string) (*entity.NodeInst, error) {
	inst, err := r.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(instID, "", ""))
	if err != nil {
		return nil, err
	}
	for _, node := range inst.SchedNodeInsts {
		if node.BasicNodeDef.RefName == refName {
			return node, nil
		}
	}
	return nil, fmt.Errorf("nodeInst refname [%s] in workflowInst: [%s] not found", refName, inst.InstID)
}

// getNodeInstByInstIDAndNodeInstID 根据节点实例 ID 获取节点实例
func (r *DefaultNodeRunner) getNodeInstByInstIDAndNodeInstID(instID string,
	nodeInstID string, refName string) (*entity.NodeInst, error) {
	nodeInst, err := r.getNodeInstByInstIDAndRefName(instID, refName)
	if err != nil {
		return nil, err
	}
	if nodeInstID != "" && nodeInst.NodeInstID != nodeInstID {
		return nil, fmt.Errorf("nodeInst refname:[%s] nodeInstID:[%d] in workflowInst: [%d] not found",
			refName, nodeInstID, instID)
	}

	return nodeInst, nil
}
