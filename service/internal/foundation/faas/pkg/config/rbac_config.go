package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	rbacGroupKey = config.NewGroupKey("faas", "RBAC") // MySQL 数据库配置
)

// GetRbacConfig 获取权限默认配置
func GetRbacConfig() *config.RbacConfig {
	conf := config.RbacConfig{
		SuperAdmins: []string{},
	}
	provider.GetConfigProvider().GetAny(context.Background(), rbacGroupKey, &conf)
	return &conf
}
