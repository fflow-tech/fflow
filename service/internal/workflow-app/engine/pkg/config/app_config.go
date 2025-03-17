package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	commonGroupKey = config.NewGroupKey("engine", "APP") // 公共配置
)

// AppConfig 公共配置项
type AppConfig struct {
	*config.AuthConfig
}

// GetAppConfig 获取应用公共的配置项
func GetAppConfig() AppConfig {
	conf := &AppConfig{
		AuthConfig: &config.AuthConfig{},
	}
	provider.GetConfigProvider().GetAny(context.Background(), commonGroupKey, &conf)
	return *conf
}
