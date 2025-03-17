package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// NodeExecutor 节点执行器接口
type NodeExecutor interface {
	Execute(context.Context, *entity.NodeInst) error
	Polling(context.Context, *entity.NodeInst) error
	Cancel(context.Context, *entity.NodeInst) error
	AsyncByTrigger(*entity.NodeInst) bool
	AsyncByPolling(*entity.NodeInst) bool
	AsyncComplete(*entity.NodeInst) bool
	Type() entity.NodeType
}

// Registry 节点执行器注册中心
type Registry interface {
	Register(e NodeExecutor)
	UnRegister(e NodeExecutor)
	GetExecutor(t entity.NodeType) (NodeExecutor, bool)
}

// NewDefaultRegistry 初始化注册中心
func NewDefaultRegistry() *DefaultRegistry {
	return &DefaultRegistry{executorMap: map[entity.NodeType]NodeExecutor{}}
}

// DefaultRegistry 节点执行器注册中心
type DefaultRegistry struct {
	executorMap map[entity.NodeType]NodeExecutor
}

// Register 注册
func (r *DefaultRegistry) Register(e NodeExecutor) {
	r.executorMap[e.Type()] = e
}

// UnRegister 取消注册
func (r *DefaultRegistry) UnRegister(e NodeExecutor) {
	delete(r.executorMap, e.Type())
}

// GetExecutor 获取执行器
func (r *DefaultRegistry) GetExecutor(t entity.NodeType) (NodeExecutor, bool) {
	e, ok := r.executorMap[t]
	return e, ok
}
