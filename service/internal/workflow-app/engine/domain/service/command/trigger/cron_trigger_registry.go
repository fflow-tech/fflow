// Package trigger 触发器实现
package trigger

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/remote"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// CronTriggerRegistry 定时触发器注册中心
type CronTriggerRegistry interface {
	Register(d *dto.RegisterTriggerDTO) error
	Unregister(d *dto.UnregisterTriggerDTO) error
}

// DefaultCronTriggerRegistry 默认定时触发器注册实现
type DefaultCronTriggerRegistry struct {
	triggerRepo  ports.TriggerRepository
	eventBusRepo ports.EventBusRepository
	remoteRepo   ports.RemoteRepository
}

// NewDefaultCronTriggerRegistry 初始化
func NewDefaultCronTriggerRegistry(repoProvider *ports.RepoProviderSet) *DefaultCronTriggerRegistry {
	return &DefaultCronTriggerRegistry{
		triggerRepo:  repoProvider.TriggerRepo(),
		eventBusRepo: repoProvider.EventBusRepo(),
		remoteRepo:   repoProvider.RemoteRepo(),
	}
}

// Register 定时触发器注册逻辑实现
func (t *DefaultCronTriggerRegistry) Register(d *dto.RegisterTriggerDTO) error {
	if err := t.validateCronTriggerIntervalTime(d.Expr); err != nil {
		return err
	}

	actions := entity.GetAllAction(d.Actions)
	for _, action := range actions {
		if entity.GetTriggerLevel(action.ActionType) != d.Level {
			log.Warnf("[%s]Action %s not %s level trigger, skip it",
				logs.GetFlowTraceID(d.DefID, d.InstID), action.ActionType, d.Level)
			continue
		}

		createTriggerDTO := &dto.CreateTriggerDTO{}
		createTriggerDTO.Trigger.BasicTriggerDef = d.TriggerDef.BasicTriggerDef
		createTriggerDTO.Type = entity.Timer
		createTriggerDTO.Action = action
		createTriggerDTO.Expr = d.Expr
		createTriggerDTO.DefID = d.DefID
		createTriggerDTO.DefVersion = d.DefVersion
		createTriggerDTO.InstID = d.InstID
		createTriggerDTO.Level = d.Level
		createTriggerDTO.Status = entity.EnabledTrigger

		triggerID, err := t.triggerRepo.Create(createTriggerDTO)
		if err != nil {
			return err
		}

		if err := t.createCronTask(d.Expr, triggerID, d.DefID, d.DefVersion); err != nil {
			return err
		}
	}
	return nil
}

// Unregister 定时触发器反注册逻辑实现
func (t *DefaultCronTriggerRegistry) Unregister(d *dto.UnregisterTriggerDTO) error {
	updateTriggerDTO := &dto.UpdateTriggerDTO{
		DefID:      d.DefID,
		DefVersion: d.DefVersion,
		InstID:     d.InstID,
		Level:      d.Level,
	}
	updateTriggerDTO.Status = entity.DisabledTrigger
	return t.triggerRepo.Update(updateTriggerDTO)
}

// validateCronTriggerIntervalTime  校验定时触发器间隔时间
func (t *DefaultCronTriggerRegistry) validateCronTriggerIntervalTime(cronExpr string) error {
	// 定时触发器触发间隔必须在一分钟和一个月之间
	intervalUnixTime, err := utils.GetIntervalTimeByCronExpr(cronExpr)
	if err != nil {
		return err
	}

	maxIntervalUnixTime, err := utils.GetIntervalTimeByCronExpr(constants.MaxIntervalTimeCronExpr)
	if err != nil {
		return err
	}
	if intervalUnixTime < constants.MinCronIntervalTime || intervalUnixTime > maxIntervalUnixTime {
		return fmt.Errorf("registry cron trigger interval time between 1min and 1month， "+
			"illegal cron expr:[%s]", cronExpr)
	}

	return nil
}

// createCronTask 创建定时任务
func (t *DefaultCronTriggerRegistry) createCronTask(expr string, triggerID, defID string, defVersion int) error {
	nowTime := time.Now()
	nextTime, err := utils.GetNextTimeByExpr(expr, nowTime)
	if err != nil {
		return err
	}

	// 消息间隔时间大于 10 天的情况下会使用定时器来实现
	if nextTime.Unix()-nowTime.Unix() >= constants.TdmqMaxCacheTime {
		addCronJobDTO := &remote.AddCronJobDTO{
			CronStr: expr,
			CronTriggerEvent: event.CronTriggerEvent{
				TriggerID:  triggerID,
				DefID:      defID,
				DefVersion: defVersion,
			},
		}

		return t.remoteRepo.AddCronJob(addCronJobDTO)
	}

	log.Infof("[%s]Send CronTriggerEvent, triggerID=%d, nextTime=%s", logs.GetFlowTraceID(defID, triggerID),
		triggerID, nextTime)
	cronTriggerEvent := &event.CronTriggerEvent{TriggerID: triggerID, DefID: defID}
	return t.eventBusRepo.SendCronPresetEvent(context.Background(), nextTime, cronTriggerEvent)
}
