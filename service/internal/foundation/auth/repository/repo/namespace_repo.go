package repo

import (
	"errors"

	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/repository/convertor"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

// NamespaceRepo 实体
type NamespaceRepo struct {
	storage.NamespaceDAO
}

// NewNamespaceRepo 实体构造函数
func NewNamespaceRepo(dao *sql.NamespaceDAO) *NamespaceRepo {
	return &NamespaceRepo{
		NamespaceDAO: dao,
	}
}

// Create 创建Namespace
func (r *NamespaceRepo) Create(req *dto.CreateNamespaceDTO) (string, error) {
	userPO, err := r.NamespaceDAO.Create(req)
	if err != nil {
		return "", err
	}

	return utils.ToString(userPO.ID), nil
}

// CreateIfNotExists 如果不存在则创建Namespace
func (r *NamespaceRepo) CreateIfNotExists(req *dto.CreateNamespaceDTO) (*entity.Namespace, error) {
	user, err := r.NamespaceDAO.Get(&dto.GetNamespaceDTO{Namespace: req.Namespace})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userPO, err := r.NamespaceDAO.Create(req)
			if err != nil {
				return nil, err
			}

			return convertor.NamespaceConvertor.ConvertPOToEntity(userPO)
		}

		return nil, err
	}

	return convertor.NamespaceConvertor.ConvertPOToEntity(user)
}

// Delete 删除Namespace
func (r *NamespaceRepo) Delete(req *dto.DeleteNamespaceDTO) error {
	return r.NamespaceDAO.Delete(req)
}

// Update 更新Namespace
func (r *NamespaceRepo) Update(req *dto.UpdateNamespaceDTO) error {
	return r.NamespaceDAO.Update(req)
}

// Get 获取Namespace
func (r *NamespaceRepo) Get(req *dto.GetNamespaceDTO) (*entity.Namespace, error) {
	userPO, err := r.NamespaceDAO.Get(req)
	if err != nil {
		return nil, err
	}
	return convertor.NamespaceConvertor.ConvertPOToEntity(userPO)
}

// PageQuery 分页查询Namespace
func (r *NamespaceRepo) PageQuery(req *dto.PageQueryNamespaceDTO) ([]*entity.Namespace, int64, error) {
	total, err := r.NamespaceDAO.Count(req)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.Namespace{}, 0, nil
	}

	userPOs, err := r.NamespaceDAO.PageQuery(req)
	if err != nil {
		return nil, 0, err
	}
	entities, err := convertor.NamespaceConvertor.ConvertPOsToEntities(userPOs)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Count 统计Namespace数量
func (r *NamespaceRepo) Count(req *dto.PageQueryNamespaceDTO) (int64, error) {
	return r.NamespaceDAO.Count(req)
}
