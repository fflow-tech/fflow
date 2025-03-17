package ports

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
)

// FunctionRepository 仓储层接口
type FunctionRepository interface {
	Get(*dto.GetFunctionReqDTO) (*entity.Function, error)
	Create(*dto.CreateFunctionDTO) (uint, error)
	CreateOrUpdateIfExists(*dto.CreateFunctionDTO) (uint, error)
	PageQuery(*dto.PageQueryFunctionDTO) ([]*entity.Function, int64, error)
	Delete(*dto.DeleteFunctionDTO) error
	CreateRunHistory(*dto.CreateRunHistoryDTO) (*entity.RunHistory, error)
	UpdateRunHistory(*dto.UpdateRunHistoryDTO) error
	PageQueryRunHistory(*dto.PageQueryRunHistoryDTO) ([]*entity.RunHistory, int64, error)
	BatchDeleteRunHistory(*dto.BatchDeleteRunHistoryDTO) error
}
