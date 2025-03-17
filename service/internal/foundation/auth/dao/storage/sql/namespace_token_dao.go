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

// NamespaceTokenDAO 数据访问对象
type NamespaceTokenDAO struct {
	db *mysql.Client
}

// NewNamespaceTokenDAO 数据访问对象构造函数
func NewNamespaceTokenDAO(db *mysql.Client) *NamespaceTokenDAO {
	return &NamespaceTokenDAO{db: db}
}

// Transaction 事务
func (dao *NamespaceTokenDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create	创建NamespaceToken
func (dao *NamespaceTokenDAO) Create(req *dto.CreateNamespaceTokenDTO) (*po.NamespaceTokenPO, error) {
	var err error
	id, err := seq.NewUint()
	if err != nil {
		return nil, err
	}

	p, err := convertor.NamespaceTokenConvertor.ConvertDTOToPO(req)
	if err != nil {
		return nil, err
	}
	p.ID = id
	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create namespace token, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Get 获取NamespaceToken信息
func (dao *NamespaceTokenDAO) Get(req *dto.GetNamespaceTokenDTO) (*po.NamespaceTokenPO, error) {
	if utils.IsZero(req.Namespace) {
		return nil, fmt.Errorf("get namespace token `namespace` must not be zero")
	}

	p, err := convertor.NamespaceTokenConvertor.ConvertGetDTOToPO(req)
	if err != nil {
		return nil, err
	}
	r := &po.NamespaceTokenPO{}
	if err := dao.db.Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get namespace token, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// Delete 删除NamespaceToken
func (dao *NamespaceTokenDAO) Delete(req *dto.DeleteNamespaceTokenDTO) error {
	if utils.IsZero(req.ID) {
		return fmt.Errorf("delete namespace token `id` must not be zero")
	}

	p, err := convertor.NamespaceTokenConvertor.ConvertDeleteDTOToPO(req)
	if err != nil {
		return err
	}
	if err := dao.db.Delete(p).Error; err != nil {
		log.Errorf("Failed to delete namespace token, caused by %s", err)
		return err
	}
	return nil
}

// Update 更新NamespaceToken信息
func (dao *NamespaceTokenDAO) Update(req *dto.UpdateNamespaceTokenDTO) error {
	if utils.IsZero(req.ID) {
		return fmt.Errorf("delete namespace token `id` must not be zero")
	}

	p, err := convertor.NamespaceTokenConvertor.ConvertUpdateDTOToPO(req)
	if err != nil {
		return err
	}
	if err := dao.db.Where("ID=? ", req.ID).Updates(p).Error; err != nil {
		log.Errorf("Failed to delete namespace token, caused by %s", err)
		return err
	}
	return nil
}

// PageQuery 分页查询NamespaceToken
func (dao *NamespaceTokenDAO) PageQuery(d *dto.PageQueryNamespaceTokenDTO) ([]*po.NamespaceTokenPO, error) {
	p, err := convertor.NamespaceTokenConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return nil, err
	}

	db := dao.db.Model(&po.NamespaceTokenPO{})
	var ps []*po.NamespaceTokenPO
	if err := db.Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).Find(&ps).Error; err != nil {
		log.Errorf("Failed to page query namespace tokens, caused by %s", err)
		return nil, err
	}
	return ps, nil
}

// Count 统计NamespaceToken数量
func (dao *NamespaceTokenDAO) Count(d *dto.PageQueryNamespaceTokenDTO) (int64, error) {
	p, err := convertor.NamespaceTokenConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return 0, err
	}

	db := dao.db.Model(&po.NamespaceTokenPO{})
	var totalCount int64
	if err := db.Where(p).Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to count namespace tokens, caused by %s", err)
		return 0, err
	}
	return totalCount, nil
}
