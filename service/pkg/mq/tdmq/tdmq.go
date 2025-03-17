// Package tdmq tdmq 相关工具包
package tdmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// Client 客户端
type Client struct {
	conf         config.TDMQConfig
	client       pulsar.Client
	producerMap  sync.Map
	producerLock sync.Mutex
	pool         concurrency.WorkerPool
}

// Consumer 消费者
type Consumer struct {
	isClosed bool
	consumer pulsar.Consumer
}

// Close 关闭消费者 这里是为了兼容kafka的close才会包一层 默认close都是有返回error的
func (c *Consumer) Close() error {
	c.isClosed = true
	time.Sleep(3 * time.Second)
	c.consumer.Close()
	return c.consumer.Unsubscribe()
}

func errLog() log.Logger {
	return log.GetDefaultLogger()
}

// GetTDMQClient 获取客户端
func GetTDMQClient(conf config.TDMQConfig) (*Client, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:            conf.URL,
		Authentication: pulsar.NewAuthenticationToken(conf.Authentication),
	})
	if err != nil {
		log.Errorf("Failed to new tdmq client, caused by %s", err)
		return nil, err
	}
	return &Client{client: client, conf: conf, pool: concurrency.GetDefaultWorkerPool()}, nil
}

// NewConsumer 初始化pulsar消费者客户端
func (c *Client) NewConsumer(ctx context.Context, topic, group string,
	handle func(context.Context, interface{}) error) (
	*Consumer, error) {
	consumer, err := c.client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            group,
		Type:                        pulsar.Shared,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionLatest,
		RetryEnable:                 c.conf.RetryEnable,
		DLQ: &pulsar.DLQPolicy{
			MaxDeliveries: c.conf.MaxDeliveries,
		},
	})
	if err != nil {
		return nil, err
	}
	pulsarConsumer := &Consumer{consumer: consumer, isClosed: false}

	if err := c.pool.Submit(func() {
		for {
			if pulsarConsumer.isClosed {
				log.Warnf("Consumer already closed")
				return
			}

			msg, err := consumer.Receive(ctx)
			if err != nil {
				log.Errorf("Failed to receive tdmq msg, caused by %s", err)
				time.Sleep(time.Second)
				continue
			}
			if err := handle(ctx, msg); err != nil {
				if errno.NeedRetryErr(err) {
					delay := c.getReconsumeDuration(msg)
					log.Infof("Set reconsume tdmq msg, msgID:%s, delay:%s", msg.ID(), delay)
					consumer.ReconsumeLater(msg, delay)
					continue
				}
			}
			consumer.Ack(msg)
		}
	}); err != nil {
		return nil, fmt.Errorf("tdmq client register consumer failed, consumer:%+v, err:%w",
			pulsarConsumer, err)
	}

	return pulsarConsumer, nil
}

func (c *Client) getReconsumeDuration(msg pulsar.Message) time.Duration {
	curDuration := time.Duration(c.conf.RetryInitialDelay*uint64(msg.RedeliveryCount())) * time.Second
	maxDuration := time.Duration(c.conf.RetryMaxDelay) * time.Second
	if curDuration.Milliseconds() > maxDuration.Milliseconds() {
		return maxDuration
	}

	return curDuration
}

// NewProducer 初始化pulsar生产者客户端
func (c *Client) NewProducer(topic string) (pulsar.Producer, error) {
	producer, err := c.client.CreateProducer(pulsar.ProducerOptions{
		Topic:           topic,
		DisableBatching: true,
	})

	if err != nil {
		return nil, err
	}

	return producer, nil
}

// SendMessage 发送消息
func (c *Client) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	pulsarMsg, ok := (msg).(pulsar.ProducerMessage)
	if !ok {
		return "", fmt.Errorf("the msg is not pulsar msg, msg:%+v", msg)
	}
	producer, err := c.GetTopicProducer(topic)
	if err != nil {
		log.Errorf("Failed to new pulsar producer, caused by %s, msg:%+v", err, msg)
		return "", fmt.Errorf("failed to new pulsar producer: %w", err)
	}
	msgID, err := producer.Send(ctx, &pulsarMsg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", msgID), nil
}

// GetTopicProducer 获取topic对应的Producer
func (c *Client) GetTopicProducer(topic string) (pulsar.Producer, error) {
	if value, ok := c.producerMap.Load(topic); ok {
		producer, ok := (value).(pulsar.Producer)
		if !ok {
			return nil, fmt.Errorf("the map value is not pulsar Producer, producer:%+v", value)
		}
		return producer, nil
	}

	c.producerLock.Lock()
	defer c.producerLock.Unlock()

	// double check.
	if value, ok := c.producerMap.Load(topic); ok {
		producer, ok := (value).(pulsar.Producer)
		if !ok {
			return nil, fmt.Errorf("the map value is not pulsar Producer, producer:%+v", value)
		}
		return producer, nil
	}

	producer, err := c.NewProducer(topic)
	if err != nil {
		log.Errorf("Failed to new pulsar producer for topic %s, caused by %s", topic, err)
		return nil, fmt.Errorf("failed to new pulsar producer: %w", err)
	}

	log.Infof("create producer of topic: %s successfully", topic)
	c.producerMap.Store(topic, producer)
	return producer, nil
}
