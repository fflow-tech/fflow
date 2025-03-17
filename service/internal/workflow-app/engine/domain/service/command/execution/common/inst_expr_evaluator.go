package common

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/expr"
)

// InstExprEvaluator 实例表达式计算器
type InstExprEvaluator struct {
	workflowInstRepo ports.WorkflowInstRepository
	exprEvaluator    expr.Evaluator
}

// NewInstExprEvaluator 初始化实例表达式计算器
func NewInstExprEvaluator(workflowInstRepo ports.WorkflowInstRepository,
	exprEvaluator expr.Evaluator) *InstExprEvaluator {
	return &InstExprEvaluator{workflowInstRepo: workflowInstRepo, exprEvaluator: exprEvaluator}
}

// EvaluateMapByInstCtx 根据实例上下文动态替换 map
func (d *InstExprEvaluator) EvaluateMapByInstCtx(query *dto.GetWorkflowInstDTO, oldMap map[string]interface{}) (
	map[string]interface{}, error) {
	ctx, err := d.workflowInstRepo.GetWorkflowInstCtx(query)
	if err != nil {
		return nil, err
	}

	return d.exprEvaluator.EvaluateMap(ctx, oldMap)
}

// MatchCondition 计算是否匹配
func (d *InstExprEvaluator) MatchCondition(query *dto.GetWorkflowInstDTO, condition string) (bool, error) {
	ctx, err := d.workflowInstRepo.GetWorkflowInstCtx(query)
	if err != nil {
		return false, err
	}

	match, err := d.exprEvaluator.Match(ctx, condition)
	if err != nil {
		return false, err
	}

	return match, nil
}
