package repo

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/jinzhu/copier"
)

// WorkflowArchiveRepo 流程归档仓储层
type WorkflowArchiveRepo struct {
	workflowInstDAO        storage.WorkflowInstDAO
	nodeInstDAO            storage.NodeInstDAO
	historyWorkflowInstDAO storage.HistoryWorkflowInstDAO
	historyNodeInstDAO     storage.HistoryNodeInstDAO
}

// NewWorkflowArchiveRepo 创建流程归档仓储层
func NewWorkflowArchiveRepo(workflowInstDAO *sql.WorkflowInstDAO,
	nodeInstDAO *sql.NodeInstDAO,
	historyWorkflowInstDAO *sql.HistoryWorkflowInstDAO,
	historyNodeInstDAO *sql.HistoryNodeInstDAO) *WorkflowArchiveRepo {
	return &WorkflowArchiveRepo{
		workflowInstDAO:        workflowInstDAO,
		nodeInstDAO:            nodeInstDAO,
		historyWorkflowInstDAO: historyWorkflowInstDAO,
		historyNodeInstDAO:     historyNodeInstDAO,
	}
}

// ArchiveWorkflowInst 归档流程实例
func (t *WorkflowArchiveRepo) ArchiveWorkflowInst(req *dto.ArchiveWorkflowInstsDTO) error {
	insts, err := t.workflowInstDAO.GetWorkflowInstsByIDs(&dto.GetWorkflowInstsByIDsDTO{
		DefID:   req.DefID,
		InstIDs: req.InstIDs,
	})
	if err != nil {
		return err
	}

	historyInsts := make([]*dto.HistoryWorkflowInstDTO, 0)
	for _, inst := range insts {
		historyInst := dto.HistoryWorkflowInstDTO{}
		copier.Copy(&historyInst, inst)
		historyInsts = append(historyInsts, &historyInst)
	}

	if err := t.historyWorkflowInstDAO.BatchCreate(historyInsts); err != nil {
		return err
	}

	return t.workflowInstDAO.DeleteWorkflowInstsByIDs(&dto.DeleteWorkflowInstsByIDsDTO{
		DefID:   req.DefID,
		InstIDs: req.InstIDs,
	})
}

// ArchiveNodeInst 归档节点实例
func (t *WorkflowArchiveRepo) ArchiveNodeInst(req *dto.ArchiveNodeInstsDTO) error {
	nodeInsts, err := t.nodeInstDAO.GetNodeInstsByIDs(&dto.GetNodeInstsByIDsDTO{
		DefID:       req.DefID,
		InstID:      req.InstID,
		NodeInstIDs: req.NodeInstIDs,
	})
	if err != nil {
		return err
	}

	historyNodeInsts := make([]*dto.HistoryNodeInstDTO, 0)
	for _, nodeInst := range nodeInsts {
		historyInst := dto.HistoryNodeInstDTO{}
		copier.Copy(&historyInst, nodeInst)
		historyNodeInsts = append(historyNodeInsts, &historyInst)
	}

	if err := t.historyNodeInstDAO.BatchCreate(historyNodeInsts); err != nil {
		return err
	}

	return t.nodeInstDAO.DeleteNodeInstsByIDs(&dto.DeleteNodeInstsByIDsDTO{
		DefID:       req.DefID,
		InstID:      req.InstID,
		NodeInstIDs: req.NodeInstIDs,
	})
}
