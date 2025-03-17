// Package kafka kafka 相关工具包
package kafka

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/patrickmn/go-cache"
	"github.com/segmentio/kafka-go"
)

// Message 消息体
type Message = kafka.Message
type Reader = kafka.Reader

// Client 客户端
type Client struct {
	config.KafkaConfig
	partitionNumCache *cache.Cache
	mutex             sync.Mutex
}

// Consumer 消费者
type Consumer struct {
	isClosed bool
	consumer *Reader
}

// Close 关闭消费者 这里是为了兼容kafka的close才会包一层 默认close都是有返回error的
func (c *Consumer) Close() error {
	c.isClosed = true
	time.Sleep(3 * time.Second)
	return c.consumer.Close()
}

// GetClient 获取客户端
func GetClient(config config.KafkaConfig) *Client {
	return &Client{
		KafkaConfig:       config,
		partitionNumCache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// SendMessage 发送 kafka 消息
func (c *Client) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	c.CreateTopicIfNotExists(topic)

	kafkaMsg, ok := (msg).(Message)
	if !ok {
		return "", fmt.Errorf("msg is not kafka msg %v", msg)
	}
	w := kafka.Writer{
		Addr:         kafka.TCP(getAddress(c.Host, c.Port)),
		Topic:        topic,
		Balancer:     &kafka.CRC32Balancer{},                                   // 相同的 key 会发到相同的 partition
		BatchTimeout: time.Duration(c.ProducerBatchTimeout) * time.Millisecond, // 发送之前的最少等待时间, 框架默认设置为1s, 我们根据实际情况配置
	}

	err := w.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		return "", fmt.Errorf("failed to write messages, topic=%s, host=%s, port=%d: %w",
			topic, c.Host, c.Port, err)
	}

	if err = w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	return buildKafkaMsgID(kafkaMsg), nil
}

// buildKafkaMsgID 构建消息ID
func buildKafkaMsgID(message kafka.Message) string {
	return strings.Join([]string{strconv.Itoa(message.Partition),
		strconv.FormatInt(message.Offset, 10)}, ":")
}

// CreateTopicIfNotExists 如果没有 topic 则创建
func (c *Client) CreateTopicIfNotExists(topic string) error {
	conn, err := kafka.Dial(c.Network, getAddress(c.Host, c.Port))
	if err != nil {
		return fmt.Errorf("failed to create controller conn: %w", err)
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     c.NumPartitions,
			ReplicationFactor: c.ReplicationFactor,
		},
	}

	partitions, err := conn.ReadPartitions(topic)
	if len(partitions) > 0 {
		return nil
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

// NewConsumer 新建一个消费者
func (c *Client) NewConsumer(ctx context.Context, topic, group string,
	handle func(context.Context, interface{}) error) (
	*Consumer, error) {
	readerConfig := getReaderConfig(topic, group, c.KafkaConfig)
	reader := kafka.NewReader(readerConfig)

	kafkaConsumer := &Consumer{consumer: reader, isClosed: false}
	go func() {
		for {
			if kafkaConsumer.isClosed {
				log.Warnf("Consumer already closed")
				return
			}

			message, err := reader.FetchMessage(ctx)
			if err != nil {
				log.ErrorContextf(ctx, "Failed to read kafka message, caused by %s", err)
				time.Sleep(time.Second)
				continue
			}
			if err = handle(ctx, message); err != nil {
				log.ErrorContextf(ctx, "Failed to consume kafka message, caused by %s", err)
				continue
			}
			if err = reader.CommitMessages(ctx, message); err != nil {
				log.ErrorContextf(ctx, "Failed to commit kafka message, caused by %s", err)
			}
		}
	}()

	return kafkaConsumer, nil
}

func getAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func getReaderConfig(topic, group string, config config.KafkaConfig) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:  []string{getAddress(config.Host, config.Port)},
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  time.Duration(config.ConsumerMaxWaitTime) * time.Millisecond,
	}
}
