package nodeexecutor

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
)

// ServiceFAASNodeExecutor FAAS 服务节点执行器实现
type ServiceFAASNodeExecutor struct {
	remoteRepo ports.RemoteRepository
}

// NewServiceFAASNodeExecutor 返回 FAAS 服务节点执行器
func NewServiceFAASNodeExecutor(remoteRepo ports.RemoteRepository) *ServiceFAASNodeExecutor {
	return &ServiceFAASNodeExecutor{remoteRepo: remoteRepo}
}

// Execute 执行节点
func (d *ServiceFAASNodeExecutor) Execute(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.FAASArgs)
	nodeInst.Input = args.Body
	rsp, err := d.call(ctx, args)
	if err != nil {
		return err
	}

	nodeInst.Output = rsp
	return nil
}

// Polling 轮询节点
func (d *ServiceFAASNodeExecutor) Polling(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.FAASArgs)
	nodeInst.PollInput = args.Body
	rsp, err := d.call(ctx, args)
	if err != nil {
		nodeInst.PollFailedCount += 1
		nodeInst.Reason.PollFailedReason = err.Error()
		return err
	}

	nodeInst.PollOutput = rsp
	return nil
}

// Cancel 取消执行节点
func (d *ServiceFAASNodeExecutor) Cancel(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.FAASArgs)
	nodeInst.CancelInput = args.Body
	rsp, err := d.call(ctx, args)
	if err != nil {
		return err
	}

	nodeInst.CancelOutput = rsp
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// call 组装 request 并发送 rpc 请求
func (d *ServiceFAASNodeExecutor) call(ctx context.Context, args *entity.FAASArgs) (map[string]interface{}, error) {
	err := d.validateArgs(args)
	if err != nil {
		return nil, fmt.Errorf("illegal args, err: %w", err)
	}

	req := convertor.AbilityArgsConvertor.ConvertEntityToCallFAASDTO(args)
	rsp, err := d.remoteRepo.CallFAAS(ctx, req)
	return rsp, err
}

// validateArgs 入参检查
func (d *ServiceFAASNodeExecutor) validateArgs(args *entity.FAASArgs) error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace must not be empty")
	}
	if args.Func == "" {
		return fmt.Errorf("func must not be empty")
	}
	return nil
}
