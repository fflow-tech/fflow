package eventbus

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
)

// DriveEventClientMock 驱动事件客户端
type DriveEventClientMock struct {
}

// SendEvent 发送流程事件
func (c *DriveEventClientMock) SendEvent(ctx context.Context, msg interface{}) error {
	return nil
}

// SendDelayEvent 发送延时事件 多久之后
func (c *DriveEventClientMock) SendDelayEvent(ctx context.Context, deliverAfter time.Duration, msg interface{}) error {
	return nil
}

// SendPresetEvent 发送定时事件 定时消息必须是10天内
func (c *DriveEventClientMock) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	return nil
}

// 发送消息
func (c *DriveEventClientMock) sendMessage(ctx context.Context, deliverAt time.Time, deliverAfter time.Duration,
	msg interface{}) error {
	return nil
}

// NewConsumer 新建一个消费者
func (c *DriveEventClientMock) NewConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return nil, nil
}

// GetEventType 获取事件类型
func (c *DriveEventClientMock) GetEventType(message interface{}) (string, error) {
	return "", nil
}

// NewDriverEventClientMock 新建事件模拟客户端
func NewDriverEventClientMock() *DriveEventClientMock {
	return &DriveEventClientMock{}
}
