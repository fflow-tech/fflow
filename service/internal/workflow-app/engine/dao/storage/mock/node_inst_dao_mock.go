package mock

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
)

// NodeInstDAOMock 模拟对象
type NodeInstDAOMock struct {
}

// NewNodeInstDAOMock 构造方法
func NewNodeInstDAOMock() *WorkflowDefDAOMock {
	return &WorkflowDefDAOMock{}
}

// Transaction 事务
func (dao *NodeInstDAOMock) Transaction(fun func() error) error {
	return nil
}

// Create 创建节点实例
func (dao *NodeInstDAOMock) Create(d *dto.CreateNodeInstDTO) (*po.NodeInstPO, error) {
	return nil, nil
}

// Get 获取节点实例
func (dao *NodeInstDAOMock) Get(d *dto.GetNodeInstDTO) (*po.NodeInstPO, error) {
	return nil, nil
}

// Delete 删除节点实例
func (dao *NodeInstDAOMock) Delete(d *dto.DeleteNodeInstDTO) error {
	return nil
}

// Update 更新节点实例
func (dao *NodeInstDAOMock) Update(d *dto.UpdateNodeInstDTO) error {
	return nil
}

// PageQuery 分页查询节点实例
func (dao *NodeInstDAOMock) PageQuery(d *dto.PageQueryNodeInstDTO) ([]*po.NodeInstPO, error) {
	return nil, nil
}

// Count 统计节点实例的数量
func (dao *NodeInstDAOMock) Count(d *dto.PageQueryNodeInstDTO) (int64, error) {
	return 0, nil
}
