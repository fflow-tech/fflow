package schedule

import (
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/redis"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// Scheduler 调度策略
type Scheduler struct {
	RedisConfig config.RedisConfig
	Key         string // 传给redis的key
}

// NewSchedule 新建 schedule 实例
func NewSchedule(key string, redisConfig config.RedisConfig) *Scheduler {
	return &Scheduler{Key: key, RedisConfig: redisConfig}
}

// Schedule 通过 redis 实现互斥任务定时器, 根据是否返回 error 来判断 timer 是否执行
func (s *Scheduler) Schedule(serviceName string,
	newNode string, holdTime time.Duration) (nowNode string, err error) {
	log.Infof("Schedule start with key:%s||%s|%s|%s", s.Key, serviceName, newNode, holdTime)
	// 测试环境和生产环境目前使用的是同一个 Redis，所以 key 值传入 Env
	lockName := fmt.Sprintf("%s.%s.%s", utils.GetEnv(), serviceName, s.Key)
	lock := redis.GetClient(s.RedisConfig).GetDistributeLock(lockName, holdTime)
	err = lock.Lock()
	if err != nil {
		return nowNode, err
	}

	if lock != nil {
		log.Infof("Set lock success, nowNode is: %s", newNode)
		nowNode = newNode
	}

	if nowNode == "" {
		return nowNode, fmt.Errorf("can not get nowNode")
	}

	return nowNode, nil
}
