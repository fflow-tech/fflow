package dto

import (
	"time"

	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// CreateRunHistoryDTO 创建执行历史
type CreateRunHistoryDTO struct {
	DefID    string `form:"def_id,omitempty"  json:"def_id,omitempty"`       // 定义ID
	Name     string `form:"name,omitempty"  json:"name,omitempty"`           // 定时器名称
	Output   string `form:"output,omitempty"  json:"output,omitempty"`       // 执行结果
	RunTimer string `form:"run_timer,omitempty"  json:"run_timer,omitempty"` // 执行时间
	CostTime int64  `form:"cost_time,omitempty"  json:"cost_time,omitempty"` // 执行耗时
	Status   string `form:"status,omitempty"  json:"status,omitempty"`       // 当前状态
}

// UpdateRunHistoryDTO 更新执行历史
type UpdateRunHistoryDTO struct {
	DefID    string `form:"def_id,omitempty"  json:"def_id,omitempty"`       // 定义ID
	RunTimer string `form:"run_timer,omitempty"  json:"run_timer,omitempty"` // 执行时间
	Output   string `form:"output,omitempty"  json:"output,omitempty"`       // 执行结果
	CostTime int64  `form:"cost_time,omitempty"  json:"cost_time,omitempty"` // 执行耗时，单位 ms
	Status   string `form:"status,omitempty"  json:"status,omitempty"`       // 当前状态
}

// GetRunHistoryDTO 获取执行历史
type GetRunHistoryDTO struct {
	DefID    string `form:"def_id,omitempty"  json:"def_id,omitempty"`       // 定义ID
	Name     string `form:"name,omitempty"  json:"name,omitempty"`           // 定时器名称
	RunTimer string `form:"run_timer,omitempty"  json:"run_timer,omitempty"` // 执行时间
}

// GetRunHistoryRspDTO 获取执行历史
type GetRunHistoryRspDTO struct {
	DefID     string    `form:"def_id,omitempty"  json:"def_id,omitempty"`         // 定义ID
	Name      string    `form:"name,omitempty"  json:"name,omitempty"`             // 定时器名称
	Output    string    `form:"output,omitempty"  json:"output,omitempty"`         // 执行结果
	RunTimer  string    `form:"run_timer,omitempty"  json:"run_timer,omitempty"`   // 执行时间
	CostTime  int64     `form:"cost_time,omitempty"  json:"cost_time,omitempty"`   // 执行耗时
	Status    string    `form:"status,omitempty"  json:"status,omitempty"`         // 当前状态             // 当前状态
	CreatedAT time.Time `form:"created_at,omitempty"  json:"created_at,omitempty"` //  创建时间
	UpdatedAt time.Time `form:"updated_at,omitempty"  json:"updated_at,omitempty"` // 完成时间
}

// DeleteRunHistoryDTO 删除执行历史
type DeleteRunHistoryDTO struct {
	DefID    string `form:"def_id,omitempty"  json:"def_id,omitempty"`       // 定义ID
	RunTimer string `form:"run_timer,omitempty"  json:"run_timer,omitempty"` // 执行时间
}

// PageQueryRunHistoryDTO 分页查询函数的请求体
type PageQueryRunHistoryDTO struct {
	DefID    string `form:"def_id,omitempty"  json:"def_id,omitempty" binding:"required"` // [必填] 定义ID
	Name     string `form:"name,omitempty"  json:"name,omitempty"`                        // 定时器名称
	RunTimer string `form:"run_timer,omitempty"  json:"run_timer,omitempty"`              // 执行时间
	*constants.PageQuery
	*constants.Order
}
