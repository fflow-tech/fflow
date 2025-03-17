package mock

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
)

// WorkflowDefDAOMock 模拟对象
type WorkflowDefDAOMock struct {
}

// NewWorkflowDefDAOMock 初始化模拟对象方法
func NewWorkflowDefDAOMock() *WorkflowDefDAOMock {
	return &WorkflowDefDAOMock{}
}

// Transaction 事务
func (dao *WorkflowDefDAOMock) Transaction(fun func() error) error {
	return nil
}

// CreateWorkflowDef 创建流程定义
func (dao *WorkflowDefDAOMock) CreateWorkflowDef(d *dto.CreateWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	return nil, nil
}

// GetWorkflowDef 获取流程定义
func (dao *WorkflowDefDAOMock) GetWorkflowDef(d *dto.GetWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	return nil, nil
}

// DeleteWorkflowDef 删除流程定义
func (dao *WorkflowDefDAOMock) DeleteWorkflowDef(d *dto.DeleteWorkflowDefDTO) error {
	return nil
}

// UpdateWorkflowDef 更新流程定义
func (dao *WorkflowDefDAOMock) UpdateWorkflowDef(d *dto.UpdateWorkflowDefDTO) error {
	return nil
}

// PageQueryWorkflowDef 分页查询流程定义
func (dao *WorkflowDefDAOMock) PageQueryWorkflowDef(d *dto.PageQueryWorkflowDefDTO) ([]*po.WorkflowDefPO, error) {
	return nil, nil
}

// Count 根据条件获取数量
func (dao *WorkflowDefDAOMock) Count(d *dto.PageQueryWorkflowDefDTO) (int64, error) {
	return 0, nil
}

// GetWorkflowLastVersion 获取流程最新版本
func (dao *WorkflowDefDAOMock) GetWorkflowLastVersion(d *dto.GetWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	return nil, nil
}
