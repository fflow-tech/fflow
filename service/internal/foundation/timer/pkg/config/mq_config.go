package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	tdmqGroupKey = config.NewGroupKey("timer", "TDMQ")
)

// GetDefaultTDMQConfig 获取TDMQ默认配置, 默认七彩石获取
func GetDefaultTDMQConfig() config.TDMQConfig {
	conf := config.TDMQConfig{
		NackRedeliveryDelay: 5,
		MaxDeliveries:       10,
		RetryInitialDelay:   1,
		RetryMaxDelay:       60,
		RetryEnable:         true,
	}
	provider.GetConfigProvider().GetAny(context.Background(), tdmqGroupKey, &conf)
	return conf
}
