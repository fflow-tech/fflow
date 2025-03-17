package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/tdmq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// DriveEventClient 驱动事件客户端
type DriveEventClient struct {
	topic  string
	client mq.Client
}

// SendEvent 发送流程事件
func (c *DriveEventClient) SendEvent(ctx context.Context, msg interface{}) error {
	return c.sendMessage(ctx, time.Time{}, 0, msg)
}

// SendDelayEvent 发送延时事件 多久之后
func (c *DriveEventClient) SendDelayEvent(ctx context.Context, deliverAfter time.Duration,
	msg interface{}) error {
	return c.sendMessage(ctx, time.Time{}, deliverAfter, msg)
}

// SendPresetEvent 发送定时事件 定时消息必须是10天内
func (c *DriveEventClient) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	return c.sendMessage(ctx, deliverAt, 0, msg)
}

// sendMessage 发送消息
func (c *DriveEventClient) sendMessage(ctx context.Context, deliverAt time.Time, deliverAfter time.Duration,
	msg interface{}) (err error) {
	startTime := time.Now()
	var msgID string
	var pulsarMsg pulsar.ProducerMessage
	defer func() {
		recordTDMQProduceLog(&logs.EventRecord{
			MsgID:     msgID,
			Message:   pulsarMsg,
			BizKey:    "DriveEventClient",
			StartTime: startTime,
			Error:     err,
		})
	}()

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to send drive message, caused by %s", err)
		return err
	}
	log.Infof("Send drive msg:%s", payload)
	routerValue, err := event.GetRouterValue(payload)
	if err != nil {
		log.Errorf("Failed to send drive message, caused by %s", err)
		return err
	}
	pulsarMsg = pulsar.ProducerMessage{
		Key:          routerValue,
		Payload:      payload,
		DeliverAt:    deliverAt,
		DeliverAfter: deliverAfter,
		EventTime:    time.Now(),
	}
	msgID, err = c.client.SendMessage(ctx, c.topic, pulsarMsg)
	return err
}

// NewConsumer 新建一个消费者
func (c *DriveEventClient) NewConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return c.client.NewConsumer(ctx, c.topic, group, handle)
}

// NewDriverEventClient 新建事件发送客户端
func NewDriverEventClient(client *tdmq.Client) *DriveEventClient {
	return &DriveEventClient{client: client, topic: config.GetEventConfig().DriveEventTopic}
}

// GetEventType 获取事件类型
func (c *DriveEventClient) GetEventType(message interface{}) (string, error) {
	msg, ok := (message).(pulsar.Message)
	if !ok {
		log.Errorf("Failed to get event type, message=%+v", message)
		return "", nil
	}
	eventType, err := event.GetEventType(msg.Payload())
	if err != nil {
		return "", err
	}
	return eventType, nil
}
