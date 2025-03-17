// Package notify 负责完成定时任务具体回调 http
package notify

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq/eventbus"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/logs"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

type notifySender interface {
	SendNotify(hashID string) error
}

// Notify 通知服务处理器定义
type Notify struct {
	sender       notifySender
	eventBusRepo ports.EventBusRepository
	consumers    []mq.Consumer
	consumerNum  int
	pool         concurrency.WorkerPool
	recorder     eventbus.Recorder
}

// NewNotifyEventProcessor 初始化通知事件服务
func NewNotifyEventProcessor(sender *command.Adapters, eventBusRepo *repo.EventBusRepo,
	recorder *eventbus.LogRecorder, workerPool *concurrency.GoWorkerPool) *Notify {
	r := &Notify{
		sender:       sender,
		eventBusRepo: eventBusRepo,
		consumerNum:  config.GetNotifyTaskConfig().ConsumerNum,
		pool:         workerPool,
		recorder:     recorder,
	}
	return r
}

// Type 处理器类型
func (p *Notify) Type() string {
	return reflect.TypeOf(p).Elem().Name()
}

// Start 启动
func (p *Notify) Start() error {
	log.Infof("start notify server")
	for i := 0; i < p.consumerNum; i++ {
		bizKey := strings.Join([]string{p.Type(), strconv.Itoa(i)}, "_")
		// 因为需要使用到当前 bizkey 所以当前方法没有单独抽出
		consumer, err := p.eventBusRepo.NewTimerTaskConsumer(context.Background(),
			config.GetEventConfig().TimerEventGroup, func(ctx context.Context, message interface{}) (err error) {
				startTime := time.Now()
				defer logs.DumpPanicStack(p.Type(), fmt.Errorf("failed to consume notify event"))
				defer func() {
					p.recorder.RecordTDMQConsumeLog(&logs.EventRecord{
						Message:   message,
						BizKey:    bizKey,
						StartTime: startTime,
						Error:     err,
					})
				}()
				return p.Working(message)
			})
		if err != nil {
			return fmt.Errorf("failed to start notify event consume: %w", err)
		}
		p.consumers = append(p.consumers, consumer)
	}
	return nil
}

// Restart 重启
func (p *Notify) Restart() error {
	if err := p.Stop(); err != nil {
		return fmt.Errorf("failed to stop notify event processor: %w", err)
	}
	return p.Start()
}

// Stop 停止
func (p *Notify) Stop() error {
	log.Infof("stop notify server")
	var wg sync.WaitGroup
	for _, consumer := range p.consumers {
		consumer := consumer
		wg.Add(1)
		if err := p.pool.Submit(func() {
			defer wg.Done()
			if err := consumer.Close(); err != nil {
				// 如果关闭不成功, 其余的还是要继续关闭
				log.Errorf("Failed to close notify event consumer, caused by %s", err)
			}
		}); err != nil {
			log.Errorf("notify service submit stop task failed, consumer: %+v,err:%w",
				consumer, err)
		}
	}
	wg.Wait()
	return nil
}

// Working 开始工作
func (p *Notify) Working(msg interface{}) error {
	hashID, err := p.GetTimeTaskEvent(msg)
	if err != nil {
		log.Errorf("Failed to Working GetTimeTaskEvent, caused by %s", err)
		return err
	}
	if err = p.pool.Submit(func() {
		if err := p.sender.SendNotify(hashID); err != nil {
			log.Errorf("Failed to Working SendNotify, caused by %s", err)
		}
	}); err != nil {
		log.Errorf("notify service submit working task failed, hashID: %s,err:%w",
			hashID, err)
	}
	return nil
}

// GetTimeTaskEvent 获取定时器任务事件数据
func (p *Notify) GetTimeTaskEvent(message interface{}) (string, error) {
	msg, ok := (message).(pulsar.Message)
	if !ok {
		log.Errorf("Failed to GetTimeTaskEvent, message=%+v", message)
		return "", nil
	}
	return string(msg.Payload()), nil
}
