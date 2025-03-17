package entity

import "time"

// NamespaceToken 实体
type NamespaceToken struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	Token     string    `json:"token,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
