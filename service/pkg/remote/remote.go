// Package remote 远程调用相关工具包
package remote

import (
	"context"
)

// AbilityCaller 能力中心客户端
type AbilityCaller interface {
	CallFAAS(context.Context, *CallFAASReqDTO) (map[string]interface{}, error)
	CallHTTP(context.Context, *CallHTTPReqDTO) (map[string]interface{}, error)
}

// CronClient 分布式定时器客户端
type CronClient interface {
	AddCronJob(*AddCronJobReqDTO) error
	CancelCronJob(*CancelCronJobReqDTO) error
}

// ChatOpsClient ChatOps 客户端
type ChatOpsClient interface {
	SendMsgToUser(userID string, msg string) error
	SendMsgToGroup(chatID string, msg string) error
}

// CloudEventClient 事件客户端
type CloudEventClient interface {
	Send(ctx context.Context, req *SendCloudEventDTO) error
}

type DefaultCronClient struct{}

func NewDefaultCronClient() *DefaultCronClient {
	return &DefaultCronClient{}
}

func (*DefaultCronClient) AddCronJob(d *AddCronJobReqDTO) error {
	return nil
}

func (*DefaultCronClient) CancelCronJob(d *CancelCronJobReqDTO) error {
	return nil
}

type DefaultChatOpsClient struct{}

func NewDefaultChatOpsClient() *DefaultChatOpsClient {
	return &DefaultChatOpsClient{}
}

func (*DefaultChatOpsClient) SendMsgToUser(userID string, msg string) error {
	return nil
}

func (*DefaultChatOpsClient) SendMsgToGroup(chatID string, msg string) error {
	return nil
}

type DefaultCloudEventClient struct{}

func NewDefaultCloudEventClient() *DefaultCloudEventClient {
	return &DefaultCloudEventClient{}
}

func (*DefaultCloudEventClient) Send(ctx context.Context, req *SendCloudEventDTO) error {
	return nil
}
