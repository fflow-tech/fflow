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
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// DriveEventProcessor 处理器定义
type DriveEventProcessor struct {
	domainService   *service.DomainService
	eventBusRepo    ports.EventBusRepository
	consumers       []mq.Consumer
	eventHandlerMap map[event.DriveEventType]func(ctx context.Context, driveEvent *dto.DriveEventDTO) error
	consumerNum     int
}

// NewDriveEventProcessor 初始化
func NewDriveEventProcessor(domainService *service.DomainService,
	eventBusRepo *repo.EventBusRepo) *DriveEventProcessor {
	p := &DriveEventProcessor{
		domainService: domainService,
		eventBusRepo:  eventBusRepo,
		consumerNum:   config.GetEventConfig().DriveEventConsumerNum,
	}
	eventHandler := newDriveEventHandler(domainService)
	p.eventHandlerMap = map[event.DriveEventType]func(ctx context.Context, driveEvent *dto.DriveEventDTO) error{
		event.WorkflowStartDrive: eventHandler.handleWorkflowStartDriveEvent,
		event.NodeScheduleDrive:  eventHandler.handleNodeScheduleDriveEvent,
		event.NodeExecuteDrive:   eventHandler.handleNodeExecuteDriveEvent,
		event.NodePollDrive:      eventHandler.handleNodePollDriveEvent,
		event.NodeCompleteDrive:  eventHandler.handleNodeCompleteDriveEvent,
		event.NodeRetryDrive:     eventHandler.handleNodeCompleteDriveEvent,
	}
	return p
}

// Type 处理器类型
func (p *DriveEventProcessor) Type() string {
	return reflect.TypeOf(p).Elem().Name()
}

// Start 启动
func (p *DriveEventProcessor) Start() error {
	for i := 0; i < p.consumerNum; i++ {
		bizKey := strings.Join([]string{p.Type(), strconv.Itoa(i)}, "_")
		consumer, err := p.eventBusRepo.NewDriveEventConsumer(context.Background(),
			config.GetEventConfig().DriveEventGroup,
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
				eventType, err := p.eventBusRepo.GetDriveEventType(message)
				if err != nil {
					return err
				}
				handle, exists := p.eventHandlerMap[event.DriveEventType(eventType)]
				if !exists {
					log.Errorf("Failed to get eventType=%s handle", eventType)
					return nil
				}

				return handle(ctx, &dto.DriveEventDTO{Message: message.(pulsar.Message)})
			})

		if err != nil {
			return fmt.Errorf("failed to start drive event consume: %w", err)
		}

		p.consumers = append(p.consumers, consumer)
	}
	return nil
}

// Restart 重启
func (p *DriveEventProcessor) Restart() error {
	err := p.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop drive event processor: %w", err)
	}

	return p.Start()
}

// Stop 停止
func (p *DriveEventProcessor) Stop() error {
	var wg sync.WaitGroup
	for _, consumer := range p.consumers {
		wg.Add(1)
		go func(consumer mq.Consumer) {
			defer wg.Done()
			if err := consumer.Close(); err != nil {
				// 如果关闭不成功, 其余的还是要继续关闭
				log.Errorf("Failed to close drive event consumer, caused by %s", err)
			}
		}(consumer)
	}
	wg.Wait()
	return nil
}

// newDriveEventHandler 实例化
func newDriveEventHandler(domainService *service.DomainService) *driveEventHandler {
	return &driveEventHandler{domainService: domainService}
}

type driveEventHandler struct {
	domainService        *service.DomainService
	driveEventHandlerMap map[event.DriveEventType]func(message interface{})
}

func (h *driveEventHandler) handleWorkflowStartDriveEvent(ctx context.Context, event *dto.DriveEventDTO) error {
	return h.domainService.Commands.ConsumeWorkflowStartDriveEvent(ctx, event)
}

func (h *driveEventHandler) handleNodeScheduleDriveEvent(ctx context.Context, event *dto.DriveEventDTO) error {
	return h.domainService.Commands.ConsumeNodeScheduleDriveEvent(ctx, event)
}

func (h *driveEventHandler) handleNodeExecuteDriveEvent(ctx context.Context, event *dto.DriveEventDTO) error {
	return h.domainService.Commands.ConsumeNodeExecuteDriveEvent(ctx, event)
}

func (h *driveEventHandler) handleNodePollDriveEvent(ctx context.Context, event *dto.DriveEventDTO) error {
	return h.domainService.Commands.ConsumeNodePollDriveEvent(ctx, event)
}

func (h *driveEventHandler) handleNodeCompleteDriveEvent(ctx context.Context, event *dto.DriveEventDTO) error {
	return h.domainService.Commands.ConsumeNodeCompleteDriveEvent(ctx, event)
}

func (h *driveEventHandler) handleNodeRetryDriveEvent(ctx context.Context, event *dto.DriveEventDTO) error {
	return h.domainService.Commands.ConsumeNodeRetryDriveEvent(ctx, event)
}
