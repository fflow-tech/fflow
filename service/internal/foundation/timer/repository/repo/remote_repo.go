package repo

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/remote"
)

// RemoteRepo 调用能力实体
type RemoteRepo struct {
	abilityCaller remote.AbilityCaller
	chatOpsClient remote.ChatOpsClient
}

// NewRemoteRepo 实体构造函数
func NewRemoteRepo(d *remote.DefaultAbilityCaller, c *remote.DefaultChatOpsClient) *RemoteRepo {
	return &RemoteRepo{abilityCaller: d, chatOpsClient: c}
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
