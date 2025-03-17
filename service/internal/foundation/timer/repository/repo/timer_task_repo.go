package repo

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/cache/redis"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// defaultKeepDays 默认的执行历史保存时间
var defaultKeepDays = 7

// TimerTaskRepo 定时器任务实体
type TimerTaskRepo struct {
	timerTaskRepo storage.TimerTaskDAO
	runHistoryDAO storage.TimerTaskRunHistoryDAO
}

// NewTimerTaskRepo 实体构造函数
func NewTimerTaskRepo(d *redis.TimerTaskDAO, t *sql.RunHistoryDAO) *TimerTaskRepo {
	return &TimerTaskRepo{timerTaskRepo: d, runHistoryDAO: t}
}

// AddTimerTask 增加定时器任务
func (t *TimerTaskRepo) AddTimerTask(d *dto.AddTimerTaskDTO) error {
	return t.timerTaskRepo.AddTimerTask(d)
}

// GetTimerTasks 获取定时器任务
func (t *TimerTaskRepo) GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error) {
	return t.timerTaskRepo.GetTimerTasks(d)
}

// DelTimerTask 删除定时器任务
func (t *TimerTaskRepo) DelTimerTask(d *dto.DelTimerTaskDTO) error {
	return t.timerTaskRepo.DelTimerTask(d)
}

// CreateHistory 创建执行记录
func (t *TimerTaskRepo) CreateHistory(d *dto.CreateRunHistoryDTO) (*entity.RunHistory, error) {
	history, err := t.runHistoryDAO.Create(d)
	if err != nil {
		return nil, err
	}
	return convertor.TaskConvertor.ConvertPOToEntity(history)
}

// UpdateHistory 更新执行记录
func (t *TimerTaskRepo) UpdateHistory(d *dto.UpdateRunHistoryDTO) error {
	return t.runHistoryDAO.Update(d)
}

// PageQueryHistory 获取执行记录列表
func (t *TimerTaskRepo) PageQueryHistory(d *dto.PageQueryRunHistoryDTO) ([]*entity.RunHistory, int64, error) {
	if d.PageQuery == nil {
		d.PageQuery = constants.NewDefaultPageQuery()
	}

	historyPOs, err := t.runHistoryDAO.PageQuery(d)
	if err != nil {
		return nil, 0, err
	}

	historyEntities, err := convertor.TaskConvertor.ConvertHistoryPOsToEntities(historyPOs)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.runHistoryDAO.Count(d)
	if err != nil {
		return nil, 0, err
	}
	return historyEntities, total, nil
}

// GetNotTriggeredTimers 获取没有触发的定时器列表
func (t *TimerTaskRepo) GetNotTriggeredTimers(bucketTime string) ([]string, error) {
	return t.timerTaskRepo.GetNotTriggeredTimers(bucketTime)
}

// GetTaskTableName 获取任务表名
func (t *TimerTaskRepo) GetTaskTableName(bucketID, timeSlice string) string {
	return t.timerTaskRepo.GetTaskTableName(bucketID, timeSlice)
}

// DeleteRunHistories 删除过期执行记录
func (t *TimerTaskRepo) DeleteRunHistories() error {
	// 获取执行理事需要保存的时间
	keepDays := config.GetAppConfig().KeepDays
	// 配置出错的兜底，避免删除了比较新的执行历史
	if keepDays < defaultKeepDays {
		keepDays = defaultKeepDays
	}
	expireTime := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -keepDays)
	return t.runHistoryDAO.DeleteByRunTime(expireTime)
}

// GetSaveTimerTask 获取保存定时器任务.
func (t *TimerTaskRepo) GetSaveTimerTask(defID string) (*dto.SaveTimerTaskDTO, error) {
	return t.timerTaskRepo.GetSaveTimerTask(defID)
}

// DeleteSaveTimerTask 删除保存定时器任务.
func (t *TimerTaskRepo) DeleteSaveTimerTask(defID string) error {
	return t.timerTaskRepo.DeleteSaveTimerTask(defID)
}

// DelPendingTimerTask 删除待执行定时器记录.
func (w *TimerTaskRepo) DelPendingTimerTask(defID string, execTime time.Time) error {
	return w.timerTaskRepo.DelPendingTimerTask(defID, execTime)
}

// CountPendingTimers 统计未执行的定时器数量.
func (w *TimerTaskRepo) CountPendingTimers(execTime time.Time) (int, error) {
	return w.timerTaskRepo.CountPendingTimers(execTime)
}
