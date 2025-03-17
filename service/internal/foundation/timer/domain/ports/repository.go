package ports

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/remote"
)

// TimerDefRepository 仓储层接口
type TimerDefRepository interface {
	CreateTimerDef(d *dto.CreateTimerDefDTO) (uint64, error)
	GetTimerDef(d *dto.GetTimerDefDTO) (*entity.TimerDef, error)
	DeleteTimerDef(d *dto.DeleteTimerDefDTO) error
	ChangeTimerStatus(d *dto.ChangeTimerStatusDTO) error
	GetTimerDefList(d *dto.PageQueryTimeDefDTO) ([]*entity.TimerDef, int64, error)
	CountTimersByStatus(status entity.TimerDefStatus) (int64, error)
	GetTimerDefByAppName(app, name string) (*entity.TimerDef, error)
}

// TimerTaskRepository 仓储层接口
type TimerTaskRepository interface {
	AddTimerTask(d *dto.AddTimerTaskDTO) error
	GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error)
	DelTimerTask(d *dto.DelTimerTaskDTO) error
	CreateHistory(d *dto.CreateRunHistoryDTO) (*entity.RunHistory, error)
	UpdateHistory(d *dto.UpdateRunHistoryDTO) error
	PageQueryHistory(d *dto.PageQueryRunHistoryDTO) ([]*entity.RunHistory, int64, error)
	GetNotTriggeredTimers(bucketTime string) ([]string, error)
	GetTaskTableName(bucketID, timeSlice string) string
	DeleteRunHistories() error
	GetSaveTimerTask(defID string) (*dto.SaveTimerTaskDTO, error)
	DeleteSaveTimerTask(defID string) error
	DelPendingTimerTask(defID string, curTime time.Time) error
	CountPendingTimers(curTime time.Time) (int, error)
}

// PollingTaskRepository 仓储层接口
type PollingTaskRepository interface {
	GetBucketNum() int
	GetTaskBucketID(hashID string) (string, error)
	SetBucketNum(num int) error
	SetTimeSlice(timeDuration string) error
	GetTimeSlice(timeDuration string) error
	SuccessTimeSlice(timeDuration string) error
}

// EventBusRepository 仓储层接口
type EventBusRepository interface {
	SendPollingEvent(ctx context.Context, msg interface{}) error
	NewPollingConsumer(ctx context.Context, group string,
		handle func(context.Context, interface{}) error) (mq.Consumer, error)
	SendTimerTaskEvent(ctx context.Context, msg interface{}) error
	NewTimerTaskConsumer(ctx context.Context, group string,
		handle func(context.Context, interface{}) error) (mq.Consumer, error)
}

// RemoteRepository 仓储层接口
type RemoteRepository interface {
	CallFAAS(ctx context.Context, req *remote.CallFAASReqDTO) (map[string]interface{}, error)
	CallHTTP(ctx context.Context, req *remote.CallHTTPReqDTO) (map[string]interface{}, error)
	SendMsgToUser(userID string, msg string) error
}

// AppRepository 仓储层接口
type AppRepository interface {
	GetAppList(d *dto.PageQueryAppDTO) ([]*entity.App, int64, error)
	CreateApp(d *dto.CreateAppDTO) (*entity.App, error)
	DeleteApp(d *dto.DeleteAppDTO) error
}
