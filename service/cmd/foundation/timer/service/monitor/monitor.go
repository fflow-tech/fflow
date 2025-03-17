package monitor

import (
	"context"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/query"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/monitor"
	"github.com/fflow-tech/fflow/service/pkg/redis"

	"github.com/fflow-tech/fflow/service/pkg/log"
)

const (
	// 因为轮询时间间隔为 10 min，所以将分布式锁的过期时间设置为 15 min.
	distributeLockExpireDuration time.Duration = 15 * time.Minute
	// 查询未触发定时器的时间范围 20 min前到 10 min 前，所以需要向前扣除 20 min.
	// TODO(@weixxxu): 根据时间点获取所处的 10 分钟区间的逻辑目前通过 CountPendingTimers api 实现
	// 后续将时间范围作为参数由调用方传入，提高方法灵活性和可读性.
	pendingTimerQueryDuration time.Duration = 20 * time.Minute
)

type reporter interface {
	ReportTimerEnabledTotalNum(total float64)
	ReportTimerFailedNum(total float64)
}

type timerCounter interface {
	CountTimersByStatus(status entity.TimerDefStatus) (int64, error)
	CountPendingTimers(curTime time.Time) (int, error)
}

type lockProvider interface {
	GetDistributeLock(name string, expireTime time.Duration) *redis.DefaultDistributeLock
}

// Monitor 监控服务.
type Monitor struct {
	redisClient lockProvider
	pool        concurrency.WorkerPool
	counter     timerCounter
	reporter    reporter
}

// NewMonitor 监控服务构造器.
func NewMonitor(counter *query.Adapters, redisClient *redis.Client,
	reporter *monitor.Reporter, workerPool *concurrency.GoWorkerPool) *Monitor {
	return &Monitor{
		pool:        workerPool,
		redisClient: redisClient,
		reporter:    reporter,
		counter:     counter,
	}
}

// ReportRecord 定时完成汇报动作.
func (m *Monitor) ReportRecord(ctx context.Context) error {
	log.Infof("start task of reporting timer record")
	now := time.Now()
	// 分钟级时间字符串.
	timeStr := now.Format(dto.TimerTaskTimeFormat)
	// 使用十分钟级时间字符串拼接成分布式锁名.
	lockName := "monitor_" + timeStr[:len(timeStr)-1]
	// 1. 争抢当前时间片下的分布式锁，以 10 分钟为粒度
	lock := m.redisClient.GetDistributeLock(lockName, distributeLockExpireDuration)
	if err := lock.Lock(); err != nil {
		return err
	}
	log.Infof("got lock: %s successfully, start to report timer record", lockName)
	// 2. 查 timer_def 表，查询当前处于激活状态下的定时器数量，进行数据上报.
	m.pool.Submit(func() {
		if enabledTimerCount, err := m.counter.CountTimersByStatus(entity.Enabled); err != nil {
			log.Errorf("count enabled timers failed, err: %v", err)
		} else {
			m.reporter.ReportTimerEnabledTotalNum(float64(enabledTimerCount))
			log.Infof("report timer enabled total num successfully, num:%d", enabledTimerCount)
		}
	})
	// 3. 查询过去十分钟的 pending 表，统计失败的 timer 总数.
	if failedTimersCount, err := m.counter.CountPendingTimers(now.Add(-pendingTimerQueryDuration)); err != nil {
		log.Errorf("count pending timers failed, err: %v", err)
	} else {
		m.reporter.ReportTimerFailedNum(float64(failedTimersCount))
		log.Infof("report timer failed num successfully, num:%d", failedTimersCount)
	}
	return nil
}
