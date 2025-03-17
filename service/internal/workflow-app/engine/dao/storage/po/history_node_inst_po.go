package po

import (
	"time"

	"gorm.io/gorm"
)

// HistoryNodeInstPO 历史节点实例
type HistoryNodeInstPO struct {
	gorm.Model
	Namespace     string    `gorm:"column:namespace;NOT NULL" json:"namespace,omitempty"`     // 命名空间
	DefID         uint64    `gorm:"column:def_id;NOT NULL" json:"def_id,omitempty"`           // 主键ID
	DefVersion    int       `gorm:"column:def_version;NOT NULL" json:"def_version,omitempty"` // 流程的版本号
	InstID        uint64    `gorm:"column:inst_id;NOT NULL" json:"inst_id,omitempty"`         // 流程实例ID
	RefName       string    `gorm:"column:ref_name;NOT NULL" json:"ref_name,omitempty"`       // 节点引用名称
	Context       string    `gorm:"column:context;NOT NULL" json:"context,omitempty"`         // 节点实例上下文
	Status        int       `gorm:"column:status;NOT NULL" json:"status,omitempty"`           // 节点实例状态 1:scheduled,2:waiting,3:running,4:completed,5:failed,6:cancelled,7:timeout
	ScheduledAt   time.Time `gorm:"column:scheduled_at;NOT NULL" json:"scheduled_at"`         // 节点开始调度时间
	WaitAt        time.Time `gorm:"column:wait_at" json:"wait_at"`                            // 节点等待开始时间
	ExecuteAt     time.Time `gorm:"column:execute_at" json:"execute_at"`                      // 节点执行开始时间
	AsynWaitResAt time.Time `gorm:"column:asyn_wait_res_at" json:"asyn_wait_res_at"`          // 异步等待结果开始时间
	CompletedAt   time.Time `gorm:"column:completed_at" json:"completed_at"`                  // 节点执行结束时间
}

// TableName 对应表名
func (m *HistoryNodeInstPO) TableName() string {
	return "history_node_inst"
}
