// Package redis 基于 redis 的缓存实现
package redis

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/cache"
	"github.com/fflow-tech/fflow/service/pkg/redis"
)

// RedisCacheDAO redis 缓存服务
type RedisCacheDAO struct {
	redisClient *redis.Client
}

// NewRedisCacheDAO 初始化 redis 缓存服务
func NewRedisCacheDAO(client *redis.Client) *RedisCacheDAO {
	return &RedisCacheDAO{redisClient: client}
}

// Set 设置 key val
func (s *RedisCacheDAO) Set(key string, value string, ttl int64) error {
	if err := s.redisClient.Set(context.Background(), key, value, ttl); err != nil {
		return err
	}
	return nil
}

// SetNX 设置 key val
func (s *RedisCacheDAO) SetNX(key, value string, ttl int64) (interface{}, error) {
	return s.redisClient.SetNX(context.Background(), key, value, ttl)
}

// Get 通过 key 获取 val
func (s *RedisCacheDAO) Get(key string) (string, error) {
	str, err := s.redisClient.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	return str, nil
}

// GetDistributeLock 获取分布式锁
func (s *RedisCacheDAO) GetDistributeLock(name string, expireTime time.Duration) cache.DistributeLock {
	return s.redisClient.GetDistributeLock(name, expireTime)
}

// GetDistributeLockWithRetry 获取分布式锁, 预期的时间范围内拿不到就返回nil
func (s *RedisCacheDAO) GetDistributeLockWithRetry(name string, expireTime time.Duration,
	trys int, retryDelay time.Duration) cache.DistributeLock {
	return s.redisClient.GetDistributeLockWithRetry(name, expireTime, trys, retryDelay)
}
