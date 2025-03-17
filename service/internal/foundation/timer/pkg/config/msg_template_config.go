package config

import (
	"bytes"
	"html/template"
)

// MsgTemplateConfig 消息模板配置
type MsgTemplateConfig struct {
	DefaultAlertTemplateTitle   string `json:"defaultAlertTemplateTitle,omitempty"`   // 告警默认标题
	DefaultAlertContentTemplate string `json:"defaultAlertContentTemplate,omitempty"` // 告警默认内容模板
}

// GetMsgTemplateConfig 获取默认配置
func GetMsgTemplateConfig() *MsgTemplateConfig {
	conf := MsgTemplateConfig{
		DefaultAlertTemplateTitle:   "【定时器服务】告警通知",
		DefaultAlertContentTemplate: "【定时器ID:{{.DefID}}, 定时器名称:{{.TimerName}} 失败原因:{{.Fail}} 请关注!】",
	}
	return &conf
}

// GetMsgTitle 获取消息标题模版
func GetMsgTitle() string {
	return GetMsgTemplateConfig().DefaultAlertTemplateTitle
}

// GetMsgContentTemplate 获取内容模版
func GetMsgContentTemplate() string {
	return GetMsgTemplateConfig().DefaultAlertContentTemplate
}

// ParseTemplate 解析模板
func ParseTemplate(text string, m interface{}) (string, error) {
	tmpl, err := template.New("t").Parse(text)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = tmpl.Option().Execute(&b, m)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
