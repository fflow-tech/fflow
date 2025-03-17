package mock

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
)

// TriggerDAOMock 模拟对象
type TriggerDAOMock struct {
}

// NewTriggerDAOMock 构造方法
func NewTriggerDAOMock() *TriggerDAOMock {
	return &TriggerDAOMock{}
}

// Transaction 事务
func (dao *TriggerDAOMock) Transaction(fun func() error) error {
	return nil
}

// CreateTrigger 创建触发器
func (dao *TriggerDAOMock) CreateTrigger(d *dto.CreateTriggerDTO) (*po.TriggerPO, error) {
	return nil, nil
}

// GetTrigger 获取触发器
func (dao *TriggerDAOMock) GetTrigger(d *dto.GetTriggerDTO) (*po.TriggerPO, error) {
	return nil, nil
}

// DeleteTrigger 删除触发器
func (dao *TriggerDAOMock) DeleteTrigger(d *dto.DeleteTriggerDTO) error {
	return nil
}

// UpdateTrigger 更新触发器
func (dao *TriggerDAOMock) UpdateTrigger(d *dto.UpdateTriggerDTO) error {
	return nil
}

// PageQuery 分页查询触发器
func (dao *TriggerDAOMock) PageQuery(d *dto.PageQueryTriggerDTO) ([]*po.TriggerPO, error) {
	return nil, nil
}

// Count 查询满足条件的触发器数量
func (dao *TriggerDAOMock) Count(d *dto.PageQueryTriggerDTO) (int64, error) {
	return 0, nil
}

// QueryTriggerByName 根据名称查询触发器列表
func (dao *TriggerDAOMock) QueryTriggerByName(d *dto.QueryTriggerDTO) ([]*po.TriggerPO, error) {
	return nil, nil
}
