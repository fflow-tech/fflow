package config

// LimiterConfig 限流器配置
type LimiterConfig struct {
	// 限流器大小.
	Burst int `json:"burst"`
	// 流量补充速率，含义为每秒补充 limit 单位流量.
	Limit int `json:"limit"`
	// 获取流量时的最大等待时长，单位: s.
	WaitingDuration int32 `json:"waitingDuration"`
	// 自刷新时间间隔，单位: s.
	RefreshInterval int32 `json:"refreshInterval"`
}
