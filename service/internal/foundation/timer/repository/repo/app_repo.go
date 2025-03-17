package repo

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// AppRepo  应用仓储层实现体
type AppRepo struct {
	appDAO storage.AppDAO
}

// NewTimerDefRepo 实体构造函数
func NewAppRepo(appDAO *sql.AppDAO) *AppRepo {
	return &AppRepo{appDAO: appDAO}
}

// GetAppList 获取 App 列表
func (r *AppRepo) GetAppList(d *dto.PageQueryAppDTO) ([]*entity.App, int64, error) {
	total, err := r.appDAO.Count(&dto.CountAppDTO{
		Name:    d.Name,
		Creator: d.Creator,
	})
	if err != nil {
		return nil, 0, err
	}

	if utils.IsZero(d.PageQuery) {
		d.PageQuery = constants.NewDefaultPageQuery()
	}
	if utils.IsZero(d.Order) {
		d.Order = constants.NewDefaultOrder()
	}

	appPOs, err := r.appDAO.PageQuery(d)
	if err != nil {
		return nil, 0, err
	}

	appEntities, err := convertor.AppConvertor.ConvertAppPOsToEntities(appPOs)
	if err != nil {
		return nil, 0, err
	}

	return appEntities, total, nil
}

// CreateApp 创建App
func (r *AppRepo) CreateApp(d *dto.CreateAppDTO) (*entity.App, error) {
	app, err := r.appDAO.Create(d)
	if err != nil {
		return nil, err
	}

	appEntity, err := convertor.AppConvertor.ConvertAppPOToEntity(app)
	if err != nil {
		return nil, err
	}
	return appEntity, nil
}

// DeleteApp 删除App
func (r *AppRepo) DeleteApp(d *dto.DeleteAppDTO) error {
	return r.appDAO.Delete(d)
}
