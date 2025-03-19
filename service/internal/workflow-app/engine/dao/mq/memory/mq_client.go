package memory

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// Client 内存实现的 Client
type Client struct {
	messages map[string]chan interface{}
	mu       sync.Mutex
}

// NewClient 创建一个新的 Client
func NewClient() *Client {
	return &Client{
		messages: make(map[string]chan interface{}),
	}
}

// SendMessage 发送消息
func (mc *Client) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if _, exists := mc.messages[topic]; !exists {
		mc.messages[topic] = make(chan interface{}, 100) // 使用缓冲通道
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	memoryMsg := NewBasicMemoryMessage(topic, msgBytes)
	mc.messages[topic] <- memoryMsg
	return memoryMsg.ID().String(), nil
}

// NewConsumer 创建一个新的消费者
func (mc *Client) NewConsumer(ctx context.Context, topic, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	mc.mu.Lock()
	ch, exists := mc.messages[topic]
	if !exists {
		ch = make(chan interface{}, 100) // 使用缓冲通道
		mc.messages[topic] = ch
	}
	mc.mu.Unlock()

	go func() {
		for msg := range ch {
			msg := msg.(*Message)
			log.Infof("Consume msg: %v", string(msg.Payload()))
			if err := handle(ctx, msg); err != nil {
				log.Errorf("Failed to handle message: %v", err)
			}
		}
	}()

	return &Consumer{}, nil
}

// Consumer 内存实现的 Consumer
type Consumer struct{}

// Close 关闭消费者
func (mc *Consumer) Close() error {
	return nil
}

// DriveEventClient 内存实现的 DriveEventClient
type DriveEventClient struct {
	client   *Client
	messages map[string]interface{}
}

// NewDriveEventClient 创建一个新的 DriveEventClient
func NewDriveEventClient() *DriveEventClient {
	return &DriveEventClient{
		client: NewClient(),
	}
}

// NewConsumer 创建一个新的消费者
func (mdec *DriveEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return mdec.client.NewConsumer(ctx, "drive-event", group, handle)
}

// SendEvent 发送事件
func (mdec *DriveEventClient) SendEvent(ctx context.Context, msg interface{}) error {
	_, err := mdec.client.SendMessage(ctx, "drive-event", msg)
	return err
}

// SendDelayEvent 发送延迟事件
func (mdec *DriveEventClient) SendDelayEvent(ctx context.Context, deliverAfter time.Duration, msg interface{}) error {
	time.Sleep(deliverAfter)
	return mdec.SendEvent(ctx, msg)
}

// SendPresetEvent 发送预设事件
func (mdec *DriveEventClient) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	time.Sleep(time.Until(deliverAt))
	return mdec.SendEvent(ctx, msg)
}

// GetEventType 获取事件类型
func (mdec *DriveEventClient) GetEventType(message interface{}) (string, error) {
	return getEventType(message)
}

func getEventType(message interface{}) (string, error) {
	msg := message.(*Message)
	msgMap := make(map[string]interface{})
	err := json.Unmarshal(msg.Payload(), &msgMap)
	if err != nil {
		return "", err
	}

	return msgMap["event_type"].(string), nil
}

// ExternalEventClient 内存实现的 ExternalEventClient
type ExternalEventClient struct {
	client *Client
}

// NewExternalEventClient 创建一个新的 ExternalEventClient
func NewExternalEventClient() *ExternalEventClient {
	return &ExternalEventClient{
		client: NewClient(),
	}
}

// NewConsumer 创建一个新的消费者
func (meec *ExternalEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return meec.client.NewConsumer(ctx, "external-event", group, handle)
}

// SendEvent 发送事件
func (meec *ExternalEventClient) SendEvent(ctx context.Context, msg interface{}) error {
	_, err := meec.client.SendMessage(ctx, "external-event", msg)
	return err
}

// GetEventType 获取事件类型
func (meec *ExternalEventClient) GetEventType(message interface{}) (string, error) {
	return getEventType(message)
}

// CronEventClient 内存实现的 CronEventClient
type CronEventClient struct {
	client *Client
}

// NewCronEventClient 创建一个新的 CronEventClient
func NewCronEventClient() *CronEventClient {
	return &CronEventClient{
		client: NewClient(),
	}
}

// NewConsumer 创建一个新的消费者
func (mcec *CronEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return mcec.client.NewConsumer(ctx, "cron-event", group, handle)
}

// SendPresetEvent 发送预设事件
func (mcec *CronEventClient) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	time.Sleep(time.Until(deliverAt))
	_, err := mcec.client.SendMessage(ctx, "cron-event", msg)
	return err
}

// GetEventType 获取事件类型
func (mcec *CronEventClient) GetEventType(message interface{}) (string, error) {
	return getEventType(message)
}

// TriggerEventClient 内存实现的 TriggerEventClient
type TriggerEventClient struct {
	client *Client
}

// NewTriggerEventClient 创建一个新的 TriggerEventClient
func NewTriggerEventClient() *TriggerEventClient {
	return &TriggerEventClient{
		client: NewClient(),
	}
}

// NewConsumer 创建一个新的消费者
func (mtec *TriggerEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return mtec.client.NewConsumer(ctx, "trigger-event", group, handle)
}

// SendEvent 发送事件
func (mtec *TriggerEventClient) SendEvent(ctx context.Context, key string, value interface{}) error {
	_, err := mtec.client.SendMessage(ctx, "trigger-event", value)
	return err
}

// GetEventType 获取事件类型
func (mtec *TriggerEventClient) GetEventType(message interface{}) (string, error) {
	return getEventType(message)
}
