package execution

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/nodeexecutor"
	"github.com/fflow-tech/fflow/service/pkg/expr"
)

// WorkflowProviderSet 流程基础依赖集合
type WorkflowProviderSet struct {
	workflowDecider      WorkflowDecider
	exprEvaluator        expr.Evaluator
	instTimeoutChecker   common.InstTimeoutChecker
	nodeExecutorRegistry nodeexecutor.Registry
	msgSender            common.MsgSender
}

// MsgSender 消息发送者
func (w *WorkflowProviderSet) MsgSender() common.MsgSender {
	return w.msgSender
}

// InstTimeoutChecker 超时检查器
func (w *WorkflowProviderSet) InstTimeoutChecker() common.InstTimeoutChecker {
	return w.instTimeoutChecker
}

// WorkflowDecider 决策器
func (w *WorkflowProviderSet) WorkflowDecider() WorkflowDecider {
	return w.workflowDecider
}

// ExprEvaluator 表达式计算
func (w *WorkflowProviderSet) ExprEvaluator() expr.Evaluator {
	return w.exprEvaluator
}

// NodeExecutorRegistry 节点执行器注册中心
func (w *WorkflowProviderSet) NodeExecutorRegistry() nodeexecutor.Registry {
	return w.nodeExecutorRegistry
}

// NewWorkflowProviderSet 实例化
func NewWorkflowProviderSet(workflowDecider *DefaultWorkflowDecider,
	exprEvaluator *expr.DefaultEvaluator,
	instTimeoutChecker *common.DefaultInstTimeoutChecker,
	nodeExecutorRegistry *nodeexecutor.DefaultRegistry,
	msgSender *common.DefaultMsgSender) *WorkflowProviderSet {
	r := &WorkflowProviderSet{
		workflowDecider:      workflowDecider,
		exprEvaluator:        exprEvaluator,
		instTimeoutChecker:   instTimeoutChecker,
		nodeExecutorRegistry: nodeExecutorRegistry,
		msgSender:            msgSender,
	}
	return r
}

// NewWorkflowProviderSetMock 实例化
func NewWorkflowProviderSetMock(workflowDecider *DefaultWorkflowDecider,
	exprEvaluator *expr.DefaultEvaluator,
	instTimeoutChecker *common.DefaultInstTimeoutChecker,
	nodeExecutorRegistry *nodeexecutor.DefaultRegistry,
	msgSender *common.DefaultMsgSender) *WorkflowProviderSet {
	return &WorkflowProviderSet{
		workflowDecider:      workflowDecider,
		exprEvaluator:        exprEvaluator,
		instTimeoutChecker:   instTimeoutChecker,
		nodeExecutorRegistry: nodeExecutorRegistry,
		msgSender:            msgSender,
	}
}
