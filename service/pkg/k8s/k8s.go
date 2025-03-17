package k8s

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	GlobalConfigName string
	GlobalConfigType string
	GlobalConfigPath string
}

// Client 客户端
type Client struct {
}

// NewClient 创建一个新的 Client 对象
func NewClient(config Config) (*Client, error) {
	viper.SetConfigName(config.GlobalConfigName)
	viper.SetConfigType(config.GlobalConfigType)
	viper.AddConfigPath(config.GlobalConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.WatchConfig()

	return &Client{}, nil
}

// GetAny 读取配置值
func (r *Client) GetAny(ctx context.Context, k config.Key, t interface{}) error {
	return utils.ToOtherInterfaceValue(t, viper.GetStringMap(k.Key))
}

// GetString 获取字符串
func (r *Client) GetString(ctx context.Context, k config.Key) (string, error) {
	return viper.GetString(k.Key), nil
}

// Register 注册服务
func (r *Client) Register(service, addr string) error {
	return nil
}

// GetConnStr 获取连接串
func (r *Client) GetConnStr(target string) (string, error) {
	return fmt.Sprintf("kubernetes:///%s", target), nil
}
