package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// JoinNodeExecutor JOIN 节点执行器实现
// 和 Fork 节点配套使用，所有的父节点都执行完成，JOIN 节点才算执行完成
type JoinNodeExecutor struct {
}

// NewJoinNodeExecutor 初始化执行器
func NewJoinNodeExecutor() *JoinNodeExecutor {
	return &JoinNodeExecutor{}
}

// Execute 执行节点
func (d *JoinNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Polling 轮询节点
func (d *JoinNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *JoinNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// AsyncComplete 是否异步完成
func (d *JoinNodeExecutor) AsyncComplete(nodeInst *entity.NodeInst) bool {
	return false
}

// AsyncByTrigger 通过触发器实现异步
func (d *JoinNodeExecutor) AsyncByTrigger(nodeInst *entity.NodeInst) bool {
	return false
}

// AsyncByPolling 通过轮询实现异步
func (d *JoinNodeExecutor) AsyncByPolling(nodeInst *entity.NodeInst) bool {
	return false
}

// Type 获取是哪种节点类型的处理器
func (d *JoinNodeExecutor) Type() entity.NodeType {
	return entity.JoinNode
}
