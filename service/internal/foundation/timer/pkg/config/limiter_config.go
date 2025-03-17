package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

const (
	defaultBurst           = 1000
	defaultLimit           = 200
	defaultWaitingSeconds  = 60
	defaultRefreshInterval = 120
)

var (
	limiterGroupKey = config.NewGroupKey("timer", "LIMITER") // EventBus 事件配置
)

// GetLimiterConfig 获取定时器服务的限流器配置
func GetLimiterConfig() *config.LimiterConfig {
	conf := config.LimiterConfig{
		Burst:           defaultBurst,
		Limit:           defaultLimit,
		WaitingDuration: defaultWaitingSeconds,
		RefreshInterval: defaultRefreshInterval,
	}

	provider.GetConfigProvider().GetAny(context.Background(), limiterGroupKey, &conf)
	return &conf
}
