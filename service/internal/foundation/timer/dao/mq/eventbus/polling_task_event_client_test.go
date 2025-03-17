package eventbus

import (
	"context"
	"testing"
)

func Test_PollingEventClient_SendEvent(t *testing.T) {
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
	mockClient := &PollingEventClient{client: &mockMQClient{}, recorder: &mockRecorder{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockClient.SendEvent(ctx, tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("SendEvent() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_PollingEventClient_NewConsumer(t *testing.T) {
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
	mockClient := &PollingEventClient{client: &mockMQClient{}, recorder: &mockRecorder{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockClient.NewConsumer(ctx, tt.topic, tt.handle); (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
