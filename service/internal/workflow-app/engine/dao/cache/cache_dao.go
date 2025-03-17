package cache

import "time"

// CacheDAO 缓存接入层
type CacheDAO interface {
	// Set 设置
	Set(key string, value string, ttl int64) error
	// SetNX 如果不存在则设置
	SetNX(key string, value string, ttl int64) (interface{}, error)
	// Get 获取值
	Get(key string) (string, error)
	// GetDistributeLock 获取分布式锁
	GetDistributeLock(name string, expireTime time.Duration) DistributeLock
	// GetDistributeLockWithRetry 获取分布式锁，如果没有拿到的话会在一定时间内重试
	GetDistributeLockWithRetry(name string, expireTime time.Duration, trys int, retryDelay time.Duration) DistributeLock
}

// DistributeLock 分布式锁
type DistributeLock interface {
	Lock() error           // 加锁
	Unlock() (bool, error) // 解锁
}
