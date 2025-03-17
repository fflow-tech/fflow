package sql

import (
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"gorm.io/gorm"
)

// RunHistoryDAO RunHistory数据访问对象
type RunHistoryDAO struct {
	db *mysql.Client
}

// NewRunHistoryDAO RunHistoryDAO数据访问对象构造函数
func NewRunHistoryDAO(db *mysql.Client) *RunHistoryDAO {
	return &RunHistoryDAO{db: db}
}

// Transaction 事务
func (dao *RunHistoryDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建
func (dao *RunHistoryDAO) Create(d *dto.CreateRunHistoryDTO) (*po.RunHistoryPO, error) {
	if utils.IsZero(d.Namespace) {
		return nil, fmt.Errorf("create RunHistory `app` or `server` must not be empty, namespace:[%s]",
			d.Namespace)
	}
	p := convertor.RunHistoryConvertor.ConvertCreateDTOToPO(d)
	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create runhistory, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Get 获取
func (dao *RunHistoryDAO) Get(d *dto.GetRunHistoryDTO) (*po.RunHistoryPO, error) {
	if utils.IsZero(d.Namespace) {
		return nil, fmt.Errorf("get runHistory `app` or `server` must not be empty, namespace:[%s]",
			d.Namespace)
	}
	p := convertor.RunHistoryConvertor.ConvertGetDTOToPO(d)
	r := &po.RunHistoryPO{}
	if err := dao.db.Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get function run history, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// Delete 根据删除
func (dao *RunHistoryDAO) Delete(d *dto.DeleteRunHistoryDTO) error {
	if utils.IsZero(d.Namespace) || utils.IsZero(d.ID) {
		return fmt.Errorf("delete runHistory `app` + `server` + `ID` must not be empty, "+
			"namespace:[%s] ID:[%d]",
			d.Namespace, d.ID)
	}

	p := convertor.RunHistoryConvertor.ConvertDeleteDTOToPO(d)
	if err := dao.db.Unscoped().Where(p).Delete(p).Error; err != nil {
		log.Errorf("Failed to delete function run history, caused by %s", err)
		return err
	}

	return nil
}

// BatchDelete 批量删除
func (dao *RunHistoryDAO) BatchDelete(d *dto.BatchDeleteRunHistoryDTO) error {
	if utils.IsZero(d.Namespace) && len(d.IDs) == 0 {
		return fmt.Errorf("batch delete run history `namespace` + `ID` must not "+
			"be empty at the same time, namespace:[%s] IDs:[%d]",
			d.Namespace, d.IDs)
	}

	db := dao.db.Unscoped()
	if len(d.IDs) > 0 {
		db.Where("id in ?", d.IDs)
	}

	p := convertor.RunHistoryConvertor.ConvertBatchDeleteDTOToPO(d)
	if err := db.Delete(p).Error; err != nil {
		return fmt.Errorf("batch delete function run history err: %w", err)
	}

	return nil
}

// DeleteByCreateTime 根据创建时间批量删除
func (dao *RunHistoryDAO) DeleteByCreateTime(t time.Time) error {
	var ids []uint
	var history po.RunHistoryPO
	if err := dao.db.Model(&history).Select("id").Where("created_at < ?", t).Find(&ids).Error; err != nil {
		return err
	}
	if err := dao.db.Unscoped().Where("id in ?", ids).Delete(&history).Error; err != nil {
		log.Errorf("Failed to delete function run history, caused by %s", err)
		return err
	}

	return nil
}

// Update 更新信息
func (dao *RunHistoryDAO) Update(d *dto.UpdateRunHistoryDTO) error {
	if utils.IsZero(d.ID) {
		return fmt.Errorf("update runHistory `id` must not be empty, id:[%d]",
			d.ID)
	}

	p := convertor.RunHistoryConvertor.ConvertUpdateDTOToPO(d)
	if err := dao.db.Where("ID=? ", d.ID).Updates(p).Error; err != nil {
		log.Errorf("Failed to update runHistory, caused by %s", err)
		return err
	}

	return nil
}

// PageQuery 分页查询函数的调用记录
func (dao *RunHistoryDAO) PageQuery(d *dto.PageQueryRunHistoryDTO) ([]*po.RunHistoryPO, error) {
	var functions []*po.RunHistoryPO

	db := dao.db.Model(&po.RunHistoryPO{})
	if !d.CreatedAt.IsZero() {
		db.Where("created_at < ?", d.CreatedAt)
	}

	if len(d.IDs) > 0 {
		db.Where("id in ?", d.IDs)
	}

	if d.MaxID != 0 {
		db.Where("id <= ?", d.MaxID)
	}

	p := convertor.RunHistoryConvertor.ConvertPageQueryDTOToPO(d)
	if err := db.Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).
		Find(&functions).Error; err != nil {
		return nil, fmt.Errorf("page query function run history err: %w", err)
	}

	return functions, nil
}

// Count 根据条件获取总数
func (dao *RunHistoryDAO) Count(d *dto.PageQueryRunHistoryDTO) (int64, error) {
	var totalCount int64
	db := dao.db.Model(&po.RunHistoryPO{})
	if !d.CreatedAt.IsZero() {
		db.Where("created_at < ?", d.CreatedAt)
	}

	if len(d.IDs) > 0 {
		db.Where("id in ?", d.IDs)
	}

	if d.MaxID != 0 {
		db.Where("id <= ?", d.MaxID)
	}

	p := convertor.RunHistoryConvertor.ConvertPageQueryDTOToPO(d)
	if err := db.Where(p).Group("id").Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to get function count, caused by %s", err)
		return 0, err
	}

	return totalCount, nil
}
