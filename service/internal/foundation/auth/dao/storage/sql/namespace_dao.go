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

// NamespaceDAO 数据访问对象
type NamespaceDAO struct {
	db *mysql.Client
}

// NewNamespaceDAO 数据访问对象构造函数
func NewNamespaceDAO(db *mysql.Client) *NamespaceDAO {
	return &NamespaceDAO{db: db}
}

// Transaction 事务
func (dao *NamespaceDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create	创建Namespace
func (dao *NamespaceDAO) Create(req *dto.CreateNamespaceDTO) (*po.NamespacePO, error) {
	var err error
	id, err := seq.NewUint()
	if err != nil {
		return nil, err
	}

	p, err := convertor.NamespaceConvertor.ConvertDTOToPO(req)
	if err != nil {
		return nil, err
	}
	p.ID = id
	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create namespace, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Get 获取Namespace信息
func (dao *NamespaceDAO) Get(req *dto.GetNamespaceDTO) (*po.NamespacePO, error) {
	p, err := convertor.NamespaceConvertor.ConvertGetDTOToPO(req)
	if err != nil {
		return nil, err
	}
	r := &po.NamespacePO{}
	if err := dao.db.Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get namespace, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// Delete 删除Namespace
func (dao *NamespaceDAO) Delete(req *dto.DeleteNamespaceDTO) error {
	if utils.IsZero(req.ID) {
		return fmt.Errorf("delete namespace `id` must not be zero")
	}

	p, err := convertor.NamespaceConvertor.ConvertDeleteDTOToPO(req)
	if err != nil {
		return err
	}
	if err := dao.db.Delete(p).Error; err != nil {
		log.Errorf("Failed to delete namespace, caused by %s", err)
		return err
	}
	return nil
}

// Update 更新Namespace信息
func (dao *NamespaceDAO) Update(req *dto.UpdateNamespaceDTO) error {
	if utils.IsZero(req.ID) {
		return fmt.Errorf("delete namespace `id` must not be zero")
	}

	p, err := convertor.NamespaceConvertor.ConvertUpdateDTOToPO(req)
	if err != nil {
		return err
	}
	if err := dao.db.Where("ID=? ", req.ID).Updates(p).Error; err != nil {
		log.Errorf("Failed to delete namespace, caused by %s", err)
		return err
	}
	return nil
}

// PageQuery 分页查询Namespace
func (dao *NamespaceDAO) PageQuery(d *dto.PageQueryNamespaceDTO) ([]*po.NamespacePO, error) {
	p, err := convertor.NamespaceConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return nil, err
	}

	db := dao.db.Model(&po.NamespacePO{})

	if !utils.IsZero(d.Namespace) {
		db = db.Where("namespace like ?", "%"+d.Namespace+"%")
	}

	var ps []*po.NamespacePO
	if err := db.Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).Find(&ps).Error; err != nil {
		log.Errorf("Failed to page query namespaces, caused by %s", err)
		return nil, err
	}
	return ps, nil
}

// Count 统计Namespace数量
func (dao *NamespaceDAO) Count(d *dto.PageQueryNamespaceDTO) (int64, error) {
	p, err := convertor.NamespaceConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return 0, err
	}

	db := dao.db.Model(&po.NamespacePO{})

	if !utils.IsZero(d.Namespace) {
		db = db.Where("namespace like ?", "%"+d.Namespace+"%")
	}

	var totalCount int64
	if err := db.Where(p).Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to count namespaces, caused by %s", err)
		return 0, err
	}
	return totalCount, nil
}
