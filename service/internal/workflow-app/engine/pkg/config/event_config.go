// Package config 提供常用的配置
package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	eventBusGroupKey = config.NewGroupKey("engine", "EVENTBUS") // 事件配置
)

// EventConfig 配置
type EventConfig struct {
	DriveEventTopic          string `json:"driveEventTopic"`
	DriveEventGroup          string `json:"driveEventGroup"`
	ExternalEventTopic       string `json:"externalEventTopic"`
	ExternalEventGroup       string `json:"externalEventGroup"`
	CronEventTopic           string `json:"cronEventTopic"`
	CronEventGroup           string `json:"cronEventGroup"`
	TriggerEventTopic        string `json:"triggerEventTopic"`
	TriggerEventGroup        string `json:"triggerEventGroup"`
	PollingEventTopic        string `json:"pollingEventTopic"`
	PollingEventGroup        string `json:"pollingEventGroup"`
	DriveEventConsumerNum    int    `json:"driveEventConsumerNum"`
	ExternalEventConsumerNum int    `json:"externalEventConsumerNum"`
	TriggerEventConsumerNum  int    `json:"triggerEventConsumerNum"`
	CronEventConsumerNum     int    `json:"cronEventConsumerNum"`
}

// GetEventConfig 获取默认配置
func GetEventConfig() *EventConfig {
	conf := EventConfig{
		DriveEventConsumerNum:    8,
		ExternalEventConsumerNum: 2,
		TriggerEventConsumerNum:  2,
		CronEventConsumerNum:     2,
		DriveEventTopic:          "driven_event",
		DriveEventGroup:          "EngineDrivenEvent.DefaultGroup",
		ExternalEventTopic:       "external_event",
		ExternalEventGroup:       "EngineExternalEvent.DefaultGroup",
		CronEventTopic:           "cron_event",
		CronEventGroup:           "EngineCronEvent.DefaultGroup",
		TriggerEventTopic:        "trigger_event",
		TriggerEventGroup:        "EngineTriggerEvent.DefaultGroup",
	}
	provider.GetConfigProvider().GetAny(context.Background(), eventBusGroupKey, &conf)
	return &conf
}
