package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var ttl = 30 * time.Second

// Config 基本的配置项
type Config struct {
	Env      string
	Address  string
	User     string
	Password string
}

// Client 客户端
type Client struct {
	config    Config
	apiClient *api.Client
}

// NewClient 创建一个新的 Client 对象
func NewClient(config Config) (*Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = config.Address
	registry, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{config: config, apiClient: registry}, nil
}

// GetAny 读取配置值
func (r *Client) GetAny(ctx context.Context, k config.Key, t interface{}) error {
	key := fmt.Sprintf("%s/%s/%s", r.config.Env, k.Group, k.Key)
	kv, _, err := r.apiClient.KV().Get(key, &api.QueryOptions{UseCache: true})
	if err != nil {
		return err
	}

	if kv == nil {
		return fmt.Errorf("get config %v value failed", k)
	}

	return json.Unmarshal(kv.Value, t)
}

// GetString 获取字符串
func (r *Client) GetString(ctx context.Context, k config.Key) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", r.config.Env, k.Group, k.Key)
	kv, _, err := r.apiClient.KV().Get(key, &api.QueryOptions{UseCache: true})
	if err != nil {
		return "", err
	}

	if kv == nil {
		return "", fmt.Errorf("get config %v value failed", k)
	}

	return string(kv.Value), nil
}

// Register 注册服务
func (r *Client) Register(service, addr string) error {
	ip, err := utils.GetOutboundIP()
	if err != nil {
		return err
	}

	port, err := utils.StrToInt(addr[1:])
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s-%s%s", service, ip, addr)

	// 注册到 Consul，包含地址、端口信息，以及健康检查
	if err = r.apiClient.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      id,
		Name:    service,
		Port:    port,
		Address: ip,
		Check: &api.AgentServiceCheck{
			TTL:     (ttl + time.Second).String(),
			Timeout: time.Minute.String(),
		},
	}); err != nil {
		return err
	}
	go func() {
		for range time.Tick(ttl) {
			if err := r.apiClient.Agent().PassTTL("service:"+id, ""); err != nil {
				log.Errorf("PassTTL failed: %v", err)
			}
		}
	}()

	return nil
}

// GetConnStr 获取连接串
func (r *Client) GetConnStr(target string) (string, error) {
	if r.config.User == "" || r.config.Password == "" {
		return fmt.Sprintf("consul://%s/%s", r.config.Address, target), nil
	}

	return fmt.Sprintf("consul://%s:%s@%s/%s", r.config.User, r.config.Password, r.config.Address, target), nil
}
