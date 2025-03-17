package sql

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// RunHistoryDAO 定时器执行历史数据访问对象
type RunHistoryDAO struct {
	db *mysql.Client
}

// NewRunHistoryDAO RunHistoryDAO 数据访问对象构造函数
func NewRunHistoryDAO(db *mysql.Client) *RunHistoryDAO {
	return &RunHistoryDAO{db: db}
}

// Transaction 事务
func (dao *RunHistoryDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建定时器执行记录
func (dao *RunHistoryDAO) Create(d *dto.CreateRunHistoryDTO) (*po.RunHistoryPO, error) {
	if utils.IsZero(d.DefID) || utils.IsZero(d.RunTimer) {
		return nil, fmt.Errorf("create RunHistory `DefID` or `RunTimer` must not be empty, "+
			"DefID:[%s] RunTimer:[%s]", d.DefID, d.RunTimer)
	}
	p, err := convertor.TaskConvertor.ConvertCreateDTOToPO(d)
	if err != nil {
		return nil, err
	}

	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create runHistory, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Get 获取定时器执行历史
func (dao *RunHistoryDAO) Get(d *dto.GetRunHistoryDTO) (*po.RunHistoryPO, error) {
	if utils.IsZero(d.DefID) || utils.IsZero(d.RunTimer) {
		return nil, fmt.Errorf("get RunHistory `DefID` or `RunTimer` must not be empty, DefID:[%s] "+
			"RunTimer:[%s]", d.DefID, d.RunTimer)
	}
	p, err := convertor.TaskConvertor.ConvertGetDTOToPO(d)
	if err != nil {
		return nil, err
	}

	r := &po.RunHistoryPO{}
	if err := dao.db.Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get runHistory, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// Delete 删除定时器执行历史
func (dao *RunHistoryDAO) Delete(d *dto.DeleteRunHistoryDTO) error {
	if utils.IsZero(d.DefID) && utils.IsZero(d.RunTimer) {
		return fmt.Errorf("get RunHistory `DefID` or `RunTimer` must not be empty, DefID:[%s] RunTimer:[%s]",
			d.DefID, d.RunTimer)
	}

	p, err := convertor.TaskConvertor.ConvertDeleteDTOToPO(d)
	if err != nil {
		return err
	}

	if err := dao.db.Unscoped().Where(p).Delete(p).Error; err != nil {
		log.Errorf("Failed to delete runHistory, caused by %s", err)
		return err
	}

	return nil
}

// Update 更新定时器执行历史
func (dao *RunHistoryDAO) Update(d *dto.UpdateRunHistoryDTO) error {
	if utils.IsZero(d.DefID) && utils.IsZero(d.RunTimer) {
		return fmt.Errorf("update RunHistory `DefID` or `RunTimer` must not be empty, defID:[%s] runTimer:[%s]",
			d.DefID, d.RunTimer)
	}

	p, err := convertor.TaskConvertor.ConvertUpdateDTOToPO(d)
	if err != nil {
		return err
	}

	if err := dao.db.Where("def_id=? and run_timer=?", d.DefID, d.RunTimer).Updates(p).Error; err != nil {
		log.Errorf("Failed to update runHistory, caused by %s", err)
		return err
	}

	return nil
}

// PageQuery 分页查询定时器的执行记录
func (dao *RunHistoryDAO) PageQuery(d *dto.PageQueryRunHistoryDTO) ([]*po.RunHistoryPO, error) {
	var runHistoryPOS []*po.RunHistoryPO
	p := convertor.TaskConvertor.ConvertPageQueryDTOToPO(d)

	db := dao.db.Model(&po.RunHistoryPO{})
	if err := db.Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).
		Find(&runHistoryPOS).Error; err != nil {
		log.Errorf("Failed to page query runHistory, caused by %s", err)
		return nil, err
	}

	return runHistoryPOS, nil
}

// Count 根据条件获取定时器总数
func (dao *RunHistoryDAO) Count(d *dto.PageQueryRunHistoryDTO) (int64, error) {
	var totalCount int64
	db := dao.db.Model(&po.RunHistoryPO{})
	p := convertor.TaskConvertor.ConvertPageQueryDTOToPO(d)
	if err := db.Where(p).Group("id").Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to get runHistory count, caused by %s", err)
		return 0, err
	}

	return totalCount, nil
}

// DeleteByRunTime 根据执行时间批量删除
func (dao *RunHistoryDAO) DeleteByRunTime(t time.Time) error {
	var ids []uint
	var history po.RunHistoryPO
	if err := dao.db.Model(&history).Select("id").Where("run_timer < ?", t).Find(&ids).Error; err != nil {
		return err
	}
	if err := dao.db.Unscoped().Where("id in ?", ids).Delete(&history).Error; err != nil {
		log.Errorf("Failed to delete history, caused by %s", err)
		return err
	}

	return nil
}
