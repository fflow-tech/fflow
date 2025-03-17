// Package command 领域 command 服务，提供增删改能力。
package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/command/validate"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/repo"

	"github.com/gorhill/cronexpr"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// TimerDefCommandService 定时器定义写服务
type TimerDefCommandService struct {
	timerDefRepo    ports.TimerDefRepository
	timerTaskRepo   ports.TimerTaskRepository
	pollingTaskRepo ports.PollingTaskRepository
	pool            concurrency.WorkerPool
}

// NewTimerDefCommandService 新建服务
func NewTimerDefCommandService(repoProviderSet *ports.RepoProviderSet,
	pool *concurrency.GoWorkerPool) *TimerDefCommandService {
	return &TimerDefCommandService{
		timerDefRepo:    repoProviderSet.TimerDefRepo(),
		timerTaskRepo:   repoProviderSet.TimerTaskRepo(),
		pollingTaskRepo: repoProviderSet.PollingTaskRepo(),
		pool:            pool,
	}
}

// CreateTimerDef 创建定时器定义.
func (m *TimerDefCommandService) CreateTimerDef(req *dto.CreateTimerDefDTO) (uint64, error) {
	// 参数校验
	if err := validate.CheckCreateDefParam(req); err != nil {
		return 0, err
	}

	return m.timerDefRepo.CreateTimerDef(req)
}

// DeleteTimerDef 删除定时器.
func (m *TimerDefCommandService) DeleteTimerDef(req *dto.DeleteTimerDefDTO) error {
	if req.DefID == "" {
		if !req.HasAppAndName() {
			return errors.New("timer app and name must not be empty when defID is empty")
		}

		timerDef, err := m.timerDefRepo.GetTimerDefByAppName(req.App, req.Name)
		// 由于要保证定时器删除幂等性，即重复删除同一个定时器仍返回删除正常响应.
		// 因此当此处的错误为根据应用和名称未找到定时器时，不返回错误.
		if errors.Is(err, repo.ErrGetTimerByAppNameNotFound) {
			log.Warnf("delete timer failed, timer app: %s, name: %s, err: %v", req.App, req.Name, err)
			return nil
		}
		if err != nil {
			return err
		}
		req.DefID = timerDef.DefID
	}

	if err := m.timerDefRepo.DeleteTimerDef(req); err != nil {
		return err
	}
	// 删除定时器时需要对 pending 表清点.
	m.pool.Submit(func() {
		if saveTask, err := m.timerTaskRepo.GetSaveTimerTask(req.DefID); err != nil {
			log.Errorf("get save timer task failed, defID: %s, err: %v", req.DefID, err)
		} else {
			if err := m.timerTaskRepo.DelPendingTimerTask(req.DefID, time.Unix(0, saveTask.UnixTime)); err != nil {
				log.Errorf("delete pending timer task failed, defID: %s, err: %v", req.DefID, err)
			}
		}
	})
	return nil
}

// ChangeTimerStatus 修改定时器定义状态
func (m *TimerDefCommandService) ChangeTimerStatus(req *dto.ChangeTimerStatusDTO) error {
	log.Infof("ChangeTimerStatus status begin defID: %s", req.DefID)
	d := &dto.GetTimerDefDTO{
		DefID: req.DefID,
	}
	timerDef, err := m.timerDefRepo.GetTimerDef(d)
	if err != nil {
		return err
	}
	// 判断状态是否没有更改
	if timerDef.Status.ToInt() == req.Status {
		log.Infof("ChangeTimerStatus status not change defID: %s", timerDef.DefID)
		return nil
	}
	// 执行更改状态
	if err := m.timerDefRepo.ChangeTimerStatus(req); err != nil {
		return err
	}

	// 激活动作.
	if req.Status == entity.Enabled.ToInt() {
		return m.registerTimerTask(req, timerDef)
	}

	// 去激活动作时需要对 pending 表进行清点.
	if req.Status == entity.Disabled.ToInt() {
		m.pool.Submit(func() {
			if saveTask, err := m.timerTaskRepo.GetSaveTimerTask(req.DefID); err != nil {
				log.Errorf("get save timer task failed, defID: %s, err: %v", req.DefID, err)
			} else {
				if err := m.timerTaskRepo.DelPendingTimerTask(req.DefID, time.Unix(0, saveTask.UnixTime)); err != nil {
					log.Errorf("delete pending timer task failed, defID: %s, err: %v", req.DefID, err)
				}
			}
		})
	}
	return nil
}

func (m *TimerDefCommandService) registerTimerTask(req *dto.ChangeTimerStatusDTO, timerDef *entity.TimerDef) error {
	if req.Status != entity.Enabled.ToInt() {
		log.Infof("registerTimerTask defID: %s Status err", timerDef.DefID)
		return nil
	}
	// 获取当前任务 当已经存在并且触发时间比当前时间要大 那么这个任务就不注册
	if saveTask, err := m.timerTaskRepo.GetSaveTimerTask(timerDef.DefID); err == nil {
		if saveTask.UnixTime > time.Now().UnixNano() {
			log.Infof("registerTimerTask task timer is ready, defID: %s", timerDef.DefID)
			return nil
		}
	}
	// 激活时注册定时任务
	addTimerTask, err := m.convertorToAddTimerTaskDTO(timerDef)
	if err != nil {
		return err
	}

	return m.timerTaskRepo.AddTimerTask(addTimerTask)
}

// getNextTimeout 获取下次触发时间
func (m *TimerDefCommandService) getNextTimeout(timerDef *entity.TimerDef) (time.Time, error) {
	switch timerDef.TimerType {
	case entity.CronTimer:
		expr, err := cronexpr.Parse(timerDef.Cron)
		if err != nil {
			return time.Time{}, err
		}
		return expr.Next(time.Now()), nil
	case entity.DelayTimer:
		return time.ParseInLocation(entity.DelayTimeFormat, timerDef.DelayTime, time.Local)
	default:
		return time.Time{}, fmt.Errorf("failed to getNextTimeout TimerType: %v", timerDef.TimerType)
	}
}

// convertorToAddTimerTaskDTO 转换成增加定时器任务DTO
func (m *TimerDefCommandService) convertorToAddTimerTaskDTO(timerDef *entity.TimerDef) (*dto.AddTimerTaskDTO, error) {
	addTimerDTO := &dto.AddTimerTaskDTO{
		HashID: timerDef.DefID,
	}
	var err error
	// 获取定时器定义的桶ID
	addTimerDTO.BucketID, err = m.pollingTaskRepo.GetTaskBucketID(timerDef.DefID)
	if err != nil {
		return nil, err
	}

	// 获取定时器下次触发时间
	addTimerDTO.TimerTime, err = m.getNextTimeout(timerDef)
	if err != nil {
		return nil, err
	}
	return addTimerDTO, nil
}
