// Package convertor 负责 po-> entity 转换
package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"

	"github.com/jinzhu/copier"
)

var (
	// DefConvertor 定时器定义转换体
	DefConvertor = &defConvertorImpl{}
	// TaskConvertor 定时器任务转换体
	TaskConvertor = &taskConvertorImpl{}
	AppConvertor  = &appConvertorImpl{}
)

type defConvertorImpl struct {
}

// ConvertPOToEntity 转换成实体
func (*defConvertorImpl) ConvertPOToEntity(p *po.TimerDefPO) (*entity.TimerDef, error) {
	def := &entity.TimerDef{}
	if err := copier.Copy(def, p); err != nil {
		return nil, err
	}
	return def, nil
}

type taskConvertorImpl struct {
}

// ConvertPOToEntity 转换成实体
func (*taskConvertorImpl) ConvertPOToEntity(p *po.RunHistoryPO) (*entity.RunHistory, error) {
	task := &entity.RunHistory{}
	if err := copier.Copy(task, p); err != nil {
		return nil, err
	}
	return task, nil
}

// ConvertHistoryPOsToEntities 将 history 的 po list 转为 entity list
func (c *taskConvertorImpl) ConvertHistoryPOsToEntities(p []*po.RunHistoryPO) ([]*entity.RunHistory, error) {
	histories := make([]*entity.RunHistory, 0, len(p))
	for _, historyPO := range p {
		history, err := c.ConvertPOToEntity(historyPO)
		if err != nil {
			return nil, err
		}

		histories = append(histories, history)
	}
	return histories, nil
}

type appConvertorImpl struct {
}

// ConvertAppPOsToEntities 转换为实体列表
func (c *appConvertorImpl) ConvertAppPOsToEntities(appPOs []*po.App) ([]*entity.App, error) {
	appEntities := make([]*entity.App, 0, len(appPOs))

	for _, appPO := range appPOs {
		appEntity, err := c.ConvertAppPOToEntity(appPO)
		if err != nil {
			return nil, err
		}
		appEntities = append(appEntities, appEntity)
	}

	return appEntities, nil
}

// ConvertAppPOToEntity 转换为实体
func (c *appConvertorImpl) ConvertAppPOToEntity(p *po.App) (*entity.App, error) {
	app := &entity.App{}
	if err := copier.Copy(app, p); err != nil {
		return nil, err
	}

	return app, nil
}
