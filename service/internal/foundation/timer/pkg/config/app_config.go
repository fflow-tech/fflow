package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	commonGroupKey = config.NewGroupKey("timer", "APP") // 公共配置
)

// AppConfig 公共配置项
type AppConfig struct {
	*config.AuthConfig
	KeepDays int               `json:"KeepDays"`
	Accounts map[string]string `json:"accounts"`
	Realm    string            `json:"realm"`
}

// GetAppConfig 获取 Timer 公共的配置项
func GetAppConfig() AppConfig {
	conf := &AppConfig{}
	provider.GetConfigProvider().GetAny(context.Background(), commonGroupKey, &conf)
	return *conf
}
