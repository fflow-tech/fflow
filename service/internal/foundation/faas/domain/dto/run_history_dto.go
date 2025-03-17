package dto

import (
	"time"

	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// CreateRunHistoryDTO 创建执行历史
type CreateRunHistoryDTO struct {
	Namespace    string `form:"namespace,omitempty" json:"namespace,omitempty"`          // 服务接口
	Version      uint   `form:"version,omitempty"  json:"version,omitempty"`             // 版本号
	FunctionName string `form:"function_name,omitempty"  json:"function_name,omitempty"` // 函数名称
	Operator     string `form:"operator,omitempty"  json:"operator,omitempty"`           // 执行者
	Input        string `form:"input,omitempty"  json:"input,omitempty"`                 // 函数入参
	Output       string `form:"output,omitempty"  json:"output,omitempty"`               // 执行结果
	Log          string `form:"log,omitempty"  json:"log,omitempty"`                     // 执行日志
	CostTime     int64  `form:"cost_time,omitempty"  json:"cost_time,omitempty"`         // 执行耗时
	Status       string `form:"status,omitempty"  json:"status,omitempty"`               // 当前状态
}

// UpdateRunHistoryDTO 更新执行历史
type UpdateRunHistoryDTO struct {
	ID       uint   `form:"id,omitempty"  json:"id,omitempty"`               // 执行ID
	Output   string `form:"output,omitempty"  json:"output,omitempty"`       // 执行结果
	Log      string `form:"log,omitempty"  json:"log,omitempty"`             // 执行日志
	CostTime int64  `form:"cost_time,omitempty"  json:"cost_time,omitempty"` // 执行耗时，单位 ms
	Status   string `form:"status,omitempty"  json:"status,omitempty"`       // 当前状态
}

// BatchDeleteExpiredRunHistoryDTO 批量删除过期的执行历史
type BatchDeleteExpiredRunHistoryDTO struct {
	IDs      []uint `form:"ids,omitempty"  json:"ids,omitempty"`             // 执行 ID 列表
	KeepDays int    `form:"keep_days,omitempty"  json:"keep_days,omitempty"` // 保留天数
}

// GetRunHistoryDTO 获取执行历史
type GetRunHistoryDTO struct {
	ID           uint   `form:"id,omitempty"  json:"id,omitempty"`                       // 执行ID
	Namespace    string `form:"namespace,omitempty" json:"namespace,omitempty"`          // 服务接口
	Version      uint   `form:"version,omitempty"  json:"version,omitempty"`             // 版本号
	FunctionName string `form:"function_name,omitempty"  json:"function_name,omitempty"` // 函数名称
	Operator     string `form:"operator,omitempty"  json:"operator,omitempty"`           // 执行者
}

// GetRunHistoryRspDTO 获取执行历史
type GetRunHistoryRspDTO struct {
	ID           uint      `form:"id,omitempty"  json:"id,omitempty"`                       // 执行ID
	Namespace    string    `form:"namespace,omitempty" json:"namespace,omitempty"`          // 服务接口
	Version      uint      `form:"version,omitempty"  json:"version,omitempty"`             // 版本号
	FunctionName string    `form:"function_name,omitempty"  json:"function_name,omitempty"` // 函数名称
	Operator     string    `form:"operator,omitempty"  json:"operator,omitempty"`           // 执行者
	Input        string    `form:"input,omitempty"  json:"input,omitempty"`                 // 函数入参
	Output       string    `form:"output,omitempty"  json:"output,omitempty"`               // 执行结果
	Log          string    `form:"log,omitempty"  json:"log,omitempty"`                     // 执行日志
	CostTime     int64     `form:"cost_time,omitempty"  json:"cost_time,omitempty"`         // 执行耗时
	Status       string    `form:"status,omitempty"  json:"status,omitempty"`               // 当前状态
	CreatedAt    time.Time `form:"created_at,omitempty"  json:"created_at,omitempty"`       //  创建时间
	UpdatedAt    time.Time `form:"updated_at,omitempty"  json:"updated_at,omitempty"`       // 完成时间
}

// DeleteRunHistoryDTO 删除执行历史
type DeleteRunHistoryDTO struct {
	ID           uint   `form:"id,omitempty"  json:"id,omitempty"`                       // 执行ID
	Namespace    string `form:"namespace,omitempty" json:"namespace,omitempty"`          // 服务接口
	Version      uint   `form:"version,omitempty"  json:"version,omitempty"`             // 版本号
	FunctionName string `form:"function_name,omitempty"  json:"function_name,omitempty"` // 函数名称
}

// PageQueryRunHistoryDTO 分页查询函数的请求体
type PageQueryRunHistoryDTO struct {
	IDs          []uint    `form:"ids,omitempty"  json:"ids,omitempty"`                     // 执行 ID 列表
	MaxID        uint      `form:"max_id,omitempty"  json:"max_id,omitempty"`               // 查询的最大 ID
	Namespace    string    `form:"namespace,omitempty" json:"namespace,omitempty"`          // 命名空间
	Version      int       `form:"version,omitempty"  json:"version,omitempty"`             // 版本号
	FunctionName string    `form:"function_name,omitempty"  json:"function_name,omitempty"` // 函数名称
	CreatedAt    time.Time `form:"created_at,omitempty"  json:"created_at,omitempty"`       // 创建时间
	*constants.PageQuery
	*constants.Order
}

// BatchDeleteRunHistoryDTO 批量删除执行历史
type BatchDeleteRunHistoryDTO struct {
	IDs       []uint `form:"ids,omitempty"  json:"ids,omitempty"`            // 执行 ID 列表
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"` // 服务接口
}
