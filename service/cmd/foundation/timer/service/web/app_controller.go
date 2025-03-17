package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/errno"
)

// AppController  http服务实现
type AppController struct {
	domainService *service.DomainService
}

// NewAppController 构造函数
func NewAppController(domainService *service.DomainService) *AppController {
	return &AppController{domainService: domainService}
}

// GetAppList 查询 App 列表
// @Summary  查询 App 列表
// @Description  查询 App 列表
// @Tags 应用相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.PageQueryAppDTO true "查询 App 列表"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/app/list [GET]
func (h *AppController) GetAppList(c *gin.Context) {
	var req dto.PageQueryAppDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	appList, total, err := h.domainService.Queries.GetAppList(&req)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(appList, total))
}

// CreateApp 创建应用
// @Summary 创建应用
// @Description 创建应用
// @Tags 应用相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.CreateAppDTO true "创建定时器定义"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/app/create [post]
func (h *AppController) CreateApp(c *gin.Context) {
	var req dto.CreateAppDTO
	// 绑定参数
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.CreateApp(&req)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSucceedWebRsp(nil))
}

// DeleteApp 删除app
// @Summary 删除app
// @Description 删除app
// @Tags 应用相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.DeleteAppDTO true "删除APP"
// @Success 200 {object} constants.WebRsp
// @Router /timer/api/v1/app/deleteApp [delete]
func (h *AppController) DeleteApp(c *gin.Context) {
	var req dto.DeleteAppDTO
	// 绑定参数
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.DeleteApp(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp("delete app success"))
}
