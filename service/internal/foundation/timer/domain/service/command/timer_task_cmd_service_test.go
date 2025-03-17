package command

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

type mockTimerTaskRepository struct{}

func (m *mockTimerTaskRepository) AddTimerTask(d *dto.AddTimerTaskDTO) error {
	if d.HashID == "add timer task fail" {
		return errors.New("fail ")
	}
	return nil
}
func (m *mockTimerTaskRepository) GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error) {
	if d.BucketTime == "error" {
		return nil, errors.New("invalid readyset")
	}
	if d.BucketTime == "send error" {
		return []string{"error"}, nil
	}
	return []string{"success", "delete fail"}, nil
}
func (m *mockTimerTaskRepository) DelTimerTask(d *dto.DelTimerTaskDTO) error {
	if d.HashID == "delete fail" {
		return errors.New(d.BucketTime)
	}
	return nil
}
func (m *mockTimerTaskRepository) CreateHistory(d *dto.CreateRunHistoryDTO) (*entity.RunHistory, error) {
	if d.DefID == "create timer history fail" {
		return nil, errors.New("invalid")
	}
	return &entity.RunHistory{}, nil
}
func (m *mockTimerTaskRepository) UpdateHistory(d *dto.UpdateRunHistoryDTO) error {
	return nil
}
func (m *mockTimerTaskRepository) PageQueryHistory(d *dto.PageQueryRunHistoryDTO) ([]*entity.RunHistory, int64, error) {
	if d.Name == "error" {
		return nil, 0, errors.New(d.Name)
	}
	return []*entity.RunHistory{{}, {}}, 2, nil
}
func (m *mockTimerTaskRepository) GetNotTriggeredTimers(bucketTime string) ([]string, error) {
	if bucketTime == "0_2006-01-02 15:05" {
		return nil, errors.New("fail")
	}
	return []string{"1"}, nil
}
func (m *mockTimerTaskRepository) GetTaskTableName(bucketID, timeSlice string) string {
	return bucketID + "_" + timeSlice
}
func (m *mockTimerTaskRepository) DeleteRunHistories() error {
	return nil
}
func (m *mockTimerTaskRepository) GetSaveTimerTask(defID string) (*dto.SaveTimerTaskDTO, error) {
	if defID == "get save fail" {
		return nil, errors.New("fail")
	}
	return &dto.SaveTimerTaskDTO{
		UnixTime: time.Now().UnixNano(),
	}, nil
}
func (m *mockTimerTaskRepository) DeleteSaveTimerTask(defID string) error {
	if strings.HasPrefix(defID, "enable") {
		return errors.New("delete fail ")
	}
	return nil
}
func (m *mockTimerTaskRepository) DelPendingTimerTask(defID string, curTime time.Time) error {
	if defID == "delete pending fail" {
		return errors.New("fail")
	}
	return nil
}
func (m *mockTimerTaskRepository) CountPendingTimers(curTime time.Time) (int, error) {
	return 0, nil
}

type mockEventBusRepository struct{}

func (m *mockEventBusRepository) SendPollingEvent(ctx context.Context, msg interface{}) error {
	if msg.(string) == "1_error" {
		return errors.New("fail")
	}
	return nil
}

func (m *mockEventBusRepository) NewPollingConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return nil, nil
}
func (m *mockEventBusRepository) SendTimerTaskEvent(ctx context.Context, msg interface{}) error {
	payload, ok := msg.(string)
	if !ok {
		return errors.New("invalid msg")
	}
	if payload == "error" {
		return errors.New("send fail")
	}
	return nil
}
func (m *mockEventBusRepository) NewTimerTaskConsumer(ctx context.Context, group string,
	handle func(context.Context, interface{}) error) (mq.Consumer, error) {
	return nil, nil
}
func Test_TimerTaskCommandService_AddTimerTask(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.AddTimerTaskDTO
		wantErr bool
	}{
		{
			name: "success",
			req:  &dto.AddTimerTaskDTO{},
		},
	}
	mockService := &TimerTaskCommandService{&mockTimerTaskRepository{}, &mockEventBusRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockService.AddTimerTask(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("AddTimerTask() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
func Test_TimerTaskCommandService_DelTimerTask(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.DelTimerTaskDTO
		wantErr bool
	}{
		{
			name: "success",
			req:  &dto.DelTimerTaskDTO{},
		},
	}
	mockService := &TimerTaskCommandService{&mockTimerTaskRepository{}, &mockEventBusRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockService.DelTimerTask(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("DelTimerTask() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
func Test_TimerTaskCommandService_SendReadyTimer(t *testing.T) {
	tests := []struct {
		name    string
		req     string
		wantErr bool
	}{
		{
			name: "success",
		},
	}
	mockService := &TimerTaskCommandService{&mockTimerTaskRepository{}, &mockEventBusRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockService.SendReadyTimer(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("SendReadyTimer() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
func Test_TimerTaskCommandService_GetReadyTimer(t *testing.T) {
	tests := []struct {
		name      string
		req       string
		startTime time.Time
		wantErr   bool
	}{
		{
			name:    "get fail",
			req:     "error",
			wantErr: true,
		},
		{
			name:    "send fail",
			req:     "send error",
			wantErr: true,
		},
		{
			name: "success",
			req:  "success",
		},
	}
	mockService := &TimerTaskCommandService{&mockTimerTaskRepository{}, &mockEventBusRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockService.GetReadyTimer(tt.req, tt.startTime); (err != nil) != tt.wantErr {
				t.Errorf("GetReadyTimer() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
