package limiter

import (
	"time"

	"github.com/fflow-tech/fflow/service/pkg/config"
)

// WaitPolicy 获取流量的策略.
type WaitPolicy struct {
	// 单次等待的超时时间.
	Timeout time.Duration
	// 重试次数.
	Retries uint32
}

// NewWaitPolicy 获取流量规则构造器.
func NewWaitPolicy(config *config.LimiterConfig) *WaitPolicy {
	return &WaitPolicy{
		Timeout: time.Duration(config.WaitingDuration) * time.Second,
	}
}

// Option 应用项.
type Option interface {
	Apply(*WaitPolicy)
}

// OptionFunc 配置 WaitPolicy 的函数.
type OptionFunc func(*WaitPolicy)

// Apply 调用应用函数.
func (f OptionFunc) Apply(r *WaitPolicy) {
	f(r)
}

// WithTimeout 设置单次等待的超时时间.
func WithTimeout(timeout time.Duration) Option {
	return OptionFunc(func(r *WaitPolicy) {
		r.Timeout = timeout
	})
}

// WithRetries 设置重试次数.
func WithRetries(retries uint32) Option {
	return OptionFunc(func(r *WaitPolicy) {
		r.Retries = retries
	})
}
