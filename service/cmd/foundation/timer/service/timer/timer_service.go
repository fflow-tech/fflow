// Package timer 负责待消费时间片具体任务派发到 notify 消费集。
package timer

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
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

type readyTimerGetter interface {
	GetReadyTimer(readySet string, startTime time.Time) error
}

// TimeTask 定时器任务处理器
type TimeTask struct {
	timerGetter       readyTimerGetter
	eventBusRepo      ports.EventBusRepository
	consumers         []mq.Consumer
	consumerNum       int
	pool              concurrency.WorkerPool
	recorder          eventbus.Recorder
	waitDuration      time.Duration
	workSleepDuration time.Duration
}

// NewTimerTaskEventProcessor 初始化定时器任务事件
func NewTimerTaskEventProcessor(getter *command.Adapters, eventBusRepo *repo.EventBusRepo,
	recorder *eventbus.LogRecorder, workerPool *concurrency.GoWorkerPool) *TimeTask {
	r := &TimeTask{
		timerGetter:  getter,
		eventBusRepo: eventBusRepo,
		consumerNum:  config.GetTimerTaskConfig().ConsumerNum,
		pool:         workerPool,
		recorder:     recorder,
		// 优雅关闭 等待一分钟 等待当前工作线程完成 因为是分钟级别的桶 所以需要1分钟来等待
		waitDuration:      time.Minute,
		workSleepDuration: time.Duration(config.GetTimerTaskConfig().WorkSleepSecond) * time.Second,
	}
	return r
}

// Type 处理器类型
func (p *TimeTask) Type() string {
	return reflect.TypeOf(p).Elem().Name()
}

// Start 启动
func (p *TimeTask) Start() error {
	log.Infof("start timeTask server")
	for i := 0; i < p.consumerNum; i++ {
		bizKey := strings.Join([]string{p.Type(), strconv.Itoa(i)}, "_")
		// 因为需要使用到当前bizkey 所以当前方法没有单独抽出
		consumer, err := p.eventBusRepo.NewPollingConsumer(context.Background(),
			config.GetEventConfig().TimerEventGroup,
			func(ctx context.Context, message interface{}) (err error) {
				startTime := time.Now()
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
			return fmt.Errorf("failed to start timerTask event consume: %w", err)
		}
		p.consumers = append(p.consumers, consumer)
	}
	return nil
}

// Restart 重启
func (p *TimeTask) Restart() error {
	// 优雅关闭 等待一分钟 等待当前工作线程完成 因为是分钟级别的桶 所以需要1分钟来等待
	if err := p.Stop(); err != nil {
		return fmt.Errorf("failed to stop timerTask event processor: %w", err)
	}
	return p.Start()
}

// Stop 停止
func (p *TimeTask) Stop() error {
	log.Infof("stop timerTask server")
	var wg sync.WaitGroup
	for _, consumer := range p.consumers {
		consumer := consumer
		wg.Add(1)
		if err := p.pool.Submit(func() {
			defer wg.Done()
			if err := consumer.Close(); err != nil {
				// 如果关闭不成功, 其余的还是要继续关闭
				log.Errorf("Failed to close timerTask event consumer, caused by %s", err)
			}
		}); err != nil {
			log.Errorf("timer service stop task failed, consumer: %+v, err:%w", consumer, err)
		}
	}
	wg.Wait()

	log.Infof("wait timer task working done")
	time.Sleep(p.waitDuration)
	return nil
}

// Working 工作
func (p *TimeTask) Working(msg interface{}) error {
	log.Infof("GoroutineID:%d timerTask working timeAt:%v", utils.GetCurrentGoroutineID(), time.Now())
	readySet, err := p.GetTimeSliceEvent(msg)
	if err != nil {
		log.Errorf("Failed to GetTimeSliceEvent ,caused by %v", err)
		return err
	}
	startTime, err := getTimeSliceStartTime(readySet)
	if err != nil {
		log.Errorf("Failed to getTimeSliceStartTime ,caused by %v", err)
		return err
	}
	if err := p.pool.Submit(func() {
		p.PollingTimeTaskSlice(startTime, readySet)
	}); err != nil {
		log.Errorf("timer service start working failed, readySet:%s, err:%w", readySet, err)
	}
	return nil
}

// PollingTimeTaskSlice 轮询时间任务切片
func (p *TimeTask) PollingTimeTaskSlice(startTime time.Time, readySet string) {
	nowTime := startTime
	for {
		if err := p.timerGetter.GetReadyTimer(readySet, nowTime); err != nil {
			log.Errorf("Failed to getTimeSliceStartTime ,goroutineID:%d caused by %v ",
				utils.GetCurrentGoroutineID(), err)
			return
		}
		time.Sleep(p.workSleepDuration)
		nowTime = nowTime.Add(p.workSleepDuration)
		if startTime.Add(time.Minute).Before(nowTime) { // 当工作了一分钟之后就不再工作了
			log.Infof("PollingTimeTaskSlice readySet %v  is end", readySet)
			return
		}
	}
}
func getTimeSliceStartTime(readySet string) (time.Time, error) {
	readySetStrings := strings.Split(readySet, "_")
	if len(readySetStrings) != 2 {
		return time.Time{}, fmt.Errorf("failed to getTimeSliceStartTime,caused error by: Split readySet")
	}
	readySetTime := readySetStrings[1]
	readyTime := readySetTime + ":00" // 时间片转换成秒级
	return time.ParseInLocation(dto.TimerTriggerTimeFormat, readyTime, time.Local)
}

// GetTimeSliceEvent 获取事件数据
func (p *TimeTask) GetTimeSliceEvent(message interface{}) (string, error) {
	msg, ok := (message).(pulsar.Message)
	if !ok {
		log.Errorf("Failed to get timer task event type, message=%+v", message)
		return "", nil
	}
	return string(msg.Payload()), nil
}
