// Package repo 仓储层负责从 DAO 层获取数据，组装为 entity 返回 service 业务层。
package repo

import (
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/cache/redis"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ErrGetTimerByAppNameNotFound  根据应用和名称获取定时器失败，错误未记录未找到.
var ErrGetTimerByAppNameNotFound = fmt.Errorf("timer with app and name is not found, err: %w", sql.ErrRecordNotFound)

// TimerDefRepo 定时器定义实体
type TimerDefRepo struct {
	timerDefRedisDAO storage.TimerDefRedisDAO
	timerDefSQLDAO   storage.TimerDefDAO
	appDAO           storage.AppDAO
}

// NewTimerDefRepo 实体构造函数
func NewTimerDefRepo(redisDAO *redis.TimerDefDAO, sqlDAO *sql.TimerDefDAO, appDAO *sql.AppDAO) *TimerDefRepo {
	return &TimerDefRepo{timerDefRedisDAO: redisDAO, timerDefSQLDAO: sqlDAO, appDAO: appDAO}
}

// CreateTimerDef 创建定时器定义
func (t *TimerDefRepo) CreateTimerDef(d *dto.CreateTimerDefDTO) (uint64, error) {
	// 判断App是否存在，必须创建已有app下的定时器
	hasAppInfo, err := t.hasAppInfo(d)
	if err != nil {
		return 0, err
	}

	if !hasAppInfo {
		return 0, fmt.Errorf("failed to CreateTimerDef , caused by app:%v not found", d.App)
	}

	// 新建时定时器状态为`未激活`
	d.Status = entity.Disabled.ToInt()
	defDTO, err := t.timerDefSQLDAO.Create(d)
	if err != nil {
		return 0, err
	}

	d.DefID = utils.UintToStr(defDTO.ID)
	if err := t.timerDefRedisDAO.AddTimerDef(d); err != nil {
		return 0, err
	}

	return uint64(defDTO.ID), nil
}

// GetTimerDef 获取定时器定义
func (t *TimerDefRepo) GetTimerDef(d *dto.GetTimerDefDTO) (*entity.TimerDef, error) {
	timerPO, err := t.timerDefRedisDAO.GetTimerDef(d)
	if err != nil {
		return nil, err
	}

	return convertor.DefConvertor.ConvertPOToEntity(timerPO)
}

// DeleteTimerDef 删除定时器定义
func (t *TimerDefRepo) DeleteTimerDef(d *dto.DeleteTimerDefDTO) error {
	if err := t.timerDefRedisDAO.DelTimerDef(d); err != nil {
		return err
	}

	return t.timerDefSQLDAO.Delete(d)
}

// ChangeTimerStatus 修改定时器状态
func (t *TimerDefRepo) ChangeTimerStatus(d *dto.ChangeTimerStatusDTO) error {
	if err := t.timerDefRedisDAO.ChangeTimerStatus(d); err != nil {
		return err
	}

	return t.timerDefSQLDAO.UpdateStatus(&dto.UpdateTimerDefDTO{DefID: d.DefID, Status: d.Status})
}

// GetTimerDefList 获取定时器定义列表
func (t *TimerDefRepo) GetTimerDefList(d *dto.PageQueryTimeDefDTO) ([]*entity.TimerDef, int64, error) {
	total, err := t.timerDefSQLDAO.Count(&dto.CountTimerDefDTO{
		App:     d.App,
		Name:    d.Name,
		Creator: d.Creator,
	})
	if err != nil {
		return nil, 0, err
	}

	timerDefPOS, err := t.timerDefSQLDAO.PageQueryTimeList(d)
	if err != nil {
		return nil, 0, err
	}

	timerEntities := make([]*entity.TimerDef, 0, len(timerDefPOS))
	for _, timerDefPO := range timerDefPOS {
		timerEntity, err := convertor.DefConvertor.ConvertPOToEntity(timerDefPO)
		if err != nil {
			return nil, 0, err
		}
		timerEntity.DefID = utils.UintToStr(timerDefPO.ID)
		timerEntities = append(timerEntities, timerEntity)
	}

	return timerEntities, total, nil
}

// hasAppInfo  是否存在App
func (t *TimerDefRepo) hasAppInfo(d *dto.CreateTimerDefDTO) (bool, error) {
	app, err := t.appDAO.Get(&dto.GetAppDTO{Name: d.App})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	}
	return !utils.IsZero(app), nil
}

// CountTimersByStatus 根据状态统计定时器数量.
func (t *TimerDefRepo) CountTimersByStatus(status entity.TimerDefStatus) (int64, error) {
	return t.timerDefSQLDAO.CountByStatus(status.ToInt())
}

// GetTimerDefByAppName 根据应用和名称获取定时器定义.
func (t *TimerDefRepo) GetTimerDefByAppName(app, name string) (*entity.TimerDef, error) {
	timerDef, err := t.timerDefSQLDAO.GetTimerDefByAppName(app, name)
	if err == sql.ErrRecordNotFound {
		return nil, ErrGetTimerByAppNameNotFound
	}
	if err != nil {
		return nil, err
	}
	return convertor.DefConvertor.ConvertPOToEntity(timerDef)
}
