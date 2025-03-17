package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	validationRulesGroupKey = config.NewGroupKey("engine", "VALIDATION_RULES") // ValidationRules 校验规则限制
)

// ValidationRulesConfig 缓存配置
type ValidationRulesConfig struct {
	DefJsonSize               int `json:"defJsonSize"`               // 定义json大小
	WorkflowInstCtxSize       int `json:"workflowInstCtxSize"`       // 流程实例上下文大小
	NodeInstCtxSize           int `json:"nodeInstCtxSize"`           // 节点实例上下文大小
	MaxNodesCount             int `json:"maxNodesCount"`             // 最大节点数量
	MinNodesCount             int `json:"minNodesCount"`             // 最小节点数量
	MaxNodeInstsForOneFlow    int `json:"maxNodeInstsForOneFlow"`    // 一个流程最大的节点实例数量
	MaxWorkflowInstsForOneDef int `json:"maxWorkflowInstsForOneDef"` // 一个流程定义最大的流程实例数量
	MaxInputArgsSize          int `json:"maxInputArgsSize"`          // 启动流程input最大字节数
}

// GetValidationRulesConfig 获取默认配置
func GetValidationRulesConfig() ValidationRulesConfig {
	conf := ValidationRulesConfig{
		DefJsonSize:               1024 * 1024,
		WorkflowInstCtxSize:       1024 * 1024,
		NodeInstCtxSize:           1024 * 1024,
		MaxNodesCount:             200,
		MinNodesCount:             1,
		MaxNodeInstsForOneFlow:    1000,
		MaxInputArgsSize:          1 * 1024 * 4,
		MaxWorkflowInstsForOneDef: 500,
	}
	provider.GetConfigProvider().GetAny(context.Background(), validationRulesGroupKey, &conf)
	return conf
}
