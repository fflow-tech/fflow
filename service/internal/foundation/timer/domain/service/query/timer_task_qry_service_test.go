package query

import (
	"errors"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

type mockTimerTaskRepository struct{}

func (m *mockTimerTaskRepository) AddTimerTask(d *dto.AddTimerTaskDTO) error {
	return nil
}

func (m *mockTimerTaskRepository) GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error) {
	return nil, nil
}
func (m *mockTimerTaskRepository) DelTimerTask(d *dto.DelTimerTaskDTO) error {
	return nil
}
func (m *mockTimerTaskRepository) CreateHistory(d *dto.CreateRunHistoryDTO) (*entity.RunHistory, error) {
	return nil, nil
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
	return nil, nil
}
func (m *mockTimerTaskRepository) DeleteSaveTimerTask(defID string) error {
	return nil
}
func (m *mockTimerTaskRepository) DelPendingTimerTask(defID string, curTime time.Time) error {
	return nil
}
func (m *mockTimerTaskRepository) CountPendingTimers(curTime time.Time) (int, error) {
	return 0, nil
}

type mockPollingTaskRepository struct {
}

func (m *mockPollingTaskRepository) GetBucketNum() int {
	return 1
}

func (m *mockPollingTaskRepository) GetTaskBucketID(hashID string) (string, error) {
	return "", nil
}
func (m *mockPollingTaskRepository) SetBucketNum(num int) error {
	return nil
}
func (m *mockPollingTaskRepository) SetTimeSlice(timeDuration string) error {
	return nil
}
func (m *mockPollingTaskRepository) GetTimeSlice(timeDuration string) error {
	return nil
}
func (m *mockPollingTaskRepository) SuccessTimeSlice(timeDuration string) error {
	return nil
}

func Test_TimerDefQueryService_GetTimerTasks(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.GetTimerTaskDTO
		wantErr bool
	}{
		{
			name: "success",
			req:  &dto.GetTimerTaskDTO{},
		},
	}

	mockService := &TimerTaskQueryService{&mockTimerTaskRepository{}, &mockPollingTaskRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockService.GetTimerTasks(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetTimerTasks() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefQueryService_PageQueryHistory(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.PageQueryRunHistoryDTO
		wantErr bool
	}{
		{
			name: "error",
			req: &dto.PageQueryRunHistoryDTO{
				Name: "error",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.PageQueryRunHistoryDTO{
				Name: "success",
			},
		},
	}

	mockService := &TimerTaskQueryService{&mockTimerTaskRepository{}, &mockPollingTaskRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := mockService.PageQueryHistory(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("PageQueryHistory() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefQueryService_GetTimeLimitTimers(t *testing.T) {
	tests := []struct {
		name      string
		startTime string
		endTime   string
		wantErr   bool
	}{
		{
			name:      "invalid start time",
			startTime: "error",
			endTime:   "error",
			wantErr:   true,
		},
		{
			name:      "invalid end time",
			startTime: "2006-01-02 15:04",
			endTime:   "error",
			wantErr:   true,
		},
		{
			name:      "get no trigger fail",
			startTime: "2006-01-02 15:04",
			endTime:   "2006-01-02 15:06",
			wantErr:   true,
		},
		{
			name:      "get no trigger fail",
			startTime: "2007-01-02 15:04",
			endTime:   "2007-01-02 15:06",
		},
	}

	mockService := &TimerTaskQueryService{&mockTimerTaskRepository{}, &mockPollingTaskRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockService.GetTimeLimitTimers(tt.startTime, tt.endTime); (err != nil) != tt.wantErr {
				t.Errorf("GetTimeLimitTimers() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefQueryService_CountPendingTimers(t *testing.T) {
	tests := []struct {
		name    string
		cur     time.Time
		wantErr bool
	}{
		{
			name: "success",
		},
	}

	mockService := &TimerTaskQueryService{&mockTimerTaskRepository{}, &mockPollingTaskRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockService.CountPendingTimers(tt.cur); (err != nil) != tt.wantErr {
				t.Errorf("CountPendingTimers() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
