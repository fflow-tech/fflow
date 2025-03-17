package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	corsGroupKey = config.NewGroupKey("timer", "CORS") // WebCors 跨域配置
)

// GetCorsConfig 获取默认配置
func GetCorsConfig() *config.CorsConfig {
	conf := config.CorsConfig{}
	provider.GetConfigProvider().GetAny(context.Background(), corsGroupKey, &conf)
	return &conf
}
