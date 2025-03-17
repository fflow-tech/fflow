package repo

import (
	"encoding/json"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// WorkflowDefRepo 流程定义仓储层
type WorkflowDefRepo struct {
	workflowDefDAO storage.WorkflowDefDAO
}

// NewWorkflowDefRepo 实体构造函数
func NewWorkflowDefRepo(d *sql.WorkflowDefDAO) *WorkflowDefRepo {
	return &WorkflowDefRepo{workflowDefDAO: d}
}

// Create 创建
func (t *WorkflowDefRepo) Create(d *dto.CreateWorkflowDefDTO) (string, error) {
	defPO, err := t.workflowDefDAO.Create(d)
	if err != nil {
		return "", err
	}
	return utils.Uint64ToStr(defPO.DefID), nil
}

// BatchCreate 批量创建
func (t *WorkflowDefRepo) BatchCreate(d []*dto.CreateWorkflowDefDTO) (err error) {
	return t.workflowDefDAO.BatchCreate(d)
}

// Get 查询单条
func (t *WorkflowDefRepo) Get(d *dto.GetWorkflowDefDTO) (*entity.WorkflowDef, error) {
	p, err := t.workflowDefDAO.Get(d)
	if err != nil {
		return nil, err
	}
	e, err := convertor.DefConvertor.ConvertPOToEntity(p)
	if err != nil {
		return nil, err
	}
	return t.setEntitySubWorkflow(e, p)
}

// setEntitySubWorkflow 设置流程定义实体的 subworkflows 字段
func (t *WorkflowDefRepo) setEntitySubWorkflow(entity *entity.WorkflowDef, p *po.WorkflowDefPO) (*entity.WorkflowDef,
	error) {
	subWorkflows, err := t.GetAllSubworkflowDefs(&dto.GetAllSubworkflowDefsDTO{
		ParentDefID:      utils.Uint64ToStr(p.DefID),
		ParentDefVersion: p.Version,
		DefJson:          p.DefJson,
	})
	if err != nil {
		return nil, err
	}

	entity.Subworkflows = subWorkflows
	return entity, nil
}

// GetAllSubworkflowDefs 获取流程的所有子流程
func (t *WorkflowDefRepo) GetAllSubworkflowDefs(d *dto.GetAllSubworkflowDefsDTO) ([]map[string]entity.WorkflowDef,
	error) {
	def := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(d.DefJson), &def); err != nil {
		return nil, err
	}
	// 保存所有子流程的定义
	var subWorkflowDefs []map[string]entity.WorkflowDef
	for _, subWorkflow := range def.Subworkflows {
		for refName := range subWorkflow {
			// 查询获取子流程定义
			getSubWorkflowDto := &dto.GetSubworkflowDefDTO{
				RefName:          refName,
				ParentDefID:      d.ParentDefID,
				ParentDefVersion: d.ParentDefVersion,
			}
			subWorkflowDef, err := t.GetSubworkflowLastVersion(getSubWorkflowDto)
			if err != nil {
				return nil, err
			}
			subWorkflowDefs = append(subWorkflowDefs, map[string]entity.WorkflowDef{refName: *subWorkflowDef})
		}
	}

	return subWorkflowDefs, nil
}

// PageQueryLastVersion 分页查询
func (t *WorkflowDefRepo) PageQueryLastVersion(req *dto.PageQueryWorkflowDefDTO) ([]*entity.WorkflowDef, error) {
	if req.PageQuery == nil {
		req.PageQuery = constants.NewDefaultPageQuery()
	}

	p, err := t.workflowDefDAO.PageQueryLastVersion(req)
	if err != nil {
		return nil, err
	}
	var defs []*entity.WorkflowDef
	for _, defPO := range p {
		def, err := convertor.DefConvertor.ConvertPOToEntity(defPO)
		if err != nil {
			return nil, err
		}
		defWithSubWorkflows, err := t.setEntitySubWorkflow(def, defPO)
		if err != nil {
			return nil, err
		}
		defs = append(defs, defWithSubWorkflows)
	}
	return defs, nil
}

// PageQueryLastVersionForWeb 分页查询
func (t *WorkflowDefRepo) PageQueryLastVersionForWeb(req *dto.PageQueryWorkflowDefDTO) (
	[]*entity.WorkflowDef, int64, error) {
	if req.PageQuery == nil {
		req.PageQuery = constants.NewDefaultPageQuery()
	}

	p, err := t.workflowDefDAO.PageQueryLastVersion(req)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.workflowDefDAO.Count(req)
	if err != nil {
		return nil, 0, err
	}

	var defs []*entity.WorkflowDef
	for _, defPO := range p {
		def, err := convertor.DefConvertor.ConvertPOToEntity(defPO)
		if err != nil {
			return nil, 0, err
		}
		defWithSubWorkflows, err := t.setEntitySubWorkflow(def, defPO)
		if err != nil {
			return nil, 0, err
		}
		defs = append(defs, defWithSubWorkflows)
	}
	return defs, total, nil
}

// GetLastVersion 获取最新版本
func (t *WorkflowDefRepo) GetLastVersion(d *dto.GetWorkflowDefDTO) (*entity.WorkflowDef, error) {
	p, err := t.workflowDefDAO.GetLastVersion(d)
	if err != nil {
		return nil, err
	}
	e, err := convertor.DefConvertor.ConvertPOToEntity(p)
	if err != nil {
		return nil, err
	}
	return t.setEntitySubWorkflow(e, p)
}

// UpdateStatus 更新流程定义状态
func (t *WorkflowDefRepo) UpdateStatus(d *dto.UpdateWorkflowDefDTO) error {
	return t.workflowDefDAO.Update(d)
}

// GetSubworkflowLastVersion 获取子流程最新版本
func (t *WorkflowDefRepo) GetSubworkflowLastVersion(d *dto.GetSubworkflowDefDTO) (*entity.WorkflowDef, error) {
	p, err := t.workflowDefDAO.GetSubWorkflowLastVersion(d)
	if err != nil {
		return nil, err
	}
	// 因为禁止子流程的嵌套，这里不需要再去获取子流程的子流程定义
	return convertor.DefConvertor.ConvertPOToEntity(p)
}
