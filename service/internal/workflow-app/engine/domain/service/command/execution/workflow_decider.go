package execution

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/pkg/errors"
)

// WorkflowDecider 流程执行决策类
// 针对逻辑节点的情况, 还是选择在执行器里面走一遍, 便于后续进行扩展
// 如果没有在执行器中执行, 会导致节点的执行顺序不连续, 容易引起使用者的误解
type WorkflowDecider interface {
	Decide(inst *entity.WorkflowInst) (*entity.DecideResult, error)
}

// DefaultWorkflowDecider 默认流程执行决策器
type DefaultWorkflowDecider struct {
	exprEvaluator      expr.Evaluator
	searchNextsFuncMap map[entity.NodeType]func(inst *entity.WorkflowInst, nodeDef interface{}) ([]string, error)
}

// NewDefaultWorkflowDecider 新建决策器
func NewDefaultWorkflowDecider(exprEvaluator *expr.DefaultEvaluator) *DefaultWorkflowDecider {
	r := &DefaultWorkflowDecider{
		exprEvaluator: exprEvaluator,
	}

	r.searchNextsFuncMap = map[entity.NodeType]func(inst *entity.WorkflowInst, nodeDef interface{}) ([]string, error){
		entity.SwitchNode: r.searchSwitchNodeNexts,
		entity.ForkNode:   r.searchForkNodeNexts,
		entity.JoinNode:   r.searchJoinNodeNexts,
	}
	return r
}

// Decide 决策
func (d *DefaultWorkflowDecider) Decide(inst *entity.WorkflowInst) (*entity.DecideResult, error) {
	log.Infof("[%s]Start to do decide, inst:%s",
		logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID), utils.StructToJsonStr(inst))
	decideResult, err := d.doDecide(inst)
	log.Infof("[%s]End to do decide, result:%s",
		logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID), decideResult)
	return decideResult, err
}

func (d *DefaultWorkflowDecider) doDecide(inst *entity.WorkflowInst) (*entity.DecideResult, error) {
	if inst == nil {
		return nil, errors.New("workflow inst must not be nil")
	}

	if notHaveCurNodeInst(inst) {
		if isRestartInst(inst) {
			return d.decideForRestart(inst)
		}

		return d.decideForStart(inst)
	}

	return d.decideByCurNodeInst(inst)
}

func isRestartInst(inst *entity.WorkflowInst) bool {
	return inst.LastRestartNode != ""
}

func notHaveCurNodeInst(inst *entity.WorkflowInst) bool {
	return inst.CurNodeInst == nil || inst.CurNodeInst.BasicNodeDef.RefName == ""
}

