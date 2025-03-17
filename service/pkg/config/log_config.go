package config

// RequestRecordConfig 请求记录配置
type RequestRecordConfig struct {
	Enable             bool     `json:"enable"`
	EnableWeb          bool     `json:"enableWeb"`
	UriPrefixWhiteList []string `json:"uriPrefixWhiteList"` // http uri 前缀白名单
	UriPrefixBlackList []string `json:"uriPrefixBlackList"` // http uri 前缀黑名单
}
