package mock

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
)

// WorkflowInstDAOMock 模拟对象
type WorkflowInstDAOMock struct {
}

// NewWorkflowInstDAOMock 构造方法
func NewWorkflowInstDAOMock() *WorkflowDefDAOMock {
	return &WorkflowDefDAOMock{}
}

// Transaction 事务
func (dao *WorkflowInstDAOMock) Transaction(fun func() error) error {
	return nil
}

// CreateWorkflowInst 创建流程实例
func (dao *WorkflowInstDAOMock) CreateWorkflowInst(d *dto.CreateWorkflowInstRepoDTO) (*po.WorkflowInstPO, error) {
	return nil, nil
}

// GetWorkflowInst 获取流程实例
func (dao *WorkflowInstDAOMock) GetWorkflowInst(d *dto.GetWorkflowInstDTO) (*po.WorkflowInstPO, error) {
	return nil, nil
}

// DeleteWorkflowInst 删除流程实例
func (dao *WorkflowInstDAOMock) DeleteWorkflowInst(d *dto.DeleteWorkflowInstsDTO) error {
	return nil
}

// UpdateWorkflowInst 更新流程实例
func (dao *WorkflowInstDAOMock) UpdateWorkflowInst(d *dto.UpdateWorkflowInstDTO) error {
	return nil
}

// PageQueryWorkflowInst 分页查询流程实例
func (dao *WorkflowInstDAOMock) PageQueryWorkflowInst(d *dto.PageQueryWorkflowInstDTO) ([]*po.WorkflowInstPO, error) {
	return nil, nil
}

// Count 查询满足条件的流程实例数量
func (dao *WorkflowInstDAOMock) Count(d *dto.PageQueryWorkflowInstDTO) (int64, error) {
	return 0, nil
}
