package trigger

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// Registry 触发器注册中心
type Registry interface {
	Register(d *dto.RegisterTriggersDTO) error     // 注册
	UnRegister(d *dto.UnRegisterTriggersDTO) error // 去注册
}

// DefaultRegistry 默认触发器实现体
type DefaultRegistry struct {
	cronTriggerRegister  CronTriggerRegistry
	eventTriggerRegister EventTriggerRegistry
}

// NewDefaultRegistry 构造方法
func NewDefaultRegistry(c *DefaultCronTriggerRegistry, e *DefaultEventTriggerRegistry) *DefaultRegistry {
	return &DefaultRegistry{
		cronTriggerRegister:  c,
		eventTriggerRegister: e,
	}
}

// Register 注册逻辑实现
func (t *DefaultRegistry) Register(d *dto.RegisterTriggersDTO) error {
	if err := validateBasicParams(d.DefID, d.InstID, d.DefVersion, d.Level); err != nil {
		return err
	}

	log.Infof("[%s]Register triggers:%s",
		logs.GetFlowTraceID(d.DefID, d.InstID),
		utils.StructToJsonStr(d.Triggers))

	// FIXME, 后面可以对这些触发器进行分组, 批量插入
	for _, triggers := range d.Triggers {
		for name, trigger := range triggers {
			trigger.RefName = name
			registerDTO := &dto.RegisterTriggerDTO{
				DefID:      d.DefID,
				DefVersion: d.DefVersion,
				InstID:     d.InstID,
				Level:      d.Level,
				TriggerDef: trigger,
			}

			if trigger.Type == entity.Timer {
				if err := t.cronTriggerRegister.Register(registerDTO); err != nil {
					return err
				}
				continue
			}

			if err := t.eventTriggerRegister.Register(registerDTO); err != nil {
				return err
			}
		}
	}

	return nil
}

// UnRegister 反注册逻辑实现
func (t *DefaultRegistry) UnRegister(d *dto.UnRegisterTriggersDTO) error {
	if err := validateBasicParams(d.DefID, d.InstID, d.DefVersion, d.Level); err != nil {
		return err
	}

	unRegisterEventTriggerDTO := &dto.UnregisterTriggerDTO{
		DefID:      d.DefID,
		DefVersion: d.DefVersion,
		InstID:     d.InstID,
		Level:      d.Level,
	}

	if err := t.cronTriggerRegister.Unregister(unRegisterEventTriggerDTO); err != nil {
		return err
	}

	return t.eventTriggerRegister.Unregister(unRegisterEventTriggerDTO)
}

// validateBasicParams  校验参数
func validateBasicParams(defID, instID string, version int, level entity.TriggerLevel) error {
	if utils.IsZero(defID) || utils.IsZero(version) || utils.IsZero(level) {
		return fmt.Errorf("Register/Unregister `DefID`、`DefVersion`、`Level` must not be zero,"+
			"`DefID`:%s ,`DefVersion`: %d, `Level`:%v", defID, version, level)
	}

	if level == entity.InstTrigger && utils.IsZero(instID) {
		return fmt.Errorf("Register/Unregister inst level trigger `InstID` must not be zero")
	}

	return nil
}
