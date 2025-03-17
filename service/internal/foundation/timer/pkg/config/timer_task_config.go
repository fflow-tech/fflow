package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

// TimerTaskConfig 定时器服务配置
type TimerTaskConfig struct {
	ConsumerNum          int   `json:"consumerNum"`          // 消费者数量
	WorkSleepSecond      int64 `json:"workSleepSecond"`      // 工作协程休眠秒数
	WorkDelayMillisecond int64 `json:"workDelayMillisecond"` // 工作协程的延时毫秒数
}

var (
	timerTaskGroupKey = config.NewGroupKey("timer", "TIMERTASK") //  定时器服务配置
)

// GetTimerTaskConfig 获取 定时器服务 默认配置, 没有特殊情况直接用默认配置就可以了
func GetTimerTaskConfig() TimerTaskConfig {
	conf := TimerTaskConfig{
		ConsumerNum:          5,
		WorkSleepSecond:      1,
		WorkDelayMillisecond: 1500,
	}
	provider.GetConfigProvider().GetAny(context.Background(), timerTaskGroupKey, &conf)
	return conf
}
