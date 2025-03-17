// Package dto 	服务各层方法参数体定义。
package dto

import "github.com/fflow-tech/fflow/service/pkg/constants"

// CreateTimerDefDTO 创建定时器定义DTO
type CreateTimerDefDTO struct {
	DefID            string          `json:"def_id,omitempty"`                         // 主键ID
	Name             string          `json:"name,omitempty" binding:"required"`        // [必填] 定时器名称
	App              string          `json:"app,omitempty" binding:"required"`         // [必填] APP 应用名
	Creator          string          `json:"creator,omitempty" binding:"required"`     // [必填] 创建人
	Status           int             `json:"status,omitempty"`                         // 定时器定义状态，1:激活, 2:未激活
	Cron             string          `json:"cron,omitempty"`                           // 定时器定时配置
	NotifyType       int             `json:"notify_type,omitempty" binding:"required"` // [必填] 通知类型 1:rpc 2:kafka
	NotifyRpcParam   NotifyRpcParam  `json:"notify_rpc_param,omitempty"`               // Rpc  回调参数
	NotifyHttpParam  NotifyHttpParam `json:"notify_http_param,omitempty"`              // Http 回调参数
	TimerType        int             `json:"timer_type,omitempty"  binding:"required"` // [必填] 定时器类型 1：延时定时器 2：cron定时器
	DelayTime        string          `json:"delay_time,omitempty"`                     // 延时定时器触发时间 格式为:"2006-01-02 15:04:05"
	EndTime          string          `json:"end_time,omitempty"`                       // 定时器停止时间 格式为:"2006-01-02 15:04:05"
	TriggerType      int             `json:"trigger_type,omitempty"`                   // 触发类型 1-触发一次 2-持续触发
	DeleteType       int             `json:"delete_type,omitempty"`                    // 自动删除机制 0-不删除 1-删除
	ExecuteTimeLimit int32           `json:"execute_time_limit,omitempty"`             // 任务单次执行时间限制，单位：s. 默认 15 s.
}

// DeleteTimerDefDTO 删除定时器定义DTO
type DeleteTimerDefDTO struct {
	DefID string `json:"def_id,omitempty" form:"def_id,omitempty"`
	App   string `json:"app,omitempty"`
	Name  string `json:"name,omitempty"`
}

// HasAppAndName 存在应用和名称信息.
func (d *DeleteTimerDefDTO) HasAppAndName() bool {
	return d.App != "" && d.Name != ""
}

// ChangeTimerStatusDTO 修改定时器定义状态DTO
type ChangeTimerStatusDTO struct {
	DefID  string `json:"def_id,omitempty"  binding:"required"` // [必填] 主键ID
	Status int    `json:"status,omitempty"  binding:"required"` // [必填] 定时器定义状态 1:激活, 2:去激活
}

// GetTimerDefDTO 获取定时器定义DTO
type GetTimerDefDTO struct {
	DefID string `json:"def_id,omitempty" form:"def_id,omitempty"  binding:"required"` // [必填] 主键ID
}

// TimerDefDTO 定时器定义DTO
type TimerDefDTO struct {
	DefID           string          `json:"def_id,omitempty"`            // 主键ID
	App             string          `json:"app,omitempty"`               // APP 应用名
	Name            string          `json:"name,omitempty"`              // 定时器名称
	Creator         string          `json:"creator,omitempty"`           // 创建人
	Status          int             `json:"status,omitempty"`            // 定时器定义状态，1:未激活, 2:已激活
	Cron            string          `json:"cron,omitempty"`              // 定时器定时配置
	NotifyType      int             `json:"notify_type,omitempty"`       // 通知类型 1:rpc 2:kafka
	TimerType       int             `json:"timer_type,omitempty"`        // 定时器类型 1：延时定时器 2：cron定时器
	DelayTime       string          `json:"delay_time,omitempty"`        // 延时定时器触发时间 格式为:"2006-01-02 15:04:05"
	NotifyRpcParam  NotifyRpcParam  `json:"notify_rpc_param,omitempty"`  // 通知 Rpc 参数
	NotifyHttpParam NotifyHttpParam `json:"notify_http_param,omitempty"` // Http 回调参数
	EndTime         string          `json:"end_time,omitempty"`          // 定时器停止时间 格式为:"2006-01-02 15:04:05"
	TriggerType     int             `json:"trigger_type,omitempty"`      // 触发类型 1-触发一次 2-持续触发
	DeleteType      int             `json:"delete_type,omitempty"`       // 自动删除机制 0-不删除 1-触发后删除
}

// NotifyRpcParam RPC 通知配置参数
type NotifyRpcParam struct {
	Service   string `json:"service,omitempty"`    // 服务名:对应 123 平台上 service.name
	Method    string `json:"method,omitempty"`     // 回调方法名
	RpcName   string `json:"rpc_name,omitempty"`   // 对应 method 别名，优先使用 RpcName 寻址
	Params    string `json:"params,omitempty"`     // 回调参数
	CalleeEnv string `json:"callee_env,omitempty"` // callee 被调服务环境
}

// NotifyHttpParam http 通知参数
type NotifyHttpParam struct {
	Method string `json:"method,omitempty" metakey:"method" ` // POST,GET 方法
	Url    string `json:"url,omitempty" metakey:"url"`        // URL 路径
	Header string `json:"header,omitempty" metakey:"header"`  // header 请求头
	Body   string `json:"body,omitempty"`                     // 参数体
}

// TimerAppNameDTO 定时器应用与名称信息.
type TimerAppNameDTO struct {
	App  string `json:"app"`
	Name string `json:"name"`
}

// PageQueryTimeDefDTO 分页查询定时器定义参数体
type PageQueryTimeDefDTO struct {
	App     string `json:"app,omitempty" form:"app,omitempty" binding:"required"` // [必填] APP 应用名
	Name    string `json:"name,omitempty" form:"name,omitempty" `                 // 定时器名称
	Creator string `json:"creator,omitempty" form:"creator,omitempty"`            // 创建人
	constants.PageQuery
	constants.Order
}

// CountTimerDefDTO 查询定时器定义总数
type CountTimerDefDTO struct {
	App     string `json:"app,omitempty"`  // APP 应用名
	Name    string `json:"name,omitempty"` // 定时器名称
	Creator string `json:",omitempty"`     // 创建人
}

// UpdateTimerDefDTO 更新定时器状态参数
type UpdateTimerDefDTO struct {
	DefID  string `json:"def_id,omitempty"  binding:"required"` // [必填] 主键ID
	Status int    `json:"status,omitempty"  binding:"required"` // [必填] 定时器定义状态，1:未激活, 2:已激活
}
