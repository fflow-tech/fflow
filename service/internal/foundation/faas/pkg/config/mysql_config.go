package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	mySQLGroupKey = config.NewGroupKey("faas", "MYSQL") // MySQL 数据库配置
)

// GetMySQLConfig 获取 mysql 默认配置, 没有特殊情况直接用默认配置就可以了
func GetMySQLConfig() config.MySQLConfig {
	conf := &config.MySQLConfig{
		SlowThreshold:             200,
		IgnoreRecordNotFoundError: true,
	}
	provider.GetConfigProvider().GetAny(context.Background(), mySQLGroupKey, &conf)
	return *conf
}
