package memory

import (
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
)

// Message 是内存版本的 pulsar.Message
type Message struct {
	id                pulsar.MessageID
	payload           []byte
	properties        map[string]string
	publishTime       time.Time
	eventTime         time.Time
	topic             string
	producerName      string
	key               string
	orderingKey       string
	redeliveryCount   uint32
	isReplicated      bool
	replicatedFrom    string
	schemaVersion     []byte
	encryptionContext *pulsar.EncryptionContext
	index             *uint64
	brokerPublishTime *time.Time
}

// NewBasicMemoryMessage 创建一个新的 MemoryMessage 并设置 payload
func NewBasicMemoryMessage(topic string, payload []byte) *Message {
	return &Message{
		id:           pulsar.NewMessageID(0, 0, 0, 0),
		payload:      payload,
		publishTime:  time.Now(),
		eventTime:    time.Now(),
		topic:        topic,
		producerName: fmt.Sprintf("producer-%s", topic),
		key:          fmt.Sprintf("key-%s", uuid.New().String()),
		orderingKey:  fmt.Sprintf("ordering-key-%s", uuid.New().String()),
		properties:   map[string]string{},
	}
}

// ID 返回消息的 ID
func (m *Message) ID() pulsar.MessageID {
	return m.id
}

// Payload 返回消息的负载
func (m *Message) Payload() []byte {
	return m.payload
}

// Properties 返回消息的属性
func (m *Message) Properties() map[string]string {
	return m.properties
}

// PublishTime 返回消息的发布时间
func (m *Message) PublishTime() time.Time {
	return m.publishTime
}

// EventTime 返回消息的事件时间
func (m *Message) EventTime() time.Time {
	return m.eventTime
}

// Topic 返回消息的主题
func (m *Message) Topic() string {
	return m.topic
}

// ProducerName 返回生产者的名称
func (m *Message) ProducerName() string {
	return m.producerName
}

// Key 返回消息的键
func (m *Message) Key() string {
	return m.key
}

// OrderingKey 返回消息的排序键
func (m *Message) OrderingKey() string {
	return m.orderingKey
}

// RedeliveryCount 返回消息的重传次数
func (m *Message) RedeliveryCount() uint32 {
	return m.redeliveryCount
}

// IsReplicated 判断消息是否从其他集群复制
func (m *Message) IsReplicated() bool {
	return m.isReplicated
}

// GetReplicatedFrom 返回消息复制的来源集群
func (m *Message) GetReplicatedFrom() string {
	return m.replicatedFrom
}

// GetSchemaValue 返回消息的反序列化值
func (m *Message) GetSchemaValue(v interface{}) error {
	// 这里需要实现反序列化逻辑
	return nil
}

// SchemaVersion 返回消息的 schema 版本
func (m *Message) SchemaVersion() []byte {
	return m.schemaVersion
}

// GetEncryptionContext 返回消息的加密上下文
func (m *Message) GetEncryptionContext() *pulsar.EncryptionContext {
	return m.encryptionContext
}

// Index 返回消息的索引
func (m *Message) Index() *uint64 {
	return m.index
}

// BrokerPublishTime 返回 broker 的发布时间
func (m *Message) BrokerPublishTime() *time.Time {
	return m.brokerPublishTime
}
