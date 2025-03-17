package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
)

// TransformNodeExecutor Transform 节点执行器实现
// 实现参数之间的转换
type TransformNodeExecutor struct {
	instExprEvaluator *common.InstExprEvaluator
}

// NewTransformNodeExecutor 初始化执行器
func NewTransformNodeExecutor(instExprEvaluator *common.InstExprEvaluator) *TransformNodeExecutor {
	return &TransformNodeExecutor{instExprEvaluator: instExprEvaluator}
}

// AsyncComplete 是否异步完成
func (d *TransformNodeExecutor) AsyncComplete(inst *entity.NodeInst) bool {
	return false
}

// AsyncByTrigger 通过触发器异步完成
func (d *TransformNodeExecutor) AsyncByTrigger(nodeInst *entity.NodeInst) bool {
	return false
}

// AsyncByPolling 通过轮询异步完成
func (d *TransformNodeExecutor) AsyncByPolling(nodeInst *entity.NodeInst) bool {
	return false
}

// Execute 执行节点
func (d *TransformNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeDef, err := entity.ToActualNodeDef(nodeInst.BasicNodeDef.Type, nodeInst.NodeDef)
	if err != nil {
		return err
	}
	actualNodeDef := nodeDef.(entity.TransformNodeDef)

	query := dto.NewGetWorkflowInstDTO(nodeInst.InstID, nodeInst.DefID, "")
	nodeInst.Output, err = d.instExprEvaluator.EvaluateMapByInstCtx(query, actualNodeDef.Output)
	return err
}

// Polling 轮询节点
func (d *TransformNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *TransformNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// Type 获取是哪种节点类型的处理器
func (d *TransformNodeExecutor) Type() entity.NodeType {
	return entity.TransformNode
}
