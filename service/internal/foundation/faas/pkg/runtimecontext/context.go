// Package runtimecontext 提供函数执行的上下文
package runtimecontext

import (
	"context"
	"net/http"

	"github.com/fflow-tech/fflow-sdk-go/faas"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/redis"
)

// RuntimeContext 运行时上下文
type RuntimeContext struct {
	ctx      context.Context
	metadata *entity.Metadata
	log      *RuntimeLogger
	storage  *Storage
	request  *http.Request
}

// NewRuntimeContext 初始化一个 context
func NewRuntimeContext(ctx context.Context, f *entity.Function, r *http.Request) RuntimeContext {
	return RuntimeContext{
		ctx:      ctx,
		metadata: &entity.Metadata{Function: f},
		log:      newRuntimeLogger(ctx, f),
		storage:  newStorage(ctx, f, redis.GetClient(config.GetRedisConfig())),
		request:  r,
	}
}

// Logger 获取 logger 实例
func (r RuntimeContext) Logger() faas.Logger {
	return r.log
}

// Storage 获取 storage 实例
func (r RuntimeContext) Storage() faas.Storage {
	return r.storage
}

// Logs 获取所有打印的 log 信息
func (r RuntimeContext) Logs() []string {
	return r.log.logs
}

// Metadata 获取函数基础信息
func (r RuntimeContext) Metadata() faas.Metadata {
	return r.metadata
}

// Context 获取当前的 context
func (r RuntimeContext) Context() context.Context {
	return r.ctx
}
