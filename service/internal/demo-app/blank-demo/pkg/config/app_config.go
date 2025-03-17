package config

import (
	"context"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	appGroupKey = config.NewGroupKey("blank-demo", "APP") // WebCors 跨域配置
)

// AppConfig 应用配置
type AppConfig struct {
	*config.AuthConfig
}

// GetAppConfig 获取应用通用配置
func GetAppConfig() *AppConfig {
	conf := AppConfig{
		AuthConfig: &config.AuthConfig{
			AdminUsername: "admin",
			AdminPassword: "rfnypTmy5kJTPtAEJGxv",
			AdminEmail:    "admin@fflow.link",
			Domain:        "localhost",
			SecretKey:     "rfnypTmy5kJTPtAEJGxv",
		},
	}
	provider.GetConfigProvider().GetAny(context.Background(), appGroupKey, &conf)
	return &conf
}
