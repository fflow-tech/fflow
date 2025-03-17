// Package redis 实现了统一的 Redis 客户端，并提供基础的分布式缓存与分布式锁功能封装
package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

const (
	connTimeoutDuration time.Duration = 30 * time.Second
)

var (
	poolMap    sync.Map
	mutex      sync.Mutex
	redsyncMap sync.Map

	// ErrFailedGetDLock 分布式锁加锁失败：该锁已经被占用。配合 GetDistributeLock 使用.
	ErrFailedGetDLock = redsync.ErrFailed
)

// Client Redis 客户端.
type Client struct {
	Pool            *redis.Pool
	RedsyncInstance *redsync.Redsync
	Config          config.RedisConfig
}

var (
	routineIDMap                    = sync.Map{}
	lockTimesForOneRoutineInitValue = 0
)

// DefaultDistributeLock 默认分布式锁.
type DefaultDistributeLock struct {
	name  string
	mutex *redsync.Mutex
}

// NewDefaultDistributeLock 实例化.
func NewDefaultDistributeLock(mutex *redsync.Mutex, name string) *DefaultDistributeLock {
	return &DefaultDistributeLock{mutex: mutex, name: name}
}

// Lock 加锁.
// 实现可重入的能力.
func (l *DefaultDistributeLock) Lock() error {
	v, ok := routineIDMap.LoadOrStore(l.getKey(), lockTimesForOneRoutineInitValue)
	if ok {
		routineIDMap.Store(l.getKey(), v.(int)+1)
		return nil
	}

	err := l.mutex.Lock()
	if err != nil {
		return err
	}

	return nil
}

// Unlock 解锁.
// 实现可重入的能力.
func (l *DefaultDistributeLock) Unlock() (bool, error) {
	v, ok := routineIDMap.LoadOrStore(l.getKey(), lockTimesForOneRoutineInitValue)
	if ok && v == lockTimesForOneRoutineInitValue {
		status, err := l.mutex.Unlock()
		routineIDMap.Delete(l.getKey())
		return status, err
	}

	routineIDMap.Store(l.getKey(), v.(int)-1)
	return true, nil
}

func (l *DefaultDistributeLock) getKey() string {
	return strings.Join([]string{l.name, strconv.Itoa(utils.GetCurrentGoroutineID())}, "_")
}

// GetClient 获取客户端.
func GetClient(config config.RedisConfig) *Client {
	return &Client{
		Pool:            getRedisPool(config),
		RedsyncInstance: getRedsyncInstance(config),
		Config:          config,
	}
}

