package dto

// RbacReqDTO 权限请求
type RbacReqDTO struct {
	User        string   `json:"user,omitempty"`
	Role        string   `json:"role,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	Object      string   `json:"object,omitempty"`
	Permission  string   `json:"permission,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Users       []string `json:"users,omitempty"`
}
