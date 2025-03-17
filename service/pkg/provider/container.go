// Package provider 提供的给领域层公用的依赖，全局统一注入，所以是领域层需要统一依赖的组件，一般为 pkg 包里面定义的接口和实现
package provider

import (
	"sync"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/registry"
)

// Container 容器
type Container struct {
	registryProvider registry.Provider
	configProvider   config.Provider
	rwMutex          sync.RWMutex
}

var (
	container = &Container{}
)

// InjectConfigProvider 注入
func InjectConfigProvider(configProvider config.Provider) {
	container.rwMutex.Lock()
	defer container.rwMutex.Unlock()
	container.configProvider = configProvider
}

// InjectRegistryProvider 注入
func InjectRegistryProvider(registryProvider registry.Provider) {
	container.rwMutex.Lock()
	defer container.rwMutex.Unlock()
	container.registryProvider = registryProvider
}

// GetConfigProvider 获取配置客户端
func GetConfigProvider() config.Provider {
	container.rwMutex.Lock()
	defer container.rwMutex.Unlock()
	return container.configProvider
}

// GetRegistryProvider 获取注册中心客户端
func GetRegistryProvider() registry.Provider {
	container.rwMutex.RLock()
	defer container.rwMutex.RUnlock()
	return container.registryProvider
}
