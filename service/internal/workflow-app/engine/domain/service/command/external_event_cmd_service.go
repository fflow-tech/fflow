package command

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/remote"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// ExternalEventCommandService 对外部事件进行处理
type ExternalEventCommandService struct {
	remoteRepo       ports.RemoteRepository
	eventBusRepo     ports.EventBusRepository
	workflowInstRepo ports.WorkflowInstRepository
	nodeInstRepo     ports.NodeInstRepository
	msgSender        common.MsgSender
	exprEvaluator    expr.Evaluator
}

// NewExternalEventCommandService 实例化 ExternalEventCommandService
func NewExternalEventCommandService(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *execution.WorkflowProviderSet) *ExternalEventCommandService {
	return &ExternalEventCommandService{
		remoteRepo:       repoProviderSet.RemoteRepo(),
		eventBusRepo:     repoProviderSet.EventBusRepo(),
		nodeInstRepo:     repoProviderSet.NodeInstRepo(),
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		msgSender:        workflowProviderSet.MsgSender(),
		exprEvaluator:    workflowProviderSet.ExprEvaluator(),
	}
}

var (
	defIDKey      = "def_id"
	instIDKey     = "inst_id"
	nodeInstIDKey = "node_inst_id"
)

// sendWebhooks 供内部使用的转发 Webhook 方法
func (s *ExternalEventCommandService) sendWebhooks(urls []string, value []byte) error {
	// 如果没有配置 Webhook，直接返回
	if len(urls) == 0 {
		return nil
	}
	body := map[string]interface{}{}
	if err := json.Unmarshal(value, &body); err != nil {
		return err
	}
	for i := 0; i < len(urls); i++ {
		go s.sendWebhook(urls[i], body)
	}

	return nil
}

func (s *ExternalEventCommandService) sendWebhook(url string, body map[string]interface{}) {
	if !utils.IsValidURL(url) {
		log.Warnf("Illegal url=[%s]", url)
		return
	}

	req := &remote.CallHTTPReqDTO{
		Method: http.MethodPost,
		URL:    url,
		Header: nil,
		Query:  nil,
		Body:   body,
	}
	rsp, err := s.remoteRepo.CallHTTP(context.Background(), req)
	if err != nil {
		// 失败了不影响后续的发送
		log.Warnf("Failed to send webhook:[%s], caused by %s", url, err)
	}
	log.Infof("Send webhook:[%s], rsp is %s", url, utils.StructToJsonStr(rsp))
}

// ConsumeForSendWebhook 消费事件
func (s *ExternalEventCommandService) ConsumeForSendWebhook(ctx context.Context, req *dto.ExternalEventDTO) error {
	externalEvent := req.Message
	basicEvent := event.BasicEvent{}
	if err := json.Unmarshal(externalEvent.Payload(), &basicEvent); err != nil {
		return err
	}
	urls, err := s.getWebhookURLs(externalEvent.Payload(), basicEvent.EventType)
	if err != nil {
		return err
	}
	return s.sendWebhooks(urls, externalEvent.Payload())
}

// ConsumeForSendChatMsg 发送聊天提醒
func (s *ExternalEventCommandService) ConsumeForSendChatMsg(ctx context.Context, req *dto.ExternalEventDTO) error {
	externalEvent := req.Message
	basicEvent := event.BasicEvent{}
	if err := json.Unmarshal(externalEvent.Payload(), &basicEvent); err != nil {
		return err
	}

	eventType := event.ExternalEventType(basicEvent.EventType)

	return s.sendChatMsg(externalEvent, eventType)
}

func (s *ExternalEventCommandService) sendChatMsg(externalEvent pulsar.Message,
	eventType event.ExternalEventType) error {
	if !event.IsWorkflowInstLevelEvent(eventType) && !event.IsNodeInstLevelEvent(eventType) {
		log.Warnf("Skip event type for send chat msg:%s", eventType)
		return nil
	}

	inst, err := s.getWorkflowInst(externalEvent.Payload())
	if err != nil {
		return err
	}

	nodeInst, err := s.getNodeInstForSendChatMsg(externalEvent.Payload(), eventType)
	if err != nil {
		return err
	}

	if err := s.sendMsgToGroup(inst, nodeInst, eventType); err != nil {
		return err
	}

	return s.sendMsgToUser(inst, nodeInst, eventType)
}

