// Package eventbus 负责事件发送、监听。
package eventbus

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq/tdmq"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// PollingEventClient 驱动事件客户端
type PollingEventClient struct {
	topic    string
	client   mq.Client
	recorder Recorder
}

// NewDriverEventClient 新建事件发送客户端
func NewDriverEventClient(client *tdmq.Client, recorder *LogRecorder) *PollingEventClient {
	return &PollingEventClient{
		client:   client,
		topic:    config.GetEventConfig().TimerEventTopic,
		recorder: recorder,
	}
}

// SendEvent 发送流程事件
func (c *PollingEventClient) SendEvent(ctx context.Context, msg interface{}) error {
	return c.sendMessage(ctx, time.Time{}, 0, msg)
}

// sendMessage 发送消息
func (c *PollingEventClient) sendMessage(ctx context.Context, deliverAt time.Time, deliverAfter time.Duration,
	msg interface{}) (err error) {
	startTime := time.Now()
	var msgID string
	var pulsarMsg pulsar.ProducerMessage
	defer func() {
		c.recorder.RecordTDMQProduceLog(&logs.EventRecord{
			MsgID:     msgID,
			Message:   pulsarMsg,
			BizKey:    "PollingTaskEventClient",
			StartTime: startTime,
			Error:     err,
		})
	}()

	payload, ok := msg.(string)
	if !ok {
		log.Errorf("no string type of msg: %T", msg)
		return fmt.Errorf("no string type of msg: %T", msg)
	}
	log.Infof("Send polling msg:%s", payload)
	pulsarMsg = pulsar.ProducerMessage{
		Payload:      []byte(payload),
		DeliverAt:    deliverAt,
		DeliverAfter: deliverAfter,
		EventTime:    time.Now(),
	}
	msgID, err = c.client.SendMessage(ctx, c.topic, pulsarMsg)
	return err
}

// NewConsumer 新建一个消费者
func (c *PollingEventClient) NewConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return c.client.NewConsumer(ctx, c.topic, group, handle)
}
