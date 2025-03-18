package memory

import (
	"context"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
)

// Client 内存实现的 Client
type Client struct {
	messages map[string][]interface{}
	mu       sync.Mutex
}

// NewClient 创建一个新的 Client
func NewClient() *Client {
	return &Client{
		messages: make(map[string][]interface{}),
	}
}

// SendMessage 发送消息
func (mc *Client) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.messages[topic] = append(mc.messages[topic], msg)
	return "message-id", nil
}

// NewConsumer 创建一个新的消费者
func (mc *Client) NewConsumer(ctx context.Context, topic, subName string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if msgs, ok := mc.messages[topic]; ok {
		for _, msg := range msgs {
			if err := handle(ctx, msg); err != nil {
				return nil, err
			}
		}
	}
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
	messages map[string][]string
}

// NewDriveEventClient 创建一个新的 DriveEventClient
func NewDriveEventClient() *DriveEventClient {
	return &DriveEventClient{
		client:   NewClient(),
		messages: make(map[string][]string),
	}
}

// NewConsumer 创建一个新的消费者
func (mdec *DriveEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return mdec.client.NewConsumer(ctx, "drive-event", group, handle)
}

// SendEvent 发送事件
func (mdec *DriveEventClient) SendEvent(ctx context.Context, msg string) error {
	if mdec.messages == nil {
		mdec.messages = make(map[string][]string)
	}
	mdec.messages["drive-event"] = append(mdec.messages["drive-event"], msg)
	return nil
}

// SendDelayEvent 发送延迟事件
func (mdec *DriveEventClient) SendDelayEvent(ctx context.Context, deliverAfter time.Duration, msg string) error {
	time.Sleep(deliverAfter)
	return mdec.SendEvent(ctx, msg)
}

// SendPresetEvent 发送预设事件
func (mdec *DriveEventClient) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg string) error {
	time.Sleep(time.Until(deliverAt))
	return mdec.SendEvent(ctx, msg)
}

// GetEventType 获取事件类型
func (mdec *DriveEventClient) GetEventType(msg interface{}) (string, error) {
	return "drive-event-type", nil
}

// ExternalEventClient 内存实现的 ExternalEventClient
type ExternalEventClient struct {
	client   *Client
	messages map[string][]string
}

// NewExternalEventClient 创建一个新的 ExternalEventClient
func NewExternalEventClient() *ExternalEventClient {
	return &ExternalEventClient{
		client:   NewClient(),
		messages: make(map[string][]string),
	}
}

// NewConsumer 创建一个新的消费者
func (meec *ExternalEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return meec.client.NewConsumer(ctx, "external-event", group, handle)
}

// SendEvent 发送事件
func (meec *ExternalEventClient) SendEvent(ctx context.Context, msg string) error {
	if meec.messages == nil {
		meec.messages = make(map[string][]string)
	}
	meec.messages["external-event"] = append(meec.messages["external-event"], msg)
	return nil
}

// GetEventType 获取事件类型
func (meec *ExternalEventClient) GetEventType(msg interface{}) (string, error) {
	return "external-event-type", nil
}

// CronEventClient 内存实现的 CronEventClient
type CronEventClient struct {
	client   *Client
	messages map[string][]string
}

// NewCronEventClient 创建一个新的 CronEventClient
func NewCronEventClient() *CronEventClient {
	return &CronEventClient{
		client:   NewClient(),
		messages: make(map[string][]string),
	}
}

// NewConsumer 创建一个新的消费者
func (mcec *CronEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return mcec.client.NewConsumer(ctx, "cron-event", group, handle)
}

// SendPresetEvent 发送预设事件
func (mcec *CronEventClient) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg string) error {
	time.Sleep(time.Until(deliverAt))
	if mcec.messages == nil {
		mcec.messages = make(map[string][]string)
	}
	mcec.messages["cron-event"] = append(mcec.messages["cron-event"], msg)
	return nil
}

// GetEventType 获取事件类型
func (mcec *CronEventClient) GetEventType(msg interface{}) (string, error) {
	return "cron-event-type", nil
}

// TriggerEventClient 内存实现的 TriggerEventClient
type TriggerEventClient struct {
	client   *Client
	messages map[string][]string
}

// NewTriggerEventClient 创建一个新的 TriggerEventClient
func NewTriggerEventClient() *TriggerEventClient {
	return &TriggerEventClient{
		client:   NewClient(),
		messages: make(map[string][]string),
	}
}

// NewConsumer 创建一个新的消费者
func (mtec *TriggerEventClient) NewConsumer(ctx context.Context, group string, handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return mtec.client.NewConsumer(ctx, "trigger-event", group, handle)
}

// SendEvent 发送事件
func (mtec *TriggerEventClient) SendEvent(ctx context.Context, key string, value string) error {
	if mtec.messages == nil {
		mtec.messages = make(map[string][]string)
	}
	mtec.messages["trigger-event"] = append(mtec.messages["trigger-event"], value)
	return nil
}

// GetEventType 获取事件类型
func (mtec *TriggerEventClient) GetEventType(msg interface{}) (string, error) {
	return "trigger-event-type", nil
}
