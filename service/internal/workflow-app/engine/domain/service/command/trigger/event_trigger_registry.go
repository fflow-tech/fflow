package trigger

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// EventTriggerRegistry 事件触发器注册中心
type EventTriggerRegistry interface {
	Register(d *dto.RegisterTriggerDTO) error
	Unregister(d *dto.UnregisterTriggerDTO) error
}

// DefaultEventTriggerRegistry 默认事件触发器注册中心实现
type DefaultEventTriggerRegistry struct {
	triggerRepo ports.TriggerRepository
}

// NewDefaultEventTriggerRegistry 构造函数
func NewDefaultEventTriggerRegistry(repoProvider *ports.RepoProviderSet) *DefaultEventTriggerRegistry {
	return &DefaultEventTriggerRegistry{triggerRepo: repoProvider.TriggerRepo()}
}

// Register 事件触发器注册逻辑实现
func (t *DefaultEventTriggerRegistry) Register(d *dto.RegisterTriggerDTO) error {
	actions := entity.GetAllAction(d.Actions)
	for _, action := range actions {
		if entity.GetTriggerLevel(action.ActionType) != d.Level {
			log.Warnf("Not %s level trigger, skip it", d.Level.String())
			continue
		}

		createTriggerDTO := &dto.CreateTriggerDTO{}
		createTriggerDTO.Trigger.BasicTriggerDef = d.TriggerDef.BasicTriggerDef
		createTriggerDTO.Type = entity.Event
		createTriggerDTO.Event = d.Event
		createTriggerDTO.Action = action
		createTriggerDTO.DefID = d.DefID
		createTriggerDTO.DefVersion = d.DefVersion
		createTriggerDTO.InstID = d.InstID
		createTriggerDTO.Level = d.Level
		createTriggerDTO.Status = entity.EnabledTrigger

		_, err := t.triggerRepo.Create(createTriggerDTO)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnRegister 事件触发器反注册逻辑实现
func (t *DefaultEventTriggerRegistry) Unregister(d *dto.UnregisterTriggerDTO) error {
	updateTriggerDTO := &dto.UpdateTriggerDTO{
		DefID:      d.DefID,
		DefVersion: d.DefVersion,
		InstID:     d.InstID,
		Level:      d.Level,
	}
	updateTriggerDTO.Status = entity.DisabledTrigger
	return t.triggerRepo.Update(updateTriggerDTO)
}
