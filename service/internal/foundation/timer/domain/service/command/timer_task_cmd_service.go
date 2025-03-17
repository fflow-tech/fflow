package command

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// TimerTaskCommandService 定时器任务写服务
type TimerTaskCommandService struct {
	timerTaskRepo ports.TimerTaskRepository
	eventRepo     ports.EventBusRepository
}

// NewTimerTaskCommandService 新建服务
func NewTimerTaskCommandService(repoProviderSet *ports.RepoProviderSet) *TimerTaskCommandService {
	return &TimerTaskCommandService{
		timerTaskRepo: repoProviderSet.TimerTaskRepo(),
		eventRepo:     repoProviderSet.EventBusRepo(),
	}
}

// AddTimerTask 创建定时器任务
func (m *TimerTaskCommandService) AddTimerTask(d *dto.AddTimerTaskDTO) error {
	return m.timerTaskRepo.AddTimerTask(d)
}

// DelTimerTask 删除定时器任务
func (m *TimerTaskCommandService) DelTimerTask(d *dto.DelTimerTaskDTO) error {
	return m.timerTaskRepo.DelTimerTask(d)
}

// GetReadyTimer 获取到期的定时器
func (m *TimerTaskCommandService) GetReadyTimer(readySet string, startTime time.Time) error {
	// 找到readySet 对应的zset 获取当前的两秒的时间
	endTime := startTime.Add(time.Duration(config.GetTimerTaskConfig().WorkDelayMillisecond) * time.Millisecond) // 增加延迟范围
	getTaskDTO := &dto.GetTimerTaskDTO{
		BucketTime: readySet,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	hashIDs, err := m.timerTaskRepo.GetTimerTasks(getTaskDTO)
	if err != nil {
		log.Errorf("GoroutineID:%d GetReadyTimer getWorkTicker%v", utils.GetCurrentGoroutineID(), err)
		return err
	}

	for _, hashID := range hashIDs {
		log.Infof("ReadySet %v GoroutineID:%d range hashID %v", readySet, utils.GetCurrentGoroutineID(), hashID)
		if err := m.SendReadyTimer(hashID); err != nil {
			return err
		}
		// 发送完成之后就删除当前时间片里面的定时任务
		delTimerTaskDTO := &dto.DelTimerTaskDTO{
			BucketTime: readySet,
			HashID:     hashID,
		}
		if err := m.DelTimerTask(delTimerTaskDTO); err != nil {
			log.Errorf("failed to GetReadyTimer DelTimerTask, caused by %v", err)
			// 这里删除失败时 继续往下执行
		}
	}
	return nil
}

// SendReadyTimer 发送到期定时器到 readyTable
func (m *TimerTaskCommandService) SendReadyTimer(hashID string) error {
	//发送DefID对应的消息到readyTable
	log.Infof("GoroutineID:%d SendReadyTimer hashID %v", utils.GetCurrentGoroutineID(), hashID)
	return m.eventRepo.SendTimerTaskEvent(context.Background(), hashID)
}
