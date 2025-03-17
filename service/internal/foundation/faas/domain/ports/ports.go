package ports

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
)

// CommandPorts 写入接口
type CommandPorts interface {
	FunctionCommandPorts
}

// QueryPorts 读入接口
type QueryPorts interface {
	FunctionQueryPorts
}

// FunctionCommandPorts 函数相关接口
type FunctionCommandPorts interface {
	CreateFunction(*dto.CreateFunctionReqDTO) (uint, error)
	UpdateFunction(*dto.UpdateFunctionDTO) (uint, error)
	DeleteFunction(*dto.DeleteFunctionDTO) error
	CallFunction(context.Context, *dto.CallFunctionReqDTO) (interface{}, error)
	DebugFunction(context.Context, *dto.DebugFunctionDTO) (*dto.DebugFunctionRspDTO, error)
	BatchDeleteExpiredRunHistory(*dto.BatchDeleteExpiredRunHistoryDTO) error
}

// FunctionQueryPorts 函数相关接口
type FunctionQueryPorts interface {
	GetFunction(*dto.GetFunctionReqDTO) (*dto.GetFunctionRspDTO, error)
	GetFunctions(*dto.PageQueryFunctionDTO) ([]*dto.GetFunctionRspDTO, int64, error)
	GetRunHistories(*dto.PageQueryRunHistoryDTO) ([]*dto.GetRunHistoryRspDTO, int64, error)
}
