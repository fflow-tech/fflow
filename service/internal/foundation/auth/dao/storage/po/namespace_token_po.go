package po

import (
	"gorm.io/gorm"
	"time"
)

// NamespaceTokenPO 命名空间token
// 一个 namespace 可以有多个 token
type NamespaceTokenPO struct {
	gorm.Model
	Name      string    `gorm:"column:name;NOT NULL"`      // name
	Namespace string    `gorm:"column:namespace;NOT NULL"` // namespace
	Token     string    `gorm:"column:token;NOT NULL"`     // token校验
	Creator   string    `gorm:"column:creator;NOT NULL"`   // creator
	ExpiredAt time.Time `gorm:"column:expired_at"`         // 失效时间
}

// TableName 表名
func (m *NamespaceTokenPO) TableName() string {
	return "namespace_token"
}
