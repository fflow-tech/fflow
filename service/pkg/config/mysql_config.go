package config

// MySQLConfig 数据库配置
type MySQLConfig struct {
	Dsn                       string `json:"dsn"`
	SlaveDsn                  string `json:"slave_dsn"`
	SlowThreshold             int    `json:"slow_threshold"`                // 慢查询的阈值, 毫秒为单位
	IgnoreRecordNotFoundError bool   `json:"ignore_record_not_found_error"` // 是否忽略未找到记录的错误
	SkipDefaultTransaction    bool   `json:"skip_default_transaction"`
}
