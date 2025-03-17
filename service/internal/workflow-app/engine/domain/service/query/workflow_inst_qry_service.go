package query

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowInstQueryService 查询服务
type WorkflowInstQueryService struct {
	workflowInstRepo ports.WorkflowInstRepository
}

// NewWorkflowInstQueryService 新建查询服务
func NewWorkflowInstQueryService(r *repo.WorkflowInstRepo) *WorkflowInstQueryService {
	return &WorkflowInstQueryService{workflowInstRepo: r}
}

// GetWorkflowInst 查询工作流实例
func (m *WorkflowInstQueryService) GetWorkflowInst(ctx context.Context,
	req *dto.GetWorkflowInstDTO) (*dto.WorkflowInstDTO, error) {
	inst, err := m.workflowInstRepo.Get(req)
	if err != nil {
		log.Errorf("Failed to GetWorkflowInst, caused by %s, req:%s", err, utils.StructToJsonStr(req))
		return nil, err
	}

	return convertor.InstConvertor.ConvertEntityToDTO(inst)
}

// GetWorkflowInstList 查询多条工作流实例
func (m *WorkflowInstQueryService) GetWorkflowInstList(ctx context.Context, req *dto.GetWorkflowInstListDTO) (
	[]*dto.WorkflowInstDTO, int64, error) {
	insts, total, err := m.workflowInstRepo.PageQuery(req)
	if err != nil {
		return nil, 0, err
	}

	var r []*dto.WorkflowInstDTO
	for _, inst := range insts {
		instDTO, err := convertor.InstConvertor.ConvertEntityToDTO(inst)
		if err != nil {
			return nil, 0, err
		}
		r = append(r, instDTO)
	}

	return r, total, nil
}
