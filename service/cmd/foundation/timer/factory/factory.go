// Package factory 依赖注入的工厂
package factory

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service/monitor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/cache/redis"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq/eventbus"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq/tdmq"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/query"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	reporter "github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/monitor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/repo"
	pconfig "github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/consul"
	"github.com/fflow-tech/fflow/service/pkg/k8s"
	"github.com/fflow-tech/fflow/service/pkg/limiter"
	ptdmq "github.com/fflow-tech/fflow/service/pkg/mq/tdmq"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/provider"
	redisclient "github.com/fflow-tech/fflow/service/pkg/redis"
	"github.com/fflow-tech/fflow/service/pkg/registry"
	"github.com/fflow-tech/fflow/service/pkg/remote"
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

// WithK8sConfig K8s 配置
func WithK8sConfig(c k8s.Config) Option {
	return func(o *Options) {
		o.K8sConfig = c
	}
}

// WithConsulConfig Consul 配置
func WithConsulConfig(c consul.Config) Option {
	return func(o *Options) {
		o.ConsulConfig = c
	}
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
	// 1. 提供配置
	provideConfig(container)
	// 2. 提供依赖的第三方库
	providePKG(container)
	// 3. 提供底层的DAO
	provideDAO(container)
	// 4. 提供事件的客户端
	provideEventClient(container)
	// 5. 提供仓储层
	provideRepo(container)
	// 6. 提供领域服务
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

func provideRepo(container *dig.Container) {
	container.Provide(repo.NewEventBusRepo)
	container.Provide(repo.NewPollingTaskRepo)
	container.Provide(repo.NewTimerDefRepo)
	container.Provide(repo.NewTimerTaskRepo)
	container.Provide(repo.NewRemoteRepo)
	container.Provide(repo.NewAppRepo)
	container.Provide(ports.NewRepoSet)
}
func provideDomainService(container *dig.Container) {
	container.Provide(command.NewPollingTaskCommandService)
	container.Provide(command.NewTimerDefCommandService)
	container.Provide(command.NewTimerTaskCommandService)
	container.Provide(command.NewCommandAdapters)
	container.Provide(command.NewNotifyCommandService)
	container.Provide(command.NewAppCommandService)
	container.Provide(query.NewTimerDefQueryService)
	container.Provide(query.NewTimerTaskQueryService)
	container.Provide(query.NewQueryAdapters)
	container.Provide(query.NewAppQueryService)
	container.Provide(service.NewDomainService)
	container.Provide(monitor.NewMonitor)
}
func provideEventClient(container *dig.Container) {
	container.Provide(eventbus.GetLogRecorder)
	container.Provide(ptdmq.GetTDMQClient)
	container.Provide(tdmq.NewClient)
	container.Provide(eventbus.NewDriverEventClient)
	container.Provide(eventbus.NewTimerTaskEventClient)
}
func provideDAO(container *dig.Container) {
	container.Provide(redisclient.GetClient)
	container.Provide(mysql.GetClient)
	container.Provide(sql.NewRunHistoryDAO)
	container.Provide(sql.NewTimerDefDAO)
	container.Provide(sql.NewAppDAO)
	container.Provide(redis.NewPollingTaskClient)
	container.Provide(redis.NewTimerTaskClient)
	container.Provide(redis.NewTimerDefClient)
	container.Provide(remote.NewDefaultChatOpsClient)
	container.Provide(remote.NewDefaultCloudEventClient)
	container.Provide(config.GetDefaultAbilityCallerConfig)
	container.Provide(remote.NewDefaultAbilityCaller)
}

func provideConfig(container *dig.Container) {
	container.Provide(config.GetRedisConfig)
	container.Provide(config.GetMySQLConfig)
	container.Provide(config.GetDefaultTDMQConfig)
	container.Provide(config.GetLimiterConfig)
}

func providePKG(container *dig.Container) {
	container.Provide(concurrency.GetDefaultWorkerPool)
	container.Provide(limiter.NewTrafficPool)
	container.Provide(reporter.GetReporter)
}
