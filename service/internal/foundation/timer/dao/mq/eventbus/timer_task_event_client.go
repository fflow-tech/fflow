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

// TimerTaskEventClient 驱动事件客户端
type TimerTaskEventClient struct {
	topic    string
	client   mq.Client
	recorder Recorder
}

// NewTimerTaskEventClient 新建事件发送客户端
func NewTimerTaskEventClient(client *tdmq.Client, recorder *LogRecorder) *TimerTaskEventClient {
	return &TimerTaskEventClient{
		client:   client,
		topic:    config.GetEventConfig().TimerEventTopic,
		recorder: recorder,
	}
}

// SendEvent 发送流程事件
func (c *TimerTaskEventClient) SendEvent(ctx context.Context, msg interface{}) error {
	return c.sendMessage(ctx, time.Time{}, 0, msg)
}

// sendMessage 发送消息
func (c *TimerTaskEventClient) sendMessage(ctx context.Context, deliverAt time.Time, deliverAfter time.Duration,
	msg interface{}) (err error) {
	startTime := time.Now()
	var msgID string
	var pulsarMsg pulsar.ProducerMessage
	defer func() {
		c.recorder.RecordTDMQProduceLog(&logs.EventRecord{
			MsgID:     msgID,
			Message:   pulsarMsg,
			BizKey:    "TimerTaskEventClient",
			StartTime: startTime,
			Error:     err,
		})
	}()

	payload, ok := msg.(string)
	if !ok {
		log.Errorf("no string type of msg: %T", msg)
		return fmt.Errorf("no string type of msg: %T", msg)
	}
	log.Infof("Send task msg:%s", payload)
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
func (c *TimerTaskEventClient) NewConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return c.client.NewConsumer(ctx, c.topic, group, handle)
}
