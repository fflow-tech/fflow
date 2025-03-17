package factory

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/remote"
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

// GetEventBusRepo 获取消费客户端
func GetEventBusRepo() (*repo.EventBusRepo, error) {
	r := &repo.EventBusRepo{}
	if err := container.Invoke(func(t *repo.EventBusRepo) {
		r = t
	}); err != nil {
		return nil, err
	}
	return r, nil
}

// GetDefaultPermissionValidator 获取校验权限服务
func GetDefaultPermissionValidator() (*remote.DefaultPermissionValidator, error) {
	r := &remote.DefaultPermissionValidator{}
	if err := container.Invoke(func(t *remote.DefaultPermissionValidator) {
		r = t
	}); err != nil {
		return nil, err
	}
	return r, nil
}

// GetCacheRepo 获取缓存客户端
func GetCacheRepo() (*repo.CacheRepo, error) {
	r := &repo.CacheRepo{}
	if err := container.Invoke(func(t *repo.CacheRepo) {
		r = t
	}); err != nil {
		return nil, err
	}
	return r, nil
}
