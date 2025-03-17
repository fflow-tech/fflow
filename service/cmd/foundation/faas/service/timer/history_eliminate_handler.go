package timer

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// EliminateTimer 超时检查处理器
type EliminateTimer struct {
	domainService *service.DomainService
}

// NewEliminateTimerHandler 初始化
func NewEliminateTimerHandler(domainService *service.DomainService) *EliminateTimer {
	return &EliminateTimer{domainService: domainService}
}

// Handle 定期淘汰老旧的执行历史的执行函数
func (h *EliminateTimer) Handle(ctx context.Context) error {
	log.InfoContext(ctx, "Do BatchDeleteExpiredRunHistory processing!")
	return h.domainService.Commands.BatchDeleteExpiredRunHistory(&dto.BatchDeleteExpiredRunHistoryDTO{
		KeepDays: config.GetAppConfig().KeepDays},
	)
}
