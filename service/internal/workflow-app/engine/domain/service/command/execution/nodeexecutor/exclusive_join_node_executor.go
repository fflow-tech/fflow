package nodeexecutor

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
)

// ExclusiveJoinNodeExecutor ExclusiveJoin 节点执行器
// 只要有一个父节点运行完成，ExclusiveJoin 则节点执行完成
type ExclusiveJoinNodeExecutor struct {
	workflowInstRepo ports.WorkflowInstRepository
}

// NewExclusiveJoinNodeExecutor 初始化执行器
func NewExclusiveJoinNodeExecutor(workflowInstRepo ports.WorkflowInstRepository) *ExclusiveJoinNodeExecutor {
	return &ExclusiveJoinNodeExecutor{
		workflowInstRepo: workflowInstRepo,
	}
}

// AsyncComplete 是否异步完成
func (d *ExclusiveJoinNodeExecutor) AsyncComplete(inst *entity.NodeInst) bool {
	return false
}

// AsyncByTrigger 通过触发器异步完成
func (d *ExclusiveJoinNodeExecutor) AsyncByTrigger(nodeInst *entity.NodeInst) bool {
	return false
}

// AsyncByPolling 通过轮询异步完成
func (d *ExclusiveJoinNodeExecutor) AsyncByPolling(nodeInst *entity.NodeInst) bool {
	return false
}

// Execute 执行节点
func (d *ExclusiveJoinNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	completedParentNode, err := d.getCompletedParentNode(nodeInst)
	if err != nil {
		return err
	}
	nodeInst.Input = completedParentNode.Input
	nodeInst.Output = completedParentNode.Output
	nodeInst.PollInput = completedParentNode.PollInput
	nodeInst.PollOutput = completedParentNode.PollOutput
	return nil
}

// getCompletedParentNode 获取已完成的父节点
func (d *ExclusiveJoinNodeExecutor) getCompletedParentNode(nodeInst *entity.NodeInst) (*entity.NodeInst, error) {
	nodeDef, err := entity.ToActualNodeDef(nodeInst.BasicNodeDef.Type, nodeInst.NodeDef)
	if err != nil {
		return nil, err
	}
	actualNodeDef := nodeDef.(entity.ExclusiveJoinNodeDef)

	inst, err := d.workflowInstRepo.Get(dto.NewGetWorkflowInstDTO(nodeInst.InstID, nodeInst.DefID, ""))
	if err != nil {
		return nil, err
	}

	for _, tmpNodeInst := range inst.SchedNodeInsts {
		if !tmpNodeInst.Status.IsTerminal() {
			continue
		}
		next, err := entity.GetNextNode(inst.WorkflowDef, &tmpNodeInst.BasicNodeDef)
		if err != nil {
			log.Warnf("[%s]Failed to get next node, caused by %s",
				logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), err)
			continue
		}
		if next == actualNodeDef.RefName {
			return tmpNodeInst, nil
		}
	}

	return nil, fmt.Errorf("[%s]not found any parent node compelted for node=%s",
		logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), actualNodeDef.RefName)
}

// Polling 轮询节点
func (d *ExclusiveJoinNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *ExclusiveJoinNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// Type 节点类型
func (d *ExclusiveJoinNodeExecutor) Type() entity.NodeType {
	return entity.ExclusiveJoinNode
}