func (s *ExternalEventCommandService) sendMsgToGroup(inst *entity.WorkflowInst,
	nodeInst *entity.NodeInst, eventType event.ExternalEventType) error {
	if inst.Owner == nil || inst.Owner.ChatGroup == "" {
		return nil
	}

	chatIDs := strings.Split(inst.Owner.ChatGroup, ";")
	msg, err := s.evaluateChatMsg(inst, nodeInst, eventType)
	if err != nil {
		return err
	}
	if msg == "" {
		return nil
	}

	for _, chatID := range chatIDs {
		if err := s.remoteRepo.SendMsgToGroup(chatID, msg); err != nil {
			log.Warnf("[%s]Failed to send msg to group [%s], caused by %s",
				logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID), chatID, err)
			continue
		}
	}

	return nil
}

func (s *ExternalEventCommandService) sendMsgToUser(inst *entity.WorkflowInst,
	nodeInst *entity.NodeInst, eventType event.ExternalEventType) error {
	if inst.Owner == nil || inst.Owner.Wechat == "" {
		return nil
	}

	userIDs := strings.Split(inst.Owner.Wechat, ";")
	msg, err := s.evaluateChatMsg(inst, nodeInst, eventType)
	if err != nil {
		return err
	}
	if msg == "" {
		return nil
	}

	for _, userID := range userIDs {
		if err := s.remoteRepo.SendMsgToUser(userID, msg); err != nil {
			log.Warnf("[%s]Failed to send msg to user [%s], caused by %s",
				logs.GetFlowTraceID(inst.WorkflowDef.DefID, inst.InstID), userID, err)
			continue
		}
	}

	return nil
}

func (s *ExternalEventCommandService) evaluateChatMsg(inst *entity.WorkflowInst,
	nodeInst *entity.NodeInst, eventType event.ExternalEventType) (string, error) {
	msgFormat := event.GetChatMsgFormat(inst, nodeInst, eventType)
	if !s.exprEvaluator.IsExpression(msgFormat) {
		return msgFormat, nil
	}

	ctx, err := entity.ConvertToCtx(inst)
	if err != nil {
		return "", err
	}

	if err := entity.AppendNodeInfoToCtxKey(ctx, nodeInst, constants.ThisNode); err != nil {
		return "", err
	}

	msg, err := s.exprEvaluator.Evaluate(ctx, msgFormat)
	if err != nil {
		return "", err
	}
	return msg.(string), nil
}

// ConsumeForWorkflowExceptionHappened 消费工作流异常事件
func (s *ExternalEventCommandService) ConsumeForWorkflowExceptionHappened(ctx context.Context,
	req *dto.ExternalEventDTO) error {
	externalEvent := req.Message
	basicEvent := event.BasicEvent{}
	if err := json.Unmarshal(externalEvent.Payload(), &basicEvent); err != nil {
		return err
	}

	workflowInst, err := s.getWorkflowInst(externalEvent.Payload())
	if err != nil {
		return err
	}
	return s.sendAlertWorkflowExceptionHappened(workflowInst, basicEvent.Reason, basicEvent.EventType)
}

// ConsumeForNodeInstExceptionHappened 消费节点实例异常事件
func (s *ExternalEventCommandService) ConsumeForNodeInstExceptionHappened(ctx context.Context,
	req *dto.ExternalEventDTO) error {
	externalEvent := req.Message
	nodeInst, err := s.getNodeInst(externalEvent.Payload())
	if err != nil {
		return err
	}

	basicEvent := event.BasicEvent{}
	if err := json.Unmarshal(externalEvent.Payload(), &basicEvent); err != nil {
		return err
	}

	return s.sendAlertForNodeInstExceptionHappened(nodeInst, basicEvent)
}

func (s *ExternalEventCommandService) getWebhookURLs(value []byte, eventType string) ([]string, error) {
	if event.IsWorkflowInstLevelEvent(event.ExternalEventType(eventType)) {
		return s.getWorkflowWebhookURLs(value)
	}

	if event.IsNodeInstLevelEvent(event.ExternalEventType(eventType)) {
		return s.getNodeWebhookURLs(value)
	}

	// 其它类型的事件，忽略掉
	log.Infof("Skip event type for send webhook: %s", eventType)
	return []string{}, nil
}

