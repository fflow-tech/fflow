// Package web 提供 http 调用的入口
package web

import (
	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/cmd/foundation/faas/factory"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/ginfilter"
	"github.com/fflow-tech/fflow/service/pkg/log"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Server Web 容器
type Server struct {
	addrs         []string
	rootRouter    *gin.Engine
	basicRouter   *gin.RouterGroup
	apiRouter     *gin.RouterGroup
	openAPIRouter *gin.RouterGroup
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
	s.basicRouter = s.rootRouter.Group("faas")
	// api 为给前端使用的接口, 需要登录
	s.apiRouter = s.basicRouter.Group("api/v1")
	// openapi 为给第三方使用的接口, 不需要登录
	s.openAPIRouter = s.basicRouter.Group("openapi/v1")

	domainService, err := factory.GetDomainService()

	if err != nil {
		panic(err)
	}

	validator, err := factory.GetDefaultPermissionValidator()
	if err != nil {
		panic(err)
	}

	s.RegisterBasicController()
	s.RegisterFaasController(NewFAASController(domainService, validator))
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

// RegisterFaasController 注册处理器
func (s *Server) RegisterFaasController(controller *FAASController) {
	s.apiRouter.Use(ginfilter.Cors(config.GetCorsConfig())).Use(ginfilter.Auth(config.GetAppConfig().AuthConfig))
	funcAPIRouter := s.apiRouter.Group("/func/")
	{
		funcAPIRouter.POST("call", controller.CallFunction)
		funcAPIRouter.POST("", controller.CreateFunction)
		funcAPIRouter.GET("", controller.GetFunction)
		funcAPIRouter.GET("list", controller.GetFunctions)
		funcAPIRouter.PUT("", controller.UpdateFunction)
		funcAPIRouter.DELETE("", controller.DeleteFunction)
		funcAPIRouter.GET("history/list", controller.GetRunHistories)
		funcAPIRouter.POST("debug", controller.DebugFunction)
		funcAPIRouter.DELETE("histories", controller.DeleteRunHistories)
	}

	s.openAPIRouter.Use(controller.CallAuth())
	funcOpenAPIRouter := s.openAPIRouter.Group("/func/")
	{
		funcOpenAPIRouter.POST("call/:namespace/:function", controller.CallFunctionForHttpPost)
		funcOpenAPIRouter.GET("call/:namespace/:function", controller.CallFunctionForHttpGet)
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
