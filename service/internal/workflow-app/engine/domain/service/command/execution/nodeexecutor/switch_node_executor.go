package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// SwitchNodeExecutor SWITCH 节点执行器实现
type SwitchNodeExecutor struct {
}

// NewSwitchNodeExecutor 初始化
func NewSwitchNodeExecutor() *SwitchNodeExecutor {
	return &SwitchNodeExecutor{}
}

// Execute 执行节点
func (d *SwitchNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Polling 轮询节点
func (d *SwitchNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *SwitchNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// AsyncComplete 是否异步完成
func (d *SwitchNodeExecutor) AsyncComplete(inst *entity.NodeInst) bool {
	return false
}

// AsyncByTrigger 通过触发器实现异步
func (d *SwitchNodeExecutor) AsyncByTrigger(inst *entity.NodeInst) bool {
	return false
}

// AsyncByPolling 通过轮询实现异步
func (d *SwitchNodeExecutor) AsyncByPolling(inst *entity.NodeInst) bool {
	return false
}

// Type 获取是哪种节点类型的处理器
func (d *SwitchNodeExecutor) Type() entity.NodeType {
	return entity.SwitchNode
}
