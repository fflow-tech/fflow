package eventbus

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// Recorder 记录者.
type Recorder interface {
	RecordTDMQProduceLog(event *logs.EventRecord)
	RecordTDMQConsumeLog(event *logs.EventRecord)
}

// GetLogRecorder 获取日志记录者.
func GetLogRecorder() *LogRecorder {
	return recorder
}

var recorder = &LogRecorder{}

// LogRecorder 日志记录者.
type LogRecorder struct{}

// RecordTDMQProduceLog 记录TDMQ事件生产日志
func (l *LogRecorder) RecordTDMQProduceLog(event *logs.EventRecord) {
	event.BizKey = constants.ServiceName + "_" + event.BizKey
	event.Type = logs.Produce
	logs.RecordTDMQLog(event)
}

// RecordTDMQConsumeLog 记录TDMQ消费日志记录
func (l *LogRecorder) RecordTDMQConsumeLog(event *logs.EventRecord) {
	event.BizKey = constants.ServiceName + "_" + event.BizKey
	event.Type = logs.Consume
	logs.RecordTDMQLog(event)
}
