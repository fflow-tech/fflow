// Package repo 实际仓储层的实现
package repo

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/cache"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/cache/redis"
)

// CacheRepo 缓存仓储层
type CacheRepo struct {
	cacheDAO cache.CacheDAO
}

// NewCacheRepo 构造方法
func NewCacheRepo(d *redis.RedisCacheDAO) *CacheRepo {
	return &CacheRepo{cacheDAO: d}
}

// Set 设置 key val
func (r *CacheRepo) Set(key string, value string, ttl int64) error {
	return r.cacheDAO.Set(key, value, ttl)
}

// SetNX 设置 key val, 如果存在就设置, 如果不存在就不设置
func (r *CacheRepo) SetNX(key string, value string, ttl int64) (interface{}, error) {
	return r.cacheDAO.SetNX(key, value, ttl)
}

// Get 通过 key 获取 val
func (r *CacheRepo) Get(key string) (string, error) {
	return r.cacheDAO.Get(key)
}

// GetDistributeLock 获取分布式锁
func (r *CacheRepo) GetDistributeLock(name string, expireTime time.Duration) cache.DistributeLock {
	return r.cacheDAO.GetDistributeLock(name, expireTime)
}

// GetDistributeLockWithRetry 获取分布式锁, 预期的时间范围内拿不到就返回nil
func (r *CacheRepo) GetDistributeLockWithRetry(name string, expireTime time.Duration,
	trys int, retryDelay time.Duration) cache.DistributeLock {
	return r.cacheDAO.GetDistributeLockWithRetry(name, expireTime, trys, retryDelay)
}
