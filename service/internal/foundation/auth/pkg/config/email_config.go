package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	emailGroupKey = config.NewGroupKey("auth", "EMAIL") // WebCors 跨域配置
)

// GetEmailConfig 获取应用邮件配置
func GetEmailConfig() *config.EmailConfig {
	conf := config.EmailConfig{}
	provider.GetConfigProvider().GetAny(context.Background(), emailGroupKey, &conf)
	return &conf
}
