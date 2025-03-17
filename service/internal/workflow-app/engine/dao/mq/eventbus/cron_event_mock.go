package eventbus

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
)

// CronEventClientMock 定时事件客户端
type CronEventClientMock struct {
	client mq.Client
}

// NewCronEventClientMock 新建定时事件客户端
func NewCronEventClientMock() *CronEventClientMock {
	return &CronEventClientMock{}
}

// NewConsumer 创建消费者
func (c *CronEventClientMock) NewConsumer(ctx context.Context,
	group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return nil, nil
}

// SendEvent 发送事件
func (c *CronEventClientMock) SendEvent(ctx context.Context, msg interface{}) error {
	return nil
}

// SendPresetEvent 发送定时消息
func (c *CronEventClientMock) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	return nil
}

// GetEventType 获取事件类型
func (c *CronEventClientMock) GetEventType(msg interface{}) (string, error) {
	return "", nil
}
