package logs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/segmentio/kafka-go"
)

// EventRecord 消息队列记录
type EventRecord struct {
	MsgID     string
	Message   interface{}
	BizKey    string
	StartTime time.Time
	Error     error
	Type      RecordType
}

// RecordType 记录类型
type RecordType string

const (
	Produce RecordType = "EventProduceRecord" // 消息生产记录
	Consume RecordType = "EventConsumeRecord" // 消息消费记录
)

// EventDetailRecord 消息详细记录
type EventDetailRecord struct {
	Type         RecordType `json:"type,omitempty"`
	MsgID        string     `json:"msg_id,omitempty"`
	DefID        string     `json:"def_id,omitempty"`
	InstID       string     `json:"inst_id,omitempty"`
	Ack          bool       `json:"ack,omitempty"`
	FailedReason string     `json:"failed_reason,omitempty"`
	BizKey       string     `json:"biz_key,omitempty"`
	Topic        string     `json:"topic,omitempty"`
	Key          string     `json:"key,omitempty"`
	Value        string     `json:"value,omitempty"`
	Costs        int64      `json:"costs,omitempty"`
}

// GetMqLogName 获取日志插件的名字
func GetMqLogName() string {
	return "mq_log"
}

// RecordTDMQLog 记录 TDMQ 的消费日志
func RecordTDMQLog(record *EventRecord) {
	eventDetail := &EventDetailRecord{}
	setBasicInfo(record, eventDetail)

	switch record.Message.(type) {
	case pulsar.Message:
		recordTDMQConsumeMsg(record.Message.(pulsar.Message), eventDetail)
	case pulsar.ProducerMessage:
		recordTDMQProducerMessage(record.Message.(pulsar.ProducerMessage), eventDetail)
	default:
		return
	}
}

func recordTDMQProducerMessage(message pulsar.ProducerMessage, eventDetail *EventDetailRecord) {
	eventDetail.Key = message.Key
	eventDetail.Value = string(message.Payload)
	eventDetail.InstID, _ = event.GetEventInstID(message.Payload)
	eventDetail.DefID, _ = event.GetEventDefID(message.Payload)

	recordMqLog(eventDetail)
}

func recordTDMQConsumeMsg(message pulsar.Message, eventDetail *EventDetailRecord) {
	eventDetail.MsgID = fmt.Sprintf("%s", message.ID())
	eventDetail.Topic = message.Topic()
	eventDetail.Key = message.Key()
	eventDetail.Value = string(message.Payload())
	eventDetail.InstID, _ = event.GetEventInstID(message.Payload())
	eventDetail.DefID, _ = event.GetEventDefID(message.Payload())

	recordMqLog(eventDetail)
}

// RecordKafkaLog 记录 Kafka 的消费日志
func RecordKafkaLog(record *EventRecord) {
	eventDetail := &EventDetailRecord{}
	setBasicInfo(record, eventDetail)

	message, ok := record.Message.(kafka.Message)
	if !ok {
		return
	}
	eventDetail.Topic = message.Topic
	eventDetail.Key = string(message.Key)
	eventDetail.Value = string(message.Value)
	eventDetail.DefID, _ = event.GetEventDefID(message.Value)
	eventDetail.InstID, _ = event.GetEventInstID(message.Value)
	eventDetail.MsgID = buildKafkaMsgID(message)

	recordMqLog(eventDetail)
}

// buildKafkaMsgID 构建消息ID
func buildKafkaMsgID(message kafka.Message) string {
	return strings.Join([]string{strconv.Itoa(message.Partition),
		strconv.FormatInt(message.Offset, 10)}, ":")
}

// setBasicInfo 设置消息基础信息
func setBasicInfo(event *EventRecord, consumeRecord *EventDetailRecord) {
	consumeRecord.Type = event.Type
	consumeRecord.MsgID = event.MsgID
	consumeRecord.Ack = true
	consumeRecord.FailedReason = ""
	if event.Error != nil {
		consumeRecord.Ack = false
		consumeRecord.FailedReason = event.Error.Error()
	}

	consumeRecord.BizKey = event.BizKey
	consumeRecord.Costs = time.Since(event.StartTime).Milliseconds()
}

// recordMqLog 记录消息队列的消费日志
func recordMqLog(m *EventDetailRecord) {
	fields := []string{
		"bizKey", m.BizKey,
		"topic", m.Topic,
		"key", m.Key,
		"value", m.Value,
		"ack", strconv.FormatBool(m.Ack),
		"failReason", m.FailedReason,
		"env", utils.GetEnv(),
		"costs", strconv.FormatInt(m.Costs, 10),
		"inst_id", m.InstID,
		"def_id", m.DefID,
		"type", string(m.Type),
		"msg_id", m.MsgID,
	}
	mqLog().Infof("[%s][%s] %s %s", GetFlowTraceID(m.DefID, m.InstID),
		m.Type, utils.StructToJsonStr(m), fields)
}

func mqLog() log.Logger {
	return log.GetDefaultLogger()
}
