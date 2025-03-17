package tdmq

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mq/tdmq"
)

// Client tdmq 客户端
type Client struct {
	client *tdmq.Client
}

// NewClient 初始化客户端
func NewClient(client *tdmq.Client) *Client {
	return &Client{client: client}
}

// SendMessage 发送消息
func (m *Client) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	msgID, err := m.client.SendMessage(ctx, topic, msg)
	if err != nil {
		// 如果第一次发送失败了，重试一次发送
		go func() {
			time.Sleep(3 * time.Second)
			msgID, err := m.client.SendMessage(ctx, topic, msg)
			if err != nil {
				log.Errorf("Failed to retry send pulsar msg:%+v, caused by %s", msg, err)
				return
			}
			log.Infof("Retry send pulsar msg:%+v, msgID:%s", msg, msgID)
		}()
	}
	return msgID, err
}

// NewConsumer 创建新的消费者
func (m *Client) NewConsumer(ctx context.Context, topic, subName string,
	handle func(context.Context, interface{}) error) (
	mq.Consumer, error) {
	return m.client.NewConsumer(ctx, topic, subName, handle)
}
