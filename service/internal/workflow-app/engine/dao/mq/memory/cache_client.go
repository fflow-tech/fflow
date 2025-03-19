package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/cache"
)

// CacheDAO 内存缓存实现
type CacheDAO struct {
	data  map[string]cacheItem
	mutex sync.RWMutex
}

type cacheItem struct {
	value      string
	expireTime time.Time
}

// NewCacheDAO 创建新的内存缓存实例
func NewCacheDAO() *CacheDAO {
	return &CacheDAO{
		data: make(map[string]cacheItem),
	}
}

// Set 设置缓存
func (m *CacheDAO) Set(key string, value string, ttl int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = cacheItem{
		value:      value,
		expireTime: time.Now().Add(time.Duration(ttl) * time.Second),
	}
	return nil
}

// SetNX 如果不存在则设置
func (m *CacheDAO) SetNX(key string, value string, ttl int64) (interface{}, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, exists := m.data[key]; exists {
		return nil, errors.New("key already exists")
	}
	m.data[key] = cacheItem{
		value:      value,
		expireTime: time.Now().Add(time.Duration(ttl) * time.Second),
	}
	return value, nil
}

// Get 获取值
func (m *CacheDAO) Get(key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, exists := m.data[key]
	if !exists || time.Now().After(item.expireTime) {
		return "", errors.New("key not found or expired")
	}
	return item.value, nil
}

// GetDistributeLock 获取分布式锁
func (m *CacheDAO) GetDistributeLock(name string, expireTime time.Duration) cache.DistributeLock {
	// 简单实现，返回一个空的锁
	return &memoryDistributeLock{}
}

// GetDistributeLockWithRetry 获取分布式锁，如果没有拿到的话会在一定时间内重试
func (m *CacheDAO) GetDistributeLockWithRetry(name string, expireTime time.Duration, trys int, retryDelay time.Duration) cache.DistributeLock {
	// 简单实现，返回一个空的锁
	return &memoryDistributeLock{}
}

// memoryDistributeLock 内存分布式锁的简单实现
type memoryDistributeLock struct{}

func (l *memoryDistributeLock) Lock() error {
	return nil
}

func (l *memoryDistributeLock) Unlock() (bool, error) {
	return true, nil
}
