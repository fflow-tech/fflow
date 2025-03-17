package entity

import "time"

// Namespace 实体
type Namespace struct {
	ID        string    `json:"id,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
