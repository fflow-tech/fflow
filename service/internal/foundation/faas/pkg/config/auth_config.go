package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
	"github.com/fflow-tech/fflow/service/pkg/remote"
)

var (
	permissionValidatorGroupKey = config.NewGroupKey("faas", "PermissionValidator") // 事件配置
)

// GetDefaultPermissionValidatorConfig 获取默认配置
func GetDefaultPermissionValidatorConfig() *remote.DefaultPermissionValidatorConfig {
	conf := remote.DefaultPermissionValidatorConfig{
		AuthTarget:          "dns:///auth-grpc-service:50042",
		LoadBalancingPolicy: "round_robin",
	}
	provider.GetConfigProvider().GetAny(context.Background(), permissionValidatorGroupKey, &conf)
	return &conf
}
