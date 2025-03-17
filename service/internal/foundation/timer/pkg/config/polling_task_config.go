package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

// PollingTaskConfig 轮询服务配置
type PollingTaskConfig struct {
	WorkNum         int64 `json:"workNum"`         // 工作协程数
	WorkSleepSecond int64 `json:"workSleepSecond"` // 工作协程休眠秒数
	BucketNum       int   `json:"bucketNum"`       // 桶数
}

var (
	pollingTaskGroupKey = config.NewGroupKey("timer", "POLLINGTASK") //  轮询服务配置
)

// GetPollingTaskConfig 获取 轮询服务 默认配置, 没有特殊情况直接用默认配置就可以了
func GetPollingTaskConfig() PollingTaskConfig {
	conf := PollingTaskConfig{
		WorkNum:         4,
		WorkSleepSecond: 10,
		BucketNum:       20,
	}
	provider.GetConfigProvider().GetAny(context.Background(), pollingTaskGroupKey, &conf)
	return conf
}
