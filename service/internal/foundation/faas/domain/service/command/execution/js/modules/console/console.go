package console

import (
	"context"

	"github.com/fflow-tech/fflow-sdk-go/faas"

	"github.com/dop251/goja"
)

type console struct {
	log faas.Logger
}

// New 初始化一个 console 实例
func New(ctx faas.Context) *console {
	return &console{log: ctx.Logger()}
}

// Log 基础的 log 方法
func (c *console) Log(context *context.Context, msg goja.Value, args ...interface{}) {
	c.log.Infof(msg.ToString().String(), args...)
}

// Logf 基础的 log 方法
func (c *console) Logf(context *context.Context, msg goja.Value, args ...interface{}) {
	c.log.Infof(msg.ToString().String(), args...)
}

// Errorf 基础的 log 方法
func (c *console) Errorf(context *context.Context, msg goja.Value, args ...interface{}) {
	c.log.Errorf(msg.ToString().String(), args...)
}

// Debugf 基础的 log 方法
func (c *console) Debugf(context *context.Context, msg goja.Value, args ...interface{}) {
	c.log.Debugf(msg.ToString().String(), args...)
}

// Warnf 基础的 log 方法
func (c *console) Warnf(context *context.Context, msg goja.Value, args ...interface{}) {
	c.log.Warnf(msg.ToString().String(), args...)
}
