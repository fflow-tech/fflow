package mq

import (
	"context"
	"time"
)

// Client  抽象接口定义 两个方法里面的 interface 还是会需要用户去适配的
type Client interface {
	SendMessage(ctx context.Context, topic string, msg interface{}) (string, error)
	NewConsumer(ctx context.Context, topic, subName string,
		handle func(context.Context, interface{}) error) (Consumer, error)
}

// Consumer Consumer接收器 当不再使用时调用关闭
type Consumer interface {
	Close() error
}

// DriveEventClient 驱动事件客户端接口
type DriveEventClient interface {
	SendEvent(ctx context.Context, msg interface{}) error
	SendDelayEvent(ctx context.Context, deliverAfter time.Duration, msg interface{}) error
	SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error
	NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (Consumer, error)
	GetEventType(msg interface{}) (string, error)
}

// ExternalEventClient 外部事件客户端接口
type ExternalEventClient interface {
	SendEvent(ctx context.Context, msg interface{}) error
	NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (Consumer, error)
	GetEventType(msg interface{}) (string, error)
}

// CronEventClient 定时事件客户端接口
type CronEventClient interface {
	NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (Consumer, error)
	// SendPresetEvent 发送定时消息
	SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error
	GetEventType(msg interface{}) (string, error)
}

// TriggerEventClient 触发器事件客户端接口
type TriggerEventClient interface {
	NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (Consumer, error)
	SendEvent(ctx context.Context, key string, value interface{}) error
}
