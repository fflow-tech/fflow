package po

import (
	"time"
)

// App 应用定义
type App struct {
	ID        int       `gorm:"column:id;primary_key"`      // 主键ID
	Name      string    `gorm:"column:name;NOT NULL"`       // 应用名
	Creator   string    `gorm:"column:creator;NOT NULL"`    // 创建人
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at"`          // 更新时间
	DeletedAt time.Time `gorm:"column:deleted_at;NOT NULL"` // 删除时间
}

// TableName APP 对应数据库表名
func (a *App) TableName() string {
	return "app"
}
