package query

import (
	"errors"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

type mockTimerDefRepository struct {
}

func (m *mockTimerDefRepository) GetTimerDefByAppName(app, name string) (*entity.TimerDef, error) {
	return nil, nil
}

func (m *mockTimerDefRepository) CreateTimerDef(d *dto.CreateTimerDefDTO) (uint64, error) {
	return 0, nil
}

func (m *mockTimerDefRepository) GetTimerDef(d *dto.GetTimerDefDTO) (*entity.TimerDef, error) {
	if d.DefID == "error" {
		return nil, errors.New(d.DefID)
	}
	if d.DefID == "invalid param" {
		return &entity.TimerDef{}, nil
	}
	return &entity.TimerDef{
		NotifyRpcParam:  "{}",
		NotifyHttpParam: "{}",
	}, nil
}

func (m *mockTimerDefRepository) DeleteTimerDef(d *dto.DeleteTimerDefDTO) error {
	return nil
}

func (m *mockTimerDefRepository) ChangeTimerStatus(d *dto.ChangeTimerStatusDTO) error {
	return nil
}

func (m *mockTimerDefRepository) GetTimerDefList(d *dto.PageQueryTimeDefDTO) ([]*entity.TimerDef, int64, error) {
	if d.Name == "list error" {
		return nil, 0, errors.New(d.Name)
	}
	if d.Name == "invalid param" {
		return []*entity.TimerDef{{}, {}}, 2, nil
	}
	return []*entity.TimerDef{{
		NotifyRpcParam:  "{}",
		NotifyHttpParam: "{}",
	}}, 1, nil
}

func (m *mockTimerDefRepository) CountTimersByStatus(status entity.TimerDefStatus) (int64, error) {
	return 0, nil
}

func Test_TimerDefQueryService_GetTimerDef(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.GetTimerDefDTO
		wantErr bool
	}{
		{
			name: "error",
			req: &dto.GetTimerDefDTO{
				DefID: "error",
			},
			wantErr: true,
		},
		{
			name: "invalid param",
			req: &dto.GetTimerDefDTO{
				DefID: "invalid param",
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

	mockService := &TimerDefQueryService{&mockTimerDefRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockService.GetTimerDef(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetTimerDef() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefQueryService_GetTimerDefList(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.PageQueryTimeDefDTO
		wantErr bool
	}{
		{
			name: "list error",
			req: &dto.PageQueryTimeDefDTO{
				Name: "list error",
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

	mockService := &TimerDefQueryService{&mockTimerDefRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := mockService.GetTimerDefList(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetTimerDefList() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_TimerDefQueryService_CountTimersByStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  entity.TimerDefStatus
		wantErr bool
	}{
		{
			name:   "list error",
			status: entity.TimerDefStatus(1),
		},
	}

	mockService := &TimerDefQueryService{&mockTimerDefRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockService.CountTimersByStatus(tt.status); (err != nil) != tt.wantErr {
				t.Errorf("CountTimersByStatus() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
