package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

const (
	defaultNetwork     = "tcp"
	defaultMaxIdle     = 3
	defaultIdleTimeout = 24
)

var (
	redisGroupKey = config.NewGroupKey("engine", "REDIS") // 连接信息
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
