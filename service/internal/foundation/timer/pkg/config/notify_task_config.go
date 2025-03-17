package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

// NotifyTaskConfig 通知服务配置
type NotifyTaskConfig struct {
	ConsumerNum int `json:"consumerNum"` // 消费者数量
}

var (
	notifyTaskGroupKey = config.NewGroupKey("timer", "NOTIFYTASK") //  通知服务配置
)

// GetNotifyTaskConfig 获取 通知服务 默认配置, 没有特殊情况直接用默认配置就可以了
func GetNotifyTaskConfig() NotifyTaskConfig {
	conf := NotifyTaskConfig{
		ConsumerNum: 5,
	}
	provider.GetConfigProvider().GetAny(context.Background(), notifyTaskGroupKey, &conf)
	return conf
}
