package web

import (
	"net/http"
	"strings"

	"github.com/fflow-tech/fflow/service/pkg/login"
	"github.com/fflow-tech/fflow/service/pkg/remote"

	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

const (
	anonymousOperator = "anonymous"
)

// FAASController foundation/faas http服务实现
type FAASController struct {
	domainService       *service.DomainService
	permissionValidator *remote.DefaultPermissionValidator
}

// NewFAASController 构造函数
func NewFAASController(domainService *service.DomainService,
	permissionValidator *remote.DefaultPermissionValidator) *FAASController {
	return &FAASController{domainService: domainService, permissionValidator: permissionValidator}
}

// CallFunction 执行函数
// @Summary 执行函数
// @Description 执行函数
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.CallFunctionReqDTO true "执行函数请求"
// @Success 200 {object} interface{}
// @Router /faas/api/v1/func/call [post]
func (h *FAASController) CallFunction(c *gin.Context) {
	var req dto.CallFunctionReqDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	req.Request = c.Request
	data, err := h.domainService.Commands.CallFunction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}
	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// CallFunctionForHttpPost 执行函数
// @Summary 执行函数
// @Description 执行函数
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param namespace path string true "命名空间"
// @Param function path string true "函数名称"
// @Param input body map[string]interface{} true "调用函数请求体"
// @Success 200 {object} interface{}
// @Router /faas/openapi/v1/func/call/{namespace}/{function} [post]
func (h *FAASController) CallFunctionForHttpPost(c *gin.Context) {
	var req dto.CallFunctionReqDTO
	if err := bindPath(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	var funcInput map[string]interface{}
	// 当 Content-type 为 application/xml 时，因 xml 的 decode 只能转到明确的类型，不能转为 map[string]interface{},
	// 这里需要跳过 bind，后续将 request body 直接传入函数的上下文，让使用方自行解析
	// 请求的 body 应该是 json 结构，接受 Content-type 为 text/json 或 application/json
	// 直接用 Bind 的话如果 Content-type 不为标准的 application/json 时，gin 似乎会尝试用 form 去解析，从而导致报错
	// 有些 webhook 回调设置的 Content-type 为 text/json，因此这里做一下兼容处理
	if c.Request.Header.Get("Content-type") != "application/xml" {
		if err := c.BindJSON(&funcInput); err != nil {
			log.Infof("The input is not valid json: %s", err.Error())
			if errBind := c.Bind(&funcInput); err != nil {
				c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, errBind.Error()))
				return
			}
			c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
			return
		}
	}

	var err error
	req.Input, err = utils.MergeMap(getQueryMap(c), funcInput)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	req.Operator = anonymousOperator
	req.Request = c.Request
	data, err := h.domainService.Commands.CallFunction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}
	// 对于返回值为 string 的请求，直接通过 String 方法返回（这里主要处理企微机器人回调无法解析返回值的问题）
	switch data.(type) {
	case string:
		c.String(http.StatusOK, "%s", data)
	default:
		c.JSON(http.StatusOK, data)
	}
}

// CallAuth 调用函数的权限校验
func (h *FAASController) CallAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.GetHeader("Namespace")
		accessToken, _ := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")

		log.Infof("Call function auth, namespace=%s, accessToken=%s", namespace, accessToken)
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

// CallFunctionForHttpGet 通过执行函数
// @Summary 通过执行函数
// @Description 执行函数
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param namespace path string true "命名空间"
// @Param function path string true "函数名称"
// @Param params query string true "调用函数请求体"
// @Success 200 {object} interface{}
// @Router /faas/openapi/v1/func/call/{namespace}/{function} [get]
func (h *FAASController) CallFunctionForHttpGet(c *gin.Context) {
	var req dto.CallFunctionReqDTO
	if err := bindPath(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	req.Input = getQueryMap(c)
	log.Infof("Call function req:%s", utils.StructToJsonStr(req))
	req.Operator = anonymousOperator
	req.Request = c.Request
	data, err := h.domainService.Commands.CallFunction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	switch data.(type) {
	case string:
		c.String(http.StatusOK, "%s", data)
	default:
		c.JSON(http.StatusOK, data)
	}
}

func getQueryMap(c *gin.Context) map[string]interface{} {
	funcInput := map[string]interface{}{}
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			funcInput[k] = v[0]
			continue
		}
		funcInput[k] = v
	}
	return funcInput
}

