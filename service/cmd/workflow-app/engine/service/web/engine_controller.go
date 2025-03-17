package web

import (
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/remote"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/cmd/workflow-app/engine/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowEngineController workflow-app/engine http服务实现
type WorkflowEngineController struct {
	domainService       *service.DomainService
	eventBusRepo        ports.EventBusRepository
	permissionValidator *remote.DefaultPermissionValidator
}

// NewWorkflowEngineController 构造函数
func NewWorkflowEngineController(domainService *service.DomainService,
	eventBusRepo *repo.EventBusRepo, permissionValidator *remote.DefaultPermissionValidator) *WorkflowEngineController {
	return &WorkflowEngineController{domainService: domainService,
		eventBusRepo:        eventBusRepo,
		permissionValidator: permissionValidator,
	}
}

// GetDefDetail 查询单条工作流定义
// @Summary 查询单条工作流定义
// @Description 查询单条工作流定义
// @Tags 工作流定义相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.GetWorkflowDefDTO true "查询流程请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/def/get [get]
func (h *WorkflowEngineController) GetDefDetail(c *gin.Context) {
	var req dto.GetWorkflowDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	// 如果是超管用户则可以查询任意流程的详情
	if utils.StrContains(config.GetRbacConfig().SuperAdmins, req.Creator) {
		req.Creator = ""
	}

	data, err := h.domainService.Queries.GetWorkflowDefByDefID(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// GetDefList 批量查询工作流定义
// @Summary 批量查询工作流定义
// @Description 批量查询工作流定义
// @Tags 工作流定义相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.PageQueryWorkflowDefDTO true "查询流程请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/def/list [get]
func (h *WorkflowEngineController) GetDefList(c *gin.Context) {
	var req dto.PageQueryWebWorkflowDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, total, err := h.domainService.Queries.GetWorkflowDefList(c.Request.Context(),
		convertor.DefConvertor.ConvertWebGetListToDTO(&req))
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(data, total))
}

// CreateDef 创建流程定义
// @Summary 创建流程定义
// @Description 创建流程定义
// @Tags 工作流定义相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.CreateWorkflowDefDTO true "创建流程定义请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/def/create [post]
func (h *WorkflowEngineController) CreateDef(c *gin.Context) {
	var req dto.CreateWorkflowDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Commands.CreateWorkflowDef(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// UpdateDef 更新流程定义
// @Summary 更新流程定义
// @Description 更新流程定义
// @Tags 工作流定义相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.CreateWorkflowDefDTO true "创建流程定义请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/def/update [post]
func (h *WorkflowEngineController) UpdateDef(c *gin.Context) {
	var req dto.CreateWorkflowDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.UpdateWorkflowDef(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// EnableDef 激活流程定义
// @Summary 激活流程定义
// @Description 激活流程定义
// @Tags 工作流定义相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.EnableWorkflowDefDTO true "激活流程定义请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/def/enable [post]
func (h *WorkflowEngineController) EnableDef(c *gin.Context) {
	var req dto.EnableWorkflowDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.EnableWorkflowDef(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// DisableDef 去激活流程定义
// @Summary 去激活流程定义
// @Description 去激活流程定义
// @Tags 工作流定义相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.DisableWorkflowDefDTO true "去激活流程定义请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/def/disable [post]
func (h *WorkflowEngineController) DisableDef(c *gin.Context) {
	var req dto.DisableWorkflowDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.DisableWorkflowDef(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// UploadDef 上传流程定义
// @Summary 上传流程定义
// @Description 上传流程定义
// @Tags 工作流定义相关接口
// @Accept multipart/form-data
// @Produce multipart/form-data
// @Param def_id formData int false "流程ID"
// @Param name formData string false "流程名称"
// @Param workflow_file formData file true "工作流文件-json格式"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/def/upload [post]
func (h *WorkflowEngineController) UploadDef(c *gin.Context) {
	var req dto.UploadWorkflowDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	var err error
	req.DefJson, err = utils.ReadFile(req.WorkflowFile)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	data, err := h.domainService.Commands.UploadWorkflowDef(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// ArchiveHistory 归档流程实例
// @Summary 归档流程实例
// @Description 归档流程实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.ArchiveHistoryWorkflowInstsDTO true "归档流程实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/archive [post]
func (h *WorkflowEngineController) ArchiveHistory(c *gin.Context) {
	var req dto.ArchiveHistoryWorkflowInstsDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.ArchiveHistoryWorkflowInsts(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// StartInst 创建流程实例
// @Summary 创建流程实例
// @Description 创建流程实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.StartWorkflowInstDTO true "启动流程实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/start [post]
func (h *WorkflowEngineController) StartInst(c *gin.Context) {
	var req dto.StartWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Commands.StartWorkflowInst(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// CancelInst 取消流程实例
// @Summary 取消流程实例
// @Description 取消流程实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.CancelWorkflowInstDTO true "取消流程实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/cancel [post]
func (h *WorkflowEngineController) CancelInst(c *gin.Context) {
	var req dto.CancelWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.CancelWorkflowInst(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// CompleteInst 标记流程实例结束
// @Summary 标记流程实例结束
// @Description 标记流程实例结束
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.CompleteWorkflowInstDTO true "标记流程实例结束请求"
// @Success 200 {object}  constants.WebRsp
// @Router /engine/api/v1/inst/complete [post]
func (h *WorkflowEngineController) CompleteInst(c *gin.Context) {
	var req dto.CompleteWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.CompleteWorkflowInst(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// PauseInst 暂停流程实例
// @Summary 暂停流程实例
// @Description 暂停流程实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.PauseWorkflowInstDTO true "暂停流程实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/pause [post]
func (h *WorkflowEngineController) PauseInst(c *gin.Context) {
	var req dto.PauseWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.PauseWorkflowInst(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// ResumeInst 恢复流程实例
// @Summary 恢复流程实例
// @Description 恢复流程实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.ResumeWorkflowInstDTO true "恢复流程实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/resume [post]
func (h *WorkflowEngineController) ResumeInst(c *gin.Context) {
	var req dto.ResumeWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.ResumeWorkflowInst(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// RestartInst 重启流程实例
// @Summary 重启流程实例
// @Description 重启流程实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.RestartWorkflowInstDTO true "重启流程实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/restart [post]
func (h *WorkflowEngineController) RestartInst(c *gin.Context) {
	var req dto.RestartWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.RestartWorkflowInst(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// Debug 调试流程
// @Summary 调试流程
// @Description 调试流程
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.DebugWorkflowInstDTO true "更新流程实例调试信息请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/debug [post]
func (h *WorkflowEngineController) Debug(c *gin.Context) {
	var req dto.DebugWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.DebugWorkflowInst(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// UpdateCtx 更新实例上下文
// @Summary 更新实例上下文
// @Description 更新实例上下文
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst body dto.UpdateWorkflowInstCtxDTO true "重启流程实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/updatectx [post]
func (h *WorkflowEngineController) UpdateCtx(c *gin.Context) {
	var req dto.UpdateWorkflowInstCtxDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.UpdateWorkflowInstCtx(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// GetInstDetail 查询单条工作流实例
// @Summary 查询单条工作流实例
// @Description 查询单条工作流实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst query dto.GetWorkflowInstDTO true "查询实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/get [get]
func (h *WorkflowEngineController) GetInstDetail(c *gin.Context) {
	var req dto.GetWorkflowInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	// 如果是超管用户则可以查询任意流程的详情
	if utils.StrContains(config.GetRbacConfig().SuperAdmins, req.Operator) {
		req.Operator = ""
	}

	data, err := h.domainService.Queries.GetWorkflowInst(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// GetInstList 查询多条工作流实例
// @Summary 查询多条工作流实例
// @Description 查询多条工作流实例
// @Tags 工作流实例相关接口
// @Accept application/json
// @Produce application/json
// @Param inst query dto.GetWebWorkflowInstListDTO true "查询多条实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/inst/list [get]
func (h *WorkflowEngineController) GetInstList(c *gin.Context) {
	var req dto.GetWebWorkflowInstListDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, total, err := h.domainService.Queries.GetWorkflowInstList(c.Request.Context(),
		convertor.InstConvertor.ConvertWebGetListToDTO(&req))
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(data, total))
}

// GetNodeInstDetail 查询节点实例信息
// @Summary 查询节点实例信息
// @Description 查询节点实例信息
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst query dto.GetNodeInstDTO true "查询节点实例请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/get [get]
func (h *WorkflowEngineController) GetNodeInstDetail(c *gin.Context) {
	var req dto.GetNodeInstDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Queries.GetNodeInstDetail(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// RerunNode 重跑节点
// @Summary 重跑节点
// @Description 重跑节点
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.RerunNodeDTO true "重跑节点请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/rerun [post]
func (h *WorkflowEngineController) RerunNode(c *gin.Context) {
	var req dto.RerunNodeDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.RerunNode(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// ResumeNode 恢复节点执行
// @Summary 恢复节点执行
// @Description 恢复节点执行
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.ResumeNodeDTO true "恢复节点执行请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/resume [post]
func (h *WorkflowEngineController) ResumeNode(c *gin.Context) {
	var req dto.ResumeNodeDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	err := h.domainService.Commands.ResumeNode(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// CancelNode 取消节点的执行
// @Summary 取消节点执行
// @Description 取消节点执行
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.CancelNodeDTO true "取消节点执行请求"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/cancel [post]
func (h *WorkflowEngineController) CancelNode(c *gin.Context) {
	var req dto.CancelNodeDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.CancelNode(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// SkipNode 跳过节点执行
// @Summary 跳过节点执行
// @Description 跳过节点执行
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.SkipNodeDTO true "跳过节点执行"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/skip [post]
func (h *WorkflowEngineController) SkipNode(c *gin.Context) {
	var req dto.SkipNodeDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.SkipNode(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// CancelSkipNode 取消跳过节点执行
// @Summary 取消跳过节点执行
// @Description 取消跳过节点执行
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.CancelSkipNodeDTO true "取消跳过节点执行"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/cancelskip [post]
func (h *WorkflowEngineController) CancelSkipNode(c *gin.Context) {
	var req dto.CancelSkipNodeDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.CancelSkipNode(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// CompleteNode 标记节点完成
// @Summary 标记节点完成
// @Description 标记节点完成
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.CompleteNodeDTO true "标记节点完成"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/complete [post]
func (h *WorkflowEngineController) CompleteNode(c *gin.Context) {
	var req dto.CompleteNodeDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.CompleteNode(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// SendCronPresetEvent 发送即时处理定时触发事件
// @Summary 发送定时触发事件
// @Description 发送定时触发事件
// @Tags 事件相关接口
// @Accept application/json
// @Produce application/json
// @Param event body dto.SendCronPresetEventDTO true "发送定时触发事件"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/event/sendcronpresetevent [post]
func (h *WorkflowEngineController) SendCronPresetEvent(c *gin.Context) {
	var req dto.SendCronPresetEventDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.eventBusRepo.SendCronPresetEvent(c.Request.Context(), time.Time{}, req.Value)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// SendDriveEvent 发送驱动事件
// @Summary 发送驱动事件
// @Description 发送驱动事件
// @Tags 事件相关接口
// @Accept application/json
// @Produce application/json
// @Param event body dto.SendDriveEventDTO true "发送驱动事件"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/event/senddriveevent [post]
func (h *WorkflowEngineController) SendDriveEvent(c *gin.Context) {
	var req dto.SendDriveEventDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.eventBusRepo.SendDriveEvent(c.Request.Context(), req.Value)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// SendExternalEvent 发送外部事件
// @Summary 发送外部事件
// @Description 发送外部事件
// @Tags 事件相关接口
// @Accept application/json
// @Produce application/json
// @Param event body dto.SendExternalEventDTO true "发送外部事件"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/event/sendexternalevent [post]
func (h *WorkflowEngineController) SendExternalEvent(c *gin.Context) {
	var req dto.SendExternalEventDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.eventBusRepo.SendExternalEvent(c.Request.Context(), req.Value)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// SendTriggerEvent 发送触发器事件
// @Summary 发送触发器事件
// @Description 发送触发器事件
// @Tags 事件相关接口
// @Accept application/json
// @Produce application/json
// @Param event body dto.SendTriggerEventDTO true "发送触发器事件"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/event/sendtriggerevent [post]
func (h *WorkflowEngineController) SendTriggerEvent(c *gin.Context) {
	var req dto.SendTriggerEventDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.eventBusRepo.SendTriggerEvent(c.Request.Context(), req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// SetNodeTimeout 标记节点超时
// @Summary 标记节点超时
// @Description 标记节点超时
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.SetNodeTimeoutDTO true "标记节点超时"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/timeout [post]
func (h *WorkflowEngineController) SetNodeTimeout(c *gin.Context) {
	var req dto.SetNodeTimeoutDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.SetTimeout(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// SetNodeNearTimeout 标记节点接近超时
// @Summary 标记节点接近超时
// @Description 标记节点接近超时
// @Tags 节点相关接口
// @Accept application/json
// @Produce application/json
// @Param nodeInst body dto.SetNodeNearTimeoutDTO true "标记节点接近超时"
// @Success 200 {object} constants.WebRsp
// @Router /engine/api/v1/node/neartimeout [post]
func (h *WorkflowEngineController) SetNodeNearTimeout(c *gin.Context) {
	var req dto.SetNodeNearTimeoutDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.SetNearTimeout(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// CallAuth 调用函数的权限校验
func (h *WorkflowEngineController) CallAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.GetHeader("Namespace")
		accessToken, _ := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")

		log.Infof("Call workflow auth, namespace=%s, accessToken=%s", namespace, accessToken)
		if err := h.permissionValidator.ValidateToken(c.Request.Context(), &remote.ValidateTokenReqDTO{
			Namespace:   namespace,
			AccessToken: accessToken,
		}); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
		return
	}
}
