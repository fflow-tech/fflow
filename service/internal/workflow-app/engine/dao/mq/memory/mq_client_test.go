package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	client := NewClient()
	ctx := context.Background()
	topic := "test-topic"
	messages := []string{"test-message-1", "test-message-2", "test-message-3"}

	var wg sync.WaitGroup
	wg.Add(len(messages))

	// 使用消费者来验证消息是否被发送
	consumer, err := client.NewConsumer(ctx, topic, "subName", func(ctx context.Context, message interface{}) error {
		defer wg.Done()
		log.Infof("TestSendMessage consumer message: %v", message)
		return nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, consumer)

	// 发送多个消息
	for _, msg := range messages {
		_, err := client.SendMessage(ctx, topic, msg)
		assert.NoError(t, err)
	}

	// 等待消费者处理完所有消息
	wg.Wait()
}

func TestNewConsumer(t *testing.T) {
	client := NewClient()
	ctx := context.Background()
	topic := "test-topic"
	msg := "test-message"

	// Send a message first
	_, err := client.SendMessage(ctx, topic, msg)
	assert.NoError(t, err)

	// Create a consumer and handle messages
	consumer, err := client.NewConsumer(ctx, topic, "subName", func(ctx context.Context, message interface{}) error {
		assert.Equal(t, msg, message)
		return nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestDriveEventClient_SendEvent(t *testing.T) {
	client := NewDriveEventClient()
	ctx := context.Background()
	msg := "drive-event-message"

	err := client.SendEvent(ctx, msg)
	assert.NoError(t, err)
}

func TestDriveEventClient_SendDelayEvent(t *testing.T) {
	client := NewDriveEventClient()
	ctx := context.Background()
	msg := "delay-event-message"
	delay := 1 * time.Second

	start := time.Now()
	err := client.SendDelayEvent(ctx, delay, msg)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, duration, delay)
}

func TestDriveEventClient_SendPresetEvent(t *testing.T) {
	client := NewDriveEventClient()
	ctx := context.Background()
	msg := "preset-event-message"
	deliverAt := time.Now().Add(1 * time.Second)

	start := time.Now()
	err := client.SendPresetEvent(ctx, deliverAt, msg)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, duration, 1*time.Second)
	consumer, err := client.NewConsumer(ctx, "drive-event", func(ctx context.Context, message interface{}) error {
		assert.Equal(t, msg, message)
		return nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestExternalEventClient_SendEvent(t *testing.T) {
	client := NewExternalEventClient()
	ctx := context.Background()
	msg := "external-event-message"

	err := client.SendEvent(ctx, msg)
	assert.NoError(t, err)
}

func TestCronEventClient_SendPresetEvent(t *testing.T) {
	client := NewCronEventClient()
	ctx := context.Background()
	msg := "cron-event-message"
	deliverAt := time.Now().Add(1 * time.Second)

	start := time.Now()
	err := client.SendPresetEvent(ctx, deliverAt, msg)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, duration, 1*time.Second)
	consumer, err := client.NewConsumer(ctx, "cron-event", func(ctx context.Context, message interface{}) error {
		assert.Equal(t, msg, message)
		return nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestTriggerEventClient_SendEvent(t *testing.T) {
	client := NewTriggerEventClient()
	ctx := context.Background()
	key := "trigger-key"
	value := "trigger-event-message"

	err := client.SendEvent(ctx, key, value)
	assert.NoError(t, err)
}
