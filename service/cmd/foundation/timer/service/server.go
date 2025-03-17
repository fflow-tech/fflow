package service

import (
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/factory"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service/notify"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service/polling"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service/timer"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/mq/eventbus"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"sync"

	"github.com/fflow-tech/fflow/service/pkg/log"
)

var (
	mutex sync.Mutex
)

// Server 默认事件服务
type Server struct {
	processorMap map[string]Processor
}

// NewServer 新建注册中心
func NewServer() *Server {
	return &Server{processorMap: map[string]Processor{}}
}

// Processor 消息控制器
type Processor interface {
	// Type 处理器类型
	Type() string
	// Start 启动
	Start() error
	// Restart 重启
	Restart() error
	// Stop 停止
	Stop() error
}

// Register 注册
func (r *Server) Register(processor Processor) {
	mutex.Lock()
	defer mutex.Unlock()
	r.processorMap[processor.Type()] = processor
}

// Server 启动所有处理器
func (r *Server) Server() error {
	mutex.Lock()
	defer mutex.Unlock()
	for _, c := range r.processorMap {
		if err := c.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭所有处理器
func (r *Server) Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	var wg sync.WaitGroup
	for _, c := range r.processorMap {
		wg.Add(1)
		go func(c Processor) {
			defer wg.Done()
			if err := c.Stop(); err != nil {
				log.Errorf("Failed to stop processor, caused by %s", err)
			}
		}(c)
	}
	wg.Wait()
	return nil
}

var (
	taskServer = NewServer()
)

// InitTimerTaskServer 初始化定时器服务
func InitTimerTaskServer() {
	log.Infof("InitTimerTaskServer")
	eventBusRepo, err := factory.GetEventBusRepo()
	if err != nil {
		panic(err)
	}

	command, err := factory.GetCommand()
	if err != nil {
		panic(err)
	}

	timerServer := timer.NewTimerTaskEventProcessor(command, eventBusRepo,
		eventbus.GetLogRecorder(), concurrency.GetDefaultWorkerPool())
	pollingServer := polling.NewServer(command, concurrency.GetDefaultWorkerPool())
	notifyServer := notify.NewNotifyEventProcessor(command, eventBusRepo,
		eventbus.GetLogRecorder(), concurrency.GetDefaultWorkerPool())

	// 注册有先后顺序 先注册通知服务 再注册定时任务服务 再注册轮询服务
	taskServer.Register(notifyServer)
	taskServer.Register(timerServer)
	taskServer.Register(pollingServer)
	if err := taskServer.Server(); err != nil {
		panic(err)
	}
}

// CloseEventServer 停止事件处理
func CloseEventServer(ch chan struct{}) error {
	log.Infof("Shutdown Task Server...")
	defer log.Infof("Task Server exit")
	defer func() {
		if ch != nil {
			ch <- struct{}{}
		}
	}()
	if err := taskServer.Close(); err != nil {
		log.Errorf("Failed to close task server, caused by %s", err)
		return err
	}
	return nil
}
