package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	schemaJsonGroupKey = config.NewGroupKey("engine", "SCHEMA_JSON") // SchemaJson json校验配置
)

// GetSchemaConfig 获取schema默认配置
func GetSchemaConfig() string {
	str, _ := provider.GetConfigProvider().GetString(context.Background(), schemaJsonGroupKey)
	return str
}
