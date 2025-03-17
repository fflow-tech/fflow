package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	eventGroupKey = config.NewGroupKey("timer", "EVENTBUS") // EventBus 事件配置
)

// EventConfig 配置
type EventConfig struct {
	TimerEventTopic string `json:"timerEventTopic"`
	TimerEventGroup string `json:"timerEventGroup"`
}

// GetEventConfig 获取默认配置
func GetEventConfig() *EventConfig {
	conf := EventConfig{}
	provider.GetConfigProvider().GetAny(context.Background(), eventGroupKey, &conf)
	return &conf
}
