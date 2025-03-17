package dto

import "time"

const (
	// TimerTaskTimeFormat 定时器任务时间格式
	TimerTaskTimeFormat = "2006-01-02 15:04"

	// TimerTriggerTimeFormat 定时器触发的时间格式 精度为秒
	TimerTriggerTimeFormat = "2006-01-02 15:04:05"
)

// AddTimerTaskDTO 增加定时器任务DTO
type AddTimerTaskDTO struct {
	BucketID  string    `json:"bucket_id,omitempty"`   // 桶名
	HashID    string    `json:"hash_id,omitempty"`     // hashID
	TimerTime time.Time `json:"timer_timer,omitempty"` // 触发时间
}

// SaveTimerTaskDTO 保存定时器任务DTO
type SaveTimerTaskDTO struct {
	BucketTimeID string `json:"bucket_time_id,omitempty"` // 桶分片ID
	TriggerTime  string `json:"trigger_time,omitempty"`   // 触发时间
	UnixTime     int64  `json:"unix_time,omitempty"`      // unix时间
}

// GetTimerTaskDTO 获取定时器任务DTO
type GetTimerTaskDTO struct {
	BucketTime string    `json:"bucket_time,omitempty"` // 桶名
	StartTime  time.Time `json:"start_time,omitempty"`  // 开始时间
	EndTime    time.Time `json:"end_time,omitempty"`    // 结束时间
}

// DelTimerTaskDTO 删除定时器任务DTO
type DelTimerTaskDTO struct {
	BucketTime string `json:"bucket_time,omitempty"` // 桶名
	HashID     string `json:"hash_id,omitempty"`     // hashID
}

// GetTimerTaskListDTO 获取定时器任务列表DTO
type GetTimerTaskListDTO struct {
	StartTime string `form:"start_time,omitempty" json:"start_time,omitempty"` // 开始时间 时间戳格式 "2021-01-01 01:01"
	EndTime   string `form:"end_time,omitempty" json:"end_time,omitempty"`     // 结束时间 时间戳格式 "2021-01-01 01:01"
}

// TimerListSendNotifyDTO 定时器列表触发DTO
type TimerListSendNotifyDTO struct {
	TimerList []string `json:"timer_list,omitempty"` // 定时器列表
}
