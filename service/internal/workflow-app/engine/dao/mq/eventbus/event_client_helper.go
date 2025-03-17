package eventbus

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// recordKafkaProduceLog 记录Kafka事件生产日志
func recordKafkaProduceLog(event *logs.EventRecord) {
	event.BizKey = constants.ServiceName + "_" + event.BizKey
	event.Type = logs.Produce
	logs.RecordKafkaLog(event)
}

// recordTDMQProduceLog 记录TDMQ事件生产日志
func recordTDMQProduceLog(event *logs.EventRecord) {
	event.BizKey = constants.ServiceName + "_" + event.BizKey
	event.Type = logs.Produce
	logs.RecordTDMQLog(event)
}
