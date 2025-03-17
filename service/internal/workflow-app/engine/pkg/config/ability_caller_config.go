// Package config 提供常用的配置
package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
	"github.com/fflow-tech/fflow/service/pkg/remote"
)

var (
	abilityCallerGroupKey = config.NewGroupKey("engine", "AbilityCaller") // 事件配置
)

// GetDefaultAbilityCallerConfig 获取默认配置
func GetDefaultAbilityCallerConfig() *remote.DefaultAbilityCallerConfig {
	conf := remote.DefaultAbilityCallerConfig{
		FaasTarget:          "dns:///faas-grpc-service:50032",
		LoadBalancingPolicy: "round_robin",
	}
	provider.GetConfigProvider().GetAny(context.Background(), abilityCallerGroupKey, &conf)
	return &conf
}
