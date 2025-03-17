package query

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/ports"
)

// FunctionQueryService 读服务
type FunctionQueryService struct {
	functionRepo ports.FunctionRepository
}

// NewFunctionQueryService 新建服务
func NewFunctionQueryService(repoProviderSet *ports.RepoProviderSet) *FunctionQueryService {
	return &FunctionQueryService{
		functionRepo: repoProviderSet.FunctionRepo(),
	}
}

// GetFunction 获取函数详情
func (m *FunctionQueryService) GetFunction(req *dto.GetFunctionReqDTO) (*dto.GetFunctionRspDTO, error) {
	// DAO层获取函数详情
	funcPO, err := m.functionRepo.Get(req)
	if err != nil {
		return nil, err
	}

	functionDTO, err := convertor.FunctionConvertor.ConvertEntityToGetDTO(funcPO)
	if err != nil {
		return nil, err
	}
	return functionDTO, nil
}

// GetFunctions 获取函数列表
func (m *FunctionQueryService) GetFunctions(req *dto.PageQueryFunctionDTO) ([]*dto.GetFunctionRspDTO,
	int64, error) {
	functions, total, err := m.functionRepo.PageQuery(req)
	if err != nil {
		return nil, 0, err
	}

	functionDTOs, err := convertor.FunctionConvertor.ConvertEntitiesToDTOs(functions)
	if err != nil {
		return nil, 0, err
	}
	return functionDTOs, total, err
}

// GetRunHistories 获取函数执行历史列表
func (m *FunctionQueryService) GetRunHistories(req *dto.PageQueryRunHistoryDTO) ([]*dto.GetRunHistoryRspDTO,
	int64, error) {
	histories, total, err := m.functionRepo.PageQueryRunHistory(req)
	if err != nil {
		return nil, 0, err
	}

	historyDTOs, err := convertor.FunctionConvertor.ConvertHistoryEntitiesToDTOs(histories)
	if err != nil {
		return nil, 0, err
	}
	return historyDTOs, total, err
}
