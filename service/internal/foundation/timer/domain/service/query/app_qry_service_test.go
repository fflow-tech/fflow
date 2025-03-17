package query

import (
	"errors"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

type mockAppRepository struct {
}

func (m *mockAppRepository) GetAppList(d *dto.PageQueryAppDTO) ([]*entity.App, int64, error) {
	if d.Name == "error" {
		return nil, 0, errors.New(d.Name)
	}
	return []*entity.App{{}, {}}, 2, nil
}

func (m *mockAppRepository) CreateApp(d *dto.CreateAppDTO) (*entity.App, error) {
	return nil, nil
}

func (m *mockAppRepository) DeleteApp(d *dto.DeleteAppDTO) error {
	return nil
}

func Test_AppQueryService_GetAppList(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.PageQueryAppDTO
		wantErr bool
	}{
		{
			name: "error",
			req: &dto.PageQueryAppDTO{
				Name: "error",
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

	mockService := &AppQueryService{&mockAppRepository{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := mockService.GetAppList(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("GetAppList() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
