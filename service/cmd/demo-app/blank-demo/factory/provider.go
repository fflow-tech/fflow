package factory

import (
	"github.com/fflow-tech/fflow/service/internal/demo-app/blank-demo/domain/service"
)

// GetDomainService 获取领域服务
func GetDomainService() (*service.DomainService, error) {
	r := &service.DomainService{}
	if err := container.Invoke(func(t *service.DomainService) {
		r = t
	}); err != nil {
		return nil, err
	}
	return r, nil
}
