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

// UserRepo 实体
type UserRepo struct {
	storage.UserDAO
}

// NewUserRepo 实体构造函数
func NewUserRepo(dao *sql.UserDAO) *UserRepo {
	return &UserRepo{
		UserDAO: dao,
	}
}

// Create 创建用户
func (r *UserRepo) Create(userDTO *dto.CreateUserDTO) (string, error) {
	userPO, err := r.UserDAO.Create(userDTO)
	if err != nil {
		return "", err
	}

	return utils.ToString(userPO.ID), nil
}

// CreateIfNotExists 如果不存在则创建用户
func (r *UserRepo) CreateIfNotExists(userDTO *dto.CreateUserDTO) (*entity.User, error) {
	user, err := r.UserDAO.Get(&dto.GetUserDTO{Username: userDTO.Username})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userPO, err := r.UserDAO.Create(userDTO)
			if err != nil {
				return nil, err
			}

			return convertor.AuthConvertor.ConvertPOToEntity(userPO)
		}

		return nil, err
	}

	return convertor.AuthConvertor.ConvertPOToEntity(user)
}

// Delete 删除用户
func (r *UserRepo) Delete(userDTO *dto.DeleteUserDTO) error {
	return r.UserDAO.Delete(userDTO)
}

// Update 更新用户
func (r *UserRepo) Update(userDTO *dto.UpdateUserDTO) error {
	return r.UserDAO.Update(userDTO)
}

// Get 获取用户
func (r *UserRepo) Get(userDTO *dto.GetUserDTO) (*entity.User, error) {
	userPO, err := r.UserDAO.Get(userDTO)
	if err != nil {
		return nil, err
	}
	return convertor.AuthConvertor.ConvertPOToEntity(userPO)
}

// PageQuery 分页查询用户
func (r *UserRepo) PageQuery(userDTO *dto.PageQueryUserDTO) ([]*entity.User, error) {
	userPOs, err := r.UserDAO.PageQuery(userDTO)
	if err != nil {
		return nil, err
	}
	return convertor.AuthConvertor.ConvertPOsToEntities(userPOs)
}
