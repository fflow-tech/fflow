package query

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// WorkflowDefQueryService 查询服务
type WorkflowDefQueryService struct {
	workflowDefRepo ports.WorkflowDefRepository
}

// NewWorkflowDefQueryService 新建查询服务
func NewWorkflowDefQueryService(r *repo.WorkflowDefRepo) *WorkflowDefQueryService {
	return &WorkflowDefQueryService{workflowDefRepo: r}
}

// GetWorkflowDefByDefID 查询工作流
func (m *WorkflowDefQueryService) GetWorkflowDefByDefID(ctx context.Context,
	d *dto.GetWorkflowDefDTO) (*dto.WorkflowDefDTO, error) {
	e, err := m.workflowDefRepo.GetLastVersion(d)
	if err != nil {
		log.Errorf("Failed to get workflow def last version, defID:[%d], caused by %s", d.DefID, err)
		return nil, err
	}

	return convertor.DefConvertor.ConvertEntityToDTO(e)
}

// GetSubworkflowByParentDefIDAndRefName 根据父工作流 ID 和 子工作流的 RefName 查询子工作流
func (m *WorkflowDefQueryService) GetSubworkflowByParentDefIDAndRefName(ctx context.Context,
	d *dto.GetSubworkflowDefDTO) (
	*dto.WorkflowDefDTO, error) {
	e, err := m.workflowDefRepo.GetSubworkflowLastVersion(d)
	if err != nil {
		log.Errorf("Failed to get workflow last version parentDefID:[%d] refName:[%s], caused by %s",
			d.ParentDefID, d.RefName, err)
		return nil, err
	}

	return convertor.DefConvertor.ConvertEntityToDTO(e)
}

// GetWorkflowDefList 批量查询工作流
func (m *WorkflowDefQueryService) GetWorkflowDefList(ctx context.Context, d *dto.PageQueryWorkflowDefDTO) (
	[]*dto.WorkflowDefDTO, int64, error) {
	defEntityList, total, err := m.workflowDefRepo.PageQueryLastVersionForWeb(d)
	if err != nil {
		log.Errorf("Failed to page query def last version, defID:%d, caused by %s", d, err)
		return nil, 0, err
	}

	defDTOList, err := convertor.DefConvertor.ConvertEntitiesToDTOs(defEntityList)
	return defDTOList, total, err
}