func (d *DefaultWorkflowDecider) decideByCurNodeInst(inst *entity.WorkflowInst) (*entity.DecideResult, error) {
	result := entity.NewDecideResult()
	curNodeInst := inst.CurNodeInst
	if d.isAtLeastOnceNode(curNodeInst) && inst.WaitSomeNodesCompleteBeforeInstComplete {
		return d.decideInstStatusByWaitAtLeastNode(inst, curNodeInst)
	}
	// 如果节点成功向后调度, 如果节点执行失败, 根据配置的失败策略决定是否需要向后调度
	if curNodeInst.Status == entity.NodeInstSucceed {
		if d.notNeedScheduleNext(curNodeInst) {
			return result, nil
		}
		if err := d.appendNextsNodeInst(inst, curNodeInst, result); err != nil {
			return nil, err
		}
	} else if curNodeInst.Status == entity.NodeInstFailed {
		if err := d.decideByFailedPolicy(inst, curNodeInst, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// notNeedScheduleNext 判断是否为异步节点或者重跑节点
// 如果为异步节点或者重跑节点, 执行成功不会向后调度, 执行失败会影响流程的结果
func (d *DefaultWorkflowDecider) notNeedScheduleNext(curNodeInst *entity.NodeInst) bool {
	return curNodeInst.BasicNodeDef.Schedule.SchedulePolicy == entity.ScheduleNextIfNotComplete || curNodeInst.FromRerun
}

func (d *DefaultWorkflowDecider) decideByFailedPolicy(inst *entity.WorkflowInst,
	curNodeInst *entity.NodeInst, result *entity.DecideResult) error {
	if curNodeInst.BasicNodeDef.Schedule.FailedPolicy == entity.Ignore {
		if d.notNeedScheduleNext(curNodeInst) {
			return nil
		}
		return d.appendNextsNodeInst(inst, curNodeInst, result)
	}
	result.InstStatus = entity.InstFailed
	result.InstFailedRootCause.FailedNodeRefNames = []string{curNodeInst.BasicNodeDef.RefName}
	if curNodeInst.Reason != nil {
		result.InstFailedRootCause.FailedReason = curNodeInst.Reason.FailedReason
	}
	return nil
}

func (d *DefaultWorkflowDecider) decideForStart(inst *entity.WorkflowInst) (*entity.DecideResult, error) {
	result := entity.NewDecideResult()
	nodeIndexMap, err := entity.GetNodeIndexDefMap(inst.WorkflowDef)
	if err != nil {
		return nil, err
	}
	node, ok := nodeIndexMap[0]
	// 如果没有节点, 直接返回成功
	if !ok {
		result.InstStatus = entity.InstSucceed
		return result, nil
	}
	basicNodeDef, err := entity.GetBasicNodeDefFromNode(inst.WorkflowDef, node)
	if err != nil {
		return nil, err
	}
	if err := d.appendNextNodeInst(inst, basicNodeDef.RefName, result); err != nil {
		return nil, err
	}
	return result, nil
}

// decideInstStatusByWaitAtLeastNode 设置等待必须执行一次节点时流程实例状态
func (d *DefaultWorkflowDecider) decideInstStatusByWaitAtLeastNode(inst *entity.WorkflowInst,
	curNodeInst *entity.NodeInst) (*entity.DecideResult, error) {
	result := entity.NewDecideResult()
	if curNodeInst.Status == entity.NodeInstFailed {
		// 只要没配忽略掉, 节点失败默认就让流程失败
		if curNodeInst.BasicNodeDef.Schedule.FailedPolicy != entity.Ignore {
			result.InstStatus = entity.InstFailed
			result.InstFailedRootCause.FailedNodeRefNames = []string{curNodeInst.BasicNodeDef.RefName}
			if curNodeInst.Reason != nil {
				result.InstFailedRootCause.FailedReason = curNodeInst.Reason.FailedReason
			}
			return result, nil
		}
	}
	hasWaitNode, err := d.hasNodeWaitExecute(inst)
	if err != nil {
		return nil, err
	}
	if hasWaitNode {
		result.InstStatus = entity.InstRunning
		return result, nil
	}
	result.InstStatus = entity.InstSucceed
	return result, nil
}

// isAtLeastOnceNode 是否为必须执行一次的节点
func (d *DefaultWorkflowDecider) isAtLeastOnceNode(nodeInst *entity.NodeInst) bool {
	return nodeInst.BasicNodeDef.Schedule.ExecuteTimesPolicy == entity.AtLeastOnce
}

// decideNodeAppendToBeScheduledResult 决策节点是否放入结果
type decideNodeAppendToBeScheduledResult struct {
	SkipCurNode bool // 跳过当前节点
	RunNextNode bool // 执行下一个节点
}

// appendNodeInst
func (d *DefaultWorkflowDecider) appendNodeInst(inst *entity.WorkflowInst, nodeInst *entity.NodeInst,
	result *entity.DecideResult) error {
	decideNodeAppendResult, err := d.decideNodeAppendToBeScheduled(inst, nodeInst, result)
	if err != nil {
		return err
	}
	// 当两种情况都不跳过时 才执行该节点
	if !decideNodeAppendResult.SkipCurNode {
		result.NodesToBeScheduled = append(result.NodesToBeScheduled, nodeInst)
	}
	if decideNodeAppendResult.RunNextNode {
		if err := d.appendNextsNodeInst(inst, nodeInst, result); err != nil {
			return err
		}
	}
	return nil
}

// decideNodeAppendToBeScheduled 决策节点是否加入等待执行队列
func (d *DefaultWorkflowDecider) decideNodeAppendToBeScheduled(inst *entity.WorkflowInst, nodeInst *entity.NodeInst,
	result *entity.DecideResult) (*decideNodeAppendToBeScheduledResult, error) {
	// 获取调度配置的决策结果
	decideScheduleResult, err := d.decideByScheduleConfig(inst, nodeInst, result)
	if err != nil {
		return nil, err
	}
	// 获取该节点是否被标记成跳过
	skipNode, err := d.needSkipNode(inst, nodeInst)
	if err != nil {
		return nil, err
	}
	decideNodeAppendResult := &decideNodeAppendToBeScheduledResult{
		SkipCurNode: decideScheduleResult.SkipCurNode || skipNode,
		RunNextNode: decideScheduleResult.RunNextNode || skipNode,
	}
	return decideNodeAppendResult, nil
}

// needSkipNode 需要跳过的节点
// 条件节点也走这个分支判断, 因为本质上来说也属于条件跳过
func (d *DefaultWorkflowDecider) needSkipNode(inst *entity.WorkflowInst, nodeInst *entity.NodeInst) (bool, error) {
	for _, skipNode := range inst.SkipNodes {
		if skipNode == nodeInst.BasicNodeDef.RefName {
			return true, nil
		}
	}

	if nodeInst.BasicNodeDef.Condition != "" {
		match, err := matchCondition(d.exprEvaluator, inst, nodeInst.BasicNodeDef.Condition)
		if err != nil {
			return false, err
		}
		return !match, nil
	}

	return false, nil
}

//  decideByScheduleConfig 获取配调度配置的结果
func (d *DefaultWorkflowDecider) decideByScheduleConfig(inst *entity.WorkflowInst, nodeInst *entity.NodeInst,
	result *entity.DecideResult) (*decideNodeAppendToBeScheduledResult, error) {
	// 执行次数决断
	decideExecuteTimesPolicyResult, err := d.decideByExecuteTimesPolicyResult(inst, nodeInst, result)
	if err != nil {
		return nil, err
	}
	// 调度策略决断
	decideSchedulePolicyResult, err := d.decideBySchedulePolicyResult(inst, nodeInst)
	if err != nil {
		return nil, err
	}
	return &decideNodeAppendToBeScheduledResult{
		SkipCurNode: decideExecuteTimesPolicyResult.SkipCurNode || decideSchedulePolicyResult.SkipCurNode,
		RunNextNode: decideExecuteTimesPolicyResult.RunNextNode || decideSchedulePolicyResult.RunNextNode}, nil
}

// decideByExecuteTimesPolicyResult 决断执行次数
func (d *DefaultWorkflowDecider) decideByExecuteTimesPolicyResult(inst *entity.WorkflowInst, nodeInst *entity.NodeInst,
	result *entity.DecideResult) (*decideNodeAppendToBeScheduledResult, error) {
	// 当有且执行一次时
	if nodeInst.BasicNodeDef.Schedule.ExecuteTimesPolicy == entity.ExactlyOnce {
		// 当执行过
		if d.checkNodeInSchedNodeInsts(inst, nodeInst.BasicNodeDef.RefName) {
			return &decideNodeAppendToBeScheduledResult{SkipCurNode: true, RunNextNode: true}, nil
		}
		// 当在准备执行数组中
		if d.checkNodeInToBeScheduledNodes(result, nodeInst.BasicNodeDef.RefName) {
			return &decideNodeAppendToBeScheduledResult{SkipCurNode: true, RunNextNode: true}, nil
		}
	}
	return &decideNodeAppendToBeScheduledResult{SkipCurNode: false, RunNextNode: false}, nil
}

// decideBySchedulePolicyResult 决断调度策略
func (d *DefaultWorkflowDecider) decideBySchedulePolicyResult(inst *entity.WorkflowInst, nodeInst *entity.NodeInst) (
	*decideNodeAppendToBeScheduledResult, error) {
	switch nodeInst.BasicNodeDef.Schedule.SchedulePolicy {
	case entity.IgnoreFirstSchedule:
		return d.ignoreFirstSchedule(inst, nodeInst)
	case entity.ScheduleNextUntilComplete:
		return &decideNodeAppendToBeScheduledResult{SkipCurNode: false, RunNextNode: false}, nil
	case entity.ScheduleNextIfNotComplete:
		return &decideNodeAppendToBeScheduledResult{SkipCurNode: false, RunNextNode: true}, nil
	default:
		return &decideNodeAppendToBeScheduledResult{SkipCurNode: false, RunNextNode: false}, nil
	}
}

// ignoreFirstSchedule 跳过第一次调度
func (d *DefaultWorkflowDecider) ignoreFirstSchedule(inst *entity.WorkflowInst, nodeInst *entity.NodeInst) (
	*decideNodeAppendToBeScheduledResult, error) {
	if d.isScheduleNode(inst, nodeInst) {
		// 被调度过则继续调度
		return &decideNodeAppendToBeScheduledResult{SkipCurNode: false, RunNextNode: false}, nil
	}
	// 不执行则放到已跳过的列表中 并执行下一个节点
	inst.IgnoreFirstScheduleNodes = append(inst.IgnoreFirstScheduleNodes, nodeInst.BasicNodeDef.RefName)
	return &decideNodeAppendToBeScheduledResult{SkipCurNode: true, RunNextNode: true}, nil
}

// isScheduleNode 检测节点是否被调度过
func (d *DefaultWorkflowDecider) isScheduleNode(inst *entity.WorkflowInst, nodeInst *entity.NodeInst) bool {
	// 已经被跳过
	for _, skipNodeRefName := range inst.IgnoreFirstScheduleNodes {
		if skipNodeRefName == nodeInst.BasicNodeDef.RefName {
			return true
		}
	}
	return false
}

func (d *DefaultWorkflowDecider) appendNextsNodeInst(inst *entity.WorkflowInst,
	curNodeInst *entity.NodeInst, result *entity.DecideResult) error {
	nexts, err := d.searchNextNodes(inst, curNodeInst)
	if err != nil {
		return err
	}

	// 0. 如果下一个节点是end节点，直接让流程结束
	if d.nextIsEndNode(curNodeInst.BasicNodeDef, nexts) {
		return d.setInstStatusForArriveEndNode(inst, result)
	}

	// 1. 将接下来的节点添加到待调度的节点中
	for _, next := range nexts {
		if err := d.appendNextNodeInst(inst, next, result); err != nil {
			return err
		}
	}

	return nil
}

// setInstStatusForArriveEndNode 结束节点的执行-对流程实例的字段设置
func (d *DefaultWorkflowDecider) setInstStatusForArriveEndNode(inst *entity.WorkflowInst,
	result *entity.DecideResult) error {
	hasWaitNode, err := d.hasNodeWaitExecute(inst)
	if err != nil {
		return err
	}
	if hasWaitNode {
		inst.WaitSomeNodesCompleteBeforeInstComplete = true
		return nil
	}
	result.InstStatus = entity.InstSucceed
	return nil
}

// hasNodeWaitExecute 检查是否有节点等待执行
func (d *DefaultWorkflowDecider) hasNodeWaitExecute(inst *entity.WorkflowInst) (bool, error) {
	for _, node := range inst.WorkflowDef.Nodes {
		basicNodeDef, err := entity.GetBasicNodeDefFromNode(inst.WorkflowDef, node)
		if err != nil {
			return false, err
		}
		// 当必须执行一次时
		if basicNodeDef.Schedule.ExecuteTimesPolicy == entity.AtLeastOnce {
			if !d.checkNodeInSchedNodeInsts(inst, basicNodeDef.RefName) {
				return true, nil
			}
		}
	}
	return false, nil
}

// checkNodeInToBeScheduledNodes 检查节点是否在准备执行队列中
func (d *DefaultWorkflowDecider) checkNodeInToBeScheduledNodes(result *entity.DecideResult, refName string) bool {
	for _, scheduleNode := range result.NodesToBeScheduled {
		if scheduleNode.BasicNodeDef.RefName == refName {
			return true
		}
	}
	return false
}

// checkNodeInSchedNodeInsts 检查节点是否在已调度的节点数组中
func (d *DefaultWorkflowDecider) checkNodeInSchedNodeInsts(inst *entity.WorkflowInst, refName string) bool {
	for _, scheduleNode := range inst.SchedNodeInsts {
		if scheduleNode.BasicNodeDef.RefName == refName {
			return true
		}
	}
	return false
}

func (d *DefaultWorkflowDecider) nextIsEndNode(basicNodeDef entity.BasicNodeDef, nexts []string) bool {
	for _, next := range nexts {
		if next == entity.EndNode {
			return true
		}
	}

	return false
}

func (d *DefaultWorkflowDecider) appendNextNodeInst(inst *entity.WorkflowInst,
	next string, result *entity.DecideResult) error {
	nodeInst, err := entity.NewNodeInstByNodeRefName(*inst, next)
	if err != nil {
		return err
	}
	// 如果是JOIN节点, 父节点还没完成就不开始执行, 等待下一次调度
	if nodeInst.BasicNodeDef.Type == entity.JoinNode {
		parentsCompleted, err := d.isParentsCompleted(inst, nodeInst.BasicNodeDef.RefName)
		if err != nil {
			return err
		}
		if !parentsCompleted {
			return nil
		}
		// 如果join节点的实例已经存在了直接返回
		if d.joinNodeInstExists(inst, nodeInst.BasicNodeDef.RefName) {
			return nil
		}
	}

	return d.appendNodeInst(inst, nodeInst, result)
}

func (d *DefaultWorkflowDecider) searchNextNodes(inst *entity.WorkflowInst,
	curNodeInst *entity.NodeInst) ([]string, error) {
	curBasicNodeDef, err := entity.GetBasicNodeDefFromNodeDef(curNodeInst.NodeDef)
	if err != nil {
		return nil, err
	}

	// 0. 如果当前节点是Switch节点, 选择第一个满足条件的节点执行
	// 1. 如果当前节点是Fork节点, 选择所有fork的节点并行执行
	// 2. 其他情况看下一个节点的类型再返回
	f, ok := d.searchNextsFuncMap[curBasicNodeDef.Type]
	if !ok {
		return d.searchDefaultNextNodes(inst, curBasicNodeDef)
	}

	nodeDef, err := entity.GetNodeDefByRefName(inst.WorkflowDef, curBasicNodeDef.RefName)
	if err != nil {
		return nil, err
	}
	return f(inst, nodeDef)
}

func (d *DefaultWorkflowDecider) searchDefaultNextNodes(inst *entity.WorkflowInst,
	curBasicNodeDef *entity.BasicNodeDef) ([]string, error) {
	next, err := entity.GetNextNode(inst.WorkflowDef, curBasicNodeDef)
	if err != nil {
		return nil, err
	}

	return []string{next}, nil
}

func (d *DefaultWorkflowDecider) searchSwitchNodeNexts(inst *entity.WorkflowInst,
	curNodeDef interface{}) ([]string, error) {
	nodeDef := curNodeDef.(entity.SwitchNodeDef)
	for _, c := range nodeDef.Switch {
		match, err := matchCondition(d.exprEvaluator, inst, c.Condition)
		if err != nil {
			return nil, err
		}

		if match {
			return []string{c.Next}, nil
		}
	}
	return []string{nodeDef.Next}, nil
}

func (d *DefaultWorkflowDecider) searchForkNodeNexts(inst *entity.WorkflowInst,
	curNodeDef interface{}) ([]string, error) {
	nodeDef := curNodeDef.(entity.ForkNodeDef)
	return nodeDef.Fork, nil
}

func (d *DefaultWorkflowDecider) searchJoinNodeNexts(inst *entity.WorkflowInst,
	curNodeDef interface{}) ([]string, error) {
	nodeDef := curNodeDef.(entity.JoinNodeDef)
	parents, err := entity.GetNodeParents(inst.WorkflowDef, nodeDef.RefName)
	if err != nil {
		return nil, err
	}
	parentsCompleted, err := d.isAllNodesCompleted(inst, parents)
	if err != nil {
		return nil, err
	}
	if parentsCompleted {
		next, err := entity.GetNextNode(inst.WorkflowDef, &nodeDef.BasicNodeDef)
		if err != nil {
			return nil, err
		}
		return []string{next}, nil
	}

	return []string{}, nil
}

func (d *DefaultWorkflowDecider) isParentsCompleted(inst *entity.WorkflowInst,
	refName string) (bool, error) {
	parents, err := entity.GetNodeParents(inst.WorkflowDef, refName)
	if err != nil {
		return false, err
	}

	return d.isAllNodesCompleted(inst, parents)
}

func (d *DefaultWorkflowDecider) isAllNodesCompleted(inst *entity.WorkflowInst,
	nodeRefNames []string) (bool, error) {
	nodeInstMap := entity.GetNodeRefNameLatestInstMap(inst.SchedNodeInsts)
	for _, parent := range nodeRefNames {
		nodeInst, exists := nodeInstMap[parent]
		if !exists {
			return false, nil
		}

		if !nodeInst.Status.IsCompleted() {
			return false, nil
		}
	}

	return true, nil
}

func (d *DefaultWorkflowDecider) decideForRestart(inst *entity.WorkflowInst) (*entity.DecideResult, error) {
	result := entity.NewDecideResult()
	if err := d.appendNextNodeInst(inst, inst.LastRestartNode, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DefaultWorkflowDecider) joinNodeInstExists(inst *entity.WorkflowInst, refName string) bool {
	joinNodeInst := entity.GetNodeInstByRefName(inst.SchedNodeInsts, refName)
	return joinNodeInst != nil && inst.CurNodeInst != nil && joinNodeInst.NodeInstID > inst.CurNodeInst.NodeInstID
}
