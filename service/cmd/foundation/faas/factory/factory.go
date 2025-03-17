// Package factory 依赖注入的工厂
package factory

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/pkg/remote"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/query"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/repository/repo"
	pconfig "github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/consul"
	"github.com/fflow-tech/fflow/service/pkg/k8s"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/provider"
	"github.com/fflow-tech/fflow/service/pkg/registry"
	"go.uber.org/dig"
)

var (
	container = dig.New()
)

// Options 选项配置
type Options struct {
	RegistryClientType registry.ProviderType
	ConfigClientType   pconfig.ProviderType
	ConsulConfig       consul.Config
	K8sConfig          k8s.Config
}

// Option 选项方法
type Option func(*Options)

// NewOptions 初始化
func NewOptions(opts ...Option) Options {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithRegistryClientType 注册中心类型
func WithRegistryClientType(t registry.ProviderType) Option {
	return func(o *Options) {
		o.RegistryClientType = t
	}
}

// WithConfigClientType 配置中心类型
func WithConfigClientType(t pconfig.ProviderType) Option {
	return func(o *Options) {
		o.ConfigClientType = t
	}
}

// WithConsulConfig Consul 配置
func WithConsulConfig(c consul.Config) Option {
	return func(o *Options) {
		o.ConsulConfig = c
	}
}

// WithK8sConfig K8s 配置
func WithK8sConfig(c k8s.Config) Option {
	return func(o *Options) {
		o.K8sConfig = c
	}
}

// New 初始化工厂
func New(opts ...Option) error {
	options := NewOptions(opts...)

	if options.RegistryClientType == registry.Consul && options.ConfigClientType == pconfig.Consul {
		if err := initConsulClient(options); err != nil {
			return err
		}
	} else if options.RegistryClientType == registry.Kubernetes && options.ConfigClientType == pconfig.Kubernetes {
		if err := initK8sClient(options); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupport registry and config client type, "+
			"RegistryClientType=%s, ConfigClientType=%s", options.RegistryClientType, options.ConfigClientType)
	}

	// 0. 提供底层的DAO
	provideDAO(container)
	// 1. 提供仓储层
	provideRepo(container)
	// 2. 提供执行层的基础能力
	provideExecutionService(container)
	// 3. 提供领域服务
	provideDomainService(container)
	return nil
}

func initConsulClient(options Options) error {
	consulConfig := options.ConsulConfig
	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return err
	}

	provider.InjectConfigProvider(consulClient)
	provider.InjectRegistryProvider(consulClient)
	return nil
}

func initK8sClient(options Options) error {
	k8sConfig := options.K8sConfig
	k8sClient, err := k8s.NewClient(k8sConfig)
	if err != nil {
		return err
	}

	provider.InjectConfigProvider(k8sClient)
	provider.InjectRegistryProvider(k8sClient)
	return nil
}

func provideRepo(container *dig.Container) {
	container.Provide(repo.NewFunctionRepo)
	container.Provide(ports.NewRepoSet)
}

func provideDomainService(container *dig.Container) {
	container.Provide(command.NewFunctionCommandService)
	container.Provide(query.NewFunctionQueryService)

	container.Provide(command.NewCommandAdapters)
	container.Provide(query.NewQueryAdapters)

	container.Provide(service.NewDomainService)
}

func provideDAO(container *dig.Container) {
	container.Provide(config.GetMySQLConfig)
	container.Provide(mysql.GetClient)
	container.Provide(sql.NewFunctionDAO)
	container.Provide(sql.NewRunHistoryDAO)

	container.Provide(config.GetDefaultPermissionValidatorConfig)
	container.Provide(remote.NewDefaultPermissionValidator)
}

func provideExecutionService(container *dig.Container) {
	container.Provide(execution.NewCodeExecutor)
}
