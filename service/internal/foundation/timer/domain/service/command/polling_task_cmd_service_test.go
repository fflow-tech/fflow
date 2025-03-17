package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PollingTaskCommandService_GetPollingTaskWorkLock(t *testing.T) {
	mockService := &PollingTaskCommandService{
		&mockPollingTaskRepository{},
		&mockTimerTaskRepository{},
		&mockEventBusRepository{},
	}

	_, err := mockService.GetPollingTaskWorkLock()
	assert.Nil(t, err)
}

func Test_PollingTaskCommandService_GetTaskBucketID(t *testing.T) {
	mockService := &PollingTaskCommandService{
		&mockPollingTaskRepository{},
		&mockTimerTaskRepository{},
		&mockEventBusRepository{},
	}

	_, err := mockService.GetTaskBucketID("test")
	assert.Nil(t, err)
}

func Test_PollingTaskCommandService_SendPollingTaskWork(t *testing.T) {
	tests := []struct {
		name      string
		timeSlice string
		wantErr   bool
	}{
		{
			name:      "send fail",
			timeSlice: "error",
			wantErr:   true,
		},
		{
			name:      "success",
			timeSlice: "success",
		},
	}

	mockService := &PollingTaskCommandService{
		&mockPollingTaskRepository{},
		&mockTimerTaskRepository{},
		&mockEventBusRepository{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockService.SendPollingTaskWork(tt.timeSlice); (err != nil) != tt.wantErr {
				t.Errorf("SendPollingTaskWork() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
