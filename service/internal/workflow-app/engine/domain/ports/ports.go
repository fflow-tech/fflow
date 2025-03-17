// Package ports 领域层和外部交互的接口
package ports

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
)

// CommandPorts 写接口
type CommandPorts interface {
	WorkflowDefCommandPorts
	WorkflowInstCommandPorts
	NodeInstCommandPorts
	ExternalEventCommandPorts
	WorkflowTriggerCommandPorts
}

// QueryPorts 查接口
type QueryPorts interface {
	WorkflowDefQueryPorts
	WorkflowInstQueryPorts
	NodeInstQueryPorts
}

// WorkflowDefCommandPorts 流程定义接口
type WorkflowDefCommandPorts interface {
	CreateWorkflowDef(context.Context, *dto.CreateWorkflowDefDTO) (string, error) // 创建流程
	UpdateWorkflowDef(context.Context, *dto.CreateWorkflowDefDTO) error           // 更新流程
	EnableWorkflowDef(context.Context, *dto.EnableWorkflowDefDTO) error           // 激活流程
	DisableWorkflowDef(context.Context, *dto.DisableWorkflowDefDTO) error         // 去激活流程
	UploadWorkflowDef(context.Context, *dto.UploadWorkflowDefDTO) (string, error) // 上传流程
}

// WorkflowDefQueryPorts 流程定义接口
type WorkflowDefQueryPorts interface {
	// GetWorkflowDefByDefID 查询流程定义
	GetWorkflowDefByDefID(context.Context, *dto.GetWorkflowDefDTO) (*dto.WorkflowDefDTO, error)
	// GetWorkflowDefList 查询流程列表
	GetWorkflowDefList(context.Context, *dto.PageQueryWorkflowDefDTO) ([]*dto.WorkflowDefDTO, int64, error)
	// GetSubworkflowByParentDefIDAndRefName 获取子流程定义
	GetSubworkflowByParentDefIDAndRefName(context.Context, *dto.GetSubworkflowDefDTO) (*dto.WorkflowDefDTO, error)
}

// WorkflowInstQueryPorts 流程实例接口
type WorkflowInstQueryPorts interface {
	// GetWorkflowInst 查询流程实例
	GetWorkflowInst(context.Context, *dto.GetWorkflowInstDTO) (*dto.WorkflowInstDTO, error)
	// GetWorkflowInstList 查询流程实例列表
	GetWorkflowInstList(context.Context, *dto.GetWorkflowInstListDTO) ([]*dto.WorkflowInstDTO, int64, error)
}

// WorkflowInstCommandPorts 流程实例接口
type WorkflowInstCommandPorts interface {
	StartWorkflowInst(context.Context, *dto.StartWorkflowInstDTO) (string, error)           // 启动流程实例
	RestartWorkflowInst(context.Context, *dto.RestartWorkflowInstDTO) error                 // 重启流程实例
	CancelWorkflowInst(context.Context, *dto.CancelWorkflowInstDTO) error                   // 取消流程实例
	CompleteWorkflowInst(context.Context, *dto.CompleteWorkflowInstDTO) error               // 标记流程实例完成
	PauseWorkflowInst(context.Context, *dto.PauseWorkflowInstDTO) error                     // 暂停流程实例
	ResumeWorkflowInst(context.Context, *dto.ResumeWorkflowInstDTO) error                   // 恢复流程实例
	UpdateWorkflowInstCtx(context.Context, *dto.UpdateWorkflowInstCtxDTO) error             // 更新流程上下文
	ConsumeWorkflowStartDriveEvent(context.Context, *dto.DriveEventDTO) error               // 消费流程启动驱动事件
	ConsumeNodeCompleteDriveEvent(context.Context, *dto.DriveEventDTO) error                // 消费节点完成驱动事件
	DebugWorkflowInst(context.Context, *dto.DebugWorkflowInstDTO) error                     // 更新流程调试信息
	CheckTimeout(context.Context) error                                                     // 检查超时情况
	ArchiveHistoryWorkflowInsts(context.Context, *dto.ArchiveHistoryWorkflowInstsDTO) error // 归档流程实例
}

// NodeInstCommandPorts 节点实例接口
type NodeInstCommandPorts interface {
	ConsumeNodeScheduleDriveEvent(context.Context, *dto.DriveEventDTO) error // 消费节点被调度驱动事件
	ConsumeNodeExecuteDriveEvent(context.Context, *dto.DriveEventDTO) error  // 消费节点执行驱动事件
	ConsumeNodePollDriveEvent(context.Context, *dto.DriveEventDTO) error     // 消费节点轮询驱动事件
	ConsumeNodeRetryDriveEvent(context.Context, *dto.DriveEventDTO) error    // 消费节点重试驱动事件
	RerunNode(context.Context, *dto.RerunNodeDTO) error                      // 重跑节点
	SkipNode(context.Context, *dto.SkipNodeDTO) error                        // 跳过节点
	CancelSkipNode(context.Context, *dto.CancelSkipNodeDTO) error            // 取消跳过节点
	ResumeNode(context.Context, *dto.ResumeNodeDTO) error                    // 恢复节点执行
	CancelNode(context.Context, *dto.CancelNodeDTO) error                    // 取消节点执行
	CompleteNode(context.Context, *dto.CompleteNodeDTO) error                // 标记节点执行完成
	SetTimeout(context.Context, *dto.SetNodeTimeoutDTO) error                // 标记节点超时
	SetNearTimeout(context.Context, *dto.SetNodeNearTimeoutDTO) error        // 标记节点接近超时
}

// NodeInstQueryPorts 节点实例接口
type NodeInstQueryPorts interface {
	GetNodeInstDetail(context.Context, *dto.GetNodeInstDTO) (*dto.NodeInstDTO, error) // 获取节点实例详情
}

// ExternalEventCommandPorts 外部事件处理接口
type ExternalEventCommandPorts interface {
	ConsumeForSendWebhook(context.Context, *dto.ExternalEventDTO) error               // 消费事件中定义的 webhooks
	ConsumeForSendChatMsg(context.Context, *dto.ExternalEventDTO) error               // 发送聊天消息
	ConsumeForWorkflowExceptionHappened(context.Context, *dto.ExternalEventDTO) error // 工作流异常消费
	ConsumeForNodeInstExceptionHappened(context.Context, *dto.ExternalEventDTO) error // 节点异常消费
}

// WorkflowTriggerCommandPorts 触发器处理接口
type WorkflowTriggerCommandPorts interface {
	ConsumeCronTriggerEvent(context.Context, *dto.CronTriggerEventDTO) error // 消费定时器触发事件
	ConsumeTriggerEvent(context.Context, *dto.TriggerEventDTO) error         // 消费触发器事件
	CronCallBack(context.Context, string) error                              // 定时器回调方法
}
