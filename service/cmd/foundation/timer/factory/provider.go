package factory

import (
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service/monitor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/repo"
)

// GetMonitor 获取监控服务.
func GetMonitor() (*monitor.Monitor, error) {
	m := &monitor.Monitor{}
	if err := container.Invoke(func(_m *monitor.Monitor) {
		m = _m
	}); err != nil {
		return nil, err
	}
	return m, nil
}

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

// GetCommand 获取 command 模块服务.
func GetCommand() (*command.Adapters, error) {
	c := &command.Adapters{}
	if err := container.Invoke(func(_c *command.Adapters) {
		c = _c
	}); err != nil {
		return nil, err
	}
	return c, nil
}
