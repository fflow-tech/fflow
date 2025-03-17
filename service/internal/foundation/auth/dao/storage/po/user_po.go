package po

import (
	"database/sql"
	"gorm.io/gorm"
)

// UserPO 用户表
type UserPO struct {
	gorm.Model
	Username string         `gorm:"column:username;NOT NULL"`  // 用户名
	NickName string         `gorm:"column:nick_name;NOT NULL"` // 用户昵称
	AuthType string         `gorm:"column:auth_type;NOT NULL"` // 认证类型
	Password string         `gorm:"column:password"`           // 用户名
	Email    sql.NullString `gorm:"column:email"`              // 邮箱
	Phone    sql.NullString `gorm:"column:phone"`              // 手机号
	Avatar   string         `gorm:"column:avatar"`             // 头像路径
	Status   int            `gorm:"column:status;default:1"`   // 用户状态，1:未激活, 2:已激活
}

// TableName 表名
func (m *UserPO) TableName() string {
	return "user"
}
