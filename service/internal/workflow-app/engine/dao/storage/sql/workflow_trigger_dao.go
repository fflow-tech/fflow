package sql

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/seq"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"gorm.io/gorm"
)

// TriggerDAO Trigger数据访问对象
type TriggerDAO struct {
	db *mysql.Client
}

// NewTriggerDAO Trigger数据访问对象构造函数
func NewTriggerDAO(db *mysql.Client) *TriggerDAO {
	return &TriggerDAO{db: db}
}

// Transaction 事务
func (dao *TriggerDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// CreateTrigger 创建触发器
func (dao *TriggerDAO) CreateTrigger(d *dto.CreateTriggerDTO) (*po.TriggerPO, error) {
	id, err := seq.NewUint()
	if err != nil {
		return nil, err
	}

	p, err := convertor.TriggerConvertor.ConvertCreateDTOToPO(d)
	p.ID = id
	if err := dao.db.Create(p).Error; err != nil {
		log.Errorf("Failed to create trigger, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// GetTrigger 获取触发器
func (dao *TriggerDAO) GetTrigger(d *dto.GetTriggerDTO) (*po.TriggerPO, error) {
	if utils.IsZero(d.DefID) || utils.IsZero(d.TriggerID) {
		return nil, fmt.Errorf("get trigger `DefID` and `TriggerID` must not be zero, "+
			"DefID:[%s] TriggerID:[%s]", d.DefID, d.TriggerID)
	}

	trigger := &po.TriggerPO{}
	p, err := convertor.TriggerConvertor.ConvertGetDTOToPO(d)
	if err != nil {
		return nil, err
	}
	if err := dao.db.ReadFromSlave(false).Where(p).Take(trigger).Error; err != nil {
		log.Errorf("Failed to get trigger, caused by %s", err)
		return nil, err
	}

	return trigger, nil
}

// DeleteTrigger 删除触发器
func (dao *TriggerDAO) DeleteTrigger(d *dto.DeleteTriggerDTO) error {
	if utils.IsZero(d.TriggerID) || utils.IsZero(d.DefID) {
		return fmt.Errorf("delete trigger `TriggerID` and `DefID` must not be zero, "+
			"TriggerID:[%s] DefID:[%s]", d.TriggerID, d.DefID)
	}

	p, err := convertor.TriggerConvertor.ConvertDeleteDTOToPO(d)
	if err != nil {
		return err
	}
	if err := dao.db.Where(p).Delete(&po.TriggerPO{}).Error; err != nil {
		log.Errorf("Failed to delete trigger, caused by %s", err)
		return err
	}

	return nil
}

// UpdateTrigger 更新触发器
func (dao *TriggerDAO) UpdateTrigger(d *dto.UpdateTriggerDTO) error {
	if utils.IsZero(d.DefID) {
		return fmt.Errorf("update trigger args must not be zero, defID:[%s] triggerID:[%s]", d.DefID, d.TriggerID)
	}

	p := convertor.TriggerConvertor.ConvertUpdateDTOToPO(d)
	db := dao.db.Where("def_id = ? and def_version = ?", d.DefID, d.DefVersion)
	if d.InstID != "" {
		db.Where("inst_id = ?", d.InstID)
	}
	if d.TriggerID != "" {
		db.Where("id = ?", d.TriggerID)
	}
	if d.Level != 0 {
		db.Where("level = ?", d.Level.IntValue())
	}
	if err := db.Updates(p).Error; err != nil {
		log.Errorf("Failed to update trigger, caused by %s", err)
		return err
	}

	return nil
}

// PageQuery 分页查询触发器
func (dao *TriggerDAO) PageQuery(d *dto.PageQueryTriggerDTO) ([]*po.TriggerPO, error) {
	if utils.IsZero(d.DefID) {
		return nil, fmt.Errorf("PageQueryLastVersion trigger `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	var triggers []*po.TriggerPO
	p, err := convertor.TriggerConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return nil, err
	}
	if err := dao.db.ReadFromSlave(false).
		Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).Find(&triggers).Error; err != nil {
		log.Errorf("Failed to page query triggers, caused by %s", err)
		return nil, err
	}

	return triggers, nil
}

// Count 根据条件获取触发器总数
func (dao *TriggerDAO) Count(d *dto.PageQueryTriggerDTO) (int64, error) {
	if utils.IsZero(d.DefID) {
		return 0, fmt.Errorf("count trigger `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	var totalCount int64
	p, err := convertor.TriggerConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return 0, err
	}
	if err := dao.db.ReadFromSlave(false).
		Model(&po.TriggerPO{}).Where(p).Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to get triggers count, caused by %s", err)
		return 0, err
	}

	return totalCount, nil
}

// QueryTriggerByName 通过name查询trigger
func (dao *TriggerDAO) QueryTriggerByName(d *dto.QueryTriggerDTO) ([]*po.TriggerPO, error) {
	if utils.IsZero(d.Event) || utils.IsZero(d.Status) {
		return nil, fmt.Errorf("query triggers by `Event`,`Status` must not be zero, Event:[%s], Status:[%d]",
			d.Event, d.Status)
	}

	var triggers []*po.TriggerPO
	p := &po.TriggerPO{Event: d.Event, Status: d.Status}
	if err := dao.db.ReadFromSlave(false).Where(p).Find(&triggers).Error; err != nil {
		log.Errorf("Failed to query triggers by name, caused by %s", err)
		return nil, err
	}

	return triggers, nil
}
