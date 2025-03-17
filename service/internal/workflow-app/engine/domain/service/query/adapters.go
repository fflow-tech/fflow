// Package query 负责处理查请求
package query

// Adapters 查询操作适配器
type Adapters struct {
	*WorkflowDefQueryService
	*WorkflowInstQueryService
	*NodeInstQueryService
}

// NewQueryAdapters 初始化查询适配器
func NewQueryAdapters(workflowDefQueryService *WorkflowDefQueryService,
	workflowInstQueryService *WorkflowInstQueryService,
	nodeInstQueryService *NodeInstQueryService) *Adapters {
	return &Adapters{
		WorkflowDefQueryService:  workflowDefQueryService,
		WorkflowInstQueryService: workflowInstQueryService,
		NodeInstQueryService:     nodeInstQueryService,
	}
}
