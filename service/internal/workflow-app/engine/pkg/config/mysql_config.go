package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	mysqlGroupKey = config.NewGroupKey("engine", "MYSQL") // Kafka 连接信息
)

// GetMySQLConfig 获取 mysql 默认配置, 没有特殊情况直接用默认配置就可以了
func GetMySQLConfig() config.MySQLConfig {
	conf := &config.MySQLConfig{
		SlowThreshold:             200,
		IgnoreRecordNotFoundError: true,
		SkipDefaultTransaction:    true,
	}
	provider.GetConfigProvider().GetAny(context.Background(), mysqlGroupKey, &conf)
	return *conf
}
