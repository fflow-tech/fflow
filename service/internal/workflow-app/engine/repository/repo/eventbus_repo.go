package repo

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/eventbus"
)

// EventBusRepo 事件总线仓储层
type EventBusRepo struct {
	driveEventClient    mq.DriveEventClient
	externalEventClient mq.ExternalEventClient
	cronEventClient     mq.CronEventClient
	triggerEventClient  mq.TriggerEventClient
}

// NewEventBusRepo 实体构造函数
func NewEventBusRepo(d *eventbus.DriveEventClient, e *eventbus.ExternalEventClient,
	c *eventbus.CronEventClient, t *eventbus.TriggerEventClient) *EventBusRepo {
	return &EventBusRepo{driveEventClient: d, externalEventClient: e, cronEventClient: c, triggerEventClient: t}
}

// SendDriveEvent 发送驱动事件
func (e *EventBusRepo) SendDriveEvent(ctx context.Context, msg interface{}) error {
	return e.driveEventClient.SendEvent(ctx, msg)
}

// SendDelayDriveEvent 发送延时驱动事件
func (e *EventBusRepo) SendDelayDriveEvent(ctx context.Context, deliverAfter time.Duration, msg interface{}) error {
	return e.driveEventClient.SendDelayEvent(ctx, deliverAfter, msg)
}

// SendPresetDriveEvent 发送定时驱动事件
func (e *EventBusRepo) SendPresetDriveEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	return e.driveEventClient.SendPresetEvent(ctx, deliverAt, msg)
}

// NewDriveEventConsumer 创建驱动事件的消费者
func (e *EventBusRepo) NewDriveEventConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return e.driveEventClient.NewConsumer(ctx, group, handle)
}

// SendExternalEvent 发送外部事件
func (e *EventBusRepo) SendExternalEvent(ctx context.Context, msg interface{}) error {
	return e.externalEventClient.SendEvent(ctx, msg)
}

// NewExternalEventConsumer 创建外部事件的消费者
func (e *EventBusRepo) NewExternalEventConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return e.externalEventClient.NewConsumer(ctx, group, handle)
}

// GetDriveEventType 创建驱动事件的事件类型
func (e *EventBusRepo) GetDriveEventType(message interface{}) (string, error) {
	return e.driveEventClient.GetEventType(message)
}

// GetExternalEventType 获取外部事件类型
func (e *EventBusRepo) GetExternalEventType(message interface{}) (string, error) {
	return e.externalEventClient.GetEventType(message)
}

// SendCronPresetEvent 发送定时事件
func (e *EventBusRepo) SendCronPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	return e.cronEventClient.SendPresetEvent(ctx, deliverAt, msg)
}

// NewCronEventConsumer 创建定时事件消费者
func (e *EventBusRepo) NewCronEventConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return e.cronEventClient.NewConsumer(ctx, group, handle)
}

// GetCronEventType 获取定时事件类型
func (e *EventBusRepo) GetCronEventType(message interface{}) (string, error) {
	return e.cronEventClient.GetEventType(message)
}

// NewTriggerEventConsumer 新建触发器事件消费者
func (e *EventBusRepo) NewTriggerEventConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (
	mq.Consumer, error) {
	return e.triggerEventClient.NewConsumer(ctx, group, handle)
}

// SendTriggerEvent 发送触发器事件
func (e *EventBusRepo) SendTriggerEvent(ctx context.Context, key string, msg interface{}) error {
	return e.triggerEventClient.SendEvent(ctx, key, msg)
}
