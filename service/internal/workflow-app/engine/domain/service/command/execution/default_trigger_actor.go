// Package execution 流程的执行层
package execution

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// DefaultTriggerActor 默认触发器 Actor
type DefaultTriggerActor struct {
	eventBusRepo   ports.EventBusRepository
	workflowRunner WorkflowRunner
	nodeRunner     NodeRunner
}

// NewDefaultActor 初始化
func NewDefaultActor(repoProvider *ports.RepoProviderSet,
	workflowRunner *DefaultWorkflowRunner,
	nodeRunner *DefaultNodeRunner) *DefaultTriggerActor {
	return &DefaultTriggerActor{
		eventBusRepo:   repoProvider.EventBusRepo(),
		workflowRunner: workflowRunner,
		nodeRunner:     nodeRunner}
}

// OnStartWorkflow 响应启动工作流
func (h *DefaultTriggerActor) OnStartWorkflow(ctx context.Context,
	trigger *entity.Trigger, actionArgs interface{}) error {
	realActionArgs := actionArgs.(*entity.StartWorkflowActionArgs)

	startWorkflowDTO := &dto.StartWorkflowInstDTO{
		DefID:   trigger.DefID,
		Name:    realActionArgs.Name,
		Input:   realActionArgs.Input,
		Creator: realActionArgs.Operator,
		Reason:  fmt.Sprintf("on start workflow, trigger:%s", utils.StructToJsonStr(trigger)),
	}
	_, err := h.workflowRunner.Start(ctx, startWorkflowDTO)
	return err
}

// OnRerunNode 响应重跑指定节点
func (h *DefaultTriggerActor) OnRerunNode(ctx context.Context, trigger *entity.Trigger, actionArgs interface{}) error {
	realActionArgs := actionArgs.(*entity.RerunNodeActionArgs)

	rerunNodeDTO := &dto.RerunNodeDTO{
		InstID:      trigger.InstID,
		DefID:       trigger.DefID,
		NodeRefName: realActionArgs.Node,
		Input:       realActionArgs.Input,
		Operator:    realActionArgs.Operator,
		Reason:      fmt.Sprintf("on rerun node, trigger:%s", utils.StructToJsonStr(trigger)),
	}
	return h.nodeRunner.Rerun(ctx, rerunNodeDTO)
}

// OnResumeNode 响应恢复指定节点
func (h *DefaultTriggerActor) OnResumeNode(ctx context.Context, trigger *entity.Trigger, actionArgs interface{}) error {
	realActionArgs := actionArgs.(*entity.ResumeNodeActionArgs)

	resumeNodeDTO := &dto.ResumeNodeDTO{
		InstID:      trigger.InstID,
		DefID:       trigger.DefID,
		NodeRefName: realActionArgs.Node,
		Input:       realActionArgs.Input,
		Operator:    realActionArgs.Operator,
		Reason:      fmt.Sprintf("on resume node, trigger:%s", utils.StructToJsonStr(trigger)),
	}
	return h.nodeRunner.Resume(ctx, resumeNodeDTO)
}

// OnCompleteNode 响应标记节点完成
func (h DefaultTriggerActor) OnCompleteNode(ctx context.Context,
	trigger *entity.Trigger, actionArgs interface{}) error {
	realActionArgs := actionArgs.(*entity.CompleteNodeActionArgs)

	completeNode := &dto.CompleteNodeDTO{
		DefID:       trigger.DefID,
		InstID:      trigger.InstID,
		NodeRefName: realActionArgs.Node,
		Status:      entity.GetNodeInstStatusByStrValue(realActionArgs.Status),
		Output:      realActionArgs.Output,
		Operator:    realActionArgs.Operator,
		Reason:      fmt.Sprintf("on compete node, trigger:%s", utils.StructToJsonStr(trigger)),
	}

	return h.nodeRunner.Complete(ctx, completeNode)
}
