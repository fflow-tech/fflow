// Package tdmq tdmq 客户端实现发送、创建消费者。
package tdmq

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/fflow-tech/fflow/service/pkg/mq/tdmq"
)

type tdmqClient interface {
	SendMessage(ctx context.Context, topic string, msg interface{}) (string, error)
	NewConsumer(ctx context.Context, topic string, group string,
		handle func(context.Context, interface{}) error) (*tdmq.Consumer, error)
}

// Client 客户端
type Client struct {
	client tdmqClient
}

// NewClient 新建MQ客户端
func NewClient(client *tdmq.Client) *Client {
	return &Client{client: client}
}

// SendMessage 发送消息
func (m *Client) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	return m.client.SendMessage(ctx, topic, msg)
}

// NewConsumer 创建新的消费者
func (m *Client) NewConsumer(ctx context.Context, topic, subName string,
	handle func(context.Context, interface{}) error) (
	mq.Consumer, error) {
	return m.client.NewConsumer(ctx, topic, subName, handle)
}
