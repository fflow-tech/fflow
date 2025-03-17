package config

// AuthConfig 权限配置
type AuthConfig struct {
	Domain        string `json:"domain,omitempty'"`
	SecretKey     string `json:"secretKey,omitempty"`
	HomePage      string `json:"homePage,omitempty"`
	AdminEmail    string `json:"adminEmail,omitempty'"`
	AdminUsername string `json:"adminUsername,omitempty'"`
	AdminPassword string `json:"adminPassword,omitempty"`
}
