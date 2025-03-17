package concurrency

import (
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"

	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/panjf2000/ants/v2"
)

// WorkerPool 协程工作池.
type WorkerPool interface {
	Submit(func()) error
}

var (
	once              = sync.Once{}
	defaultWorkerPool *GoWorkerPool
)

// GoWorkerPool golang 协程工作池.
type GoWorkerPool struct {
	pool *ants.Pool
}

// Submit 提交任务.
func (g *GoWorkerPool) Submit(f func()) error {
	return g.pool.Submit(f)
}

// GetDefaultWorkerPool 获取默认的协程工作池.
func GetDefaultWorkerPool() *GoWorkerPool {
	once.Do(func() {
		conf := config.GetWorkerPoolConfig()
		log.Infof("Ready to init worker pool, size: %d", conf.Size)

		pool, err := ants.NewPool(
			conf.Size,
			ants.WithExpiryDuration(time.Duration(conf.ExpireDuration)*time.Second),
		)
		if err != nil {
			log.Errorf("Init worker pool failed, err: %w", err)
		}
		defaultWorkerPool = &GoWorkerPool{
			pool: pool,
		}
	})
	return defaultWorkerPool
}
