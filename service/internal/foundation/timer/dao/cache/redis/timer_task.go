package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/redis"
)

const (
	// saveTaskTable 保存定时任务的表名
	saveTaskTable = "task_save"
	// 因为 pending 表的时间范围为 10 min，所以将分布式锁过期设置为 15 min
	pendingTableLockExpireDuration time.Duration = 15 * time.Minute
	// pending 表的过期时间设置为 7 d.
	pendingTableExpireDuration time.Duration = 7 * 24 * time.Hour
)

// TimerTaskDAO 定时器任务
type TimerTaskDAO struct {
	redisClient *redis.Client
	pool        concurrency.WorkerPool
}

// NewTimerTaskClient 新建定时器任务客户端
func NewTimerTaskClient(db *redis.Client, workerPool *concurrency.GoWorkerPool) *TimerTaskDAO {
	return &TimerTaskDAO{
		redisClient: db,
		pool:        workerPool,
	}
}

// GetTaskTableName 获取任务表名
func (w *TimerTaskDAO) GetTaskTableName(bucketID, timeSlice string) string {
	return fmt.Sprintf("%s_%s", bucketID, timeSlice)
}

// AddTimerTask 增加定时器任务到时间切片中
// TODO(@weixxxu): 当前在 dao 层下沉了过于复杂的逻辑，后续将 dao 层方法逻辑优化得简单一些
// 复合逻辑尽量在 repo 或 service 聚合
func (w *TimerTaskDAO) AddTimerTask(d *dto.AddTimerTaskDTO) error {
	// 获取所属时间分片
	timeSlice := d.TimerTime.Format(dto.TimerTaskTimeFormat)
	tableName := w.GetTaskTableName(d.BucketID, timeSlice)
	log.Infof("AddTimerTask TableName:%v Time:%v, HashID:%v", tableName, d.TimerTime.Unix(), d.HashID)
	// 将 timer 的 DefID 放入对应的 tableName score 就是其真正的触发时间 秒级为单位
	if err := w.redisClient.ZSet(context.Background(), tableName, d.TimerTime.Unix(), d.HashID); err != nil {
		return err
	}
	// 定时器任务添加成功，接下来需要在 save 表中存储定时器对应的下一次执行时间信息，以及根据定时器执行时间，在 pending 表中进行打点.
	w.pool.Submit(func() {
		if err := w.SaveTimerTask(d, tableName); err != nil {
			log.Errorf("Failed to SaveTimerTask error, caused by %v", err)
		}
		pendingTableName := getPendingTaskTableName(d.TimerTime)
		if err := w.redisClient.HSet(context.Background(),
			pendingTableName, getPendingTaskKey(d.HashID, d.TimerTime), ""); err != nil {
			log.Errorf("failed to save pending timer task, timer: %+v error:%v", d, err)
			return
		}
		log.Infof("save pending timer successfully, tableName: %s,timer: %+v", pendingTableName, d)
		// 获取分布式锁，进行 pending 表的过期时间设置.
		lock := w.redisClient.GetDistributeLock("pending_lock_"+pendingTableName, pendingTableLockExpireDuration)
		if err := lock.Lock(); err != nil {
			return
		}
		if err := w.redisClient.Expire(context.Background(),
			pendingTableName, getPendingTableExpireSeconds(d.TimerTime)); err != nil {
			// 过期时间设置失败，则提前解开分布式锁，由其他协称完成过期时间设置
			lock.Unlock()
			log.Errorf("set expire time failed, pendingTable: %s, err: %v", pendingTableName, err)
		}
	})
	return nil
}

// getPendingTableExpireSeconds pending 表的过期时间点为 timer 执行时间往后顺延 7 天, 单位为秒.
func getPendingTableExpireSeconds(execTime time.Time) int64 {
	timeDiff := time.Until(execTime)
	if timeDiff < 0 {
		timeDiff = 0
	}
	// 单位由纳秒转为秒.
	return int64(timeDiff+pendingTableExpireDuration) / int64(time.Second)
}

// DelPendingTimerTask 删除待执行定时器记录.
func (w *TimerTaskDAO) DelPendingTimerTask(defID string, execTime time.Time) error {
	tableName := getPendingTaskTableName(execTime)
	return w.redisClient.HDel(context.Background(), tableName, getPendingTaskKey(defID, execTime))
}

// CountPendingTimers 统计未执行的定时器数量.
func (w *TimerTaskDAO) CountPendingTimers(execTime time.Time) (int, error) {
	return w.redisClient.HLen(context.Background(), getPendingTaskTableName(execTime))
}

