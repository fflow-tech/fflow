package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
)

var (
	msgTemplateGroupKey = config.NewGroupKey("engine", "MSG_TEMPLATE") // MsgTemplate 消息模板配置
)

// MsgTemplateConfig 消息模板配置
type MsgTemplateConfig struct {
	DefaultAlertTemplateTitle    string `json:"defaultAlertTemplateTitle,omitempty"`    // 告警默认标题
	DefaultAlertContentTemplate  string `json:"defaultAlertContentTemplate,omitempty"`  // 告警默认内容模板
	WorkflowAlertContentTemplate string `json:"workflowAlertContentTemplate,omitempty"` // 流程告警异常内容模板
	NodeAlertContentTemplate     string `json:"nodeAlertContentTemplate,omitempty"`     // 节点告警异常内容模板
}

// GetMsgTemplateConfig 获取默认配置
func GetMsgTemplateConfig() *MsgTemplateConfig {
	conf := MsgTemplateConfig{
		DefaultAlertTemplateTitle:   "【流程引擎2.0】告警通知",
		DefaultAlertContentTemplate: "【流程实例ID:{{.InstID}}, 节点实例名称:{{.NodeRefName}} 请关注!】",
		WorkflowAlertContentTemplate: "【流程实例ID:{{.InstID}}\n" +
			"异常类型:{{.ExceptionalType}} 异常原因:{{.Reason}} 请关注!】",
		NodeAlertContentTemplate: "【流程实例ID:{{.InstID}}, 节点实例名称:{{.NodeRefName}} \n" +
			"异常类型:{{.ExceptionalType}} 异常原因:{{.Reason}} 请关注!】",
	}
	provider.GetConfigProvider().GetAny(context.Background(), msgTemplateGroupKey, &conf)
	return &conf
}
