// Package config 提供常用的配置
package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	historyArchiverGroupKey = config.NewGroupKey("engine", "HISTORY_ARCHIVER") // 超时检查器配置的 group key
)

// HistoryArchiveConfig 超时检查配置
type HistoryArchiveConfig struct {
	EnableArchive     bool `json:"enableArchive"`    // 是否启用超时检查
	GoroutinePoolSize int  `json:"poolSize"`         // 归档协程池大小
	KeepDataDuration  int  `json:"keepDataDuration"` // 保留数据时长，单位天
}

// GetHistoryArchiverConfig 获取默认配置
func GetHistoryArchiverConfig() *HistoryArchiveConfig {
	conf := HistoryArchiveConfig{GoroutinePoolSize: 50, KeepDataDuration: 90, EnableArchive: true}
	provider.GetConfigProvider().GetAny(context.Background(), historyArchiverGroupKey, &conf)
	return &conf
}
