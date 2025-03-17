package command

import (
	"errors"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

type mockPollingTaskRepository struct{}

func (m *mockPollingTaskRepository) GetBucketNum() int {
	return 2
}
func (m *mockPollingTaskRepository) GetTaskBucketID(hashID string) (string, error) {
	if hashID == "get bucket id fail" {
		return "", errors.New("fail")
	}
	if hashID == "enable get bucket id fail" {
		return "", errors.New("fail")
	}
	return "1", nil
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

type mockTimerDefRepository struct{}

func (m *mockTimerDefRepository) GetTimerDefByAppName(app, name string) (*entity.TimerDef, error) {
	return nil, nil
}

func (m *mockTimerDefRepository) CreateTimerDef(d *dto.CreateTimerDefDTO) (uint64, error) {
	if d.Name == "create error" {
		return 0, errors.New("create fail")
	}
	return 1, nil
}
func (m *mockTimerDefRepository) GetTimerDef(d *dto.GetTimerDefDTO) (*entity.TimerDef, error) {
	if d.DefID == "get def fail" {
		return nil, errors.New("fail")
	}
	if d.DefID == "status no change" {
		return &entity.TimerDef{
			DefID: d.DefID,
		}, nil
	}
	if d.DefID == "enable invalid cron type" {
		return &entity.TimerDef{
			DefID:     d.DefID,
			TimerType: 10,
			Status:    3,
		}, nil
	}
	if d.DefID == "enable parse cron timer fail" {
		return &entity.TimerDef{
			DefID:     d.DefID,
			Status:    3,
			TimerType: 2,
		}, nil
	}
	if d.DefID == "enable parse end time fail" {
		return &entity.TimerDef{
			DefID:     d.DefID,
			Status:    1,
			TimerType: 2,
			EndTime:   "111",
		}, nil
	}
	if d.DefID == "enable invalid end time fail" {
		return &entity.TimerDef{
			DefID:     d.DefID,
			Status:    1,
			TimerType: 2,
			EndTime:   "2006-01-02 15:04:05",
		}, nil
	}
	if d.DefID == "enable delay timer" {
		return &entity.TimerDef{
			DefID:     d.DefID,
			Status:    1,
			TimerType: 1,
		}, nil
	}
	if d.DefID == "enable once timer" {
		return &entity.TimerDef{
			DefID:       d.DefID,
			Status:      1,
			TriggerType: 1,
		}, nil
	}
	if d.DefID == "enable timer invalid cron" {
		return &entity.TimerDef{
			DefID:  d.DefID,
			Status: 1,
		}, nil
	}
	if d.DefID == "enable get bucket id fail" {
		return &entity.TimerDef{
			DefID:  d.DefID,
			Status: 1,
			Cron:   "0 0 */1 * * ? *",
		}, nil
	}
	if d.DefID == "enable success" {
		return &entity.TimerDef{
			DefID:  d.DefID,
			Status: 1,
			Cron:   "0 0 */1 * * ? *",
		}, nil
	}
	return &entity.TimerDef{
		DefID:     d.DefID,
		Status:    3,
		TimerType: 2,
		Cron:      "0 0 */1 * * ? *",
	}, nil
}
func (m *mockTimerDefRepository) DeleteTimerDef(d *dto.DeleteTimerDefDTO) error {
	if d.DefID == "error" {
		return errors.New("delete fail")
	}
	return nil
}
func (m *mockTimerDefRepository) ChangeTimerStatus(d *dto.ChangeTimerStatusDTO) error {
	if d.DefID == "change fail" {
		return errors.New("fail")
	}
	return nil
}
func (m *mockTimerDefRepository) GetTimerDefList(d *dto.PageQueryTimeDefDTO) ([]*entity.TimerDef, int64, error) {
	return nil, 0, nil
}
func (m *mockTimerDefRepository) CountTimersByStatus(status entity.TimerDefStatus) (int64, error) {
	return 0, nil
}

type mockWorkerPool struct{}

func (m *mockWorkerPool) Submit(f func()) error {
	f()
	return nil
}
func Test_TimerDefCommandService_CreateTimerDef(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.CreateTimerDefDTO
		wantErr bool
	}{
		{
			name:    "validate fail",
			req:     &dto.CreateTimerDefDTO{},
			wantErr: true,
		},
		{
			name: "create fail",
			req: &dto.CreateTimerDefDTO{
				Name:        "create error",
				TimerType:   2,
				TriggerType: 2,
				NotifyType:  3,
				NotifyHttpParam: dto.NotifyHttpParam{
					Method: "GET",
					Url:    "www.baidu.com",
				},
				Cron: "0 0 */1 * * ? *",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.CreateTimerDefDTO{
				Name:        "success",
				TimerType:   2,
				TriggerType: 2,
				NotifyType:  3,
				NotifyHttpParam: dto.NotifyHttpParam{
					Method: "GET",
					Url:    "www.baidu.com",
				},
				Cron: "0 0 */1 * * ? *",
			},
		},
	}
	mockService := &TimerDefCommandService{&mockTimerDefRepository{}, &mockTimerTaskRepository{}, &mockPollingTaskRepository{}, &mockWorkerPool{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockService.CreateTimerDef(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateTimerDef() got err:%v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
func Test_TimerDefCommandService_DeleteTimerDef(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.DeleteTimerDefDTO
		wantErr bool
	}{
		{
			name: "delete fail",
			req: &dto.DeleteTimerDefDTO{
				DefID: "error",
			},
			wantErr: true,
		},
		{
			name: "get save fail",
			req: &dto.DeleteTimerDefDTO{
				DefID: "get save fail",
			},
		},
		{
			name: "delete pending fail",
			req: &dto.DeleteTimerDefDTO{
				DefID: "delete pending fail",
			},
		},
	}
	mockService := &TimerDefCommandService{&mockTimerDefRepository{}, &mockTimerTaskRepository{}, &mockPollingTaskRepository{}, &mockWorkerPool{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockService.DeleteTimerDef(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTimerDef() got err:%v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
func Test_TimerDefCommandService_ChangeTimerStatus(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.ChangeTimerStatusDTO
		wantErr bool
	}{
		{
			name: "get def fail",
			req: &dto.ChangeTimerStatusDTO{
				DefID: "get def fail",
			},
			wantErr: true,
		},
		{
			name: "status no change",
			req: &dto.ChangeTimerStatusDTO{
				DefID: "status no change",
			},
		},
		{
			name: "change fail",
			req: &dto.ChangeTimerStatusDTO{
				DefID: "change fail",
			},
			wantErr: true,
		},
		{
			name: "unable get save fail",
			req: &dto.ChangeTimerStatusDTO{
				DefID:  "get save fail",
				Status: 2,
			},
		},
		{
			name: "unable delete pending fail",
			req: &dto.ChangeTimerStatusDTO{
				DefID:  "delete pending fail",
				Status: 2,
			},
		},
		{
			name: "enable get bucketID fail",
			req: &dto.ChangeTimerStatusDTO{
				DefID:  "get bucket id fail",
				Status: 1,
			},
			wantErr: true,
		},
		{
			name: "enable invalid cron type",
			req: &dto.ChangeTimerStatusDTO{
				DefID:  "enable invalid cron type",
				Status: 1,
			},
			wantErr: true,
		},
		{
			name: "enable parse cron timer fail",
			req: &dto.ChangeTimerStatusDTO{
				DefID:  "enable parse cron timer fail",
				Status: 1,
			},
			wantErr: true,
		},
		{
			name: "add timer task fail",
			req: &dto.ChangeTimerStatusDTO{
				DefID:  "add timer task fail",
				Status: 1,
			},
			wantErr: true,
		},
		{
			name: "add timer task success",
			req: &dto.ChangeTimerStatusDTO{
				DefID:  "add timer task success",
				Status: 1,
			},
		},
	}
	mockService := &TimerDefCommandService{&mockTimerDefRepository{}, &mockTimerTaskRepository{}, &mockPollingTaskRepository{}, &mockWorkerPool{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockService.ChangeTimerStatus(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("ChangeTimerStatus() got err:%v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
