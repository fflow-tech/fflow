// Package factory 依赖注入的工厂
package factory

import (
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/memory"
	memorymq "github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/memory"
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
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/provider"
	"github.com/fflow-tech/fflow/service/pkg/registry"
	"github.com/fflow-tech/fflow/service/pkg/remote"
	localsqlite "github.com/fflow-tech/fflow/service/pkg/sqlite"
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

// New 初始化工厂
func New(opts ...Option) error {
	options := NewOptions(opts...)

	if options.RegistryClientType == registry.Kubernetes && options.ConfigClientType == pconfig.Kubernetes {
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

// CreateTables 创建表结构
func CreateTables() error {
	mysqlClient, err := GetMySQLClient()
	if err != nil {
		return fmt.Errorf("failed to get MySQL client: %w", err)
	}

	sqliteClient := &localsqlite.MySQLClient{
		Client: mysqlClient,
	}

	// 创建表
	if err := sqliteClient.CreateTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func provideRepo(container *dig.Container) {
	container.Provide(repo.NewCacheRepoWithMemory)
	container.Provide(repo.NewWorkflowDefRepo)
	container.Provide(repo.NewWorkflowInstRepo)
	container.Provide(repo.NewWorkflowArchiveRepo)
	container.Provide(repo.NewNodeInstRepo)
	container.Provide(repo.NewEventBusRepoWithMemory)
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
	// 使用内存消息队列实现
	container.Provide(memorymq.NewClient)
	container.Provide(memorymq.NewExternalEventClient)
	container.Provide(memorymq.NewTriggerEventClient)
	container.Provide(memorymq.NewCronEventClient)
	container.Provide(memorymq.NewDriveEventClient)
}

func provideDAO(container *dig.Container) {
	container.Provide(config.GetMySQLConfig)
	container.Provide(localsqlite.GetMySQLClient)

	// 注册DAO
	container.Provide(sql.NewWorkflowDefDAO)
	container.Provide(sql.NewWorkflowInstDAO)
	container.Provide(sql.NewHistoryWorkflowInstDAO)
	container.Provide(sql.NewHistoryNodeInstDAO)
	container.Provide(sql.NewNodeInstDAO)
	container.Provide(sql.NewTriggerDAO)

	// 使用内存缓存替代Redis
	container.Provide(memory.NewCacheDAO)

	// 替换远程调用为本地mock
	container.Provide(config.GetDefaultAbilityCallerConfig)
	container.Provide(remote.NewDefaultAbilityCaller)
	container.Provide(remote.NewDefaultCronClient)
	container.Provide(remote.NewDefaultChatOpsClient)
	container.Provide(remote.NewDefaultCloudEventClient)
	container.Provide(remote.NewDefaultPermissionValidator)
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

// GetMySQLClient 获取MySQL客户端
func GetMySQLClient() (*mysql.Client, error) {
	r := &mysql.Client{}
	if err := container.Invoke(func(t *mysql.Client) {
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

// GetWorkflowDefRepo 获取工作流定义仓储
func GetWorkflowDefRepo() (*repo.WorkflowDefRepo, error) {
	r := &repo.WorkflowDefRepo{}
	if err := container.Invoke(func(t *repo.WorkflowDefRepo) {
		r = t
	}); err != nil {
		return nil, err
	}
	return r, nil
}

// GetWorkflowInstRepo 获取工作流实例仓储
func GetWorkflowInstRepo() (*repo.WorkflowInstRepo, error) {
	r := &repo.WorkflowInstRepo{}
	if err := container.Invoke(func(t *repo.WorkflowInstRepo) {
		r = t
	}); err != nil {
		return nil, err
	}
	return r, nil
}

// GetNodeInstRepo 获取节点实例仓储
func GetNodeInstRepo() (*repo.NodeInstRepo, error) {
	r := &repo.NodeInstRepo{}
	if err := container.Invoke(func(t *repo.NodeInstRepo) {
		r = t
	}); err != nil {
		return nil, err
	}
	return r, nil
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
