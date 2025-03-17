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

// NamespaceTokenRepo 实体
type NamespaceTokenRepo struct {
	storage.NamespaceTokenDAO
}

// NewNamespaceTokenRepo 实体构造函数
func NewNamespaceTokenRepo(dao *sql.NamespaceTokenDAO) *NamespaceTokenRepo {
	return &NamespaceTokenRepo{
		NamespaceTokenDAO: dao,
	}
}

// Create 创建NamespaceToken
func (r *NamespaceTokenRepo) Create(req *dto.CreateNamespaceTokenDTO) (string, error) {
	userPO, err := r.NamespaceTokenDAO.Create(req)
	if err != nil {
		return "", err
	}

	return utils.ToString(userPO.ID), nil
}

// CreateIfNotExists 如果不存在则创建NamespaceToken
func (r *NamespaceTokenRepo) CreateIfNotExists(req *dto.CreateNamespaceTokenDTO) (*entity.NamespaceToken, error) {
	user, err := r.NamespaceTokenDAO.Get(&dto.GetNamespaceTokenDTO{Namespace: req.Namespace})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userPO, err := r.NamespaceTokenDAO.Create(req)
			if err != nil {
				return nil, err
			}

			return convertor.NamespaceTokenConvertor.ConvertPOToEntity(userPO)
		}

		return nil, err
	}

	return convertor.NamespaceTokenConvertor.ConvertPOToEntity(user)
}

// Delete 删除NamespaceToken
func (r *NamespaceTokenRepo) Delete(req *dto.DeleteNamespaceTokenDTO) error {
	return r.NamespaceTokenDAO.Delete(req)
}

// Update 更新NamespaceToken
func (r *NamespaceTokenRepo) Update(req *dto.UpdateNamespaceTokenDTO) error {
	return r.NamespaceTokenDAO.Update(req)
}

// Get 获取NamespaceToken
func (r *NamespaceTokenRepo) Get(req *dto.GetNamespaceTokenDTO) (*entity.NamespaceToken, error) {
	userPO, err := r.NamespaceTokenDAO.Get(req)
	if err != nil {
		return nil, err
	}
	return convertor.NamespaceTokenConvertor.ConvertPOToEntity(userPO)
}

// PageQuery 分页查询NamespaceToken
func (r *NamespaceTokenRepo) PageQuery(req *dto.PageQueryNamespaceTokenDTO) ([]*entity.NamespaceToken, int64, error) {
	total, err := r.NamespaceTokenDAO.Count(req)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.NamespaceToken{}, 0, nil
	}

	userPOs, err := r.NamespaceTokenDAO.PageQuery(req)
	if err != nil {
		return nil, 0, err
	}
	entities, err := convertor.NamespaceTokenConvertor.ConvertPOsToEntities(userPOs)
	if err != nil {
		return nil, 0, err
	}
	return entities, total, nil
}

// Count 统计数量
func (r *NamespaceTokenRepo) Count(req *dto.PageQueryNamespaceTokenDTO) (int64, error) {
	return r.NamespaceTokenDAO.Count(req)
}
