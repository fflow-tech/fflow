package repo

import (
	"context"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/stretchr/testify/assert"
)

type mockTimerTaskEventClient struct {
}

func (m *mockTimerTaskEventClient) SendEvent(ctx context.Context, msg interface{}) error {
	return nil
}
func (m *mockTimerTaskEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return nil, nil
}

type mockPollingTaskEventClient struct {
}

func (m *mockPollingTaskEventClient) SendEvent(ctx context.Context, msg interface{}) error {
	return nil
}
func (m *mockPollingTaskEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return nil, nil
}

func Test_EventBusRepo_SendPollingEvent(t *testing.T) {
	mockRepo := &EventBusRepo{
		timerTaskEventClient:   &mockTimerTaskEventClient{},
		pollingTaskEventClient: &mockPollingTaskEventClient{},
	}

	assert.Nil(t, mockRepo.SendPollingEvent(context.Background(), "test"))
}

func Test_EventBusRepo_NewPollingConsumer(t *testing.T) {
	mockRepo := &EventBusRepo{
		timerTaskEventClient:   &mockTimerTaskEventClient{},
		pollingTaskEventClient: &mockPollingTaskEventClient{},
	}

	_, err := mockRepo.NewPollingConsumer(context.Background(), "test", nil)
	assert.Nil(t, err)
}

func Test_EventBusRepo_SendTimerTaskEvent(t *testing.T) {
	mockRepo := &EventBusRepo{
		timerTaskEventClient:   &mockTimerTaskEventClient{},
		pollingTaskEventClient: &mockPollingTaskEventClient{},
	}

	assert.Nil(t, mockRepo.SendTimerTaskEvent(context.Background(), "test"))
}

func Test_EventBusRepo_NewTimerTaskConsumer(t *testing.T) {
	mockRepo := &EventBusRepo{
		timerTaskEventClient:   &mockTimerTaskEventClient{},
		pollingTaskEventClient: &mockPollingTaskEventClient{},
	}

	_, err := mockRepo.NewTimerTaskConsumer(context.Background(), "test", nil)
	assert.Nil(t, err)
}
