// Package web 提供 http 调用的入口
package web

import (
	"github.com/fflow-tech/fflow/service/cmd/foundation/auth/factory"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/ginfilter"
	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Server Web 容器
type Server struct {
	addrs       []string
	rootRouter  *gin.Engine
	basicRouter *gin.RouterGroup
	apiRouter   *gin.RouterGroup
}

// NewServer 新建 Web 容器
func NewServer(addr ...string) *Server {
	s := &Server{}
	s.addrs = addr
	s.rootRouter = gin.Default()
	s.rootRouter.MaxMultipartMemory = maxMultipartMemory
	s.rootRouter.Use(
		gin.Recovery(),
	)
	s.basicRouter = s.rootRouter.Group("auth")
	// api 为给前端使用的接口, 需要登录
	s.apiRouter = s.basicRouter.Group("api/v1")

	domainService, err := factory.GetDomainService()
	if err != nil {
		panic(err)
	}
	s.RegisterBasicController()
	s.RegisterAuthController(NewAuthController(domainService))
	return s
}

// 限制单次传输大小
const maxMultipartMemory = 8 << 20 // 8 MiB

// RegisterBasicController 注册基础处理器
func (s *Server) RegisterBasicController() {
	// ping handler
	s.basicRouter.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 注册gin swagger服务
	s.basicRouter.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
}

// RegisterAuthController 注册处理器
func (s *Server) RegisterAuthController(controller *AuthController) {
	loginAPIRouter := s.apiRouter.Group("/login/")
	{
		loginAPIRouter.POST("account", controller.Login)
		loginAPIRouter.POST("outLogin", controller.OutLogin)
		loginAPIRouter.POST("captcha", controller.GetCaptcha)
	}

	oauth2APIRouter := s.apiRouter.Group("/oauth2/")
	{
		oauth2APIRouter.GET("callback", controller.Oauth2Callback)
	}

	userNamespaceAPIRouter := s.apiRouter.Group("/user/").
		Use(ginfilter.Cors(config.GetCorsConfig())).
		Use(ginfilter.Auth(config.GetAppConfig().AuthConfig))
	{
		userNamespaceAPIRouter.GET("namespace", controller.GetUserNamespaces)
		userNamespaceAPIRouter.POST("namespace", controller.SetUserCurrentNamespace)
	}

	namespaceAPIRouter := s.apiRouter.Group("/namespaces/").
		Use(ginfilter.Cors(config.GetCorsConfig())).
		Use(ginfilter.Auth(config.GetAppConfig().AuthConfig))
	{
		namespaceAPIRouter.GET("", controller.GetNamespaces)
		namespaceAPIRouter.POST("", controller.CreateNamespace)
		namespaceAPIRouter.POST(":namespace/tokens", controller.RegisterNamespaceAPIToken)
		namespaceAPIRouter.GET(":namespace/tokens", controller.GetNamespaceTokens)
	}

	userAPIRouter := s.apiRouter.Group("/").
		Use(ginfilter.Cors(config.GetCorsConfig())).
		Use(ginfilter.Auth(config.GetAppConfig().AuthConfig))
	{
		userAPIRouter.GET("currentUser", controller.CurrentUser)
	}
}

// Serve 启动监听
func (s *Server) Serve() error {
	log.Infof("Start to bind http")
	return s.rootRouter.Run(s.addrs...)
}

// Close 关闭 HTTP 服务
func (s *Server) Close(ch chan struct{}) error {
	log.Infof("Shutdown HTTP Server...")
	defer log.Infof("HTTP Server exit")
	defer func() {
		if ch != nil {
			ch <- struct{}{}
		}
	}()
	return nil
}
