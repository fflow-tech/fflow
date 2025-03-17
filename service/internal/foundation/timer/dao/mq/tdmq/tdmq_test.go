package tdmq

import (
	"context"
	"testing"

	"github.com/fflow-tech/fflow/service/pkg/mq/tdmq"

	"github.com/stretchr/testify/assert"
)

type mockTDMQClient struct{}

func (m *mockTDMQClient) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	return "", nil
}

func (m *mockTDMQClient) NewConsumer(ctx context.Context, topic string,
	group string, handle func(context.Context, interface{}) error) (*tdmq.Consumer, error) {
	return nil, nil
}

func Test_Client_SendMessage(t *testing.T) {
	mockClient := Client{&mockTDMQClient{}}
	_, err := mockClient.SendMessage(context.Background(), "test", nil)
	assert.Nil(t, err)
}

func Test_Client_NewConsumer(t *testing.T) {
	mockClient := Client{&mockTDMQClient{}}
	_, err := mockClient.NewConsumer(context.Background(), "test", "test", nil)
	assert.Nil(t, err)
}
