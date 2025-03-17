// Package entity 业务领域模型定义。
package entity

// TimerDef 定时器定义
type TimerDef struct {
	Name             string          `json:"name,omitempty"`               // 名字
	App              string          `json:"app,omitempty"`                // APP 应用名
	DefID            string          `json:"def_id,omitempty"`             // 对应唯一hash值 这个ID会每次更新版本的时候更新
	NotifyType       TimerNotifyType `json:"notify_type,omitempty"`        // 通知类型 rpc / kafka
	Cron             string          `json:"cron,omitempty"`               // 定时器设置格式
	Creator          string          `json:"creator,omitempty"`            // 用户名称 每个定时器都有归属
	Status           TimerDefStatus  `json:"status,omitempty"`             // 状态 激活/未激活
	Version          uint64          `json:"version,omitempty"`            // 版本号
	TimerType        TimerType       `json:"timer_type,omitempty"`         // 定时器类型
	DelayTime        string          `json:"delay_time,omitempty"`         // 延时触发时间
	NotifyRpcParam   string          `json:"notify_rpc_param,omitempty"`   // 通知 Rpc 参数
	NotifyHttpParam  string          `json:"notify_http_param,omitempty"`  // Http 回调参数
	EndTime          string          `json:"end_time,omitempty"`           // 定时器停止时间 格式为:"2006-01-02 15:04:05"
	TriggerType      TriggerType     `json:"trigger_type,omitempty"`       // 触发类型 1-触发一次 2-持续触发
	DeleteType       DeleteType      `json:"delete_type,omitempty"`        // 自动删除机制 0-不删除 1-触发后删除
	ExecuteTimeLimit int32           `json:"execute_time_limit,omitempty"` // 定时任务单次执行时间限制，单位：s. 默认 15s.
}

// TriggerType 触发类型
type TriggerType int

const (
	TriggerOnce TriggerType = 1 // 触发一次
	TriggerMany TriggerType = 2 // 触发多次
)

// ToInt 转成数字
func (t TriggerType) ToInt() int {
	return int(t)
}

// TimerType 定时器类型
type TimerType int

const (
	DelayTimer      TimerType = 1
	CronTimer       TimerType = 2
	DelayTimeFormat           = "2006-01-02 15:04:05"
)

// ToInt 转成数字
func (t TimerType) ToInt() int {
	return int(t)
}

// TimerNotifyType 定时器通知类型
type TimerNotifyType int

// ToInt 转成字符串
func (t TimerNotifyType) ToInt() int {
	return int(t)
}

// 定时器通知类型
const (
	RPC   TimerNotifyType = 1
	KAFKA TimerNotifyType = 2
	HTTP  TimerNotifyType = 3
)

// TimerDefStatus 定时器定义状态
type TimerDefStatus int

// ToInt 转成数字
func (t TimerDefStatus) ToInt() int {
	return int(t)
}

// 定时器状态
const (
	Enabled  TimerDefStatus = 1
	Disabled TimerDefStatus = 2
)

// DeleteType 删除类型
type DeleteType int

const (
	NotDelete     DeleteType = 0 // 不删除
	TriggerDelete DeleteType = 1 // 触发后删除
)

// ToInt 转成数字
func (t DeleteType) ToInt() int {
	return int(t)
}
