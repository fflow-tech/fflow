package repo

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq/eventbus"
)

// EventBusRepo 事件总线实体
type EventBusRepo struct {
	timerTaskEventClient   mq.TimerTaskEventClient
	pollingTaskEventClient mq.PollingEventClient
}

// NewEventBusRepo 实体构造函数
func NewEventBusRepo(d *eventbus.PollingEventClient, e *eventbus.TimerTaskEventClient) *EventBusRepo {
	return &EventBusRepo{timerTaskEventClient: e, pollingTaskEventClient: d}
}

// SendPollingEvent 发送轮询事件
func (e *EventBusRepo) SendPollingEvent(ctx context.Context, msg interface{}) error {
	return e.pollingTaskEventClient.SendEvent(ctx, msg)
}

// NewPollingConsumer 新建轮询监听者
func (e *EventBusRepo) NewPollingConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return e.pollingTaskEventClient.NewConsumer(ctx, group, handle)
}

// SendTimerTaskEvent 发送定时器任务事件
func (e *EventBusRepo) SendTimerTaskEvent(ctx context.Context, msg interface{}) error {
	return e.timerTaskEventClient.SendEvent(ctx, msg)
}

// NewTimerTaskConsumer 新建定时器任务监听者
func (e *EventBusRepo) NewTimerTaskConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return e.timerTaskEventClient.NewConsumer(ctx, group, handle)
}
