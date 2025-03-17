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

// CronEventClient 定时事件客户端
type CronEventClient struct {
	client mq.Client
}

// NewCronEventClient 新建定时事件客户端
func NewCronEventClient(client *tdmq.Client) *CronEventClient {
	return &CronEventClient{client: client}
}

// NewConsumer 创建消费者
func (c *CronEventClient) NewConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return c.client.NewConsumer(ctx, config.GetEventConfig().CronEventTopic,
		group, handle)
}

// SendPresetEvent 发送定时消息
func (c *CronEventClient) SendPresetEvent(ctx context.Context, deliverAt time.Time, msg interface{}) error {
	return c.sendMessage(ctx, deliverAt, 0, msg)
}

// GetEventType 获取事件类型
func (c *CronEventClient) GetEventType(message interface{}) (string, error) {
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

// 发送消息
func (c *CronEventClient) sendMessage(ctx context.Context, deliverAt time.Time, deliverAfter time.Duration,
	msg interface{}) (err error) {
	var msgID string
	startTime := time.Now()
	var pulsarMsg pulsar.ProducerMessage
	defer func() {
		recordTDMQProduceLog(&logs.EventRecord{
			MsgID:     msgID,
			Message:   pulsarMsg,
			BizKey:    "CronEventClient",
			StartTime: startTime,
			Error:     err,
		})
	}()

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to send cron message, caused by %s", err)
		return err
	}
	pulsarMsg = pulsar.ProducerMessage{
		Payload:      payload,
		DeliverAt:    deliverAt,
		DeliverAfter: deliverAfter,
	}
	msgID, err = c.client.SendMessage(ctx, config.GetEventConfig().CronEventTopic, pulsarMsg)
	return err
}
