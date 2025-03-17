// Package polling 负责时间片调度，派送到 timer 消费集。
package polling

import (
	"fmt"
	"reflect"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/fflow-tech/fflow/service/pkg/log"
)

const (
	workTimeInterval = 2 * time.Second // 每个工作线程的扫描间隔时间
)

type pollingProxy interface {
	GetPollingTaskWorkLock() (string, error)
	SendPollingTaskWork(timeSlice string) error
}

// Polling 轮询器
type Polling struct {
	command  pollingProxy
	workChan []chan bool
	pool     concurrency.WorkerPool
}

// NewServer 新建轮询器
func NewServer(command *command.Adapters, workerPool *concurrency.GoWorkerPool) *Polling {
	s := &Polling{
		command: command,
		pool:    workerPool,
	}

	workNum := config.GetPollingTaskConfig().WorkNum
	s.workChan = make([]chan bool, 0, workNum)
	for ; workNum > 0; workNum-- {
		tempChan := make(chan bool)
		s.workChan = append(s.workChan, tempChan)
	}
	return s
}

// Start 启动
func (p *Polling) Start() error {
	log.Infof("start polling server")
	for _, ch := range p.workChan {
		if err := p.pool.Submit(func() {
			startWork(ch, p.command)
		}); err != nil {
			log.Errorf("polling service start failed, err:%w", err)
		}
		// 避免同一时刻启动多个工作线程 分批启动减小误差精度
		time.Sleep(workTimeInterval)
	}
	return nil
}

// Type 处理器类型
func (p *Polling) Type() string {
	return reflect.TypeOf(p).Elem().Name()
}

// Restart 重启
func (p *Polling) Restart() error {
	if err := p.Stop(); err != nil {
		return fmt.Errorf("failed to restart external polling processor: %w", err)
	}
	return p.Start()
}

// Stop 停止
func (p *Polling) Stop() error {
	log.Infof("stop polling server")
	for _, stop := range p.workChan {
		stop <- true
	}
	return nil
}

// startWork 执行工作协程 通过chan控制轮询停止
func startWork(stop chan bool, command pollingProxy) {
	for {
		work(command)
		select {
		case <-stop:
			log.Infof("polling service work stop")
			close(stop)
			return
		default:
			time.Sleep(time.Duration(config.GetPollingTaskConfig().WorkSleepSecond) * time.Second)
		}
	}
}

// work 获取时间片并发送时间任务
func work(command pollingProxy) {
	// 获取时间片任务所
	timeSlice, err := command.GetPollingTaskWorkLock()
	if err != nil {
		log.Debugf("GetPollingTaskWorkLock GoroutineID:%d err %v", utils.GetCurrentGoroutineID(), err)
		return
	}
	log.Infof("GoroutineID:%d GetTimeSlicingWorkLock success timeAt:%v timeSlice is: %v",
		utils.GetCurrentGoroutineID(), time.Now(), timeSlice)
	if err := command.SendPollingTaskWork(timeSlice); err != nil {
		log.Errorf("Failed to SendPollingTaskWork GoroutineID:%d caused by %v", utils.GetCurrentGoroutineID(),
			err)
		return
	}
	log.Infof("GoroutineID:%d SendPollingTaskWork success timeAt:%v timeSlice is: %v",
		utils.GetCurrentGoroutineID(), time.Now(), timeSlice)
}
