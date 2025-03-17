package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	commonGroupKey = config.NewGroupKey("faas", "APP") // 公共配置
)

// AppConfig 公共配置项
type AppConfig struct {
	*config.AuthConfig
	KeepDays  int               `json:"keepDays"`
	Accounts  map[string]string `json:"accounts"`
	Realm     string            `json:"realm"`
	CallToken string            `json:"callToken"`
}

// GetAppConfig 获取 FAAS 公共的配置项
func GetAppConfig() *AppConfig {
	conf := &AppConfig{}
	provider.GetConfigProvider().GetAny(context.Background(), commonGroupKey, &conf)
	return conf
}
