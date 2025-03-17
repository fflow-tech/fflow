package repo

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// TriggerRepo 触发器仓储层
type TriggerRepo struct {
	triggerDAO storage.TriggerDAO
}

// NewTriggerRepo 初始化触发器仓储层
func NewTriggerRepo(t *sql.TriggerDAO) *TriggerRepo {
	return &TriggerRepo{triggerDAO: t}
}

// Create 创建
func (t *TriggerRepo) Create(d *dto.CreateTriggerDTO) (string, error) {
	trigger, err := t.triggerDAO.CreateTrigger(d)
	if err != nil {
		return "", err
	}

	return utils.UintToStr(trigger.ID), nil
}

// Update 更新
func (t *TriggerRepo) Update(d *dto.UpdateTriggerDTO) error {
	return t.triggerDAO.UpdateTrigger(d)
}

// Get 获取
func (t *TriggerRepo) Get(d *dto.GetTriggerDTO) (*entity.Trigger, error) {
	triggerPO, err := t.triggerDAO.GetTrigger(d)
	if err != nil {
		return nil, err
	}

	return convertor.TriggerConvertor.ConvertPOToEntity(triggerPO)
}

// Count 统计
func (t *TriggerRepo) Count(d *dto.PageQueryTriggerDTO) (int64, error) {
	return t.triggerDAO.Count(d)
}

// PageQuery 分页查询
func (t *TriggerRepo) PageQuery(req *dto.PageQueryTriggerDTO) ([]*entity.Trigger, error) {
	if req.PageQuery == nil {
		req.PageQuery = constants.NewDefaultPageQuery()
	}

	triggerPOs, err := t.triggerDAO.PageQuery(req)
	if err != nil {
		return nil, err
	}

	r := []*entity.Trigger{}
	for _, triggerPO := range triggerPOs {
		trigger, err := convertor.TriggerConvertor.ConvertPOToEntity(triggerPO)
		if err != nil {
			return nil, err
		}
		r = append(r, trigger)
	}

	return r, nil
}

// QueryByName 根据名称批量获取触发器
func (t *TriggerRepo) QueryByName(d *dto.QueryTriggerDTO) ([]*entity.Trigger, error) {
	triggerPOS, err := t.triggerDAO.QueryTriggerByName(d)
	if err != nil {
		return nil, err
	}

	triggers := []*entity.Trigger{}
	for _, triggerPO := range triggerPOS {
		triggerEntity, err := convertor.TriggerConvertor.ConvertPOToEntity(triggerPO)
		if err != nil {
			return nil, nil
		}
		triggers = append(triggers, triggerEntity)
	}

	return triggers, nil
}
