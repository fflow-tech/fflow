package eventbus

import (
	"context"
	"fmt"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/fflow-tech/fflow/service/pkg/logs"

	"github.com/apache/pulsar-client-go/pulsar"
)

type mockRecorder struct{}

func (m *mockRecorder) RecordTDMQProduceLog(event *logs.EventRecord) {

}

func (m *mockRecorder) RecordTDMQConsumeLog(event *logs.EventRecord) {
}

type mockMQClient struct{}

func (m *mockMQClient) SendMessage(ctx context.Context, topic string, msg interface{}) (string, error) {
	if string(msg.(pulsar.ProducerMessage).Payload) == "fail" {
		return "", fmt.Errorf("fail")
	}
	return "", nil
}
func (m *mockMQClient) NewConsumer(ctx context.Context, topic, subName string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return nil, nil
}

func Test_TimerTaskEventClient_SendEvent(t *testing.T) {
	tests := []struct {
		name    string
		msg     interface{}
		wantErr bool
	}{
		{
			name:    "invalid msg type",
			msg:     1,
			wantErr: true,
		},
		{
			name:    "send fail",
			msg:     "fail",
			wantErr: true,
		},
		{
			name: "send success",
			msg:  "success",
		},
	}

	ctx := context.Background()
	mockClient := &TimerTaskEventClient{client: &mockMQClient{}, recorder: &mockRecorder{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockClient.SendEvent(ctx, tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("SendEvent() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskEventClient_NewConsumer(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		handle  func(context.Context, interface{}) error
		wantErr bool
	}{
		{
			name:  "invalid msg type",
			topic: "test",
		},
	}

	ctx := context.Background()
	mockClient := &TimerTaskEventClient{client: &mockMQClient{}, recorder: &mockRecorder{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockClient.NewConsumer(ctx, tt.topic, tt.handle); (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