// SetIntKey 执行 Redis SET 命令，传入的 value 为整型 不带过期时间.
func (c *Client) SetIntKey(ctx context.Context, key string, value int) error {
	if key == "" {
		return fmt.Errorf("redis SetInt key can't be empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	return err
}

// Set 执行 Redis SET 命令，expireTime 时间单位为秒.
func (c *Client) Set(ctx context.Context, key, value string, expireTime int64) error {
	if key == "" || value == "" {
		return fmt.Errorf("redis SET key or value can't be empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", key, value, "EX", expireTime)
	return err
}

// SetNX 执行 Redis SetNX 命令，expireTime 时间单位为秒.
func (c *Client) SetNX(ctx context.Context, key, value string, expireTime int64) (interface{}, error) {
	if key == "" || value == "" {
		return nil, fmt.Errorf("redis SETNX key or value can't be empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	r, err := conn.Do("SETNX", key, value)
	if err != nil {
		return nil, err
	}
	if r.(int64) == 1 {
		return conn.Do("EXPIRE", key, expireTime)
	}

	return r, nil
}

// SetInt 执行 Redis SET 命令，传入的 value 为整型，expireTime 时间单位为秒.
func (c *Client) SetInt(ctx context.Context, key string, value int, expireTime uint64) error {
	if key == "" {
		return fmt.Errorf("redis SetInt key can't be empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", key, value, "EX", expireTime)
	return err
}

// MSet 执行 Redis MSET 命令.
func (c *Client) MSet(ctx context.Context, args ...interface{}) error {
	// redigo 对为 nil 或 empty 的参数报错信息很模糊，因此手动添加错误信息
	if len(args) == 0 {
		return fmt.Errorf("redis MSET key or value can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("MSET", args...)
	return err
}

// Get 执行 Redis GET 命令.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return redis.String(conn.Do("GET", key))
}

// MGet 执行 Redis MGET 命令.
func (c *Client) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	// redigo 对为 nil 或 empty 的参数报错信息很模糊，因此手动添加错误信息
	if len(keys) == 0 {
		return nil, fmt.Errorf("redis MSET args can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	args := make([]interface{}, len(keys))
	for i := range keys {
		args[i] = keys[i]
	}
	return redis.Values(conn.Do("MGET", args...))
}

// Scan 封装了 redigo.Scan() 方法，配合 MGet 使用，传入的 values 为指针类型
// 该方法将从 Redis 获得的 interface slice 转换为期望的数据类型，并返回未处理的数据
// 如果 src 中的数据全部被处理，返回的是 []interface{}{}，而不是 nil
func (c *Client) Scan(src []interface{}, values ...interface{}) ([]interface{}, error) {
	// redigo 没有对 nil 值进行处理
	if src == nil || values == nil {
		return nil, fmt.Errorf("redis Scan src and values can't be nil")
	}
	unprocessed, err := redis.Scan(src, values...)
	if err != nil {
		return nil, err
	}
	return unprocessed, nil
}

// INCR 实现了 Redis INCR 指令，其实质上是执行了 INCRBY Key 1 指令
func (c *Client) INCR(ctx context.Context, key string) error {
	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("INCR", key)
	return err
}

// GetDistributeLock 获取一个分布式可重入锁.
func (c *Client) GetDistributeLock(name string, expireTime time.Duration) *DefaultDistributeLock {
	return NewDefaultDistributeLock(c.RedsyncInstance.NewMutex(name, redsync.WithExpiry(expireTime),
		redsync.WithTries(1)), name)
}

// GetDistributeLockWithRetry 获取一个分布式可重入锁.
func (c *Client) GetDistributeLockWithRetry(name string, expireTime time.Duration,
	trys int, retryDelay time.Duration) *DefaultDistributeLock {
	return NewDefaultDistributeLock(c.RedsyncInstance.NewMutex(
		name, redsync.WithExpiry(expireTime), redsync.WithTries(trys), redsync.WithRetryDelay(retryDelay)), name)
}

// Del 执行 Redis DEL 命令.
func (c *Client) Del(ctx context.Context, key string) error {
	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	return err
}

// Expire 执行 Redis EXPIRE 命令.
func (c *Client) Expire(ctx context.Context, key string, expireTime int64) error {
	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("EXPIRE", key, expireTime)
	return err
}

func getRedisPool(config config.RedisConfig) *redis.Pool {
	if db, ok := poolMap.Load(config.Address); ok {
		return db.(*redis.Pool)
	}
	mutex.Lock()
	defer mutex.Unlock()

	pool := newRedisPool(config)
	poolMap.Store(config.Address, pool)
	return pool
}

func newRedisPool(conf config.RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := newRedisConn(conf)
			if err != nil {
				log.Errorf("Failed to get redis connection, caused by %s", err)
				return nil, err
			}
			return c, nil
		},
		MaxActive: conf.MaxActive,
		Wait:      conf.Wait,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Errorf("Failed to ping redis server, caused by %s", err)
			}
			return err
		},
	}
}

func getRedsyncInstance(config config.RedisConfig) *redsync.Redsync {
	if db, ok := redsyncMap.Load(config.Address); ok {
		return db.(*redsync.Redsync)
	}
	mutex.Lock()
	defer mutex.Unlock()
	redsyncInstance := redsync.New(redigo.NewPool(getRedisPool(config)))
	redsyncMap.Store(config.Address, redsyncInstance)
	return redsyncInstance
}

func newRedisConn(conf config.RedisConfig) (redis.Conn, error) {
	if conf.Address == "" {
		panic("Cannot get redis address from config")
	}

	conn, err := redis.Dial(conf.Network, conf.Address,
		redis.DialPassword(conf.Password),
		redis.DialConnectTimeout(connTimeoutDuration))
	if err != nil {
		log.Errorf("Failed to connect to redis, caused by %s", err)
		return nil, err
	}
	return conn, nil
}

// HSet 执行Redis HSet 命令.
func (c *Client) HSet(ctx context.Context, args ...interface{}) error {
	// redigo 对为 nil 或 empty 的参数报错信息很模糊，因此手动添加错误信息
	if len(args) == 0 {
		return fmt.Errorf("redis HSET key or value can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("HSET", args...)
	return err
}

// HLen 执行 Redis HLen 命令.
func (c *Client) HLen(ctx context.Context, key string) (int, error) {
	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return redis.Int(conn.Do("HLEN", key))
}

// ZSet 执行Redis ZSet 命令.
func (c *Client) ZSet(ctx context.Context, args ...interface{}) error {
	// redigo 对为 nil 或 empty 的参数报错信息很模糊，因此手动添加错误信息
	if len(args) == 0 {
		return fmt.Errorf("redis ZSET key or value can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("ZADD", args...)
	return err
}

// ZRange 执行redis ZRange 命令.
func (c *Client) ZRange(ctx context.Context, keys ...interface{}) ([]string, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("redis ZSET key or value can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	args := make([]interface{}, 0, len(keys))
	args = append(args, keys...)
	values, err := redis.Values(conn.Do("ZRANGEBYSCORE", args...))
	if err != nil {
		return nil, err
	}
	var rspValues []string
	for _, value := range values {
		tempValue := (value).([]byte)
		rspValues = append(rspValues, string(tempValue))
	}
	return rspValues, nil
}

// ZREM 执行redis ZREM 命令.
func (c *Client) ZRem(ctx context.Context, keys ...interface{}) error {
	if len(keys) == 0 {
		return fmt.Errorf("redis ZSET key or value can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	args := make([]interface{}, 0, len(keys))
	args = append(args, keys...)
	_, err = conn.Do("ZRem", args...)
	return err
}

// Exists 执行Redis Exists 命令.
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	// redigo 对为 nil 或 empty 的参数报错信息很模糊，因此手动添加错误信息
	if len(keys) == 0 {
		return 0, fmt.Errorf("redis Exists args can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	args := make([]interface{}, len(keys))
	for i := range keys {
		args[i] = keys[i]
	}
	return redis.Int64(conn.Do("exists", args...))
}

// Hexists 执行Redis Hexists 命令.
func (c *Client) Hexists(ctx context.Context, keys ...string) (int64, error) {
	// redigo 对为 nil 或 empty 的参数报错信息很模糊，因此手动添加错误信息
	if len(keys) == 0 {
		return 0, fmt.Errorf("redis HGET args can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	args := make([]interface{}, len(keys))
	for i := range keys {
		args[i] = keys[i]
	}
	return redis.Int64(conn.Do("Hexists", args...))
}

// SetEx 执行 Redis SET 命令，expireTime 时间单位为秒.
func (c *Client) SetEx(ctx context.Context, key, value string, expireTime int64) error {
	if key == "" || value == "" {
		return fmt.Errorf("redis SET key or value can't be empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", key, value, "EX", expireTime)
	return err
}

// SetIntEx 执行 Redis SET 命令，传入的 value 为整型，expireTime 时间单位为秒.
func (c *Client) SetIntEx(ctx context.Context, key string, value int, expireTime uint64) error {
	if key == "" {
		return fmt.Errorf("redis SetInt key can't be empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", key, value, "EX", expireTime)
	return err
}

// HGet 执行Redis HGet 命令.
func (c *Client) HGet(ctx context.Context, keys ...string) (string, error) {
	// redigo 对为 nil 或 empty 的参数报错信息很模糊，因此手动添加错误信息
	if len(keys) == 0 {
		return "", fmt.Errorf("redis HGET args can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	args := make([]interface{}, len(keys))
	for i := range keys {
		args[i] = keys[i]
	}
	return redis.String(conn.Do("HGET", args...))
}

// HDel 执行 Redis HDel 命令.
func (c *Client) HDel(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return fmt.Errorf("redis HDel args can't be nil or empty")
	}

	tContext, cancel := context.WithTimeout(ctx, connTimeoutDuration)
	defer cancel()

	conn, err := c.Pool.GetContext(tContext)
	if err != nil {
		return err
	}
	defer conn.Close()

	args := make([]interface{}, len(keys))
	for i := range keys {
		args[i] = keys[i]
	}
	_, err = conn.Do("HDEL", args...)
	return err
}
