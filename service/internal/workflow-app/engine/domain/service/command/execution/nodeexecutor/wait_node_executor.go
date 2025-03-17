package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// WaitNodeExecutor 等待节点执行器实现
type WaitNodeExecutor struct {
}

// NewWaitNodeExecutor 初始化
func NewWaitNodeExecutor() *WaitNodeExecutor {
	return &WaitNodeExecutor{}
}

// Execute 执行节点
func (d *WaitNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Polling 轮询节点
func (d *WaitNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *WaitNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// AsyncComplete 是否异步完成
func (d *WaitNodeExecutor) AsyncComplete(inst *entity.NodeInst) bool {
	return false
}

// AsyncByTrigger 通过触发器实现异步
func (d *WaitNodeExecutor) AsyncByTrigger(inst *entity.NodeInst) bool {
	return false
}

// AsyncByPolling 通过轮询实现异步
func (d *WaitNodeExecutor) AsyncByPolling(inst *entity.NodeInst) bool {
	return false
}

// Type 获取是哪种节点类型的处理器
func (d *WaitNodeExecutor) Type() entity.NodeType {
	return entity.WaitNode
}
