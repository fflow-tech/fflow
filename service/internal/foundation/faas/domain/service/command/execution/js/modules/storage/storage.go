package storage

import (
	"context"

	"github.com/fflow-tech/fflow-sdk-go/faas"

	"github.com/dop251/goja"
)

type storage struct {
	storage faas.Storage
}

// New 初始化一个 console 实例
func New(ctx faas.Context) *storage {
	return &storage{storage: ctx.Storage()}
}

// Get 获取 key 的值
func (s *storage) Get(context *context.Context, msg goja.Value) []interface{} {
	val, err := s.storage.Get(msg.ToString().String())
	return []interface{}{val, err}
}

// Set 设置 key 的值，如果没有设置过期时间，默认半年过期
func (s *storage) Set(context *context.Context, key string, value goja.Value, expireTime int64) string {
	err := s.storage.Set(key, value.ToString().String(), expireTime)
	if err != nil {
		return err.Error()
	}
	return ""
}

// Del 删除
func (s *storage) Del(key string) string {
	err := s.storage.Del(key)
	if err != nil {
		return err.Error()
	}
	return ""
}
