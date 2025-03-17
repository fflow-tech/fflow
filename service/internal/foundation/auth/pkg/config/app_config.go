package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	appGroupKey = config.NewGroupKey("auth", "APP") // WebCors 跨域配置
)

// AppConfig 应用配置
type AppConfig struct {
	*config.AuthConfig
}

// GetAppConfig 获取应用通用配置
func GetAppConfig() *AppConfig {
	conf := AppConfig{
		AuthConfig: &config.AuthConfig{},
	}
	provider.GetConfigProvider().GetAny(context.Background(), appGroupKey, &conf)
	return &conf
}
