// Package mock Mock数据缓存层实现
package mock

// CacheDAOMock 缓存的 MOCK 能力
type CacheDAOMock struct {
}

// NewMockCacheDAO 初始化 MOCK 能力
func NewMockCacheDAO() *CacheDAOMock {
	return &CacheDAOMock{}
}

// Set 设置缓存值
func (s CacheDAOMock) Set(key string, value string, ttl int64) error {
	return nil
}

// Get 获取缓存值
func (s CacheDAOMock) Get(key string) (string, error) {
	return "mock cache data", nil
}
