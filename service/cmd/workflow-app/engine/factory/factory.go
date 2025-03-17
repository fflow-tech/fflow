// Package factory 依赖注入的工厂
package factory

import (
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/cache/redis"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/eventbus"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/tdmq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/nodeexecutor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/trigger"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/query"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	pconfig "github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/consul"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/k8s"
	pkafka "github.com/fflow-tech/fflow/service/pkg/mq/kafka"
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

	// 0. 提供底层的 DAO
	provideDAO(container)
	// 1. 提供事件的客户端
	provideEventClient(container)
	// 2. 提供仓储层
	provideRepo(container)
	// 3. 提供执行层的基础能力
	provideExecutionService(container)
	// 4. 初始化触发器层基础能力
	provideTriggerService(container)
	// 5. 提供领域服务
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
	container.Provide(repo.NewCacheRepo)
	container.Provide(repo.NewWorkflowDefRepo)
	container.Provide(repo.NewWorkflowInstRepo)
	container.Provide(repo.NewWorkflowArchiveRepo)
	container.Provide(repo.NewNodeInstRepo)
	container.Provide(repo.NewEventBusRepo)
	container.Provide(repo.NewRemoteRepo)
	container.Provide(repo.NewTriggerRepo)
	container.Provide(ports.NewRepoSet)
}

func provideExecutionService(container *dig.Container) {
	container.Provide(expr.NewDefaultEvaluator)
	container.Provide(nodeexecutor.NewDefaultRegistry)
	container.Provide(execution.NewDefaultWorkflowDecider)
	container.Provide(common.NewDefaultInstTimeoutChecker)
	container.Provide(common.NewDefaultMsgSender)
	container.Provide(execution.NewDefaultWorkflowExecutor)
	container.Provide(execution.NewDefaultWorkflowUpdater)
	container.Provide(execution.NewDefaultNodePoller)
	container.Provide(execution.NewDefaultNodeRunner)
	container.Provide(execution.NewDefaultWorkflowRunner)
	container.Provide(execution.NewDefaultTimeoutChecker)
	container.Provide(execution.NewDefaultHistoryArchiver)
	container.Provide(execution.NewWorkflowProviderSet)
}

func provideTriggerService(container *dig.Container) {
	container.Provide(execution.NewDefaultActor)
	container.Provide(trigger.NewDefaultEventTriggerRegistry)
	container.Provide(trigger.NewDefaultCronTriggerRegistry)
	container.Provide(trigger.NewDefaultRegistry)
}

func provideDomainService(container *dig.Container) {
	container.Provide(command.NewWorkflowDefCommandService)
	container.Provide(command.NewWorkflowInstCommandService)
	container.Provide(command.NewNodeInstCommandService)
	container.Provide(command.NewExternalEventCommandService)
	container.Provide(command.NewWorkflowTriggerCommandService)
	container.Provide(query.NewWorkflowDefQueryService)
	container.Provide(query.NewNodeInstQueryService)
	container.Provide(query.NewWorkflowInstQueryService)

	container.Provide(command.NewCommandAdapters)
	container.Provide(query.NewQueryAdapters)

	container.Provide(service.NewDomainService)
}

func provideEventClient(container *dig.Container) {
	container.Provide(pkafka.GetClient)
	container.Provide(eventbus.NewExternalEventClient)
	container.Provide(eventbus.NewTriggerEventClient)

	container.Provide(config.GetTDMQConfig)
	container.Provide(ptdmq.GetTDMQClient)
	container.Provide(tdmq.NewClient)
	container.Provide(eventbus.NewDriverEventClient)
	container.Provide(eventbus.NewCronEventClient)
}

func provideDAO(container *dig.Container) {
	container.Provide(mysql.GetClient)
	container.Provide(sql.NewWorkflowDefDAO)
	container.Provide(sql.NewWorkflowInstDAO)
	container.Provide(sql.NewHistoryWorkflowInstDAO)
	container.Provide(sql.NewHistoryNodeInstDAO)
	container.Provide(sql.NewNodeInstDAO)
	container.Provide(config.GetRedisConfig)
	container.Provide(config.GetMySQLConfig)
	container.Provide(redisclient.GetClient)
	container.Provide(redis.NewRedisCacheDAO)

	container.Provide(config.GetDefaultAbilityCallerConfig)
	container.Provide(remote.NewDefaultAbilityCaller)
	container.Provide(remote.NewDefaultCronClient)
	container.Provide(remote.NewDefaultChatOpsClient)
	container.Provide(remote.NewDefaultCloudEventClient)
	container.Provide(sql.NewTriggerDAO)

	container.Provide(config.GetDefaultPermissionValidatorConfig)
	container.Provide(remote.NewDefaultPermissionValidator)
}
