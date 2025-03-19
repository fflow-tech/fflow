package nodeexecutor

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// ServiceMCPNodeExecutor implements a node executor for MCP protocol
type ServiceMCPNodeExecutor struct {
	remoteRepo ports.RemoteRepository
}

// NewServiceMCPNodeExecutor creates a new instance of ServiceMCPNodeExecutor
func NewServiceMCPNodeExecutor(remoteRepo ports.RemoteRepository) *ServiceMCPNodeExecutor {
	return &ServiceMCPNodeExecutor{remoteRepo: remoteRepo}
}

// Execute 执行节点
func (d *ServiceMCPNodeExecutor) Execute(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.MCPArgs)
	if args.Body != nil {
		nodeInst.Input = args.Body
	} else {
		nodeInst.Input = utils.StringMapToInterfaceMap(args.Parameters)
	}

	rsp, err := d.call(ctx, args)
	if err != nil {
		return err
	}

	nodeInst.Output = rsp
	return nil
}

// Polling 轮询节点
func (d *ServiceMCPNodeExecutor) Polling(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.MCPArgs)
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
func (d *ServiceMCPNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.MCPArgs)
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
func (d *ServiceMCPNodeExecutor) call(ctx context.Context, args *entity.MCPArgs) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
