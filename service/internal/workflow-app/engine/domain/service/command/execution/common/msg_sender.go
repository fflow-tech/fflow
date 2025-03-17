package common

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// MsgSender 告警发送者
type MsgSender interface {
	SendWeChatMsg(receivers string, template MsgTemplate, params map[string]interface{}) error     // 发送个人消息
	SendChatGroupMsg(chatGroups string, template MsgTemplate, params map[string]interface{}) error // 发送群聊消息
}

// MsgTemplate 消息模板
type MsgTemplate string

const (
	InstTimeoutAlert         MsgTemplate = "INST_TIMEOUT_ALERT"           // 流程实例超时告警模板
	NodeInstTimeoutAlert     MsgTemplate = "NODE_INST_TIMEOUT_ALERT"      // 节点超时告警模板
	NodeInstNearTimeoutAlert MsgTemplate = "NODE_INST_NEAR_TIMEOUT_ALERT" // 节点接近超时告警模板
	InstExceptionalAlert     MsgTemplate = "INST_EXCEPTIONAL_ALERT"       // 流程实例异常告警
	NodeInstExceptionalAlert MsgTemplate = "NODE_INST_EXCEPTIONAL_ALERT"  // 节点实例异常告警
)

// DefaultMsgSender 默认消息发送者
type DefaultMsgSender struct {
	remoteRepo ports.RemoteRepository
}

// NewDefaultMsgSender 实例化
func NewDefaultMsgSender(remoteRepo *repo.RemoteRepo) *DefaultMsgSender {
	return &DefaultMsgSender{
		remoteRepo: remoteRepo,
	}
}

// SendWeChatMsg 发送个人消息
// 多用户使用分号";"进行分隔
func (s *DefaultMsgSender) SendWeChatMsg(receivers string, template MsgTemplate, params map[string]interface{}) error {
	contentTemplate := getMsgContentTemplate(template)
	content, err := ParseTemplate(contentTemplate, params)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s - %v \n%s", getMsgTitle(template), template, content)
	users := strings.Split(receivers, ";")
	log.Infof("SendWeChatMsg user:%v", users)
	for _, user := range users {
		if err := s.remoteRepo.SendMsgToUser(user, msg); err != nil {
			return err
		}
	}
	return nil
}

// SendChatGroupMsg 发送群聊消息
// 多群聊使用分号";"进行分隔
func (s *DefaultMsgSender) SendChatGroupMsg(chatGroups string,
	template MsgTemplate, params map[string]interface{}) error {
	contentTemplate := getMsgContentTemplate(template)
	content, err := ParseTemplate(contentTemplate, params)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s\n%s", getMsgTitle(template), content)
	chatGroupNames := strings.Split(chatGroups, ";")
	log.Infof("SendChatGroupMsg user:%v", chatGroupNames)
	for _, chatGroup := range chatGroupNames {
		if err := s.remoteRepo.SendMsgToGroup(chatGroup, msg); err != nil {
			return err
		}
	}
	return nil
}

// getMsgContentTemplate 获取内容模版
func getMsgContentTemplate(template MsgTemplate) string {
	switch template {
	case InstExceptionalAlert:
		return config.GetMsgTemplateConfig().WorkflowAlertContentTemplate
	case NodeInstExceptionalAlert:
		return config.GetMsgTemplateConfig().NodeAlertContentTemplate
	default:
		return config.GetMsgTemplateConfig().DefaultAlertTemplateTitle
	}
}

// getMsgTitle 获取消息标题模版
func getMsgTitle(template MsgTemplate) string {
	return config.GetMsgTemplateConfig().DefaultAlertTemplateTitle
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
