// Package event 提供事件相关的对象
package event

import (
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// ExternalEventType 外部事件类型
type ExternalEventType string

// String 转成字符串
func (t ExternalEventType) String() string {
	return string(t)
}

// 外部事件类型枚举，这些事件主要提供给外部第三方来使用
const (
	DefCreate           ExternalEventType = "DefCreateEvent"           // 流程创建事件
	DefUpdate           ExternalEventType = "DefUpdateEvent"           // 流程更新事件
	DefEnable           ExternalEventType = "DefEnableEvent"           // 流程激活事件
	DefDisable          ExternalEventType = "DefDisableEvent"          // 流程去激活事件
	WorkflowStart       ExternalEventType = "WorkflowStartEvent"       // 流程启动时间
	WorkflowSuccess     ExternalEventType = "WorkflowSuccessEvent"     // 流程执行成功事件
	WorkflowFail        ExternalEventType = "WorkflowFailEvent"        // 流程执行失败事件
	WorkflowCancel      ExternalEventType = "WorkflowCancelEvent"      // 流程执行取消事件
	WorkflowPause       ExternalEventType = "WorkflowPauseEvent"       // 流程执行暂停事件
	WorkflowResume      ExternalEventType = "WorkflowResumeEvent"      // 流程执行恢复事件
	WorkflowTimeout     ExternalEventType = "WorkflowTimeoutEvent"     // 流程执行超时事件
	WorkflowNearTimeout ExternalEventType = "WorkflowNearTimeoutEvent" // 流程执行接近超时事件
	NodeStart           ExternalEventType = "NodeStartEvent"           // 节点开始执行事件
	NodeSuccess         ExternalEventType = "NodeSuccessEvent"         // 节点成功执行事件
	NodeCancel          ExternalEventType = "NodeCancelEvent"          // 节点取消执行事件
	NodeSkip            ExternalEventType = "NodeSkipEvent"            // 节点被跳过事件
	NodeCancelSkip      ExternalEventType = "NodeCancelSkipEvent"      // 节点被取消跳过事件
	NodeWait            ExternalEventType = "NodeWaitEvent"            // 节点等待事件
	NodeAsynWait        ExternalEventType = "NodeAsynWaitEvent"        // 节点异步等待事件
	NodeFail            ExternalEventType = "NodeFailEvent"            // 节点执行失败事件
	NodeTimeout         ExternalEventType = "NodeTimeoutEvent"         // 节点超时事件
	NodeNearTimeout     ExternalEventType = "NodeNearTimeoutEvent"     // 节点接近超时事件
)

var (
	instLevelEvent = map[ExternalEventType]bool{
		WorkflowStart:       true,
		WorkflowSuccess:     true,
		WorkflowFail:        true,
		WorkflowCancel:      true,
		WorkflowPause:       true,
		WorkflowTimeout:     true,
		WorkflowNearTimeout: true,
	}
	nodeInstLevelEvent = map[ExternalEventType]bool{
		NodeStart:       true,
		NodeSuccess:     true,
		NodeCancel:      true,
		NodeSkip:        true,
		NodeCancelSkip:  true,
		NodeWait:        true,
		NodeFail:        true,
		NodeTimeout:     true,
		NodeNearTimeout: true,
		NodeAsynWait:    true,
	}
)

// IsWorkflowInstLevelEvent 判断是否属于流程实例级别的事件
func IsWorkflowInstLevelEvent(eventType ExternalEventType) bool {
	_, ok := instLevelEvent[eventType]
	return ok
}

// IsNodeInstLevelEvent 判断是否属于节点实例级别的事件
func IsNodeInstLevelEvent(eventType ExternalEventType) bool {
	_, ok := nodeInstLevelEvent[eventType]
	return ok
}

// DriveEventType 驱动事件类型
type DriveEventType string

// String 转成字符串
func (t DriveEventType) String() string {
	return string(t)
}

// 内部驱动事件类型枚举
const (
	WorkflowStartDrive DriveEventType = "WorkflowStartDriveEvent" // 流程启动驱动事件
	NodeScheduleDrive  DriveEventType = "NodeScheduleDriveEvent"  // 节点被调度驱动事件
	NodeExecuteDrive   DriveEventType = "NodeExecuteDriveEvent"   // 节点执行驱动事件
	NodePollDrive      DriveEventType = "NodePollDriveEvent"      // 节点轮询驱动事件
	NodeCompleteDrive  DriveEventType = "NodeCompleteDriveEvent"  // 节点完成驱动事件
	NodeRetryDrive     DriveEventType = "NodeRetryDriveEvent"     // 节点重试驱动事件
)

// BasicEvent 基础信息
type BasicEvent struct {
	Namespace   string                 `json:"namespace,omitempty"`
	EventType   string                 `json:"event_type,omitempty"`   // 事件类型 对应各个具体事件枚举
	RouterValue string                 `json:"router_value,omitempty"` // 路由数值 根据这个数据计算路由转发，相同的值发送到相同的分区
	Operator    string                 `json:"operator,omitempty"`     // 操作人
	Reason      string                 `json:"reason,omitempty"`       // 原因
	EventTime   time.Time              `json:"event_time,omitempty"`   // 事件时间, 添加这个的原因是重新消费的时候, pulsar 原生的 EventTime 会丢失
	Biz         map[string]interface{} `json:"biz,omitempty"`          // 业务字段
}

// NewDriveBasicEvent 新建驱动基础事件
// 为了提高并发能力, 将 defID 和 instID 合成为 key
func NewDriveBasicEvent(eventType DriveEventType, namespace string, defID string, instID string) BasicEvent {
	return BasicEvent{
		EventType:   eventType.String(),
		Namespace:   namespace,
		RouterValue: strings.Join([]string{defID, instID}, "_"),
		EventTime:   time.Now(),
	}
}

// NewExternalBasicEvent 新建基础事件
// 如果是流程定义级别的事件，instID 传 0
func NewExternalBasicEvent(eventType ExternalEventType, namespace string, defID string,
	instID string, biz map[string]interface{}) BasicEvent {
	return BasicEvent{
		EventType:   eventType.String(),
		Namespace:   namespace,
		Biz:         biz,
		RouterValue: strings.Join([]string{defID, instID}, "_"),
		EventTime:   time.Now(),
	}
}

// NewExternalBasicEventWithReason 新建基础事件
func NewExternalBasicEventWithReason(eventType ExternalEventType, namespace string, defID string,
	instID string, operator, reason string, biz map[string]interface{}) BasicEvent {
	return BasicEvent{
		EventType:   eventType.String(),
		Namespace:   namespace,
		RouterValue: strings.Join([]string{defID, instID}, "_"),
		Operator:    operator,
		Reason:      reason,
		Biz:         biz,
		EventTime:   time.Now(),
	}
}

// CronEventType 定时时间类型
type CronEventType string

// String 转化为字符串
func (c CronEventType) String() string {
	return string(c)
}

// GetEventType 获取消息的事件类型
func GetEventType(msg []byte) (string, error) {
	return utils.GetStrFromJson(msg, "event_type")
}

// GetEventInstID 获取消息里面的实例 ID
func GetEventInstID(msg []byte) (string, error) {
	return utils.GetStrFromJson(msg, "inst_id")
}

// GetEventDefID 获取消息里面的定义 ID
func GetEventDefID(msg []byte) (string, error) {
	return utils.GetStrFromJson(msg, "def_id")
}

// GetRouterValue 获取消息的路由数值
func GetRouterValue(msg []byte) (string, error) {
	return utils.GetStrFromJson(msg, "router_value")
}

// GetChatMsgFormat 根据事件类型获取聊天消息格式
func GetChatMsgFormat(inst *entity.WorkflowInst, nodeInst *entity.NodeInst, eventType ExternalEventType) string {
	if IsNodeInstLevelEvent(eventType) {
		return getNodeChatMsgFormat(nodeInst, eventType)
	}

	return getWorkflowChatMsgFormat(inst, eventType)
}

func getNodeChatMsgFormat(nodeInst *entity.NodeInst, eventType ExternalEventType) string {
	msg := nodeInst.BasicNodeDef.Msg
	switch eventType {
	case NodeStart:
		return msg.StartMsg
	case NodeSuccess:
		return msg.SuccessMsg
	case NodeFail:
		return msg.FailMsg
	case NodeTimeout:
		return msg.TimeoutMsg
	case NodeNearTimeout:
		return msg.NearTimeoutMsg
	case NodeCancel:
		return msg.CancelMsg
	case NodeWait:
		return msg.WaitMsg
	case NodeAsynWait:
		return msg.AsynWaitMsg
	default:
		return ""
	}
}

func getWorkflowChatMsgFormat(inst *entity.WorkflowInst, eventType ExternalEventType) string {
	msg := inst.WorkflowDef.Msg
	switch eventType {
	case WorkflowStart:
		return msg.StartMsg
	case WorkflowSuccess:
		return msg.SuccessMsg
	case WorkflowFail:
		return msg.FailMsg
	case WorkflowCancel:
		return msg.CancelMsg
	case WorkflowPause:
		return msg.PauseMsg
	case WorkflowResume:
		return msg.ResumeMsg
	case NodeSkip:
		return msg.SkipNodeMsg
	case NodeCancelSkip:
		return msg.CancelSkipNodeMsg
	case WorkflowTimeout:
		return msg.TimeoutMsg
	case WorkflowNearTimeout:
		return msg.NearTimeoutMsg
	default:
		return ""
	}
}
