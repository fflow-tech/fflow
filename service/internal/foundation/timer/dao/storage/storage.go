// Package storage 存储层接口定义
package storage

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
)

// Transaction 执行事务
type Transaction interface {
	Transaction(fun func(*mysql.Client) error) error
}

// TimerDefDAO 定时器定义存储层接口
type TimerDefDAO interface {
	Transaction
	Create(d *dto.CreateTimerDefDTO) (*po.TimerDefPO, error)
	Delete(d *dto.DeleteTimerDefDTO) error
	PageQueryTimeList(d *dto.PageQueryTimeDefDTO) ([]*po.TimerDefPO, error)
	Count(d *dto.CountTimerDefDTO) (int64, error)
	UpdateStatus(d *dto.UpdateTimerDefDTO) error
	CountByStatus(status int) (int64, error)
	GetTimerDefByAppName(app string, name string) (*po.TimerDefPO, error)
}

// PollingTaskDAO 轮询任务存储层接口
type PollingTaskDAO interface {
	GetBucketNum() int
	GetTaskBucketID(hashID string) (string, error)
	SetBucketNum(num int) error
	SetTimeSlice(timeDuration string) error
	GetTimeSlice(timeDuration string) error
	SuccessTimeSlice(timeDuration string) error
}

// TimerDefRedisDAO 定时器定义 redis 存储层接口
type TimerDefRedisDAO interface {
	AddTimerDef(def *dto.CreateTimerDefDTO) error
	GetTimerDef(d *dto.GetTimerDefDTO) (*po.TimerDefPO, error)
	DelTimerDef(d *dto.DeleteTimerDefDTO) error
	ChangeTimerStatus(d *dto.ChangeTimerStatusDTO) error
}

// TimerTaskDAO 定时器任务存储层接口
type TimerTaskDAO interface {
	AddTimerTask(d *dto.AddTimerTaskDTO) error
	GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error)
	DelTimerTask(d *dto.DelTimerTaskDTO) error
	GetNotTriggeredTimers(bucketTime string) ([]string, error)
	GetTaskTableName(bucketID, timeSlice string) string
	GetSaveTimerTask(defID string) (*dto.SaveTimerTaskDTO, error)
	DeleteSaveTimerTask(defID string) error
	DelPendingTimerTask(defID string, curTime time.Time) error
	CountPendingTimers(curTime time.Time) (int, error)
}

// TimerTaskRunHistoryDAO 定时器任务执行历史存储层接口
type TimerTaskRunHistoryDAO interface {
	Create(d *dto.CreateRunHistoryDTO) (*po.RunHistoryPO, error)
	Get(d *dto.GetRunHistoryDTO) (*po.RunHistoryPO, error)
	Delete(d *dto.DeleteRunHistoryDTO) error
	Update(d *dto.UpdateRunHistoryDTO) error
	PageQuery(d *dto.PageQueryRunHistoryDTO) ([]*po.RunHistoryPO, error)
	Count(d *dto.PageQueryRunHistoryDTO) (int64, error)
	DeleteByRunTime(t time.Time) error
}

// AppDAO 应用 DAO 层接口
type AppDAO interface {
	Create(d *dto.CreateAppDTO) (*po.App, error)
	Get(d *dto.GetAppDTO) (*po.App, error)
	PageQuery(d *dto.PageQueryAppDTO) ([]*po.App, error)
	Delete(d *dto.DeleteAppDTO) error
	Count(d *dto.CountAppDTO) (int64, error)
}
