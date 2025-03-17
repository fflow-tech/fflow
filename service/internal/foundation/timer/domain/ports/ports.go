// Package ports 领域接口定义，包括 command&query、 仓储层接口。
package ports

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

// CommandPorts 写入接口
type CommandPorts interface {
	TimerDefCommandPorts
	TimerTaskCommandPorts
	PollingTaskCommandPorts
	NotifyCommandPorts
	AppCommandPorts
}

// QueryPorts 查询接口
type QueryPorts interface {
	TimerDefQueryPorts
	TimerTaskQueryPorts
	AppQueryPorts
}

// TimerDefCommandPorts  定时器定义命令接口
type TimerDefCommandPorts interface {
	CreateTimerDef(req *dto.CreateTimerDefDTO) (uint64, error)
	DeleteTimerDef(req *dto.DeleteTimerDefDTO) error
	ChangeTimerStatus(req *dto.ChangeTimerStatusDTO) error
}

// TimerDefQueryPorts 定时器定义查询接口
type TimerDefQueryPorts interface {
	GetTimerDef(d *dto.GetTimerDefDTO) (*dto.TimerDefDTO, error)
	GetTimerDefList(d *dto.PageQueryTimeDefDTO) ([]*dto.TimerDefDTO, int64, error)
	CountTimersByStatus(status entity.TimerDefStatus) (int64, error)
}

// TimerTaskCommandPorts 定时器任务命令接口
type TimerTaskCommandPorts interface {
	AddTimerTask(d *dto.AddTimerTaskDTO) error
	DelTimerTask(d *dto.DelTimerTaskDTO) error
	GetReadyTimer(readySet string, startTime time.Time) error
	DeleteRunHistories() error
}

// TimerTaskQueryPorts 定时器任务查询接口
type TimerTaskQueryPorts interface {
	GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error)
	PageQueryHistory(d *dto.PageQueryRunHistoryDTO) ([]*dto.GetRunHistoryRspDTO, int64, error)
	GetTimeLimitTimers(startTime, endTime string) ([]string, error)
	CountPendingTimers(curTime time.Time) (int, error)
}

// PollingTaskCommandPorts  轮询任务命令接口
type PollingTaskCommandPorts interface {
	GetPollingTaskWorkLock() (string, error)
	SendPollingTaskWork(timeSlice string) error
	GetTaskBucketID(hashID string) (string, error)
}

// NotifyCommandPorts 通知任务命令接口
type NotifyCommandPorts interface {
	SendNotify(hashID string) error
	TimerListSendNotify(defIDs []string) error
	ManualTriggerSend(hashID string) error
	ManualTriggerSendList(defIDs []string) error
}

// AppQueryPorts 应用查询接口
type AppQueryPorts interface {
	GetAppList(d *dto.PageQueryAppDTO) ([]*dto.App, int64, error)
}

// AppCommandPorts 应用操作接口
type AppCommandPorts interface {
	CreateApp(d *dto.CreateAppDTO) error
	DeleteApp(d *dto.DeleteAppDTO) error
}
