package query

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
)

// AppQueryService App 定义查询服务
type AppQueryService struct {
	appRepo ports.AppRepository
}

// NewAppQueryService 新建查询服务
func NewAppQueryService(r *ports.RepoProviderSet) *AppQueryService {
	return &AppQueryService{appRepo: r.AppRepo()}
}

// GetAppList 获取 App 列表
func (q *AppQueryService) GetAppList(d *dto.PageQueryAppDTO) ([]*dto.App, int64, error) {
	appEntities, total, err := q.appRepo.GetAppList(d)
	if err != nil {
		return nil, 0, err
	}

	appList, err := convertor.AppConvertor.ConvertAppEntitiesToDTOs(appEntities)
	if err != nil {
		return nil, 0, err
	}

	return appList, total, nil
}
