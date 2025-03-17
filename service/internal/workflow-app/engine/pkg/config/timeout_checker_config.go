// Package config 提供常用的配置
package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	timeoutCheckerGroupKey = config.NewGroupKey("engine", "TIMEOUT_CHECKER") // 超时检查器配置的 group key
)

// TimeoutCheckerConfig 超时检查配置
type TimeoutCheckerConfig struct {
	GoroutinePoolSize int `json:"poolSize"` // 超时检查协程池大小
}

// GetTimeoutCheckerConfig 获取默认配置
func GetTimeoutCheckerConfig() *TimeoutCheckerConfig {
	conf := TimeoutCheckerConfig{GoroutinePoolSize: 50}
	provider.GetConfigProvider().GetAny(context.Background(), timeoutCheckerGroupKey, &conf)
	return &conf
}
