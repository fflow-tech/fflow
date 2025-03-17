package config

// RbacConfig 权限配置
type RbacConfig struct {
	SuperAdmins []string `json:"superAdmins"` // 超管人员
}
