package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// ForkNodeExecutor Fork 节点执行器实现
// 实现节点并行的能力，从 Fork 节点拉出来的分支会并行执行
type ForkNodeExecutor struct {
}

// NewForkNodeExecutor 初始化执行器
func NewForkNodeExecutor() *ForkNodeExecutor {
	return &ForkNodeExecutor{}
}

// AsyncComplete 是否异步完成
func (d *ForkNodeExecutor) AsyncComplete(inst *entity.NodeInst) bool {
	return false
}

// AsyncByTrigger 通过触发器异步完成
func (d *ForkNodeExecutor) AsyncByTrigger(nodeInst *entity.NodeInst) bool {
	return false
}

// AsyncByPolling 通过轮询异步完成
func (d *ForkNodeExecutor) AsyncByPolling(nodeInst *entity.NodeInst) bool {
	return false
}

// Execute 执行节点
func (d *ForkNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Polling 轮询节点
func (d *ForkNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *ForkNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// Type 获取是哪种节点类型的处理器
func (d *ForkNodeExecutor) Type() entity.NodeType {
	return entity.ForkNode
}
