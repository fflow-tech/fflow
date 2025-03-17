package registry

// ProviderType 注册中心客户端类型
type ProviderType string

const (
	Consul     ProviderType = "consul"
	Kubernetes ProviderType = "kubernetes"
)

// Provider 注册服务客户端
type Provider interface {
	// Register 注册服务
	Register(service, addr string) error
	// GetConnStr 连接串
	GetConnStr(target string) (string, error)
}
