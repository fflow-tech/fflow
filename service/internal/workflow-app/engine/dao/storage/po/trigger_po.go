package po

import (
	"gorm.io/gorm"
)

// TriggerPO 触发器
type TriggerPO struct {
	gorm.Model
	Namespace  string `gorm:"column:namespace;NOT NULL" json:"namespace,omitempty"` // 命名空间
	Type       string `gorm:"column:type;NOT NULL"`                                 // 触发器类型
	Event      string `gorm:"column:event;NOT NULL"`                                // 事件名称
	Expr       string `gorm:"column:expr;NOT NULL"`                                 // 定时触发器的时间表达式
	Attribute  string `gorm:"column:attribute;NOT NULL"`                            // 触发器属性
	Level      int    `gorm:"column:level;NOT NULL"`                                // 触发器类型 1:流程级别触发器, 2:流程实例级别触发器, 3:流程实例节点级别触发器
	DefID      uint64 `gorm:"column:def_id;NOT NULL"`                               // 主键ID
	DefVersion int    `gorm:"column:def_version;NOT NULL"`                          // 流程的版本号
	InstID     uint64 `gorm:"column:inst_id"`                                       // 流程实例ID
	Status     int    `gorm:"column:status;NOT NULL"`                               // 触发器状态，1:未激活, 2:已激活
}

// TableName 触发器表名
func (m *TriggerPO) TableName() string {
	return "trigger"
}
