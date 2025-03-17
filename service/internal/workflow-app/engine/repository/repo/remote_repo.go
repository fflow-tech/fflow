package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/remote"
)

// RemoteRepo 远程调用仓储层
type RemoteRepo struct {
	abilityCaller    remote.AbilityCaller
	cronClient       remote.CronClient
	chatOpsClient    remote.ChatOpsClient
	cloudEventClient remote.CloudEventClient
}

// NewRemoteRepo 实体构造函数
func NewRemoteRepo(abilityCaller *remote.DefaultAbilityCaller,
	cronClient *remote.DefaultCronClient,
	chatOpsClient *remote.DefaultChatOpsClient,
	cloudEventClient *remote.DefaultCloudEventClient) *RemoteRepo {
	return &RemoteRepo{
		abilityCaller:    abilityCaller,
		cronClient:       cronClient,
		chatOpsClient:    chatOpsClient,
		cloudEventClient: cloudEventClient,
	}
}

// CallFAAS 调用 faas 能力
func (t *RemoteRepo) CallFAAS(ctx context.Context, req *remote.CallFAASReqDTO) (map[string]interface{}, error) {
	return t.abilityCaller.CallFAAS(ctx, req)
}

// CallHTTP 调用 http 能力
func (t *RemoteRepo) CallHTTP(ctx context.Context, req *remote.CallHTTPReqDTO) (map[string]interface{}, error) {
	return t.abilityCaller.CallHTTP(ctx, req)
}

// SendMsgToUser 发送消息给用户
func (t *RemoteRepo) SendMsgToUser(userID string, msg string) error {
	return t.chatOpsClient.SendMsgToUser(userID, msg)
}

// SendMsgToGroup 发送消息到群聊
func (t *RemoteRepo) SendMsgToGroup(chatID string, msg string) error {
	return t.chatOpsClient.SendMsgToGroup(chatID, msg)
}

// AddCronJob 添加定时任务
func (t *RemoteRepo) AddCronJob(addCronJobDTO *remote.AddCronJobDTO) error {
	jobName := fmt.Sprintf("%s:%d", constants.PxCronJobNamePrefix, addCronJobDTO.TriggerID)
	params, err := json.Marshal(addCronJobDTO.CronTriggerEvent)
	if err != nil {
		return err
	}

	addCronJobReqDTO := &remote.AddCronJobReqDTO{
		Name:    jobName,
		Params:  string(params),
		CronStr: addCronJobDTO.CronStr,
	}
	return t.cronClient.AddCronJob(addCronJobReqDTO)
}

// CancelCronJob  取消定时任务
func (t *RemoteRepo) CancelCronJob(jobName string) error {
	return t.cronClient.CancelCronJob(&remote.CancelCronJobReqDTO{Name: jobName})
}

// SendCloudEvent 发送事件
func (t *RemoteRepo) SendCloudEvent(ctx context.Context, req *remote.SendCloudEventDTO) error {
	return t.cloudEventClient.Send(ctx, req)
}
