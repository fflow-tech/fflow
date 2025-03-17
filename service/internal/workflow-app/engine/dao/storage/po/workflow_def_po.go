package po

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// WorkflowDefPO 流程定义
type WorkflowDefPO struct {
	gorm.Model
	Namespace   string          `gorm:"column:namespace;NOT NULL" json:"namespace,omitempty"`           // 命名空间
	DefID       uint64          `gorm:"column:def_id;NOT NULL" json:"def_id,omitempty"`                 // 主键 ID
	ParentDefID uint64          `gorm:"column:parent_def_id;NOT NULL" json:"parent_def_id,omitempty"`   // 父流程定义 ID
	Attribute   WorkflowDefAttr `gorm:"column:attribute;type:json;NOT NULL" json:"attribute,omitempty"` // 流程定义的其他属性
	Version     int             `gorm:"column:version;NOT NULL" json:"version,omitempty"`               // 流程的版本号
	Name        string          `gorm:"column:name;NOT NULL" json:"name,omitempty"`                     // 流程定义名称
	DefJson     string          `gorm:"column:def_json;NOT NULL" json:"def_json,omitempty"`             // 流程定义的内容
	Creator     string          `gorm:"column:creator;NOT NULL" json:"creator,omitempty"`               // 创建人
	Status      int             `gorm:"column:status;NOT NULL" json:"status,omitempty"`                 // 流程定义状态，1:未激活, 2:已激活
	Description string          `gorm:"column:description" json:"description,omitempty"`                // 流程定义描述
}

// WorkflowDefAttr 流程额外属性
type WorkflowDefAttr struct {
	RefName          string `json:"ref_name,omitempty"`           // 子流程对应的 RefName
	ParentDefVersion int    `json:"parent_def_version,omitempty"` // 子流程对应的父流程的版本号
}

// Value 实现方法
func (w WorkflowDefAttr) Value() (driver.Value, error) {
	return json.Marshal(w)
}

// Scan 实现方法
func (w *WorkflowDefAttr) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), w)
}

// TableName 对应表名
func (m *WorkflowDefPO) TableName() string {
	return "workflow_def"
}
