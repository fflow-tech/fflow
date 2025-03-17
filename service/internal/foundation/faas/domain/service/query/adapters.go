package query

// Adapters 查询操作适配器
type Adapters struct {
	*FunctionQueryService
}

// NewQueryAdapters 创建一个查询器
func NewQueryAdapters(functionQueryService *FunctionQueryService) *Adapters {
	return &Adapters{
		FunctionQueryService: functionQueryService,
	}
}
