package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	tdMQGroupKey = config.NewGroupKey("engine", "TDMQ") // TDMQ MQ配置
)

// GetTDMQConfig 获取 TDMQ 默认配置, 默认七彩石获取
func GetTDMQConfig() config.TDMQConfig {
	conf := config.TDMQConfig{
		NackRedeliveryDelay: 5,
		MaxDeliveries:       10,
		RetryInitialDelay:   1,
		RetryMaxDelay:       60,
		RetryEnable:         true,
	}
	provider.GetConfigProvider().GetAny(context.Background(), tdMQGroupKey, &conf)
	return conf
}
