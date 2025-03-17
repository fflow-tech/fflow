package dto

import (
	"github.com/apache/pulsar-client-go/pulsar"
)

// SendCronPresetEventDTO 发送定时事件
type SendCronPresetEventDTO struct {
	Key   string                 `json:"key,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

// SendTriggerEventDTO 发送触发器事件
type SendTriggerEventDTO struct {
	Key   string                 `json:"key,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

// SendDriveEventDTO 发送触发器事件
type SendDriveEventDTO struct {
	Key   string                 `json:"key,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

// SendExternalEventDTO 发送外部事件
type SendExternalEventDTO struct {
	Key   string                 `json:"key,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

// ExternalEventDTO 消费外部事件
type ExternalEventDTO struct {
	Message pulsar.Message
}

// DriveEventDTO 内部驱动事件
type DriveEventDTO struct {
	Message pulsar.Message
}

// CronTriggerEventDTO 定时器触发器事件
type CronTriggerEventDTO struct {
	Message pulsar.Message
}

// TriggerEventDTO 触发器事件
type TriggerEventDTO struct {
	Message pulsar.Message // 事件内容
}
