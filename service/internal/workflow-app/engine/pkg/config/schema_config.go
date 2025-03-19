package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	schemaJsonGroupKey = config.NewGroupKey("engine", "SCHEMA_JSON") // SchemaJson json 校验配置
)

// GetSchemaConfig 获取 schema 默认配置
func GetSchemaConfig() string {
	defaultSchemaJSON := `{"title":"content","type":"object","description":"流程定义规则","required":["name"],"properties":{"name":{"type":"string","description":"流程名称","minLength":1},"desc":{"type":"string","description":"流程描述","minLength":1},"timeout":{"type":"object","description":"流程超时配置","properties":{"duration":{"type":"string","description":"流程超时时间"}}},"biz":{"type":"object","description":"业务配置，用户可以自定义"},"variables":{"type":"object","description":"全局变量","properties":{}},"owner":{"type":"object","properties":{"wechat":{"type":"string","description":"流程的拥有者的企业微信号"},"groupChat":{"type":"string","description":"流程相关的群聊"}}},"webhooks":{"type":"array","description":"流程事件webhook地址列表","items":{"type":"string"}}}}`
	str, _ := provider.GetConfigProvider().GetString(context.Background(), schemaJsonGroupKey)
	if str == "" {
		return defaultSchemaJSON
	}
	return str
}
