package po

import (
	"time"

	"gorm.io/gorm"
)

// HistoryWorkflowInstPO 历史流程实例
type HistoryWorkflowInstPO struct {
	gorm.Model
	Namespace   string    `gorm:"column:namespace;NOT NULL" json:"namespace,omitempty"`     // 命名空间
	DefID       uint64    `gorm:"column:def_id;NOT NULL" json:"def_id,omitempty"`           // 主键ID
	DefVersion  int       `gorm:"column:def_version;NOT NULL" json:"def_version,omitempty"` // 流程的版本号
	Name        string    `gorm:"column:name;NOT NULL" json:"name,omitempty"`               // 流程实例名称
	Creator     string    `gorm:"column:creator;NOT NULL" json:"creator,omitempty"`         // 创建人
	Context     string    `gorm:"column:context;NOT NULL" json:"context,omitempty"`         // 流程实例上下文
	Status      int       `gorm:"column:status;NOT NULL" json:"status,omitempty"`           // 流程实例状态 1:running,2:completed,3:failed,4:cancelled,5:timeout
	StartAt     time.Time `gorm:"column:start_at;NOT NULL" json:"start_at"`                 // 流程执行开始时间
	CompletedAt time.Time `gorm:"column:completed_at" json:"completed_at"`                  // 流程执行结束时间
}

// TableName 对应表名
func (m *HistoryWorkflowInstPO) TableName() string {
	return "history_workflow_inst"
}
