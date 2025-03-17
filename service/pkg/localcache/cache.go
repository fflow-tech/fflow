// Package localcache 实现了本地缓存
package localcache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
)

type Client interface {
	Set(key string, value interface{}) error
	Get(key string, value interface{}) error
	Append(key string, value interface{}) error
	Close() error
}

// DefaultClient 本地缓存客户端
type DefaultClient struct {
	cache *bigcache.BigCache
}

var (
	// ErrEntryNotFound Get Key 不存在时会返回该错误
	ErrEntryNotFound = bigcache.ErrEntryNotFound
)

// Options 选项配置
type Options struct {
	// HardMaxCacheSize 设置缓存最大值，单位为 MB, 0 表示无限制
	HardMaxCacheSize int
	// LifeWindow 每条数据的存活时间
	LifeWindow time.Duration
	// CleanWindow 后，会删除被认为不活跃的对象，<=0 代表不操作
	CleanWindow time.Duration
}

// Option 选项方法
type Option func(*Options)

// NewOptions 初始化
func NewOptions(opts ...Option) Options {
	options := Options{
		HardMaxCacheSize: 64 * 1024 * 1024,
		LifeWindow:       20 * time.Second,
		CleanWindow:      30 * time.Second,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithHardMaxCacheSize 缓存最大大小
func WithHardMaxCacheSize(size int) Option {
	return func(o *Options) {
		o.HardMaxCacheSize = size
	}
}

// WithLifeWindow 每条数据的存活时间
func WithLifeWindow(w time.Duration) Option {
	return func(o *Options) {
		o.LifeWindow = w
	}
}

// WithCleanWindow CleanWindow 后，会删除被认为不活跃的对象，<=0 代表不操作
func WithCleanWindow(w time.Duration) Option {
	return func(o *Options) {
		o.CleanWindow = w
	}
}

// NewDefaultClient 初始化
func NewDefaultClient(opts ...Option) (*DefaultClient, error) {
	options := NewOptions(opts...)
	bigCache, err := bigcache.NewBigCache(getConfig(options))
	if err != nil {
		panic(fmt.Errorf("failed to init local cache: %w", err))
	}
	return &DefaultClient{
		cache: bigCache,
	}, nil
}

// Set 向 Client 插入新数据
func (l *DefaultClient) Set(key string, value interface{}) error {
	if value == nil || key == "" {
		return fmt.Errorf("client set key or value can't be empty")
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return l.cache.Set(key, jsonValue)
}

// Get 如果 Key 不存在会返回 ErrEntryNotFound
func (l *DefaultClient) Get(key string, value interface{}) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	data, err := l.cache.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// Append 在原有 value 后面追加新数据，如果 key 不存在则新建 key
func (l *DefaultClient) Append(key string, value interface{}) error {
	if value == nil {
		return fmt.Errorf("append value can't be empty")
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return l.cache.Append(key, jsonValue)
}

// Close 程序退出时调用 Close()方法，平滑关闭 Client
func (l *DefaultClient) Close() error {
	return l.cache.Close()
}

func getConfig(options Options) bigcache.Config {
	return bigcache.Config{
		// 每条数据的存活时间
		LifeWindow: options.LifeWindow,
		// CleanWindow 后，会删除被认为不活跃的对象，<=0 代表不操作
		CleanWindow: options.CleanWindow,
		// 设置缓存最大值，单位为MB,0表示无限制
		HardMaxCacheSize: options.HardMaxCacheSize,
		// shards 数量，必须为 2 的倍数
		Shards: 1024,
		// 设置最大存储对象数量，仅在初始化时可以设置
		// MaxEntriesInWindow = rps * lifeWindow
		MaxEntriesInWindow: 1000 * 10,
		// 缓存对象的最大字节数，仅在初始化时可以设置，设置为一个操作系统页面大小
		MaxEntrySize: 4 * 1024,
		// 是否打印内存分配信息
		Verbose: true,
		// 在缓存过期或者被删除时,可设置回调函数，参数是(key, val)，默认是nil不设置
		OnRemove: nil,
		// 在缓存过期或者被删除时,可设置回调函数，参数是(key, val, reason)，默认是nil不设置
		OnRemoveWithReason: nil,
	}
}
