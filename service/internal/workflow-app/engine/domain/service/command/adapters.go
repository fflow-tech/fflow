// Package command 负责处理写请求
package command

// Adapters 适配器
type Adapters struct {
	*WorkflowDefCommandService
	*WorkflowInstCommandService
	*NodeInstCommandService
	*ExternalEventCommandService
	*WorkflowTriggerCommandService
}

// NewCommandAdapters 初始化适配器
func NewCommandAdapters(defService *WorkflowDefCommandService,
	instService *WorkflowInstCommandService,
	nodeInstService *NodeInstCommandService,
	externalEventService *ExternalEventCommandService,
	triggerCommandService *WorkflowTriggerCommandService,
) *Adapters {
	return &Adapters{WorkflowDefCommandService: defService,
		WorkflowInstCommandService:    instService,
		NodeInstCommandService:        nodeInstService,
		ExternalEventCommandService:   externalEventService,
		WorkflowTriggerCommandService: triggerCommandService,
	}
}
