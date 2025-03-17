package event

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// recordTDMQConsumeLog 记录TDMQ消费日志记录
func recordTDMQConsumeLog(event *logs.EventRecord) {
	event.BizKey = constants.ServiceName + "_" + event.BizKey
	event.Type = logs.Consume
	logs.RecordTDMQLog(event)
}

// recordKafkaConsumeLog 记录 kafka 消费日志记录
func recordKafkaConsumeLog(event *logs.EventRecord) {
	event.BizKey = constants.ServiceName + "_" + event.BizKey
	event.Type = logs.Consume
	logs.RecordKafkaLog(event)
}
