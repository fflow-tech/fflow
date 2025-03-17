package config

// RedisConfig 缓存配置
type RedisConfig struct {
	Network     string `json:"network"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	MaxIdle     int    `json:"maxIdle"`
	IdleTimeout int    `json:"idleTimeout"`
	// 连接池最大存活的连接数.
	MaxActive int `json:"maxActive"`
	// 当连接数达到上限时，新的请求是等待还是立即报错.
	Wait bool `json:"wait"`
}
