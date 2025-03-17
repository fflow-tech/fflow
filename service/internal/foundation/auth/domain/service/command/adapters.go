package command

// Adapters 适配器
type Adapters struct {
	*AuthCommandService
	*RbacCommandService
}

// NewCommandAdapters 初始化适配器
func NewCommandAdapters(authCommandService *AuthCommandService, rbacCommandService *RbacCommandService) *Adapters {
	return &Adapters{AuthCommandService: authCommandService, RbacCommandService: rbacCommandService}
}
