package web

import (
	"github.com/fflow-tech/fflow/service/pkg/login"
	"gorm.io/gorm/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/service"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/config"
	ac "github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/errno"
)

// AuthController foundation/auth http 服务实现
type AuthController struct {
	domainService *service.DomainService
}

// NewAuthController 构造函数
func NewAuthController(domainService *service.DomainService) *AuthController {
	return &AuthController{domainService: domainService}
}

// Login 基础用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.LoginReqDTO true "登录请求"
// @Success 200 {object} dto.LoginRspDTO
// @Router /auth/api/v1/login/account [post]
func (h *AuthController) Login(c *gin.Context) {
	var req dto.LoginReqDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.NewFailedLoginRsp(err.Error()))
		return
	}

	if len(req.EmailReceiver) > 0 && len(req.Captcha) > 0 {
		h.verifyCaptcha(c, req, dto.VerifyCaptchaReqDTO{
			EmailReceiver: req.EmailReceiver,
			Captcha:       req.Captcha,
		})
		return
	}

	data, err := h.domainService.Commands.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewFailedLoginRsp(err.Error()))
		return
	}

	if err := login.SetUserInfoToCookie(c, &dto.CurrentUserData{
		Namespace: h.getDefaultNamespace(c, config.GetAppConfig().AdminUsername),
		Username:  config.GetAppConfig().AdminUsername,
		Email:     config.GetAppConfig().AdminEmail,
		Avatar:    ac.DefaultAvatar,
	}, config.GetAppConfig().SecretKey); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewFailedLoginRsp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AuthController) getDefaultNamespace(c *gin.Context, username string) string {
	namespaces, err := h.domainService.Commands.GetDomainsForUser(c.Request.Context(), &dto.RbacReqDTO{User: username})
	if err != nil {
		return ""
	}

	if len(namespaces) > 0 {
		return namespaces[0]
	}

	return ac.DefaultNamespace
}

// CurrentUser 获取当前用户
// @Summary 获取当前用户
// @Description 获取当前用户
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.CurrentUserReqDTO true "当前用户请求"
// @Success 200 {object} dto.CurrentUserRspDTO
// @Router /auth/api/v1/currentUser [get]
func (h *AuthController) CurrentUser(c *gin.Context) {
	user, err := login.GetUserInfoFromCookie(c, config.GetAppConfig().SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	// 允许所有域名跨域访问
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")

	c.JSON(http.StatusOK, &dto.CurrentUserRspDTO{
		Success: true,
		Data: &dto.CurrentUserData{
			Namespace: user.Namespace,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
		},
	})
}

// OutLogin 用户取消登录
// @Summary 用户取消登录
// @Description 用户取消登录
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.OutLoginReqDTO true "取消登录请求"
// @Success 200 {object} dto.OutLoginRspDTO
// @Router /auth/api/v1/login/outLogin [post]
func (h *AuthController) OutLogin(c *gin.Context) {
	var req dto.OutLoginReqDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	data, err := h.domainService.Commands.OutLogin(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}
	c.SetCookie(login.SessionCookieName, "", -1, "/", config.GetAppConfig().Domain, false, false)
	c.JSON(http.StatusOK, data)
}

// Oauth2Callback Oauth2 用户登录回调
// @Summary 用户登录回调
// @Description 用户登录回调
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.Oauth2CallbackReqDTO true "登录请求"
// @Success 200 {object} dto.Oauth2CallbackRspDTO
// @Router /auth/api/v1/oauth2/callback [get]
func (h *AuthController) Oauth2Callback(c *gin.Context) {
	var req dto.Oauth2CallbackReqDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	req.Code = c.Query("code")
	data, err := h.domainService.Commands.Oauth2Callback(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	if data == nil || len(data.Username) <= 0 {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, "illegal user"))
		return
	}

	if err := login.SetUserInfoToCookie(c, &dto.CurrentUserData{
		Namespace: h.getDefaultNamespace(c, data.Username),
		Username:  data.Username,
		Email:     data.Email,
		Avatar:    data.Avatar,
	}, config.GetAppConfig().SecretKey); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, config.GetAppConfig().HomePage)
}

// GetCaptcha 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.GetCaptchaReqDTO true "发送邮箱验证码请求"
// @Success 200 {object} dto.GetCaptchaRspDTO
// @Router /auth/api/v1/login/captcha [post]
func (h *AuthController) GetCaptcha(c *gin.Context) {
	var req dto.GetCaptchaReqDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	data, err := h.domainService.Commands.GetCaptcha(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(data))
}

func (h *AuthController) verifyCaptcha(c *gin.Context, loginReq dto.LoginReqDTO, verifyCaptchaReq dto.VerifyCaptchaReqDTO) {
	data, err := h.domainService.Commands.VerifyCaptcha(c.Request.Context(), &verifyCaptchaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewFailedLoginRsp(err.Error()))
		return
	}

	if !data.IsValidCaptcha {
		c.JSON(http.StatusOK, dto.NewFailedLoginRsp("illegal verification code"))
		return
	}

	if err := login.SetUserInfoToCookie(c, &dto.CurrentUserData{
		Namespace: h.getDefaultNamespace(c, data.Username),
		Username:  data.Username,
		Email:     data.Email,
		Avatar:    data.Avatar,
	}, config.GetAppConfig().SecretKey); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewFailedLoginRsp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, &dto.LoginRspDTO{
		Status:           "ok",
		Type:             loginReq.Type,
		CurrentAuthority: "admin",
	})
}

