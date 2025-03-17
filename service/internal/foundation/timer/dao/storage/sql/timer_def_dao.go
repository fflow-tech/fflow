// Package sql 数据库交互层，负责 MySQL 交互完成基本增删改查
package sql

import (
	"fmt"
	"gorm.io/gorm"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// ErrRecordNotFound 数据库记录未找到.
var ErrRecordNotFound = gorm.ErrRecordNotFound

// TimerDefDAO TimerDef数据访问对象
type TimerDefDAO struct {
	db *mysql.Client
}

// NewTimerDefDAO NewTimerDef数据访问对象构造函数
func NewTimerDefDAO(db *mysql.Client) *TimerDefDAO {
	return &TimerDefDAO{db: db}
}

// Transaction 事务
func (dao *TimerDefDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建定时器定义
func (dao *TimerDefDAO) Create(def *dto.CreateTimerDefDTO) (*po.TimerDefPO, error) {
	p, err := convertor.DefConvertor.ConvertCreateDTOToPO(def)
	if err != nil {
		return nil, err
	}

	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create timer def, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Delete 删除定时器定义
func (dao *TimerDefDAO) Delete(d *dto.DeleteTimerDefDTO) error {
	if utils.IsZero(d.DefID) {
		return fmt.Errorf("delete def `DefID` must not be zero, DefID:[%s]", d.DefID)
	}
	p, err := convertor.DefConvertor.ConvertDeleteDTOToPO(d)
	if err != nil {
		return err
	}

	if err := dao.db.Unscoped().Model(&po.TimerDefPO{}).Delete(p).Error; err != nil {
		log.Errorf("Failed to delete timer def, caused by %s", err)
		return err
	}
	return nil
}

// PageQueryTimeList 分页查询定时器.
func (dao *TimerDefDAO) PageQueryTimeList(d *dto.PageQueryTimeDefDTO) ([]*po.TimerDefPO, error) {
	db := dao.db.Model(&po.TimerDefPO{})
	if !utils.IsZero(d.Name) {
		db = db.Where("name like ?", "%"+d.Name+"%")
	}

	if !utils.IsZero(d.Creator) {
		db = db.Where("creator like ?", "%"+d.Creator+"%")
	}

	var timerPOs []*po.TimerDefPO
	if err := db.Where("app = ?", d.App).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).
		Find(&timerPOs).Error; err != nil {
		return nil, err
	}

	return timerPOs, nil
}

// GetTimerDefByAppName 根据应用和名称获取定时器.
func (dao *TimerDefDAO) GetTimerDefByAppName(app string, name string) (*po.TimerDefPO, error) {
	var timerDef po.TimerDefPO
	return &timerDef, dao.db.Model(&po.TimerDefPO{}).Where("app = ? AND name = ?", app, name).First(&timerDef).Error
}

// Count 查询定时器总数
func (dao *TimerDefDAO) Count(d *dto.CountTimerDefDTO) (int64, error) {
	db := dao.db.Model(&po.TimerDefPO{})
	if !utils.IsZero(d.Name) {
		db = db.Where("name like ?", "%"+d.Name+"%")
	}

	if !utils.IsZero(d.Creator) {
		db = db.Where("creator like ?", "%"+d.Creator+"%")
	}

	var total int64
	if err := db.Where("app = ?", d.App).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// CountByStatus 根据状态统计定时器总数.
func (dao *TimerDefDAO) CountByStatus(status int) (int64, error) {
	db := dao.db.Model(&po.TimerDefPO{})
	var total int64
	if err := db.Where("status = ?", status).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// UpdateStatus 更新定时器状态
func (dao *TimerDefDAO) UpdateStatus(d *dto.UpdateTimerDefDTO) error {
	p, err := convertor.DefConvertor.ConvertUpdateDTOToPO(d)
	if err != nil {
		return err
	}

	return dao.db.Where("id = ?", d.DefID).Updates(p).Error
}
