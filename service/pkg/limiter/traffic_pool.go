package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"golang.org/x/time/rate"
)

// TrafficPool 限流池.
type TrafficPool struct {
	ctx           context.Context
	cancel        func()
	refresher     *time.Ticker
	limiter       *rate.Limiter
	limiterConfig *config.LimiterConfig
}

// NewTrafficPool 流量池构造器.
func NewTrafficPool(limiterConfig *config.LimiterConfig) *TrafficPool {
	t := TrafficPool{
		limiter:       rate.NewLimiter(rate.Limit(limiterConfig.Limit), int(limiterConfig.Burst)),
		limiterConfig: limiterConfig,
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	// 执行自刷新操作.
	t.refresher = t.refresh(time.Duration(limiterConfig.RefreshInterval) * time.Second)
	return &t
}

// Stop 停止限流池.
func (t *TrafficPool) Stop() {
	t.cancel()
	t.refresher.Stop()
}

// Get 获取一个单位的流量, 此方法会根据注入的 configProvider 提供的 config 进行超时限制，并且失败后不会进行重试.
func (t *TrafficPool) Get(opts ...Option) error {
	policy := NewWaitPolicy(t.limiterConfig)
	for _, opt := range opts {
		opt.Apply(policy)
	}
	tContext, cancel := context.WithTimeout(t.ctx, policy.Timeout)
	defer cancel()
	var err error
	for i := 0; i <= int(policy.Retries); i++ {
		err = t.limiter.Wait(tContext)
		if err == nil {
			return nil
		}
		// 可重试的错误，则发起下一次的尝试.
		if isRetryableErr(err) {
			continue
		}
		// 不可重试的错误，直接返回.
		return fmt.Errorf("get token from traffic pool failed, unretryable err: %v", err)
	}
	// 重试次数用尽.
	return fmt.Errorf("get token from traffic pool failed, retries: %d, err: %v", policy.Retries, err)
}

// TryGetNoRetry 尝试获取流量，设置单次请求的超时时间，不做重试.
func (t *TrafficPool) TryGetNoRetry(timeout time.Duration) error {
	return t.Get(WithTimeout(timeout))
}

// TryGet 尝试获取流量，设置单次请求的超时时间和重试次数.
func (t *TrafficPool) TryGet(timeout time.Duration, retries uint32) error {
	return t.Get(WithTimeout(timeout), WithRetries(retries))
}

// errLimiterTimeout 限流器等待流量超时错误.
var errLimiterTimeout = fmt.Errorf("rate: Wait(n=1) would exceed context deadline")

func isRetryableErr(err error) bool {
	if err == context.DeadlineExceeded {
		return true
	}
	return err.Error() == errLimiterTimeout.Error()
}

// refresh 自刷新，动态调整池子容量和补充速率.
func (t *TrafficPool) refresh(interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for cur := range ticker.C {
			t.limiter.SetBurst(t.limiterConfig.Burst)
			t.limiter.SetLimit(rate.Limit(t.limiterConfig.Limit))
			log.Infof("set burst: %d, limit: %d, time: %v", t.limiterConfig.Burst, t.limiterConfig.Limit, cur)
		}
	}()
	return ticker
}
