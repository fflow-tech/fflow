package query

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// TimerTaskQueryService 定时器任务查询服务
type TimerTaskQueryService struct {
	timerTaskRepo   ports.TimerTaskRepository
	pollingTaskRepo ports.PollingTaskRepository
}

// NewTimerTaskQueryService 新建查询服务
func NewTimerTaskQueryService(r *ports.RepoProviderSet) *TimerTaskQueryService {
	return &TimerTaskQueryService{timerTaskRepo: r.TimerTaskRepo(), pollingTaskRepo: r.PollingTaskRepo()}
}

// GetTimerTasks 获取定时器任务
func (m *TimerTaskQueryService) GetTimerTasks(d *dto.GetTimerTaskDTO) ([]string, error) {
	return m.timerTaskRepo.GetTimerTasks(d)
}

// PageQueryHistory 分页获取执行历史
func (m *TimerTaskQueryService) PageQueryHistory(d *dto.PageQueryRunHistoryDTO) ([]*dto.GetRunHistoryRspDTO, int64,
	error) {
	runHistory, total, err := m.timerTaskRepo.PageQueryHistory(d)
	if err != nil {
		return nil, 0, err
	}
	runHistoryRspList, err := convertor.RunHistoryConvertor.ConvertEntitiesToDTOs(runHistory)
	if err != nil {
		return nil, 0, err
	}
	return runHistoryRspList, total, err
}

// GetTimeLimitTimers 获取时间范围的定时器列表
func (m *TimerTaskQueryService) GetTimeLimitTimers(startTime, endTime string) ([]string, error) {
	bucketNum := m.pollingTaskRepo.GetBucketNum()
	bucketTimes, err := m.getTimeLimitBuckets(startTime, endTime)
	if err != nil {
		return nil, err
	}
	var allTimers []string
	for i := 0; i < bucketNum; i++ {
		for _, bucketTime := range bucketTimes {
			bucketName := m.timerTaskRepo.GetTaskTableName(utils.UintToStr(uint(i)), bucketTime)
			timers, err := m.timerTaskRepo.GetNotTriggeredTimers(bucketName)
			if err != nil {
				return nil, err
			}
			allTimers = append(allTimers, timers...)
		}
	}
	return allTimers, nil
}

func (m *TimerTaskQueryService) getTimeLimitBuckets(startTimeString, endTimeString string) ([]string, error) {
	startTime, err := time.ParseInLocation(dto.TimerTaskTimeFormat, startTimeString, time.Local)
	if err != nil {
		return nil, err
	}

	endTime, err := time.ParseInLocation(dto.TimerTaskTimeFormat, endTimeString, time.Local)
	if err != nil {
		return nil, err
	}

	var timerList []string
	for startTime.Before(endTime) {
		timerList = append(timerList, startTime.Format(dto.TimerTaskTimeFormat))
		startTime = startTime.Add(time.Minute)
	}
	return timerList, nil
}

// CountPendingTimers 统计未执行的定时器数量.
func (m *TimerTaskQueryService) CountPendingTimers(curTime time.Time) (int, error) {
	return m.timerTaskRepo.CountPendingTimers(curTime)
}
