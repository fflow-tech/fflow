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
	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// ExternalEventProcessor 处理器结构体定义
type ExternalEventProcessor struct {
	domainService    *service.DomainService
	eventBusRepo     ports.EventBusRepository
	eventHandlersMap map[event.ExternalEventType][]func(context.Context, *dto.ExternalEventDTO) error
	consumers        []mq.Consumer
	consumerNum      int
}

// NewExternalEventProcessor 初始化
func NewExternalEventProcessor(domainService *service.DomainService,
	eventBusRepo *repo.EventBusRepo) *ExternalEventProcessor {
	r := &ExternalEventProcessor{
		domainService: domainService,
		eventBusRepo:  eventBusRepo,
		consumerNum:   config.GetEventConfig().ExternalEventConsumerNum,
	}
	eventHandler := &externalEventHandler{domainService: domainService}
	r.eventHandlersMap = map[event.ExternalEventType][]func(context.Context, *dto.ExternalEventDTO) error{
		event.WorkflowFail:    {eventHandler.handleForWorkflowExceptionHappened},
		event.WorkflowCancel:  {eventHandler.handleForWorkflowExceptionHappened},
		event.WorkflowTimeout: {eventHandler.handleForWorkflowExceptionHappened},
		event.NodeFail:        {eventHandler.handleForNodeInstExceptionHappened},
		event.NodeCancel:      {eventHandler.handleForNodeInstExceptionHappened},
		event.NodeTimeout:     {eventHandler.handleForNodeInstExceptionHappened},
	}
	return r
}

// Type 处理器类型
func (p *ExternalEventProcessor) Type() string {
	return reflect.TypeOf(p).Elem().Name()
}

// Start 启动
func (p *ExternalEventProcessor) Start() error {
	for i := 0; i < p.consumerNum; i++ {
		bizKey := strings.Join([]string{p.Type(), strconv.Itoa(i)}, "_")
		consumer, err := p.eventBusRepo.NewExternalEventConsumer(context.Background(),
			config.GetEventConfig().ExternalEventGroup,
			func(ctx context.Context, message interface{}) (err error) {
				startTime := time.Now()
				defer logs.DumpPanicStack(p.Type(), fmt.Errorf("consume external event failed"))
				defer func() {
					recordTDMQConsumeLog(&logs.EventRecord{
						Message:   message,
						BizKey:    bizKey,
						StartTime: startTime,
						Error:     err,
					})
				}()
				eventType, err := p.eventBusRepo.GetExternalEventType(message)
				if err != nil {
					return err
				}
				// 处理 Webhook
				req := &dto.ExternalEventDTO{Message: message.(pulsar.Message)}
				if err = p.domainService.Commands.ConsumeForSendWebhook(ctx, req); err != nil {
					log.Warnf("Failed to consume [%s] for send webhook, caused by %s", eventType, err)
				}
				if err = p.domainService.Commands.ConsumeForSendChatMsg(ctx, req); err != nil {
					log.Warnf("Failed to consume [%s] for send wechat remind, caused by %s", eventType, err)
				}

				eventHandlers, exists := p.eventHandlersMap[event.ExternalEventType(eventType)]
				if !exists {
					log.Debugf("The eventType=%s handler not exists", eventType)
				}
				for _, handler := range eventHandlers {
					if err := handler(ctx, req); err != nil {
						log.Warnf("Failed to consume external event, caused by %s", err)
					}
				}
				return nil
			})

		if err != nil {
			return fmt.Errorf("start external event consume failed, err: %w", err)
		}

		p.consumers = append(p.consumers, consumer)
	}
	return nil
}

// Restart 重启
func (p *ExternalEventProcessor) Restart() error {
	err := p.Stop()
	if err != nil {
		return fmt.Errorf("failed to restart external event processor: %w", err)
	}
	return p.Start()
}

// Stop 停止
func (p *ExternalEventProcessor) Stop() error {
	var wg sync.WaitGroup
	for _, consumer := range p.consumers {
		wg.Add(1)
		go func(consumer mq.Consumer) {
			defer wg.Done()
			if err := consumer.Close(); err != nil {
				// 如果关闭不成功, 其余的还是要继续关闭
				log.Errorf("Failed to close external event consumer, caused by %s", err)
			}
		}(consumer)
	}
	wg.Wait()
	return nil
}

// externalEventHandler 外部事件 handler
type externalEventHandler struct {
	domainService           *service.DomainService
	externalEventHandlerMap map[event.ExternalEventType]func(message interface{})
}

// handleForWorkflowExceptionHappened 消费流程级别异常
func (h *externalEventHandler) handleForWorkflowExceptionHappened(ctx context.Context,
	req *dto.ExternalEventDTO) error {
	return h.domainService.Commands.ConsumeForWorkflowExceptionHappened(ctx, req)
}

// handleForNodeInstExceptionHappened 消费流程节点级别异常
func (h *externalEventHandler) handleForNodeInstExceptionHappened(ctx context.Context,
	req *dto.ExternalEventDTO) error {
	return h.domainService.Commands.ConsumeForNodeInstExceptionHappened(ctx, req)
}
