package event

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// CronEventProcessor 处理器
type CronEventProcessor struct {
	eventBusRepo  ports.EventBusRepository
	domainService *service.DomainService
	consumers     []mq.Consumer
	consumerNum   int
}

// NewCronEventProcessor 初始化
func NewCronEventProcessor(domainService *service.DomainService,
	eventBusRepo *repo.EventBusRepo) *CronEventProcessor {
	return &CronEventProcessor{
		eventBusRepo:  eventBusRepo,
		domainService: domainService,
		consumerNum:   config.GetEventConfig().CronEventConsumerNum,
	}
}

// Type 处理器类型
func (p *CronEventProcessor) Type() string {
	return reflect.TypeOf(p).Elem().Name()
}

// Start 启动
func (p *CronEventProcessor) Start() error {
	for i := 0; i < p.consumerNum; i++ {
		bizKey := strings.Join([]string{p.Type(), strconv.Itoa(i)}, "_")
		consumer, err := p.eventBusRepo.NewCronEventConsumer(context.Background(),
			config.GetEventConfig().CronEventGroup,
			func(ctx context.Context, message interface{}) (err error) {
				startTime := time.Now()
				defer logs.DumpPanicStack(p.Type(), fmt.Errorf("failed to consume drive event"))
				defer func() {
					recordTDMQConsumeLog(&logs.EventRecord{
						Message:   message,
						BizKey:    bizKey,
						StartTime: startTime,
						Error:     err,
					})
				}()
				return p.domainService.Commands.ConsumeCronTriggerEvent(ctx,
					&dto.CronTriggerEventDTO{Message: message.(pulsar.Message)})
			})

		if err != nil {
			return fmt.Errorf("failed to start cron event consume: %w", err)
		}

		p.consumers = append(p.consumers, consumer)
	}
	return nil
}

// Restart 重启
func (p *CronEventProcessor) Restart() error {
	err := p.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop cron event processor: %w", err)
	}

	return p.Start()
}

// Stop 停止
func (p *CronEventProcessor) Stop() error {
	var wg sync.WaitGroup
	for _, consumer := range p.consumers {
		wg.Add(1)
		go func(consumer mq.Consumer) {
			defer wg.Done()
			if err := consumer.Close(); err != nil {
				// 如果关闭不成功, 其余的还是要继续关闭
				log.Errorf("Failed to close cron event consumer, caused by %s", err)
			}
		}(consumer)
	}
	wg.Wait()
	return nil
}
