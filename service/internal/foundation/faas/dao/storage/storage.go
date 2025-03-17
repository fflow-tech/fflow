package storage

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
)

// Transaction 执行事务
type Transaction interface {
	Transaction(fun func(*mysql.Client) error) error
}

// FunctionDAO 存储层接口
type FunctionDAO interface {
	Transaction
	Create(*dto.CreateFunctionDTO) (*po.FunctionPO, error)
	Get(*dto.GetFunctionReqDTO) (*po.FunctionPO, error)
	Delete(*dto.DeleteFunctionDTO) error
	Update(*dto.CreateFunctionDTO) error
	PageQueryLastVersion(*dto.PageQueryFunctionDTO) ([]*po.FunctionPO, error)
	Count(*dto.PageQueryFunctionDTO) (int64, error)
}

// RunHistoryDAO 存储层接口
type RunHistoryDAO interface {
	Transaction
	Create(*dto.CreateRunHistoryDTO) (*po.RunHistoryPO, error)
	Get(*dto.GetRunHistoryDTO) (*po.RunHistoryPO, error)
	Delete(*dto.DeleteRunHistoryDTO) error
	BatchDelete(d *dto.BatchDeleteRunHistoryDTO) error
	Update(*dto.UpdateRunHistoryDTO) error
	PageQuery(*dto.PageQueryRunHistoryDTO) ([]*po.RunHistoryPO, error)
	Count(*dto.PageQueryRunHistoryDTO) (int64, error)
	DeleteByCreateTime(time.Time) error
}
