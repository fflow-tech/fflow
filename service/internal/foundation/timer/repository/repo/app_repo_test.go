package repo

import (
	"errors"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

type mockAppDAO struct{}

func (m *mockAppDAO) Create(d *dto.CreateAppDTO) (*po.App, error) {
	if d.Name == "create fail" {
		return nil, errors.New(d.Name)
	}
	return &po.App{}, nil
}
func (m *mockAppDAO) Get(d *dto.GetAppDTO) (*po.App, error) {
	if d.Name == "error app" {
		return nil, errors.New(d.Name)
	}
	if d.Name == "empty app" {
		return nil, nil
	}
	return &po.App{
		ID:   1,
		Name: "test",
	}, nil
}
func (m *mockAppDAO) PageQuery(d *dto.PageQueryAppDTO) ([]*po.App, error) {
	if d.Name == "list fail" {
		return nil, errors.New(d.Name)
	}
	return []*po.App{{}, {}}, nil
}
func (m *mockAppDAO) Delete(d *dto.DeleteAppDTO) error {
	if d.Name == "delete fail" {
		return errors.New(d.Name)
	}
	return nil
}
func (m *mockAppDAO) Count(d *dto.CountAppDTO) (int64, error) {
	if d.Name == "count fail" {
		return 0, errors.New(d.Name)
	}
	return 0, nil
}

func Test_AppRepo_GetAppList(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.PageQueryAppDTO
		wantErr bool
	}{
		{
			name: "count fail",
			req: &dto.PageQueryAppDTO{
				Name: "count fail",
			},
			wantErr: true,
		},
		{
			name: "list fail",
			req: &dto.PageQueryAppDTO{
				Name: "list fail",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.PageQueryAppDTO{
				Name: "success",
			},
		},
	}

	mockRepo := &AppRepo{&mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := mockRepo.GetAppList(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetAppList() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_AppRepo_CreateApp(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.CreateAppDTO
		wantErr bool
	}{
		{
			name: "create fail",
			req: &dto.CreateAppDTO{
				Name: "create fail",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.CreateAppDTO{
				Name: "success",
			},
		},
	}

	mockRepo := &AppRepo{&mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.CreateApp(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateApp() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_AppRepo_DeleteApp(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.DeleteAppDTO
		wantErr bool
	}{
		{
			name: "delete fail",
			req: &dto.DeleteAppDTO{
				Name: "delete fail",
			},
			wantErr: true,
		},
		{
			name: "success",
			req: &dto.DeleteAppDTO{
				Name: "success",
			},
		},
	}

	mockRepo := &AppRepo{&mockAppDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.DeleteApp(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("DeleteApp() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
