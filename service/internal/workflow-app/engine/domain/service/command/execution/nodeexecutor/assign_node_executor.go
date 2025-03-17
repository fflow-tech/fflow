// Package nodeexecutor 提供节点实例的执行能力
package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// AssignNodeExecutor 赋值节点执行器实现
// 实现给全局变量赋值的功能
type AssignNodeExecutor struct {
	workflowInstRepo ports.WorkflowInstRepository
	exprEvaluator    expr.Evaluator
}

// AsyncComplete 是否异步完成
func (d *AssignNodeExecutor) AsyncComplete(inst *entity.NodeInst) bool {
	return false
}

// AsyncByTrigger 通过触发器实现异步
func (d *AssignNodeExecutor) AsyncByTrigger(inst *entity.NodeInst) bool {
	return false
}

// AsyncByPolling 通过轮询实现异步
func (d *AssignNodeExecutor) AsyncByPolling(inst *entity.NodeInst) bool {
	return false
}

// NewAssignNodeExecutor 初始化执行器
func NewAssignNodeExecutor(workflowInstRepo ports.WorkflowInstRepository,
	exprEvaluator expr.Evaluator) *AssignNodeExecutor {
	return &AssignNodeExecutor{workflowInstRepo: workflowInstRepo,
		exprEvaluator: exprEvaluator}
}

// Execute 执行节点
func (d *AssignNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeDef, err := entity.ToActualNodeDef(nodeInst.BasicNodeDef.Type, nodeInst.NodeDef)
	if err != nil {
		return err
	}
	actualNodeDef := nodeDef.(entity.AssignNodeDef)
	assignKey := &entity.AssignKey{}
	for _, key := range actualNodeDef.Assign {
		if len(key.Variables) > 0 {
			assignKey.Variables = key.Variables
		}
		if len(key.Biz) > 0 {
			assignKey.Biz = key.Biz
		}
		if len(key.Owner) > 0 {
			assignKey.Owner = key.Owner
		}
	}
	return d.handleAssignKey(nodeInst, assignKey)
}

func (d *AssignNodeExecutor) handleAssignKey(nodeInst *entity.NodeInst, key *entity.AssignKey) error {
	inst, err := d.workflowInstRepo.Get(&dto.GetWorkflowInstDTO{
		InstID: nodeInst.InstID,
		DefID:  nodeInst.DefID,
	})
	if err != nil {
		return err
	}

	if err := d.assign(inst, key); err != nil {
		return err
	}

	return d.updateInst(inst)
}

func (d *AssignNodeExecutor) updateInst(inst *entity.WorkflowInst) error {
	updateDTO, err := convertor.InstConvertor.ConvertEntityToUpdateDTO(inst)
	if err != nil {
		log.Errorf("Failed to convert entity to update dto, caused by %s", err)
		return err
	}
	return d.workflowInstRepo.UpdateWithDefID(updateDTO)
}

// evaluateMap 根据上下文计算实际的值
func (d *AssignNodeExecutor) evaluateMap(inst *entity.WorkflowInst, oldMap map[string]interface{}) (
	map[string]interface{}, error) {
	ctx, err := entity.ConvertToCtx(inst)
	if err != nil {
		return nil, err
	}
	newMap, err := d.exprEvaluator.EvaluateMap(ctx, oldMap)
	if err != nil {
		return nil, err
	}
	return newMap, nil
}

// assign 赋值
func (d *AssignNodeExecutor) assign(inst *entity.WorkflowInst, key *entity.AssignKey) error {
	// 1. 给全局变量赋值
	if err := d.assignVariables(inst, key); err != nil {
		return err
	}

	// 2. 给业务字段赋值
	if err := d.assignBiz(inst, key); err != nil {
		return err
	}

	// 3. 给所有者字段赋值
	return d.assignOwner(inst, key)
}

func (d *AssignNodeExecutor) assignBiz(inst *entity.WorkflowInst, key *entity.AssignKey) error {
	biz, err := d.evaluateMap(inst, key.Biz)
	if err != nil {
		return err
	}
	inst.Biz, err = utils.MergeMap(inst.Biz, biz)
	if err != nil {
		return err
	}

	return nil
}

func (d *AssignNodeExecutor) assignVariables(inst *entity.WorkflowInst, key *entity.AssignKey) error {
	variables, err := d.evaluateMap(inst, key.Variables)
	if err != nil {
		return err
	}
	inst.Variables, err = utils.MergeMap(inst.Variables, variables)
	if err != nil {
		return err
	}
	return nil
}

func (d *AssignNodeExecutor) assignOwner(inst *entity.WorkflowInst, key *entity.AssignKey) error {
	owner, err := d.evaluateMap(inst, key.Owner)
	if err != nil {
		return err
	}

	if len(owner) > 0 {
		chatGroup := owner["chatGroup"]
		if chatGroup != nil {
			inst.Owner.ChatGroup = chatGroup.(string)
		}
		wechat := owner["wechat"]
		if wechat != nil {
			inst.Owner.Wechat = wechat.(string)
		}
	}
	return nil
}

// Polling 轮询节点
func (d *AssignNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	return nil
}

// Cancel 取消执行节点
func (d *AssignNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// Type 获取是哪种节点类型的处理器
func (d *AssignNodeExecutor) Type() entity.NodeType {
	return entity.AssignNode
}
