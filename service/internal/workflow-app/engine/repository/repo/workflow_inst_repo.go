package repo

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowInstRepo 流程实例仓储层
type WorkflowInstRepo struct {
	workflowInstDAO storage.WorkflowInstDAO
	nodeInstDAO     storage.NodeInstDAO
}

// NewWorkflowInstRepo 实体构造函数
func NewWorkflowInstRepo(workflowInstDAO *sql.WorkflowInstDAO, nodeInstDAO *sql.NodeInstDAO) *WorkflowInstRepo {
	return &WorkflowInstRepo{workflowInstDAO: workflowInstDAO, nodeInstDAO: nodeInstDAO}
}

// Create 创建流程实例
func (t *WorkflowInstRepo) Create(d *dto.CreateWorkflowInstRepoDTO) (string, error) {
	if err := ValidateWorkflowInstCtxSize(d.Context); err != nil {
		return "", err
	}
	inst, err := t.workflowInstDAO.Create(d)
	if err != nil {
		return "", err
	}

	return utils.UintToStr(inst.ID), nil
}

// Get 查询单条数据并转换成领域模型
func (t *WorkflowInstRepo) Get(d *dto.GetWorkflowInstDTO) (*entity.WorkflowInst, error) {
	instPO, err := t.workflowInstDAO.Get(d)
	if err != nil {
		return nil, err
	}
	d.DefID = utils.Uint64ToStr(instPO.DefID)
	curNodeInsts, err := t.getTotalNodeInsts(d)
	if err != nil {
		return nil, err
	}

	r, err := convertor.InstConvertor.ConvertPOToEntity(instPO, curNodeInsts, d.CurNodeInstID)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PageQuery 分页查询
func (t *WorkflowInstRepo) PageQuery(req *dto.GetWorkflowInstListDTO) ([]*entity.WorkflowInst, int64, error) {
	if req.PageQuery == nil {
		req.PageQuery = constants.NewDefaultPageQuery()
	}

	query := &dto.PageQueryWorkflowInstDTO{
		Namespace:     req.Namespace,
		DefID:         req.DefID,
		PageQuery:     constants.NewPageQuery(req.PageIndex, req.PageSize),
		Status:        req.Status,
		Name:          req.Name,
		Creator:       req.Operator,
		CreatedBefore: req.CreatedBefore,
		ReadFromSlave: req.ReadFromSlave,
	}
	instPOs, err := t.workflowInstDAO.PageQuery(query)
	if err != nil {
		return nil, 0, err
	}
	var insts []*entity.WorkflowInst
	for _, instPO := range instPOs {
		inst, err := convertor.InstConvertor.ConvertPOToEntity(instPO, []*entity.NodeInst{}, "")
		if err != nil {
			return nil, 0, err
		}
		insts = append(insts, inst)
	}
	total, err := t.workflowInstDAO.Count(query)
	if err != nil {
		return nil, 0, err
	}

	return insts, total, nil
}

// Count 统计一个定义下的实例数量
func (t *WorkflowInstRepo) Count(req *dto.GetWorkflowInstListDTO) (int64, error) {
	query := &dto.PageQueryWorkflowInstDTO{DefID: req.DefID}
	return t.workflowInstDAO.Count(query)
}

// GetWorkflowInstCtx 获取流程实例上下文
func (t *WorkflowInstRepo) GetWorkflowInstCtx(d *dto.GetWorkflowInstDTO) (map[string]interface{}, error) {
	inst, err := t.Get(d)
	if err != nil {
		return nil, fmt.Errorf("[%s]failed to get workflow inst: %w", d.InstID, err)
	}

	return entity.ConvertToCtx(inst)
}

func (t *WorkflowInstRepo) getTotalNodeInsts(d *dto.GetWorkflowInstDTO) ([]*entity.NodeInst, error) {
	queryDTO := &dto.PageQueryNodeInstDTO{
		DefID:     d.DefID,
		InstID:    d.InstID,
		PageQuery: constants.NewPageQuery(1, config.GetValidationRulesConfig().MaxNodeInstsForOneFlow),
	}
	nodePOs, err := t.nodeInstDAO.PageQuery(queryDTO)

	if err != nil {
		return nil, err
	}
	var nodeInsts []*entity.NodeInst
	for _, nodePO := range nodePOs {
		nodeInst, err := convertor.NodeInstConvertor.ConvertPOToEntity(nodePO)
		if err != nil {
			return nil, err
		}
		nodeInsts = append(nodeInsts, nodeInst)
	}
	return nodeInsts, nil
}

// UpdateWithDefID 更新流程实例
func (t *WorkflowInstRepo) UpdateWithDefID(d *dto.UpdateWorkflowInstDTO) error {
	if err := ValidateWorkflowInstCtxSize(d.Context); err != nil {
		return err
	}
	return t.workflowInstDAO.Update(d)
}

// UpdateWorkflowInstFailed 更新流程实例为失败
func (t *WorkflowInstRepo) UpdateWorkflowInstFailed(d *dto.UpdateWorkflowInstFailedDTO) error {
	return t.workflowInstDAO.UpdateFailed(d)
}
