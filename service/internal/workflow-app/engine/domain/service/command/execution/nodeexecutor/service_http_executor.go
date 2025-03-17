package nodeexecutor

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// ServiceHTTPNodeExecutor HTTP 服务节点执行器实现
type ServiceHTTPNodeExecutor struct {
	remoteRepo ports.RemoteRepository
}

// NewServiceHTTPNodeExecutor 返回 HTTP 服务节点执行器
func NewServiceHTTPNodeExecutor(remoteRepo ports.RemoteRepository) *ServiceHTTPNodeExecutor {
	return &ServiceHTTPNodeExecutor{remoteRepo: remoteRepo}
}

// Execute 执行节点
func (d *ServiceHTTPNodeExecutor) Execute(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.HTTPArgs)
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
func (d *ServiceHTTPNodeExecutor) Polling(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.HTTPArgs)
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
func (d *ServiceHTTPNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.HTTPArgs)
	nodeInst.CancelInput = args.Body
	rsp, err := d.call(ctx, args)
	if err != nil {
		return err
	}

	nodeInst.CancelOutput = rsp
	nodeInst.Status = entity.NodeInstCancelled
	return nil
}

// call 组装 request 并发送 http 请求
func (d *ServiceHTTPNodeExecutor) call(ctx context.Context, args *entity.HTTPArgs) (map[string]interface{}, error) {
	err := d.validateArgs(args)
	if err != nil {
		return nil, fmt.Errorf("illegal args: %w", err)
	}

	req := convertor.AbilityArgsConvertor.ConvertEntityToCallHTTPDTO(args)

	rsp, err := d.remoteRepo.CallHTTP(ctx, req)
	return rsp, err
}

// validateArgs 入参检查
func (d *ServiceHTTPNodeExecutor) validateArgs(args *entity.HTTPArgs) error {
	if args.Method == "" {
		return errors.New("method is empty")
	}
	if args.URL == "" {
		return errors.New("url is empty")
	}
	return nil
}
