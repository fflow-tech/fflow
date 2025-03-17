package po

import (
	"gorm.io/gorm"
)

// NamespacePO 命名空间
type NamespacePO struct {
	gorm.Model
	Namespace   string `gorm:"column:namespace;NOT NULL"`   // namespace
	Description string `gorm:"column:description;NOT NULL"` // description
	Creator     string `gorm:"column:creator;NOT NULL"`     // creator
}

// TableName 表名
func (m *NamespacePO) TableName() string {
	return "namespace"
}
