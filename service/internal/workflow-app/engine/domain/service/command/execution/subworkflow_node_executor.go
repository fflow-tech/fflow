package execution

import (
	"context"
	"fmt"
	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
)

// SubWorkflowExecutor 等待节点执行器实现
type SubWorkflowExecutor struct {
	workflowRunner    WorkflowRunner
	workflowDefRepo   ports.WorkflowDefRepository
	workflowInstRepo  ports.WorkflowInstRepository
	instExprEvaluator *common.InstExprEvaluator
}

// NewSubWorkflowExecutor 新建
func NewSubWorkflowExecutor(workflowRunner WorkflowRunner, workflowDefRepo ports.WorkflowDefRepository,
	repoProviderSet *ports.RepoProviderSet, instExprEvaluator *common.InstExprEvaluator) *SubWorkflowExecutor {
	return &SubWorkflowExecutor{
		workflowRunner:    workflowRunner,
		workflowDefRepo:   workflowDefRepo,
		workflowInstRepo:  repoProviderSet.WorkflowInstRepo(),
		instExprEvaluator: instExprEvaluator,
	}
}

// Execute 执行节点
func (d *SubWorkflowExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeDef, err := entity.ToActualNodeDef(nodeInst.BasicNodeDef.Type, nodeInst.NodeDef)
	if err != nil {
		return err
	}
	actualNodeDef := nodeDef.(entity.SubworkflowNodeDef)

	// 获取子流程定义并启动流程实例
	subWorkflowDef, err := d.getSubworkflowDef(nodeInst, actualNodeDef)
	if err != nil {
		return err
	}
	instID, err := d.startSubworkflow(nodeInst, subWorkflowDef, actualNodeDef.Args)
	if err != nil {
		return err
	}

	// nodeInst 中更新子流程信息
	nodeInst.SubworkflowDefID = subWorkflowDef.DefID
	nodeInst.SubworkflowInstID = instID
	log.Infof("execute subworkflow node success, start a subWorkflow with instID:[%d]", instID)
	return nil
}

// startSubworkflow 启动子流程
func (d *SubWorkflowExecutor) startSubworkflow(nodeInst *entity.NodeInst,
	subWorkflowDef *entity.WorkflowDef, args entity.SubworkflowArgs) (string, error) {
	startReq, err := d.getStartReqDTO(nodeInst, subWorkflowDef, args)
	if err != nil {
		return "", err
	}
	instID, err := d.workflowRunner.Start(context.Background(), startReq)
	if err != nil {
		return "", err
	}

	return instID, nil
}

// getStartReqDTO 获取启动流程参数
func (d *SubWorkflowExecutor) getStartReqDTO(nodeInst *entity.NodeInst, subWorkflowDef *entity.WorkflowDef,
	args entity.SubworkflowArgs) (*dto.StartWorkflowInstDTO, error) {
	// 拼装得到流程启动参数
	startReason := fmt.Sprintf("Start subworkflow [subworkflowDefID:%s, refName:%s] for defID:[%s]|instID:[%s]|"+
		"nodeInstID:[%s]", subWorkflowDef.DefID, subWorkflowDef.RefName,
		nodeInst.DefID, nodeInst.InstID, nodeInst.NodeInstID)
	// 从 ctx 中解析流程入参中的变量
	query := dto.NewGetWorkflowInstDTO(nodeInst.InstID, nodeInst.DefID, "")
	input, err := d.instExprEvaluator.EvaluateMapByInstCtx(query, args.Input)
	if err != nil {
		return nil, err
	}
	startReq := &dto.StartWorkflowInstDTO{
		DefID:            subWorkflowDef.DefID,
		ParentInstID:     nodeInst.InstID,
		ParentNodeInstID: nodeInst.NodeInstID,
		Name:             args.Name,
		Creator:          args.Operator,
		Input:            input,
		Reason:           startReason,
	}

	// 从节点启动时的 input 作为子流程的 input
	if len(nodeInst.Input) > 0 {
		inputOfNodeInst, err := d.instExprEvaluator.EvaluateMapByInstCtx(query, nodeInst.Input)
		if err != nil {
			return nil, err
		}
		startReq.Input = inputOfNodeInst
	}

	return startReq, nil
}

// getSubworkflowDef 获取子流程定义
func (d *SubWorkflowExecutor) getSubworkflowDef(nodeInst *entity.NodeInst, actualNodeDef entity.SubworkflowNodeDef) (
	*entity.WorkflowDef, error) {
	// 如果配置了子流程的 ID，则直接根据 ID 查询；同时配置 ID 和 subworkflow 以 ID 为准
	if actualNodeDef.ID != "" {
		return d.getSubworkflowDefByIDAndVersion(actualNodeDef.ID, actualNodeDef.Version)
	}
	// 查询得到子流程定义
	getDefReq := &dto.GetSubworkflowDefDTO{
		RefName:          actualNodeDef.Subworkflow,
		ParentDefID:      nodeInst.DefID,
		ParentDefVersion: nodeInst.DefVersion,
	}
	subworkflowDef, err := d.workflowDefRepo.GetSubworkflowLastVersion(getDefReq)
	if err != nil {
		return nil, err
	}

	return subworkflowDef, nil
}

// getSubworkflowDefByIDAndVersion 根据 id 和 version 查询流程定义
func (d *SubWorkflowExecutor) getSubworkflowDefByIDAndVersion(id string, version int) (*entity.WorkflowDef, error) {
	if version > 0 {
		return d.workflowDefRepo.Get(&dto.GetWorkflowDefDTO{DefID: id, Version: version})
	}
	// 未填写 version 时则使用最新版本
	return d.workflowDefRepo.GetLastVersion(&dto.GetWorkflowDefDTO{DefID: id})
}

// Polling 轮询节点
func (d *SubWorkflowExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *SubWorkflowExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	cancelReqDTO := &dto.CancelWorkflowInstDTO{
		DefID:    nodeInst.SubworkflowDefID,
		InstID:   nodeInst.SubworkflowInstID,
		Operator: nodeInst.Operator.CancelledOperator,
	}
	err := d.workflowRunner.Cancel(ctx, cancelReqDTO)
	if err != nil {
		return err
	}

	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// AsyncComplete 是否异步完成
func (d *SubWorkflowExecutor) AsyncComplete(nodeInst *entity.NodeInst) bool {
	return d.AsyncByTrigger(nodeInst) || d.AsyncByPolling(nodeInst)
}

// AsyncByTrigger 通过触发器实现异步
func (d *SubWorkflowExecutor) AsyncByTrigger(nodeInst *entity.NodeInst) bool {
	// 子流程节点的完成是根据流程的完成事件来异步确定的
	return true
}

// AsyncByPolling 通过轮询实现异步
func (d *SubWorkflowExecutor) AsyncByPolling(nodeInst *entity.NodeInst) bool {
	return false
}

// Type 获取是哪种节点类型的处理器
func (d *SubWorkflowExecutor) Type() entity.NodeType {
	return entity.SubWorkflowNode
}
