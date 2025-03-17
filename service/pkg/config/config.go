package config

import "context"

// ProviderType 配置中心客户端类型
type ProviderType string

const (
	Consul     ProviderType = "consul"
	Kubernetes ProviderType = "kubernetes"
)

// Key 把 Group 和 Key 包含在一起
type Key struct {
	Group string
	Key   string
}

// NewGroupKey 新建配置 key
func NewGroupKey(group, key string) Key {
	return Key{group, key}
}

// Provider 对配置获取的一个简单的封装
type Provider interface {
	// GetAny 将配置序列化到一个结构体里面
	GetAny(ctx context.Context, k Key, t interface{}) error
	// GetString 获取一个字符串
	GetString(ctx context.Context, k Key) (string, error)
}
