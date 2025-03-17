package command

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
)

// AppCommandService App 指令服务
type AppCommandService struct {
	appRepo ports.AppRepository
}

// NewAppCommandService 新建服务
func NewAppCommandService(r *ports.RepoProviderSet) *AppCommandService {
	return &AppCommandService{appRepo: r.AppRepo()}
}

// CreateApp 新建App
func (a *AppCommandService) CreateApp(d *dto.CreateAppDTO) error {
	_, err := a.appRepo.CreateApp(d)
	return err
}

// DeleteApp 删除App
func (a *AppCommandService) DeleteApp(d *dto.DeleteAppDTO) error {
	return a.appRepo.DeleteApp(d)
}