// CreateNamespace 添加命名空间
// @Summary 添加命名空间
// @Description 添加命名空间
// @Tags 命名空间相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.CreateNamespaceDTO true "添加命名空间请求"
// @Success 200 {object} constants.WebRsp
// @Router /auth/api/v1/namespaces [post]
func (h *AuthController) CreateNamespace(c *gin.Context) {
	var req dto.CreateNamespaceDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	user, err := login.GetUserInfoFromCookie(c, config.GetAppConfig().SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	req.Creator = user.Username
	err = h.domainService.Commands.CreateNamespace(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

// GetNamespaces 获取所有的命名空间
// @Summary 获取所有的命名空间
// @Description 获取所有的命名空间
// @Tags 命名空间相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.GetNamespacesReqDTO true "查询命名空间请求"
// @Success 200 {object} constants.WebRsp
// @Router /auth/api/v1/namespaces [get]
func (h *AuthController) GetNamespaces(c *gin.Context) {
	var req dto.GetNamespacesReqDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	user, err := login.GetUserInfoFromCookie(c, config.GetAppConfig().SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	req.Username = user.Username
	data, total, err := h.domainService.Commands.GetNamespaces(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(data, total))
}

// RegisterNamespaceAPIToken 获取命名空间的 Token
// @Summary 获取命名空间的 Token
// @Description 获取命名空间的 Token
// @Tags 命名空间相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.RegisterNamespaceAPITokenReqDTO true "获取命名空间的 Token 请求"
// @Success 200 {object} constants.WebRsp
// @Router /auth/api/v1/namespaces/{namespace}/tokens [get]
func (h *AuthController) RegisterNamespaceAPIToken(c *gin.Context) {
	var req dto.RegisterNamespaceAPITokenReqDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	user, err := login.GetUserInfoFromCookie(c, config.GetAppConfig().SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	req.Namespace = c.Param("namespace")
	req.Creator = user.Username
	token, err := h.domainService.Commands.RegisterNamespaceAPIToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(token))
}

// GetNamespaceTokens 获取命名空间的 Token
// @Summary 获取命名空间的 Token 请求
// @Description 获取命名空间的 Token
// @Tags 命名空间相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.PageQueryNamespaceTokenDTO true "获取命名空间的 Token 请求"
// @Success 200 {object} constants.WebRsp
// @Router /auth/api/v1/namespaces/{namespace}/tokens [get]
func (h *AuthController) GetNamespaceTokens(c *gin.Context) {
	var req dto.PageQueryNamespaceTokenDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	// 从 Cookie 中获取用户信息，避免用户信息被纂改
	user, err := login.GetUserInfoFromCookie(c, config.GetAppConfig().SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	tokens, total, err := h.domainService.Commands.GetNamespaceTokens(c.Request.Context(),
		&dto.PageQueryNamespaceTokenDTO{
			Namespace: c.Param("namespace"),
			Creator:   user.Username,
			PageQuery: req.PageQuery,
			Order:     req.Order,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(tokens, total))
}

// GetUserNamespaces 获取用户命名空间
// @Summary 获取用户命名空间
// @Description 获取用户命名空间
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.GetUserNamespacesReqDTO true "获取用户命名空间请求"
// @Success 200 {object} constants.WebRsp
// @Router /auth/api/v1/user/namespace [get]
func (h *AuthController) GetUserNamespaces(c *gin.Context) {
	// 从 Cookie 中获取用户信息，避免用户信息被纂改
	user, err := login.GetUserInfoFromCookie(c, config.GetAppConfig().SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	namespaces, err := h.domainService.Commands.GetDomainsForUser(c.Request.Context(), &dto.RbacReqDTO{User: user.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(&dto.GetUserNamespacesRspDTO{
		CurrentNamespace: user.Namespace,
		Namespaces:       namespaces,
	}))
}

// SetUserCurrentNamespace 设置用户当前的命名空间
// @Summary 设置用户当前的命名空间
// @Description 设置用户当前的命名空间
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.SetUserCurrentNamespaceReqDTO true "设置用户当前的命名空间"
// @Success 200 {object} constants.WebRsp
// @Router /auth/api/v1/user/namespace [post]
func (h *AuthController) SetUserCurrentNamespace(c *gin.Context) {
	var req dto.SetUserCurrentNamespaceReqDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	user, err := login.GetUserInfoFromCookie(c, config.GetAppConfig().SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}
	namespaces, err := h.domainService.Commands.GetDomainsForUser(c.Request.Context(), &dto.RbacReqDTO{User: user.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}
	if !utils.Contains(namespaces, req.CurrentNamespace) {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, "user is not have permission for namespace: "+req.CurrentNamespace))
		return
	}

	user.Namespace = req.CurrentNamespace
	if err := login.SetUserInfoToCookie(c, user, config.GetAppConfig().SecretKey); err != nil {
		c.JSON(http.StatusInternalServerError, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
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

	if c.Request.Method == "GET" {
		login.SetNamespace(req, currentUser)
	} else {
		login.SetNamespace(req, currentUser)
		login.SetOperator(req, currentUser)
	}

	return nil
}
