package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/tdmq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// ExternalEventClient 外部事件客户端
type ExternalEventClient struct {
	topic  string
	client mq.Client
}

// SendEvent 发送流程外部事件
func (c *ExternalEventClient) SendEvent(ctx context.Context, msg interface{}) (err error) {
	startTime := time.Now()
	var msgID string
	var pulsarMsg pulsar.ProducerMessage
	defer func() {
		recordTDMQProduceLog(&logs.EventRecord{
			MsgID:     msgID,
			Message:   pulsarMsg,
			BizKey:    "ExternalEventClient",
			StartTime: startTime,
			Error:     err,
		})
	}()

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to send external message, caused by %s", err)
		return err
	}
	log.Infof("Send external msg:%s", payload)
	routerValue, err := event.GetRouterValue(payload)
	if err != nil {
		log.Errorf("Failed to send external message, caused by %s", err)
		return err
	}
	pulsarMsg = pulsar.ProducerMessage{
		Key:       routerValue,
		Payload:   payload,
		EventTime: time.Now(),
	}
	msgID, err = c.client.SendMessage(ctx, c.topic, pulsarMsg)
	return err
}

// NewConsumer 新建一个消费者
func (c *ExternalEventClient) NewConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return c.client.NewConsumer(ctx, c.topic, group, handle)
}

// NewExternalEventClient 新建事件发送客户端
func NewExternalEventClient(client *tdmq.Client) *ExternalEventClient {
	return &ExternalEventClient{client: client, topic: config.GetEventConfig().ExternalEventTopic}
}

// GetEventType 获取事件类型
func (c *ExternalEventClient) GetEventType(message interface{}) (string, error) {
	msg, ok := (message).(pulsar.Message)
	if !ok {
		return "", fmt.Errorf("not kafka message %v", message)
	}
	eventType, err := event.GetEventType(msg.Payload())
	if err != nil {
		return "", err
	}
	return eventType, nil
}
