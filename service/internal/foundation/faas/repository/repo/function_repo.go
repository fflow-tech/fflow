package repo

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// FunctionRepo 实体
type FunctionRepo struct {
	functionDAO   storage.FunctionDAO
	runHistoryDAO storage.RunHistoryDAO
}

// NewFunctionRepo 实体构造函数
func NewFunctionRepo(function *sql.FunctionDAO, history *sql.RunHistoryDAO) *FunctionRepo {
	return &FunctionRepo{functionDAO: function, runHistoryDAO: history}
}

// Get 查询函数
func (t *FunctionRepo) Get(d *dto.GetFunctionReqDTO) (*entity.Function, error) {
	functionPO, err := t.functionDAO.Get(d)
	if err != nil {
		return nil, err
	}
	return convertor.FunctionConvertor.ConvertPOToEntity(functionPO)
}

// PageQuery 查询函数
func (t *FunctionRepo) PageQuery(d *dto.PageQueryFunctionDTO) ([]*entity.Function, int64, error) {
	if d.PageQuery == nil {
		d.PageQuery = constants.NewDefaultPageQuery()
	}

	functionPOs, err := t.functionDAO.PageQueryLastVersion(d)
	if err != nil {
		return nil, 0, err
	}
	functionEntities, err := convertor.FunctionConvertor.ConvertPOsToEntities(functionPOs)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.functionDAO.Count(d)
	if err != nil {
		return nil, 0, err
	}
	return functionEntities, total, nil
}

// Create 创建函数
func (t *FunctionRepo) Create(d *dto.CreateFunctionDTO) (uint, error) {
	functionPO, err := t.functionDAO.Create(d)
	if err != nil {
		return 0, err
	}
	return functionPO.ID, nil
}

// CreateOrUpdateIfExists 创建函数
func (t *FunctionRepo) CreateOrUpdateIfExists(d *dto.CreateFunctionDTO) (uint, error) {
	functionPO, err := t.functionDAO.Get(&dto.GetFunctionReqDTO{
		Namespace: d.Namespace,
		Function:  d.Function,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			createdFunctionPO, err := t.functionDAO.Create(d)
			if err != nil {
				return 0, err
			}

			return createdFunctionPO.ID, nil
		}

		return 0, err
	}

	createReq := &dto.CreateFunctionDTO{
		Namespace:   functionPO.Namespace,
		Function:    functionPO.Name,
		Code:        d.Code,
		Description: d.Description,
		Token:       functionPO.Token,
		Version:     functionPO.Version + 1,
		Updater:     d.Updater,
		Creator:     functionPO.Creator,
		Language:    d.Language,
	}
	return t.Create(createReq)
}

// Delete 查询函数
func (t *FunctionRepo) Delete(d *dto.DeleteFunctionDTO) error {
	return t.functionDAO.Delete(d)
}

// CreateRunHistory 查询执行记录
func (t *FunctionRepo) CreateRunHistory(d *dto.CreateRunHistoryDTO) (*entity.RunHistory, error) {
	history, err := t.runHistoryDAO.Create(d)
	if err != nil {
		return nil, err
	}
	return convertor.FunctionConvertor.ConvertHistoryPOToEntity(history)
}

// UpdateRunHistory 更新执行记录
func (t *FunctionRepo) UpdateRunHistory(d *dto.UpdateRunHistoryDTO) error {
	return t.runHistoryDAO.Update(d)
}

// PageQueryRunHistory 获取执行记录列表
func (t *FunctionRepo) PageQueryRunHistory(d *dto.PageQueryRunHistoryDTO) ([]*entity.RunHistory, int64, error) {
	if d.PageQuery == nil {
		d.PageQuery = constants.NewDefaultPageQuery()
	}

	historyPOs, err := t.runHistoryDAO.PageQuery(d)
	if err != nil {
		return nil, 0, err
	}
	historyEntities, err := convertor.FunctionConvertor.ConvertHistoryPOsToEntities(historyPOs)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.runHistoryDAO.Count(d)
	if err != nil {
		return nil, 0, err
	}
	return historyEntities, total, nil
}

// BatchDeleteRunHistory 批量删除
func (t *FunctionRepo) BatchDeleteRunHistory(d *dto.BatchDeleteRunHistoryDTO) error {
	return t.runHistoryDAO.BatchDelete(d)
}
