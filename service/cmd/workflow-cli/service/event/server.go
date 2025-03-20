// Package event 提供消费事件的入口
package event

import (
	"sync"

	"github.com/fflow-tech/fflow/service/cmd/workflow-app/engine/service/event"
	"github.com/fflow-tech/fflow/service/cmd/workflow-cli/factory"
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
	r.processorMap[processor.Type()] = processor
}

// Serve 启动所有处理器
func (r *Server) Serve() error {
	mutex.Lock()
	defer mutex.Unlock()
	domainService, err := factory.GetDomainService()
	if err != nil {
		panic(err)
	}
	eventBusRepo, err := factory.GetEventBusRepo()
	if err != nil {
		panic(err)
	}

	r.Register(event.NewDriveEventProcessor(domainService, eventBusRepo))
	r.Register(event.NewTriggerEventProcessor(domainService, eventBusRepo))
	r.Register(event.NewCronEventProcessor(domainService, eventBusRepo))
	r.Register(event.NewExternalEventProcessor(domainService, eventBusRepo))

	for _, c := range r.processorMap {
		log.Infof("Start processor: %s", c.Type())
		if err := c.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭所有处理器
func (r *Server) Close(ch chan struct{}) error {
	log.Infof("Shutdown Event Server...")
	defer log.Infof("Event Server exit")
	defer func() {
		if ch != nil {
			ch <- struct{}{}
		}
	}()

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
