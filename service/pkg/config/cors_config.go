package config

// CorsConfig 配置
type CorsConfig struct {
	AllowOrigins []string `json:"allowOrigins"` // 允许的 origin
	AllowHeaders []string `json:"allowHeaders"` // 允许的 request header
	AllowMethods []string `json:"allowMethods"` // 允许的 request method
}
