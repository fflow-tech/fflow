package command

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// PollingTaskCommandService 轮询服务写服务
type PollingTaskCommandService struct {
	pollingTaskRepo ports.PollingTaskRepository
	timerTaskRepo   ports.TimerTaskRepository
	eventRepo       ports.EventBusRepository
}

// NewPollingTaskCommandService 新建服务
func NewPollingTaskCommandService(repoProviderSet *ports.RepoProviderSet) *PollingTaskCommandService {
	return &PollingTaskCommandService{
		pollingTaskRepo: repoProviderSet.PollingTaskRepo(),
		timerTaskRepo:   repoProviderSet.TimerTaskRepo(),
		eventRepo:       repoProviderSet.EventBusRepo(),
	}
}

// GetPollingTaskWorkLock 获取轮询时间片的锁
func (p *PollingTaskCommandService) GetPollingTaskWorkLock() (string, error) {
	timeNow := time.Now().Format(dto.TimerTaskTimeFormat)
	if err := p.pollingTaskRepo.GetTimeSlice(timeNow); err != nil {
		return "", err
	}
	return timeNow, nil
}

// SendPollingTaskWork 发送时间分片任务
func (p *PollingTaskCommandService) SendPollingTaskWork(timeSlice string) error {
	// 获取所有桶的数值
	bucketNum := p.pollingTaskRepo.GetBucketNum()
	// 拼装发送时间分片到消息队列 这里使用异步的方式来投递
	for i := 0; i < bucketNum; i++ {
		// 这里当失败时 可以考虑如何来做处理 避免消息重发或漏发
		if err := p.sendPollingTask(fmt.Sprintf("%d_%s", i, timeSlice)); err != nil {
			return err
		}
		log.Infof(" GoroutineID:%d  sendTimeSlice success timeAt:%v", utils.GetCurrentGoroutineID(), time.Now())
	}
	// 标记事件片任务完成
	return p.pollingTaskRepo.SuccessTimeSlice(timeSlice)
}

// sendTimeSlice 发送时间切片 这里需要集中考虑事件定义和内容设计
func (p *PollingTaskCommandService) sendPollingTask(timeSlice string) error {
	log.Infof(" GoroutineID:%d sendPollingTask timeSlice %v", utils.GetCurrentGoroutineID(), timeSlice)
	return p.eventRepo.SendPollingEvent(context.Background(), timeSlice)
}

// GetTaskBucketID 获取任务所属的桶ID
func (p *PollingTaskCommandService) GetTaskBucketID(hashID string) (string, error) {
	return p.pollingTaskRepo.GetTaskBucketID(hashID)
}
