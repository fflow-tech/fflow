package ports

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/cache"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/remote"
)

// WorkflowDefRepository 流程定义仓储层接口
type WorkflowDefRepository interface {
	// Create 创建流程
	Create(*dto.CreateWorkflowDefDTO) (string, error)
	// BatchCreate 批量创建流程
	BatchCreate([]*dto.CreateWorkflowDefDTO) error
	// Get 获取流程定义
	Get(*dto.GetWorkflowDefDTO) (*entity.WorkflowDef, error)
	// PageQueryLastVersion 分页查询最新版本的流程定义
	PageQueryLastVersion(*dto.PageQueryWorkflowDefDTO) ([]*entity.WorkflowDef, error)
	// PageQueryLastVersionForWeb 提供给页面的分页查询最新版本的流程定义
	PageQueryLastVersionForWeb(*dto.PageQueryWorkflowDefDTO) ([]*entity.WorkflowDef, int64, error)
	// GetLastVersion 获取当前最新版本
	GetLastVersion(*dto.GetWorkflowDefDTO) (*entity.WorkflowDef, error)
	// GetSubworkflowLastVersion 获取子流程最新版本
	GetSubworkflowLastVersion(*dto.GetSubworkflowDefDTO) (*entity.WorkflowDef, error)
	// UpdateStatus 更新流程状态
	UpdateStatus(*dto.UpdateWorkflowDefDTO) error
}

// WorkflowInstRepository 流程实例信息仓储层接口
type WorkflowInstRepository interface {
	Create(*dto.CreateWorkflowInstRepoDTO) (string, error)                        // 创建流程实例
	UpdateWithDefID(*dto.UpdateWorkflowInstDTO) error                             // 根据流程实例
	UpdateWorkflowInstFailed(*dto.UpdateWorkflowInstFailedDTO) error              // 更新流程实例，专门提供给实例失败使用
	Get(*dto.GetWorkflowInstDTO) (*entity.WorkflowInst, error)                    // 获取单个流程实例
	GetWorkflowInstCtx(*dto.GetWorkflowInstDTO) (map[string]interface{}, error)   // 获取流程实例上下文
	PageQuery(*dto.GetWorkflowInstListDTO) ([]*entity.WorkflowInst, int64, error) // 分页查询流程实例
	Count(req *dto.GetWorkflowInstListDTO) (int64, error)                         // 统计流程实例数量
}

// NodeInstRepository 节点实例信息仓储层接口
type NodeInstRepository interface {
	Create(*dto.CreateNodeInstDTO) (string, error)                   // 创建节点实例
	UpdateWithDefID(*dto.UpdateNodeInstDTO) error                    // 更新节点实例
	Get(*dto.GetNodeInstDTO) (*entity.NodeInst, error)               // 获取节点实例
	PageQuery(*dto.PageQueryNodeInstDTO) ([]*entity.NodeInst, error) // 分页查询节点实例
}

// WorkflowArchiveRepository 归档流程仓储层接口
type WorkflowArchiveRepository interface {
	ArchiveWorkflowInst(*dto.ArchiveWorkflowInstsDTO) error // 归档流程实例
	ArchiveNodeInst(*dto.ArchiveNodeInstsDTO) error         // 归档节点实例
}

// EventBusRepository 事件总线仓储层接口
type EventBusRepository interface {
	SendDriveEvent(ctx context.Context, msg interface{}) error                                  // 发送驱动事件
	SendDelayDriveEvent(ctx context.Context, deliverAfter time.Duration, msg interface{}) error // 发送延迟驱动事件
	SendPresetDriveEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error       // 发送定时驱动事件
	GetDriveEventType(message interface{}) (string, error)                                      // 获取驱动事件类型
	SendExternalEvent(ctx context.Context, msg interface{}) error                               // 发送外部事件
	GetExternalEventType(message interface{}) (string, error)                                   // 获取外部事件类型
	SendCronPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error        // 发送定时消息
	GetCronEventType(message interface{}) (string, error)                                       // 获取定时事件类型
	SendTriggerEvent(ctx context.Context, key string, value interface{}) error                  // 发送触发器事件
	NewExternalEventConsumer(ctx context.Context, group string,
		handle func(context.Context, interface{}) error) (mq.Consumer, error) // 初始化外部事件消费者
	NewDriveEventConsumer(ctx context.Context, group string,
		handle func(context.Context, interface{}) error) (mq.Consumer, error) // 初始化驱动事件消费者
	NewCronEventConsumer(ctx context.Context, group string,
		handle func(context.Context, interface{}) error) (mq.Consumer, error) // 初始化定时事件消费者
	NewTriggerEventConsumer(ctx context.Context, group string,
		handle func(context.Context, interface{}) error) (mq.Consumer, error) // 初始化触发器事件消费者
}

// CacheRepository 缓存仓储层接口
type CacheRepository interface {
	Set(key string, value string, ttl int64) error                                // 设置值
	SetNX(key, value string, ttl int64) (interface{}, error)                      // 如果 key 不存在则设置值
	Get(key string) (string, error)                                               // 获取值
	GetDistributeLock(name string, expireTime time.Duration) cache.DistributeLock // 可重入锁
	GetDistributeLockWithRetry(name string, expireTime time.Duration,
		trys int, retryDelay time.Duration) cache.DistributeLock // 可重入锁
}

// RemoteRepository 远程调用仓储层接口
type RemoteRepository interface {
	CallFAAS(context.Context, *remote.CallFAASReqDTO) (map[string]interface{}, error) // 调用 faas 接口
	CallHTTP(context.Context, *remote.CallHTTPReqDTO) (map[string]interface{}, error) // 调用 http 接口
	AddCronJob(*remote.AddCronJobDTO) error                                           // 添加定时任务
	CancelCronJob(jobName string) error                                               // 取消定时任务
	SendMsgToUser(userID, msg string) error                                           // 发送企微消息给用户
	SendMsgToGroup(chatID, msg string) error                                          // 发送企微消息给群聊
	SendCloudEvent(ctx context.Context, req *remote.SendCloudEventDTO) error          // 发送事件
}

// TriggerRepository 触发器仓储层接口
type TriggerRepository interface {
	Create(*dto.CreateTriggerDTO) (string, error)                  // 创建触发器
	Update(*dto.UpdateTriggerDTO) error                            // 更新触发器
	Get(*dto.GetTriggerDTO) (*entity.Trigger, error)               // 获取触发器
	Count(*dto.PageQueryTriggerDTO) (int64, error)                 // 统计触发器数量
	PageQuery(*dto.PageQueryTriggerDTO) ([]*entity.Trigger, error) // 分页查询触发器
	QueryByName(*dto.QueryTriggerDTO) ([]*entity.Trigger, error)   // 根据名称查询触发器
}
