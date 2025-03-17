package repo

import (
	"errors"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

type mockTimerTaskDAO struct{}

func (m *mockTimerTaskDAO) AddTimerTask(d *dto.AddTimerTaskDTO) error {
	return nil
}

func (m *mockTimerTaskDAO) GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error) {
	return nil, nil
}

func (m *mockTimerTaskDAO) DelTimerTask(d *dto.DelTimerTaskDTO) error {
	return nil
}

func (m *mockTimerTaskDAO) GetNotTriggeredTimers(bucketTime string) ([]string, error) {
	return nil, nil
}

func (m *mockTimerTaskDAO) GetTaskTableName(bucketID, timeSlice string) string {
	return ""
}

func (m *mockTimerTaskDAO) GetSaveTimerTask(defID string) (*dto.SaveTimerTaskDTO, error) {
	return nil, nil
}

func (m *mockTimerTaskDAO) DeleteSaveTimerTask(defID string) error {
	return nil
}

func (m *mockTimerTaskDAO) DelPendingTimerTask(defID string, curTime time.Time) error {
	return nil
}

func (m *mockTimerTaskDAO) CountPendingTimers(curTime time.Time) (int, error) {
	return 0, nil
}

type mockTimerTaskRunHistoryDAO struct{}

func (m *mockTimerTaskRunHistoryDAO) Create(d *dto.CreateRunHistoryDTO) (*po.RunHistoryPO, error) {
	if d.Name == "error" {
		return nil, errors.New(d.Name)
	}
	return &po.RunHistoryPO{}, nil
}

func (m *mockTimerTaskRunHistoryDAO) Get(d *dto.GetRunHistoryDTO) (*po.RunHistoryPO, error) {
	return nil, nil
}

func (m *mockTimerTaskRunHistoryDAO) Delete(d *dto.DeleteRunHistoryDTO) error {
	return nil
}

func (m *mockTimerTaskRunHistoryDAO) Update(d *dto.UpdateRunHistoryDTO) error {
	return nil
}

func (m *mockTimerTaskRunHistoryDAO) PageQuery(d *dto.PageQueryRunHistoryDTO) ([]*po.RunHistoryPO, error) {
	if d.DefID == "list fail" {
		return nil, errors.New(d.DefID)
	}
	return []*po.RunHistoryPO{{}, {}}, nil
}

func (m *mockTimerTaskRunHistoryDAO) Count(d *dto.PageQueryRunHistoryDTO) (int64, error) {
	if d.DefID == "count fail" {
		return 0, errors.New(d.DefID)
	}
	return 1, nil
}

func (m *mockTimerTaskRunHistoryDAO) DeleteByRunTime(t time.Time) error {
	return nil
}

func Test_TimerTaskRepo_AddTimerTask(t *testing.T) {
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

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.AddTimerTask(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("AddTimerTask() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_GetTimerTasks(t *testing.T) {
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

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.GetTimerTasks(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetTimerTasks() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_DelTimerTask(t *testing.T) {
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

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.DelTimerTask(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("DelTimerTask() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_CreateHistory(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.CreateRunHistoryDTO
		wantErr bool
	}{
		{
			name: "error",
			req: &dto.CreateRunHistoryDTO{
				Name: "error",
			},
			wantErr: true,
		},
		{
			name: "success",
			req:  &dto.CreateRunHistoryDTO{},
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.CreateHistory(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateHistory() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_UpdateHistory(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.UpdateRunHistoryDTO
		wantErr bool
	}{
		{
			name: "success",
			req:  &dto.UpdateRunHistoryDTO{},
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.UpdateHistory(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("UpdateHistory() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_PageQueryHistory(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.PageQueryRunHistoryDTO
		wantErr bool
	}{
		{
			name: "list fail",
			req: &dto.PageQueryRunHistoryDTO{
				DefID: "list fail",
			},
			wantErr: true,
		},
		{
			name: "count fail",
			req: &dto.PageQueryRunHistoryDTO{
				DefID: "count fail",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.PageQueryRunHistoryDTO{
				DefID: "success",
			},
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := mockRepo.PageQueryHistory(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("PageQueryHistory() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_GetNotTriggeredTimers(t *testing.T) {
	tests := []struct {
		name       string
		bucketTime string
		wantErr    bool
	}{
		{
			name: "success",
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.GetNotTriggeredTimers(tt.bucketTime); (err != nil) != tt.wantErr {
				t.Errorf("GetNotTriggeredTimers() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_GetTaskTableName(t *testing.T) {
	tests := []struct {
		name      string
		bucketID  string
		timeSlice string
		want      string
	}{
		{
			name: "success",
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockRepo.GetTaskTableName(tt.bucketID, tt.timeSlice); got != tt.want {
				t.Errorf("GetTaskTableName() got : %s, expect : %s", got, tt.want)
			}
		})
	}
}

func Test_TimerTaskRepo_GetSaveTimerTask(t *testing.T) {
	tests := []struct {
		name    string
		defID   string
		wantErr bool
	}{
		{
			name: "success",
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.GetSaveTimerTask(tt.defID); (err != nil) != tt.wantErr {
				t.Errorf("GetSaveTimerTask() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_DeleteSaveTimerTask(t *testing.T) {
	tests := []struct {
		name    string
		defID   string
		wantErr bool
	}{
		{
			name: "success",
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.DeleteSaveTimerTask(tt.defID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSaveTimerTask() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerTaskRepo_CountPendingTimers(t *testing.T) {
	tests := []struct {
		name     string
		execTime time.Time
		wantErr  bool
	}{
		{
			name: "success",
		},
	}

	mockRepo := &TimerTaskRepo{&mockTimerTaskDAO{}, &mockTimerTaskRunHistoryDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.CountPendingTimers(tt.execTime); (err != nil) != tt.wantErr {
				t.Errorf("CountPendingTimers() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
