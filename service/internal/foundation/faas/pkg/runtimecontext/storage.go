package runtimecontext

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/redis"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

const timeForHalfYear = time.Hour * 24 * 30 * 6

// Storage 利用 redis 实现持久化能力
type Storage struct {
	ctx         context.Context
	function    *entity.Function
	redisClient *redis.Client
}

// 获取一个 Storage 实例
func newStorage(ctx context.Context, function *entity.Function, client *redis.Client) *Storage {
	return &Storage{
		ctx:         ctx,
		function:    function,
		redisClient: client,
	}
}

// Get 获取 key 的值
func (s *Storage) Get(key string) (any, error) {
	return s.redisClient.Get(context.Background(), s.getRealKey(key))
}

// Set 设置 key 的值，如果没有设置过期时间，默认半年过期
func (s *Storage) Set(key string, value any, expireTime int64) error {
	t := int64(timeForHalfYear.Seconds())
	if expireTime > 0 {
		t = expireTime
	}

	var realValue string
	switch value.(type) {
	case string:
		realValue = value.(string)
	default:
		realValue = utils.StructToJsonStr(value)
	}
	return s.redisClient.Set(context.Background(), s.getRealKey(key), realValue, t)
}

// Expire 标记过期
func (s *Storage) Expire(key string, expireTime int64) error {
	return s.redisClient.Expire(context.Background(), s.getRealKey(key), expireTime)
}

// Del 删除
func (s *Storage) Del(key string) error {
	return s.redisClient.Del(context.Background(), s.getRealKey(key))
}

// getRealKey 获取拼装了默认的 key 前缀的真实 key
func (s *Storage) getRealKey(key string) string {
	namespace := s.function
	keyPrefix := fmt.Sprintf("%s.%s", utils.GetEnv(), namespace.Namespace)

	debugMode := s.ctx.Value("debugMode")
	// KEY 不存在或不为调试模式
	if debugMode == nil || !debugMode.(bool) {
		return fmt.Sprintf("%s_%s", keyPrefix, key)
	}
	return fmt.Sprintf("%s_%s_%s", "debug", keyPrefix, key)
}