func (s *ExternalEventCommandService) getWorkflowWebhookURLs(value []byte) ([]string, error) {
	workflowInst, err := s.getWorkflowInst(value)
	if err != nil {
		return []string{}, err
	}
	return workflowInst.WorkflowDef.Webhooks, nil
}

func (s *ExternalEventCommandService) getWorkflowInst(value []byte) (*entity.WorkflowInst, error) {
	defID, err := utils.GetStrFromJson(value, defIDKey)
	if err != nil {
		return nil, err
	}
	instID, err := utils.GetStrFromJson(value, instIDKey)
	if err != nil {
		return nil, err
	}
	// 从数据库中获取 Webhook URLs
	getDTO := &dto.GetWorkflowInstDTO{
		DefID:  defID,
		InstID: instID,
	}
	workflowInst, err := s.workflowInstRepo.Get(getDTO)
	if err != nil {
		return nil, err
	}
	return workflowInst, nil
}

func (s *ExternalEventCommandService) getNodeWebhookURLs(value []byte) ([]string, error) {
	defID, err := utils.GetStrFromJson(value, defIDKey)
	if err != nil {
		return nil, err
	}
	nodeInstID, err := utils.GetStrFromJson(value, nodeInstIDKey)
	if err != nil {
		return nil, err
	}
	// 从数据库中获取 Webhook URLs
	getDTO := &dto.GetNodeInstDTO{
		DefID:      defID,
		NodeInstID: nodeInstID,
	}
	nodeInst, err := s.nodeInstRepo.Get(getDTO)
	if err != nil {
		return nil, err
	}
	return nodeInst.BasicNodeDef.Webhooks, nil
}

// sendAlertWorkflowExceptionHappened 流程异常时发送告警
func (s *ExternalEventCommandService) sendAlertWorkflowExceptionHappened(workflowInst *entity.WorkflowInst,
	reason, eventType string) error {
	if workflowInst.Owner == nil {
		return nil
	}

	msgInfo := map[string]interface{}{
		"InstID":          workflowInst.InstID,
		"ExceptionalType": eventType,
		"Reason":          template.HTML(reason),
	}
	if workflowInst.Owner.ChatGroup != "" {
		return s.msgSender.SendChatGroupMsg(workflowInst.Owner.ChatGroup, common.InstExceptionalAlert, msgInfo)
	}

	if workflowInst.Owner.Wechat != "" {
		return s.msgSender.SendWeChatMsg(workflowInst.Owner.Wechat, common.InstExceptionalAlert, msgInfo)
	}

	return nil
}

func (s *ExternalEventCommandService) getNodeInstForSendChatMsg(value []byte,
	eventType event.ExternalEventType) (*entity.NodeInst, error) {
	if !event.IsNodeInstLevelEvent(eventType) {
		return nil, nil
	}

	return s.getNodeInst(value)
}

func (s *ExternalEventCommandService) getNodeInst(value []byte) (*entity.NodeInst, error) {
	defID, err := utils.GetStrFromJson(value, defIDKey)
	if err != nil {
		return nil, err
	}

	nodeInstID, err := utils.GetStrFromJson(value, nodeInstIDKey)
	if err != nil {
		return nil, err
	}
	getDTO := &dto.GetNodeInstDTO{
		NodeInstID: nodeInstID,
		DefID:      defID,
	}
	return s.nodeInstRepo.Get(getDTO)
}

// sendAlertForNodeInstExceptionHappened 节点异常时发送告警
func (s *ExternalEventCommandService) sendAlertForNodeInstExceptionHappened(nodeInst *entity.NodeInst,
	basicEvent event.BasicEvent) error {
	if nodeInst.Owner == nil {
		return nil
	}

	params := map[string]interface{}{
		"InstID":          nodeInst.InstID,
		"NodeRefName":     nodeInst.BasicNodeDef.RefName,
		"ExceptionalType": basicEvent.EventType,
		"Reason":          template.HTML(basicEvent.Reason),
	}

	if nodeInst.Owner.ChatGroup != "" {
		return s.msgSender.SendChatGroupMsg(nodeInst.Owner.ChatGroup, common.NodeInstExceptionalAlert, params)
	}

	if nodeInst.Owner.Wechat != "" {
		return s.msgSender.SendWeChatMsg(nodeInst.Owner.Wechat, common.NodeInstExceptionalAlert, params)
	}

	return nil
}
