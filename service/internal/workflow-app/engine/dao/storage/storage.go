package storage

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
)

// Transaction 执行事务
type Transaction interface {
	Transaction(fun func(*mysql.Client) error) error
}

// WorkflowDefDAO 存储层接口
type WorkflowDefDAO interface {
	Transaction
	Create(d *dto.CreateWorkflowDefDTO) (*po.WorkflowDefPO, error)
	BatchCreate(d []*dto.CreateWorkflowDefDTO) error
	Get(d *dto.GetWorkflowDefDTO) (*po.WorkflowDefPO, error)
	Delete(d *dto.DeleteWorkflowDefDTO) error
	Update(d *dto.UpdateWorkflowDefDTO) error
	PageQueryLastVersion(d *dto.PageQueryWorkflowDefDTO) ([]*po.WorkflowDefPO, error)
	Count(d *dto.PageQueryWorkflowDefDTO) (int64, error)
	GetLastVersion(d *dto.GetWorkflowDefDTO) (*po.WorkflowDefPO, error)
	GetSubWorkflowLastVersion(d *dto.GetSubworkflowDefDTO) (*po.WorkflowDefPO, error)
}

// WorkflowInstDAO 存储层接口
type WorkflowInstDAO interface {
	Transaction
	Create(d *dto.CreateWorkflowInstRepoDTO) (*po.WorkflowInstPO, error)
	Get(d *dto.GetWorkflowInstDTO) (*po.WorkflowInstPO, error)
	Delete(d *dto.DeleteWorkflowInstsDTO) error
	Update(d *dto.UpdateWorkflowInstDTO) error
	UpdateFailed(d *dto.UpdateWorkflowInstFailedDTO) error
	PageQuery(d *dto.PageQueryWorkflowInstDTO) ([]*po.WorkflowInstPO, error)
	Count(d *dto.PageQueryWorkflowInstDTO) (int64, error)
	GetWorkflowInstsByIDs(d *dto.GetWorkflowInstsByIDsDTO) ([]*po.WorkflowInstPO, error)
	DeleteWorkflowInstsByIDs(d *dto.DeleteWorkflowInstsByIDsDTO) error
}

// NodeInstDAO 存储层接口
type NodeInstDAO interface {
	Transaction
	Create(d *dto.CreateNodeInstDTO) (*po.NodeInstPO, error)
	Get(d *dto.GetNodeInstDTO) (*po.NodeInstPO, error)
	Delete(d *dto.DeleteNodeInstDTO) error
	Update(d *dto.UpdateNodeInstDTO) error
	PageQuery(d *dto.PageQueryNodeInstDTO) ([]*po.NodeInstPO, error)
	Count(d *dto.PageQueryNodeInstDTO) (int64, error)
	GetNodeInstsByIDs(d *dto.GetNodeInstsByIDsDTO) ([]*po.NodeInstPO, error)
	DeleteNodeInstsByIDs(d *dto.DeleteNodeInstsByIDsDTO) error
}

// HistoryWorkflowInstDAO 存储层接口
type HistoryWorkflowInstDAO interface {
	Transaction
	BatchCreate([]*dto.HistoryWorkflowInstDTO) error
}

// HistoryNodeInstDAO 存储层接口
type HistoryNodeInstDAO interface {
	Transaction
	BatchCreate(req []*dto.HistoryNodeInstDTO) error
}

// TriggerDAO 存储层接口
type TriggerDAO interface {
	Transaction
	CreateTrigger(d *dto.CreateTriggerDTO) (*po.TriggerPO, error)
	GetTrigger(d *dto.GetTriggerDTO) (*po.TriggerPO, error)
	DeleteTrigger(d *dto.DeleteTriggerDTO) error
	UpdateTrigger(d *dto.UpdateTriggerDTO) error
	PageQuery(d *dto.PageQueryTriggerDTO) ([]*po.TriggerPO, error)
	Count(d *dto.PageQueryTriggerDTO) (int64, error)
	QueryTriggerByName(d *dto.QueryTriggerDTO) ([]*po.TriggerPO, error)
}
