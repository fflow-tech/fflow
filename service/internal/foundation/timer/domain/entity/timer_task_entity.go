package entity

import "time"

// RunHistory 执行历史实体
type RunHistory struct {
	ID        uint      `json:"id,omitempty"`        // ID
	DefID     string    `json:"def_id,omitempty"`    // 定义ID
	Name      string    `json:"name,omitempty"`      // 定时器名称
	Output    string    `json:"output,omitempty"`    // 执行结果
	RunTimer  string    `json:"run_timer,omitempty"` // 执行时间
	CostTime  int64     `json:"cost_time,omitempty"` // 执行耗时
	Status    string    `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"` // 更新时间
	CreatedAt time.Time `json:"created_at,omitempty"` // 创建时间
}

// RunStatus 运行状态
type RunStatus string

const (
	Running RunStatus = "running"
	Succeed RunStatus = "succeed"
	Failed  RunStatus = "failed"
	// Timeout 任务执行超时.
	Timeout RunStatus = "timeout"
)

// String 转换成string
func (r RunStatus) String() string {
	return string(r)
}
