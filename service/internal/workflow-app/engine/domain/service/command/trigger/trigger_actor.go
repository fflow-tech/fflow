package trigger

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// Actor 触发器的响应接口
type Actor interface {
	OnStartWorkflow(ctx context.Context, trigger *entity.Trigger, actionArgs interface{}) error // 触发启动流程
	OnRerunNode(ctx context.Context, trigger *entity.Trigger, actionArgs interface{}) error     // 触发重跑节点
	OnResumeNode(ctx context.Context, trigger *entity.Trigger, actionArgs interface{}) error    // 触发恢复节点
	OnCompleteNode(ctx context.Context, trigger *entity.Trigger, actionArgs interface{}) error  // 触发标记节点完成
}
