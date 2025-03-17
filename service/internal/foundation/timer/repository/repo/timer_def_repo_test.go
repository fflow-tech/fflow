package repo

import (
	"errors"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"gorm.io/gorm"
)

type mockTimerDefRedisDAO struct{}

func (m *mockTimerDefRedisDAO) AddTimerDef(def *dto.CreateTimerDefDTO) error {
	if def.DefID == "0" {
		return errors.New("error")
	}
	return nil
}

func (m *mockTimerDefRedisDAO) GetTimerDef(d *dto.GetTimerDefDTO) (*po.TimerDefPO, error) {
	if d.DefID == "error" {
		return nil, errors.New(d.DefID)
	}
	return &po.TimerDefPO{}, nil
}

func (m *mockTimerDefRedisDAO) DelTimerDef(d *dto.DeleteTimerDefDTO) error {
	if d.DefID == "error" {
		return errors.New(d.DefID)
	}
	return nil
}

func (m *mockTimerDefRedisDAO) ChangeTimerStatus(d *dto.ChangeTimerStatusDTO) error {
	if d.DefID == "error" {
		return errors.New(d.DefID)
	}
	return nil
}

type mockTimerDefDAO struct{}

func (m *mockTimerDefDAO) Create(d *dto.CreateTimerDefDTO) (*po.TimerDefPO, error) {
	if d.Name == "create error" {
		return nil, errors.New(d.Name)
	}
	if d.Name == "invalid defID" {
		return &po.TimerDefPO{
			App:   "test",
			Name:  "test",
			DefID: "error",
		}, nil
	}

	return &po.TimerDefPO{
		App:   "test",
		Name:  "test",
		DefID: "test",
		Model: gorm.Model{ID: 1},
	}, nil
}

func (m *mockTimerDefDAO) Transaction(fun func(*mysql.Client) error) error {
	return nil
}

func (m *mockTimerDefDAO) Delete(d *dto.DeleteTimerDefDTO) error {
	return nil
}

func (m *mockTimerDefDAO) PageQueryTimeList(d *dto.PageQueryTimeDefDTO) ([]*po.TimerDefPO, error) {
	if d.Name == "list fail" {
		return nil, errors.New(d.Name)
	}
	return []*po.TimerDefPO{{}, {}}, nil
}

func (m *mockTimerDefDAO) Count(d *dto.CountTimerDefDTO) (int64, error) {
	if d.Name == "count fail" {
		return 0, errors.New(d.Name)
	}
	return 10, nil
}

func (m *mockTimerDefDAO) UpdateStatus(d *dto.UpdateTimerDefDTO) error {
	return nil
}

func (m *mockTimerDefDAO) CountByStatus(status int) (int64, error) {
	return 0, nil
}

func (m *mockTimerDefDAO) GetTimerDefByAppName(app string, name string) (*po.TimerDefPO, error) {
	return nil, nil
}

func Test_TimerDefRepo_CreateTimerDef(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.CreateTimerDefDTO
		wantErr bool
	}{
		{
			name: "invalid app",
			req: &dto.CreateTimerDefDTO{
				App: "error app",
			},
			wantErr: true,
		},
		{
			name: "empty app",
			req: &dto.CreateTimerDefDTO{
				App: "empty app",
			},
			wantErr: true,
		},
		{
			name: "create error",
			req: &dto.CreateTimerDefDTO{
				App:  "test",
				Name: "create error",
			},
			wantErr: true,
		},
		{
			name: "add def error",
			req: &dto.CreateTimerDefDTO{
				App:  "test",
				Name: "invalid defID",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.CreateTimerDefDTO{
				App:  "test",
				Name: "success",
			},
		},
	}

	mockRepo := &TimerDefRepo{&mockTimerDefRedisDAO{}, &mockTimerDefDAO{}, &mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.CreateTimerDef(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateTimerDef() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefRepo_GetTimerDef(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.GetTimerDefDTO
		wantErr bool
	}{
		{
			name: "invalid app",
			req: &dto.GetTimerDefDTO{
				DefID: "error",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.GetTimerDefDTO{
				DefID: "success",
			},
		},
	}

	mockRepo := &TimerDefRepo{&mockTimerDefRedisDAO{}, &mockTimerDefDAO{}, &mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.GetTimerDef(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetTimerDef() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefRepo_DeleteTimerDef(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.DeleteTimerDefDTO
		wantErr bool
	}{
		{
			name: "invalid app",
			req: &dto.DeleteTimerDefDTO{
				DefID: "error",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.DeleteTimerDefDTO{
				DefID: "success",
			},
		},
	}

	mockRepo := &TimerDefRepo{&mockTimerDefRedisDAO{}, &mockTimerDefDAO{}, &mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.DeleteTimerDef(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTimerDefDTO() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefRepo_ChangeTimerStatus(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.ChangeTimerStatusDTO
		wantErr bool
	}{
		{
			name: "invalid app",
			req: &dto.ChangeTimerStatusDTO{
				DefID: "error",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.ChangeTimerStatusDTO{
				DefID: "success",
			},
		},
	}

	mockRepo := &TimerDefRepo{&mockTimerDefRedisDAO{}, &mockTimerDefDAO{}, &mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.ChangeTimerStatus(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("ChangeTimerStatus() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefRepo_GetTimerDefList(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.PageQueryTimeDefDTO
		wantErr bool
	}{
		{
			name: "count fail",
			req: &dto.PageQueryTimeDefDTO{
				Name: "count fail",
			},
			wantErr: true,
		},
		{
			name: "list fail",
			req: &dto.PageQueryTimeDefDTO{
				Name: "list fail",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.PageQueryTimeDefDTO{
				Name: "success",
			},
		},
	}

	mockRepo := &TimerDefRepo{&mockTimerDefRedisDAO{}, &mockTimerDefDAO{}, &mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := mockRepo.GetTimerDefList(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetTimerDefList() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefRepo_CountTimersByStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  entity.TimerDefStatus
		wantErr bool
	}{
		{
			name:   "success",
			status: entity.TimerDefStatus(1),
		},
	}

	mockRepo := &TimerDefRepo{&mockTimerDefRedisDAO{}, &mockTimerDefDAO{}, &mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.CountTimersByStatus(tt.status); (err != nil) != tt.wantErr {
				t.Errorf("CountTimersByStatus() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
