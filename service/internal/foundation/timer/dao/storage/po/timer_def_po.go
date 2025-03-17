package po

import "gorm.io/gorm"

// TimerDefPO 定时器定义
type TimerDefPO struct {
	gorm.Model
	DefID            string `gorm:"column:def_id;NOT NULL" json:"def_id,omitempty"`                         // 对应唯一hash值 这个ID会每次更新版本的时候更新
	App              string `gorm:"column:app;NOT NULL" json:"app,omitempty"`                               // APP 应用名
	Name             string `gorm:"column:name;NOT NULL" json:"name,omitempty"`                             // 定时器定义名称
	Creator          string `gorm:"column:creator;NOT NULL" json:"creator,omitempty"`                       // 创建人
	Status           int    `gorm:"column:status;NOT NULL" json:"status,omitempty"`                         // 定时器定义状态，1:未激活, 2:已激活
	Cron             string `gorm:"column:cron;NOT NULL" json:"cron,omitempty"`                             // 定时器定时配置
	NotifyType       int    `gorm:"column:notify_type;NOT NULL" json:"notify_type,omitempty"`               // 通知类型 1:rpc 2:kafka
	NotifyRpcParam   string `gorm:"column:notify_rpc_param;NOT NULL" json:"notify_rpc_param,omitempty"`     // 通知 Rpc 参数
	NotifyHttpParam  string `gorm:"column:notify_http_param;NOT NULL" json:"notify_http_param,omitempty"`   // Http 回调参数
	TimerType        int    `gorm:"column:timer_type;NOT NULL" json:"timer_type,omitempty"`                 // 定时器类型
	DelayTime        string `gorm:"column:delay_time;NOT NULL" json:"delay_time,omitempty"`                 // 延时触发时间
	EndTime          string `json:"end_time,omitempty"`                                                     // 定时器停止时间 格式为:"2006-01-02 15:04:05"
	TriggerType      int    `json:"trigger_type,omitempty"`                                                 // 触发类型 1-触发一次 2-持续触发
	DeleteType       int    `json:"delete_type,omitempty"`                                                  // 自动删除机制 0-不删除 1-触发后删除
	ExecuteTimeLimit int32  `gorm:"column:execute_time_limit;NOT NULL" json:"execute_time_limit,omitempty"` // 任务执行时间限制，单位: s.
}

// TableName 对应表名
func (m *TimerDefPO) TableName() string {
	return "timer_def"
}
