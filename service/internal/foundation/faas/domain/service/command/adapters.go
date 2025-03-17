package command

// Adapters 适配器
type Adapters struct {
	*FunctionCommandService
}

// NewCommandAdapters 初始化适配器
func NewCommandAdapters(funcService *FunctionCommandService) *Adapters {
	return &Adapters{FunctionCommandService: funcService}
}
