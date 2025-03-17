package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	workerPoolGroupKey = config.NewGroupKey("timer", "WORKERPOOL") // timer 全局协程池配置.
)

// WorkerPool 协程池配置.
type WorkerPool struct {
	// 协程池容量上限.
	Size int `json:"size"`
	// 空闲协程回收时间，单位：秒.
	ExpireDuration int `json:"expireDuration"`
}

// GetWorkerPoolConfig 获取 Timer 协程池配置.
func GetWorkerPoolConfig() WorkerPool {
	conf := WorkerPool{
		Size:           10000,
		ExpireDuration: 10,
	}
	provider.GetConfigProvider().GetAny(context.Background(), workerPoolGroupKey, &conf)
	return conf
}
