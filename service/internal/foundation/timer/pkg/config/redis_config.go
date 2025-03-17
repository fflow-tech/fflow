package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	redisGroupKey = config.NewGroupKey("timer", "REDIS") // Redis 缓存配置
)

const (
	defaultNetwork     = "tcp"
	defaultMaxIdle     = 2500
	defaultIdleTimeout = 24
	defaultMaxActive   = 2500
	defaultWaitFlag    = true
)

// GetRedisConfig 获取默认配置
func GetRedisConfig() config.RedisConfig {
	conf := config.RedisConfig{
		Network:     defaultNetwork,
		MaxIdle:     defaultMaxIdle,
		IdleTimeout: defaultIdleTimeout,
		MaxActive:   defaultMaxActive,
		Wait:        defaultWaitFlag,
	}
	provider.GetConfigProvider().GetAny(context.Background(), redisGroupKey, &conf)
	return conf
}
