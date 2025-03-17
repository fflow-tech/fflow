package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/tdmq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// TriggerEventClient 触发器事件客户端
type TriggerEventClient struct {
	topic  string
	client mq.Client
}

// NewTriggerEventClient 新建触发器事件客户端
func NewTriggerEventClient(client *tdmq.Client) *TriggerEventClient {
	return &TriggerEventClient{client: client, topic: config.GetEventConfig().TriggerEventTopic}
}

// NewConsumer 新建一个消费者
func (c *TriggerEventClient) NewConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return c.client.NewConsumer(ctx, c.topic, group, handle)
}

// SendEvent 发送事件
func (c *TriggerEventClient) SendEvent(ctx context.Context, key string, value interface{}) (err error) {
	startTime := time.Now()
	var msgID string
	var pulsarMsg pulsar.ProducerMessage
	defer func() {
		recordTDMQProduceLog(&logs.EventRecord{
			MsgID:     msgID,
			Message:   pulsarMsg,
			BizKey:    "TriggerEventClient",
			StartTime: startTime,
			Error:     err,
		})
	}()

	payload, err := json.Marshal(value)
	if err != nil {
		log.Errorf("Failed to send trigger message, caused by %s", err)
		return err
	}

	pulsarMsg = pulsar.ProducerMessage{
		Key:       key,
		Payload:   payload,
		EventTime: time.Now(),
	}
	msgID, err = c.client.SendMessage(ctx, c.topic, pulsarMsg)
	return err
}
