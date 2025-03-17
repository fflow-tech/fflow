// Package mq 消息队列对外接口定义。
package mq

import (
	"context"
)

// Client  抽象接口定义 两个方法里面的interface还是会需要用户去适配的
type Client interface {
	SendMessage(ctx context.Context, topic string, msg interface{}) (string, error)
	NewConsumer(ctx context.Context, topic, subName string,
		handle func(context.Context, interface{}) error) (Consumer, error)
}

// Consumer Consumer接收器 当不再使用时调用关闭
type Consumer interface {
	Close() error
}

// PollingEventClient 驱动事件客户端接口
type PollingEventClient interface {
	SendEvent(ctx context.Context, msg interface{}) error
	NewConsumer(ctx context.Context, group string,
		handle func(context.Context, interface{}) error) (Consumer, error)
}

// TimerTaskEventClient 定时器任务事件客户端接口
type TimerTaskEventClient interface {
	SendEvent(ctx context.Context, msg interface{}) error
	NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (Consumer, error)
}
