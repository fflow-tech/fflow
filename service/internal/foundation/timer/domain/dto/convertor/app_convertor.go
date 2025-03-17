package convertor

import (
	"github.com/jinzhu/copier"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

var (
	// AppConvertor app 转换体
	AppConvertor = &appConvertorImpl{}
)

type appConvertorImpl struct {
}

// ConvertAppEntitiesToDTOs 转换 app entities->dtos
func (c *appConvertorImpl) ConvertAppEntitiesToDTOs(appEntities []*entity.App) ([]*dto.App, error) {
	appDtos := make([]*dto.App, 0, len(appEntities))
	for _, appEntity := range appEntities {
		appDTO, err := c.ConvertAppEntityToDTO(appEntity)
		if err != nil {
			return nil, err
		}
		appDtos = append(appDtos, appDTO)
	}

	return appDtos, nil
}

// ConvertAppEntityToDTO 转换 app entity->dto
func (c *appConvertorImpl) ConvertAppEntityToDTO(e *entity.App) (*dto.App, error) {
	app := &dto.App{}
	if err := copier.Copy(app, e); err != nil {
		return nil, err
	}

	return app, nil
}
