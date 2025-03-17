package entity

import "time"

// App 应用实体
type App struct {
	ID        int       `json:"id,omitempty"`         // 应用 ID
	Name      string    `json:"name,omitempty"`       // 应用名
	Creator   string    `json:"creator,omitempty"`    // 创建人
	CreatedAt time.Time `json:"created_at,omitempty"` // 创建时间
	UpdatedAt time.Time `json:"updated_at,omitempty"` // 更新时间
}
