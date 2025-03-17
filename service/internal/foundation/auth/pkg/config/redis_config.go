package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	redisGroupKey = config.NewGroupKey("auth", "REDIS") // Redis 缓存配置
)

const (
	defaultNetwork     = "tcp"
	defaultMaxIdle     = 3
	defaultIdleTimeout = 24
)

// GetRedisConfig 获取默认配置
func GetRedisConfig() config.RedisConfig {
	conf := config.RedisConfig{
		Network:     defaultNetwork,
		MaxIdle:     defaultMaxIdle,
		IdleTimeout: defaultIdleTimeout,
	}
	provider.GetConfigProvider().GetAny(context.Background(), redisGroupKey, &conf)
	return conf
}
