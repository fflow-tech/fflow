// Package execution 提供节点实例的执行能力
package execution

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/remote"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// EventNodeExecutor 事件节点执行器实现
// 事件节点执行器
type EventNodeExecutor struct {
	remoteRepo       ports.RemoteRepository
	workflowInstRepo ports.WorkflowInstRepository
	exprEvaluator    expr.Evaluator
}

// AsyncComplete 是否异步完成
func (d *EventNodeExecutor) AsyncComplete(inst *entity.NodeInst) bool {
	return inst.BasicNodeDef.AsyncComplete
}

// AsyncByTrigger 通过触发器实现异步
func (d *EventNodeExecutor) AsyncByTrigger(inst *entity.NodeInst) bool {
	return inst.BasicNodeDef.AsyncComplete
}

// AsyncByPolling 通过轮询实现异步
func (d *EventNodeExecutor) AsyncByPolling(inst *entity.NodeInst) bool {
	return false
}

// NewEventNodeExecutor 初始化执行器
func NewEventNodeExecutor(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet) *EventNodeExecutor {
	return &EventNodeExecutor{
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		exprEvaluator:    workflowProviderSet.ExprEvaluator(),
		remoteRepo:       repoProviderSet.RemoteRepo(),
	}
}

var (
	defaultEventSource = "workflow"
	defaultEventType   = "workflow"
)

// Execute 执行节点
func (d *EventNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeDef, err := entity.ToActualNodeDef(nodeInst.BasicNodeDef.Type, nodeInst.NodeDef)
	if err != nil {
		return err
	}
	actualNodeDef := nodeDef.(entity.EventNodeDef)

	if err := d.validateArgs(actualNodeDef.Args); err != nil {
		return err
	}

	req := &remote.SendCloudEventDTO{
		Target: actualNodeDef.Args.Target,
		Event: &remote.CloudEvent{
			ID:              generateCloudEventID(nodeInst),
			Source:          defaultEventSource,
			Type:            defaultEventType,
			DataContentType: actualNodeDef.Args.Event.DataContentType,
		},
	}
	if req.Event.Data, err = d.buildEventDataByCtx(nodeInst, actualNodeDef.Args.Event.Data); err != nil {
		return err
	}
	if nodeInst.Input, err = utils.StructToMap(req); err != nil {
		return err
	}

	return d.remoteRepo.SendCloudEvent(ctx, req)
}

func (d *EventNodeExecutor) buildEventDataByCtx(nodeInst *entity.NodeInst,
	oldEventData map[string]interface{}) (map[string]interface{}, error) {
	workflowInst, err := d.workflowInstRepo.Get(&dto.GetWorkflowInstDTO{
		InstID: nodeInst.InstID,
		DefID:  nodeInst.DefID,
	})
	if err != nil {
		return nil, err
	}

	result, err := evaluateMapForCurNodeInst(d.exprEvaluator, workflowInst, nodeInst, oldEventData)
	if err != nil {
		return nil, fmt.Errorf("[%s]failed to evaluate event data: %w",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), err)
	}

	return result, nil
}

// generateCloudEventID 通过节点相关信息生成云事件 ID
func generateCloudEventID(nodeInst *entity.NodeInst) string {
	return fmt.Sprintf("%d-%d-%d-%s", nodeInst.DefID, nodeInst.InstID,
		nodeInst.NodeInstID, nodeInst.BasicNodeDef.RefName)
}

// validateArgs 入参检查
func (d *EventNodeExecutor) validateArgs(args entity.EventArgs) error {
	if !utils.IsValidURL(args.Target) {
		return fmt.Errorf("event target must not be empty")
	}
	if args.Event.Source == "" {
		return fmt.Errorf("event source must not be empty")
	}
	if args.Event.Type == "" {
		return fmt.Errorf("event type must not be empty")
	}
	if len(args.Event.Data) == 0 {
		return fmt.Errorf("event data type must not be empty")
	}
	return nil
}

// Polling 轮询节点
func (d *EventNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *EventNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// Type 获取是哪种节点类型的处理器
func (d *EventNodeExecutor) Type() entity.NodeType {
	return entity.EventNode
}
