package sql

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/seq"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"gorm.io/gorm"
)

// UserDAO 数据访问对象
type UserDAO struct {
	db *mysql.Client
}

// NewUserDAO 数据访问对象构造函数
func NewUserDAO(db *mysql.Client) *UserDAO {
	return &UserDAO{db: db}
}

// Transaction 事务
func (dao *UserDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create	创建用户
func (dao *UserDAO) Create(req *dto.CreateUserDTO) (*po.UserPO, error) {
	var err error
	id, err := seq.NewUint()
	if err != nil {
		return nil, err
	}

	p, err := convertor.UserConvertor.ConvertDTOToPO(req)
	if err != nil {
		return nil, err
	}
	p.ID = id
	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create user, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Get 获取用户信息
func (dao *UserDAO) Get(userDTO *dto.GetUserDTO) (*po.UserPO, error) {
	if utils.IsZero(userDTO.Username) {
		return nil, fmt.Errorf("get user `username` must not be zero")
	}

	p, err := convertor.UserConvertor.ConvertGetDTOToPO(userDTO)
	if err != nil {
		return nil, err
	}
	r := &po.UserPO{}
	if err := dao.db.Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get user, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// Delete 删除用户
func (dao *UserDAO) Delete(userDTO *dto.DeleteUserDTO) error {
	if utils.IsZero(userDTO.ID) {
		return fmt.Errorf("delete user `id` must not be zero")
	}

	p, err := convertor.UserConvertor.ConvertDeleteDTOToPO(userDTO)
	if err != nil {
		return err
	}
	if err := dao.db.Delete(p).Error; err != nil {
		log.Errorf("Failed to delete user, caused by %s", err)
		return err
	}
	return nil
}

// Update 更新用户信息
func (dao *UserDAO) Update(userDTO *dto.UpdateUserDTO) error {
	if utils.IsZero(userDTO.ID) {
		return fmt.Errorf("delete user `id` must not be zero")
	}

	p, err := convertor.UserConvertor.ConvertUpdateDTOToPO(userDTO)
	if err != nil {
		return err
	}
	if err := dao.db.Where("ID=? ", userDTO.ID).Updates(p).Error; err != nil {
		log.Errorf("Failed to delete user, caused by %s", err)
		return err
	}
	return nil
}

// PageQuery 分页查询用户
func (dao *UserDAO) PageQuery(d *dto.PageQueryUserDTO) ([]*po.UserPO, error) {
	p, err := convertor.UserConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return nil, err
	}

	var ps []*po.UserPO
	if err := dao.db.Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).Find(&ps).Error; err != nil {
		log.Errorf("Failed to page query users, caused by %s", err)
		return nil, err
	}
	return ps, nil
}

// Count 统计用户数量
func (dao *UserDAO) Count(d *dto.PageQueryUserDTO) (int64, error) {
	p, err := convertor.UserConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return 0, err
	}

	var totalCount int64
	if err := dao.db.Model(po.UserPO{}).Where(p).Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to count users, caused by %s", err)
		return 0, err
	}
	return totalCount, nil
}
