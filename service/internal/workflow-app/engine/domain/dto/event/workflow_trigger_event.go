package event

// CronTriggerEvent 定时器触发器事件
type CronTriggerEvent struct {
	Namespace  string `json:"namespace,omitempty"`
	TriggerID  string `json:"trigger_id"`  // 触发器唯一标识
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion int    `json:"def_version"` // 流程定义版本号
}