func bindPath(c *gin.Context, req *dto.CallFunctionReqDTO) error {
	req.Function = c.Param("function")
	req.Namespace = c.Param("namespace")
	return nil
}

// DebugFunction  调试函数
// @Summary  调试函数
// @Description  调试函数
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param debugReq body dto.DebugFunctionDTO true " 调试函数请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func/debug [post]
func (h *FAASController) DebugFunction(c *gin.Context) {
	var req dto.DebugFunctionDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Commands.DebugFunction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}
	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// CreateFunction 创建函数
// @Summary 创建函数
// @Description 创建函数
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param createReq body dto.CreateFunctionReqDTO true "创建函数请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func [post]
func (h *FAASController) CreateFunction(c *gin.Context) {
	var req dto.CreateFunctionReqDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Commands.CreateFunction(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// GetFunction 查询函数
// @Summary 查询函数详情
// @Description 查询函数详情
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param function query dto.GetFunctionReqDTO true "查询函数请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func [get]
func (h *FAASController) GetFunction(c *gin.Context) {
	var req dto.GetFunctionReqDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Queries.GetFunction(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// GetFunctions 查询函数列表
// @Summary 查询函数列表
// @Description 查询函数列表
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param function query dto.PageQueryFunctionDTO true "查询函数列表请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func/list [get]
func (h *FAASController) GetFunctions(c *gin.Context) {
	var req dto.PageQueryFunctionDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, total, err := h.domainService.Queries.GetFunctions(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(data, total))
}

// UpdateFunction 更新函数
// @Summary 更新函数
// @Description 更新函数
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param updateReq body dto.UpdateFunctionDTO true "更新函数请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func [put]
func (h *FAASController) UpdateFunction(c *gin.Context) {
	var req dto.UpdateFunctionDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Commands.UpdateFunction(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

// DeleteFunction 删除函数
// @Summary 删除函数
// @Description 删除函数
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param deleteReq body dto.DeleteFunctionDTO true "删除函数请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func [delete]
func (h *FAASController) DeleteFunction(c *gin.Context) {
	var req dto.DeleteFunctionDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.DeleteFunction(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp("delete success"))
}

// GetRunHistories 查询函数执行列表
// @Summary 查询函数执行列表
// @Description 查询函数执行列表
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param runHistory query dto.PageQueryRunHistoryDTO true "查询函数执行列表请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func/history/list [get]
func (h *FAASController) GetRunHistories(c *gin.Context) {
	var req dto.PageQueryRunHistoryDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, total, err := h.domainService.Queries.GetRunHistories(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(data, total))
}

// DeleteRunHistories 删除历史执行记录
// @Summary 删除历史执行记录
// @Description 删除历史执行记录
// @Tags 函数相关接口
// @Accept application/json
// @Produce application/json
// @Param deleteReq body dto.BatchDeleteExpiredRunHistoryDTO true "删除历史执行记录请求"
// @Success 200 {object} constants.WebRsp
// @Router /faas/api/v1/func/histories [delete]
func (h *FAASController) DeleteRunHistories(c *gin.Context) {
	var req dto.BatchDeleteExpiredRunHistoryDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if req.KeepDays == 0 {
		req.KeepDays = config.GetAppConfig().KeepDays
	}

	err := h.domainService.Commands.BatchDeleteExpiredRunHistory(&req)
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp("batch delete success"))
}

func bindReq(c *gin.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	secretKey := config.GetAppConfig().SecretKey
	currentUser, err := login.GetUserInfoFromCookie(c, secretKey)
	if err != nil {
		return err
	}

	log.Infof("Current user: %s", utils.StructToJsonStr(currentUser))
	if c.Request.Method == "GET" {
		login.SetNamespace(req, currentUser)
	} else {
		login.SetNamespace(req, currentUser)
		login.SetOperator(req, currentUser)
	}

	return nil
}
