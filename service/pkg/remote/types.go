package remote

import "github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"

// CallRPCReqDTO RPC 请求配置和请求体
type CallRPCReqDTO struct {
	Version   string                 `json:"version" metakey:"version"`
	Source    string                 `json:"source" metakey:"source"`
	SourceKey string                 `json:"sourceKey" metakey:"sourceKey"`
	MockMode  bool                   `json:"mockMode" metakey:"mockMode"`
	Service   string                 `json:"service" metakey:"service"`
	Method    string                 `json:"method" metakey:"method"`
	RPCName   string                 `json:"rpcName" metakey:"rpcName"`
	Protocol  string                 `json:"protocol" metakey:"protocol"`
	CalleeEnv string                 `json:"calleeEnv" metakey:"calleeEnv"`
	Target    string                 `json:"target" metakey:"target"`
	Network   string                 `json:"network" metakey:"network"`
	Body      map[string]interface{} `json:"body,omitempty"`
}

// CallFAASReqDTO FAAS 请求配置和请求体
type CallFAASReqDTO struct {
	MockMode  bool                   `json:"mockMode" metakey:"mockMode"`
	Namespace string                 `json:"namespace" metakey:"namespace"`
	Function  string                 `json:"function" metakey:"function"`
	Body      map[string]interface{} `json:"body,omitempty"`
}

// CallHTTPReqDTO HTTP 请求配置和请求体
type CallHTTPReqDTO struct {
	MockMode bool                   `json:"mockMode" metakey:"mockMode"`
	Method   string                 `json:"method" metakey:"method"`
	URL      string                 `json:"url" metakey:"url"`
	Header   map[string]string      `json:"header,omitempty" metakey:"header"`
	Query    map[string]string      `json:"query,omitempty" metakey:"query"`
	Body     map[string]interface{} `json:"body,omitempty"`
}

// ValidateTokenReqDTO 请求配置和请求体
type ValidateTokenReqDTO struct {
	Namespace   string `json:"namespace" metakey:"namespace"`
	AccessToken string `json:"accessToken" metakey:"accessToken"`
}

// AddCronJobReqDTO 创建定时任务请求体
type AddCronJobReqDTO struct {
	Name    string // job name，需保证唯一
	Params  string // rpc 回调参数
	CronStr string // cron表达式，unix 时间戳，秒级
}

// CancelCronJobReqDTO 取消定时任务请求体
type CancelCronJobReqDTO struct {
	Name string // job name，需保证唯一
}

// AddCronJobDTO 创建定时任务DTO
type AddCronJobDTO struct {
	CronStr string // cron表达式，unix 时间戳，秒级
	event.CronTriggerEvent
}

// SendCloudEventDTO 发送事件对象
type SendCloudEventDTO struct {
	Target    string      `json:"target,omitempty"`
	Event     *CloudEvent `json:"event,omitempty"`
	AppID     string      `json:"appID,omitempty"`
	EventType string      `json:"eventType,omitempty"`
}

// CloudEvent 事件定义，遵循 Cloud Event 协议
type CloudEvent struct {
	Source          string                 `json:"source,omitempty"`
	Type            string                 `json:"type,omitempty"`
	ID              string                 `json:"id,omitempty"`
	DataContentType string                 `json:"dataContentType,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
}
