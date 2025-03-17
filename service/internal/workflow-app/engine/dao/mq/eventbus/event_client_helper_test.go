package eventbus

import (
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// Test_recordKafkaProduceLog 测试记录 kafka 日志
func Test_recordKafkaProduceLog(t *testing.T) {
	type args struct {
		event *logs.EventRecord
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"正常情况", args{event: &logs.EventRecord{
				MsgID:     "testMsgID",
				Message:   nil,
				BizKey:    "testBiz",
				StartTime: time.Time{},
				Error:     nil,
				Type:      "testType",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordKafkaProduceLog(tt.args.event)
		})
	}
}

// Test_recordTDMQProduceLog 测试记录 tdmq 日志
func Test_recordTDMQProduceLog(t *testing.T) {
	type args struct {
		event *logs.EventRecord
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"正常情况", args{event: &logs.EventRecord{
				MsgID:     "testMsgID",
				Message:   nil,
				BizKey:    "testBiz",
				StartTime: time.Time{},
				Error:     nil,
				Type:      "testType",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordTDMQProduceLog(tt.args.event)
		})
	}
}
