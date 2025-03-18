package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	client := NewClient()
	ctx := context.Background()
	topic := "test-topic"
	msg := "test-message"

	messageID, err := client.SendMessage(ctx, topic, msg)
	assert.NoError(t, err)
	assert.Equal(t, "message-id", messageID)
	assert.Contains(t, client.messages[topic], msg)
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
	assert.Contains(t, client.messages["drive-event"], msg)
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
	assert.Contains(t, client.messages["drive-event"], msg)
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
	assert.Contains(t, client.messages["drive-event"], msg)
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
	assert.Contains(t, client.messages["external-event"], msg)
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
	assert.Contains(t, client.messages["cron-event"], msg)
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
	assert.Contains(t, client.messages["trigger-event"], value)
}