// getPendingTaskTableName 获取待执行定时器表名.
func getPendingTaskTableName(execTime time.Time) string {
	// 表名规则为，根据 time 获取所处的 10 分钟区间范围
	// 以 2022-05-26 09:21 为例，则表名为 pending_2022-05-26 09:2_2022-05-26 09:3
	begin := execTime.Format(dto.TimerTaskTimeFormat)
	end := execTime.Add(time.Minute * 10).Format(dto.TimerTaskTimeFormat)
	return fmt.Sprintf("pending_%s_%s", begin[:len(begin)-1], end[:len(end)-1])
}
func getPendingTaskKey(defID string, execTime time.Time) string {
	return fmt.Sprintf("%s_%d", defID, execTime.UnixNano())
}

// GetTimerTasks 获取时间切片中的定时器任务
func (w *TimerTaskDAO) GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error) {
	return w.getBucketTimers(d.BucketTime, d.StartTime.Unix(), d.EndTime.Unix())
}

// DelTimerTask 删除TimerTask
func (w *TimerTaskDAO) DelTimerTask(d *dto.DelTimerTaskDTO) error {
	return w.redisClient.ZRem(context.Background(), d.BucketTime, d.HashID)
}

// SaveTimerTask 保存时间任务 会替换成最新的任务记录
func (w *TimerTaskDAO) SaveTimerTask(d *dto.AddTimerTaskDTO, tableName string) error {
	saveTimerTaskDTO := &dto.SaveTimerTaskDTO{
		BucketTimeID: tableName,
		UnixTime:     d.TimerTime.UnixNano(),
		TriggerTime:  d.TimerTime.Format(dto.TimerTriggerTimeFormat),
	}
	jsonValue, err := json.Marshal(saveTimerTaskDTO)
	if err != nil {
		return err
	}
	return w.redisClient.HSet(context.Background(), saveTaskTable, d.HashID, string(jsonValue))
}

// GetSaveTimerTask 获取保存定时器任务.
func (w *TimerTaskDAO) GetSaveTimerTask(defID string) (*dto.SaveTimerTaskDTO, error) {
	taskDTOInfo, err := w.redisClient.HGet(context.Background(), saveTaskTable, defID)
	if err != nil {
		return nil, err
	}
	var taskDTO dto.SaveTimerTaskDTO
	if err = json.Unmarshal([]byte(taskDTOInfo), &taskDTO); err != nil {
		return nil, err
	}
	return &taskDTO, nil
}

// DeleteSaveTimerTask 删除保存定时器任务.
func (w *TimerTaskDAO) DeleteSaveTimerTask(defID string) error {
	return w.redisClient.HDel(context.Background(), saveTaskTable, defID)
}

// GetNotTriggeredTimers 获取未触发的定时器列表
func (w *TimerTaskDAO) GetNotTriggeredTimers(bucketTime string) ([]string, error) {
	startTime, err := w.getBucketStartTimeUnix(bucketTime)
	if err != nil {
		return nil, err
	}
	endTime, err := w.getBucketEndTimeUnix(bucketTime)
	if err != nil {
		return nil, err
	}
	return w.getBucketTimers(bucketTime, startTime.Unix(), endTime.Unix())
}
func (w *TimerTaskDAO) getBucketTimers(bucketTime string, startTime, endTime int64) ([]string, error) {
	return w.redisClient.ZRange(context.Background(), bucketTime, startTime, endTime)
}
func (w *TimerTaskDAO) getBucketTimeString(bucketTime string) (string, error) {
	bucketTimeStrings := strings.Split(bucketTime, "_")
	if len(bucketTimeStrings) != 2 {
		return "", fmt.Errorf("failed to getBucketTimeTimeUnix,caused error by: Split readySet")
	}
	return bucketTimeStrings[1], nil
}
func (w *TimerTaskDAO) getBucketStartTimeUnix(bucketTime string) (time.Time, error) {
	timeString, err := w.getBucketTimeString(bucketTime)
	if err != nil {
		return time.Time{}, err
	}
	startTimeString := timeString + ":00" // 开始时间片
	return time.ParseInLocation(dto.TimerTriggerTimeFormat, startTimeString, time.Local)
}
func (w *TimerTaskDAO) getBucketEndTimeUnix(bucketTime string) (time.Time, error) {
	timeString, err := w.getBucketTimeString(bucketTime)
	if err != nil {
		return time.Time{}, err
	}
	endTimeString := timeString + ":59" // 结束时间片
	return time.ParseInLocation(dto.TimerTriggerTimeFormat, endTimeString, time.Local)
}
