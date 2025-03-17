package rpc

import (
	"context"
	"fmt"
	pb "github.com/fflow-tech/fflow/api/workflow-app/engine"
	"github.com/fflow-tech/fflow/service/cmd/workflow-app/engine/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/errno"
)

// WorkflowEngineService workflow-app/engine后端服务实现
type WorkflowEngineService struct {
	pb.UnimplementedWorkflowServer
	domainService *service.DomainService
}

// NewWorkflowEngineService WorkFlowEngineService构造函数
func NewWorkflowEngineService(domainService *service.DomainService) *WorkflowEngineService {
	return &WorkflowEngineService{domainService: domainService}
}

// CreateDef 创建工作流定义
func (w *WorkflowEngineService) CreateDef(ctx context.Context, req *pb.CreateDefReq) (*pb.CreateDefRsp, error) {
	rsp := &pb.CreateDefRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	createDefDTO, err := convertor.DefConvertor.ConvertCreatePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	defID, err := w.domainService.Commands.CreateWorkflowDef(ctx, createDefDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	rsp.DefID = defID
	return rsp, nil
}

// UpdateDef 更新工作流定义
func (w *WorkflowEngineService) UpdateDef(ctx context.Context, req *pb.UpdateDefReq) (*pb.UpdateDefRsp, error) {
	rsp := &pb.UpdateDefRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	updateDefDTO, err := convertor.DefConvertor.ConvertUpdatePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.UpdateWorkflowDef(ctx, updateDefDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetDefDetail 获取工作流定义
func (w *WorkflowEngineService) GetDefDetail(ctx context.Context, req *pb.GetDefDetailReq) (
	*pb.GetDefDetailRsp, error) {
	rsp := &pb.GetDefDetailRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	getDefDTO, err := convertor.DefConvertor.ConvertGetPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	defDTO, err := w.domainService.Queries.GetWorkflowDefByDefID(ctx, getDefDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	rsp.WorkflowDef = convertor.DefConvertor.ConvertDTOToPb(defDTO)
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// EnableDef 激活工作流
func (w *WorkflowEngineService) EnableDef(ctx context.Context, req *pb.EnableDefReq) (*pb.EnableDefRsp, error) {
	rsp := &pb.EnableDefRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	enableDTO, err := convertor.DefConvertor.ConvertEnablePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.EnableWorkflowDef(ctx, enableDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// DisableDef 去激活工作流
func (w *WorkflowEngineService) DisableDef(ctx context.Context, req *pb.DisableDefReq) (*pb.DisableDefRsp, error) {
	rsp := &pb.DisableDefRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	disableDTO, err := convertor.DefConvertor.ConvertDisablePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.DisableWorkflowDef(ctx, disableDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// StartInst 启动流程实例
func (w *WorkflowEngineService) StartInst(ctx context.Context, req *pb.StartInstReq) (*pb.StartInstRsp, error) {
	rsp := &pb.StartInstRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	createInstDTO, err := convertor.InstConvertor.ConvertStartPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	defID, err := w.domainService.Commands.StartWorkflowInst(ctx, createInstDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	rsp.InstID = defID
	return rsp, nil
}

// CompleteInst 标记流程完成
func (w *WorkflowEngineService) CompleteInst(ctx context.Context, req *pb.CompleteInstReq) (
	*pb.CompleteInstRsp, error) {
	rsp := &pb.CompleteInstRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	completeDTO, err := convertor.InstConvertor.ConvertCompletePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.CompleteWorkflowInst(ctx, completeDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// CancelInst 终止流程
func (w *WorkflowEngineService) CancelInst(ctx context.Context, req *pb.CancelInstReq) (*pb.CancelInstRsp, error) {
	rsp := &pb.CancelInstRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	cancelInstDTO, err := convertor.InstConvertor.ConvertCancelPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.CancelWorkflowInst(ctx, cancelInstDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// PauseInst 暂停流程
func (w *WorkflowEngineService) PauseInst(ctx context.Context, req *pb.PauseInstReq) (*pb.PauseInstRsp, error) {
	rsp := &pb.PauseInstRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	pauseInstDTO, err := convertor.InstConvertor.ConvertPausePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.PauseWorkflowInst(ctx, pauseInstDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// ResumeInst 恢复流程
func (w *WorkflowEngineService) ResumeInst(ctx context.Context, req *pb.ResumeInstReq) (*pb.ResumeInstRsp, error) {
	rsp := &pb.ResumeInstRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	resumeInstDTO, err := convertor.InstConvertor.ConvertResumePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.ResumeWorkflowInst(ctx, resumeInstDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// UpdateInstContext 更新实例环境参数
func (w *WorkflowEngineService) UpdateInstContext(ctx context.Context, req *pb.UpdateInstContextReq) (
	*pb.UpdateInstContextRsp, error) {
	rsp := &pb.UpdateInstContextRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	updateCtxDTO, err := convertor.InstConvertor.ConvertUpdateCtxPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.UpdateWorkflowInstCtx(ctx, updateCtxDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetInstDetail 获取流程实例运行状态
func (w *WorkflowEngineService) GetInstDetail(ctx context.Context, req *pb.GetInstDetailReq) (
	*pb.GetInstDetailRsp, error) {
	rsp := &pb.GetInstDetailRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	getDTO, err := convertor.InstConvertor.ConvertGetPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	data, err := w.domainService.Queries.GetWorkflowInst(ctx, getDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	err = convertor.InstConvertor.ConvertDTOToPb(data, rsp)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetInstList 获取流程实例列表
func (w *WorkflowEngineService) GetInstList(ctx context.Context, req *pb.GetInstListReq) (*pb.GetInstListRsp, error) {
	rsp := &pb.GetInstListRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if req.DefID == "" {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, "defID must not be empty")
		return rsp, nil
	}

	getListDTO, err := convertor.InstConvertor.ConvertGetListPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	insts, total, err := w.domainService.Queries.GetWorkflowInstList(ctx, getListDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	err = convertor.InstConvertor.ConvertDTOListToPb(insts, total, rsp)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	rsp.PageIndex = req.PageIndex
	rsp.PageSize = req.PageSize
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// RestartInst 从某个节点重跑实例
func (w *WorkflowEngineService) RestartInst(ctx context.Context, req *pb.RestartInstReq) (*pb.RestartInstRsp, error) {
	rsp := &pb.RestartInstRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	restartDTO, err := convertor.InstConvertor.ConvertRestartPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.RestartWorkflowInst(ctx, restartDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetNodeInstDetail 获取节点实例信息
func (w *WorkflowEngineService) GetNodeInstDetail(ctx context.Context, req *pb.GetNodeInstDetailReq) (
	*pb.GetNodeInstDetailRsp, error) {
	rsp := &pb.GetNodeInstDetailRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	getDTO, err := convertor.NodeInstConvertor.ConvertGetPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	nodeInstDTO, err := w.domainService.Queries.GetNodeInstDetail(ctx, getDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	data, err := convertor.NodeInstConvertor.ConvertDTOToPb(nodeInstDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}
	rsp.NodeItem = data
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// SkipNode 跳过某节点
func (w *WorkflowEngineService) SkipNode(ctx context.Context, req *pb.SkipNodeReq) (
	*pb.SkipNodeRsp, error) {
	rsp := &pb.SkipNodeRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	skipDTO, err := convertor.NodeInstConvertor.ConvertSkipPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.SkipNode(ctx, skipDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// CancelSkipNode 恢复某节点(取消跳过)
func (w *WorkflowEngineService) CancelSkipNode(ctx context.Context, req *pb.CancelSkipNodeReq) (
	*pb.CancelSkipNodeRsp, error) {
	rsp := &pb.CancelSkipNodeRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	cancelSkipNode, err := convertor.NodeInstConvertor.ConvertCancelSkipPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.CancelSkipNode(ctx, cancelSkipNode)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// CompleteNode 结束节点
func (w *WorkflowEngineService) CompleteNode(ctx context.Context, req *pb.CompleteNodeReq) (
	*pb.CompleteNodeRsp, error) {
	rsp := &pb.CompleteNodeRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	completeDTO, err := convertor.NodeInstConvertor.ConvertCompletePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.CompleteNode(ctx, completeDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// CancelNode 终止指定已运行节点
func (w *WorkflowEngineService) CancelNode(ctx context.Context, req *pb.CancelNodeReq) (
	*pb.CancelNodeRsp, error) {
	rsp := &pb.CancelNodeRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	cancelDTO, err := convertor.NodeInstConvertor.ConvertCancelPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.CancelNode(ctx, cancelDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// RerunNode 重跑指定已运行节点
func (w *WorkflowEngineService) RerunNode(ctx context.Context, req *pb.RerunNodeReq) (
	*pb.RerunNodeRsp, error) {
	rsp := &pb.RerunNodeRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	rerunDTO, err := convertor.NodeInstConvertor.ConvertRerunPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.RerunNode(ctx, rerunDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// ResumeNode 继续指定等待中的节点
func (w *WorkflowEngineService) ResumeNode(ctx context.Context, req *pb.ResumeNodeReq) (
	*pb.ResumeNodeRsp, error) {
	rsp := &pb.ResumeNodeRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	resumeDTO, err := convertor.NodeInstConvertor.ConvertResumePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	err = w.domainService.Commands.ResumeNode(ctx, resumeDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// NewSucceedRsp 生成成功返回
func NewSucceedRsp() *pb.BasicRsp {
	return &pb.BasicRsp{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
	}
}

// NewFailedRsp 通过自定义的错误码生成请求返回
func NewFailedRsp(code int32, message string) *pb.BasicRsp {
	return &pb.BasicRsp{
		Code:    code,
		Message: message,
	}
}

func validateBasicReq(req *pb.BasicReq) error {
	if req == nil || req.Operator == "" {
		return fmt.Errorf("the req or req's operator must not be empty")
	}

	return nil
}
