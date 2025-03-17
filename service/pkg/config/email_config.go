package config

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	From     string `json:"from,omitempty"`
	Password string `json:"password,omitempty"`
}
