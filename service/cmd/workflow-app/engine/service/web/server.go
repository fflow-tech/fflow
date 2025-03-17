// Package web 提供 http 调用的入口
package web

import (
	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/cmd/workflow-app/engine/factory"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/ginfilter"
	"github.com/fflow-tech/fflow/service/pkg/log"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Server Web容器
type Server struct {
	addrs         []string
	rootRouter    *gin.Engine
	basicRouter   *gin.RouterGroup
	engineRouter  *gin.RouterGroup
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
	s.basicRouter = s.rootRouter.Group("engine")
	s.engineRouter = s.basicRouter.Group("api/v1")
	s.openAPIRouter = s.basicRouter.Group("openapi/v1")

	domainService, err := factory.GetDomainService()
	if err != nil {
		panic(err)
	}
	eventBusRepo, err := factory.GetEventBusRepo()
	if err != nil {
		panic(err)
	}
	validator, err := factory.GetDefaultPermissionValidator()
	if err != nil {
		panic(err)
	}
	s.RegisterBasicController()
	s.RegisterEngineController(NewWorkflowEngineController(domainService, eventBusRepo, validator))
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

// RegisterEngineController 注册处理器
func (s *Server) RegisterEngineController(controller *WorkflowEngineController) {
	s.engineRouter.Use(ginfilter.Cors(config.GetCorsConfig())).Use(ginfilter.Auth(config.GetAppConfig().AuthConfig))
	defRouter := s.engineRouter.Group("/def/").Use()
	{
		defRouter.GET("get", controller.GetDefDetail)
		defRouter.GET("list", controller.GetDefList)
		defRouter.POST("create", controller.CreateDef)
		defRouter.POST("update", controller.UpdateDef)
		defRouter.POST("enable", controller.EnableDef)
		defRouter.POST("disable", controller.DisableDef)
		defRouter.POST("upload", controller.UploadDef)
	}
	instRouter := s.engineRouter.Group("/inst/").Use()
	{
		instRouter.GET("get", controller.GetInstDetail)
		instRouter.GET("list", controller.GetInstList)
		instRouter.POST("start", controller.StartInst)
		instRouter.POST("restart", controller.RestartInst)
		instRouter.POST("cancel", controller.CancelInst)
		instRouter.POST("pause", controller.PauseInst)
		instRouter.POST("resume", controller.ResumeInst)
		instRouter.POST("complete", controller.CompleteInst)
		instRouter.POST("updatectx", controller.UpdateCtx)
		instRouter.POST("debug", controller.Debug)
		instRouter.POST("archive", controller.ArchiveHistory)
	}
	nodeInstRouter := s.engineRouter.Group("/node/").Use()
	{
		nodeInstRouter.GET("get", controller.GetNodeInstDetail)
		nodeInstRouter.POST("rerun", controller.RerunNode)
		nodeInstRouter.POST("resume", controller.ResumeNode)
		nodeInstRouter.POST("cancel", controller.CancelNode)
		nodeInstRouter.POST("complete", controller.CompleteNode)
		nodeInstRouter.POST("skip", controller.SkipNode)
		nodeInstRouter.POST("cancelskip", controller.CancelSkipNode)
		nodeInstRouter.POST("timeout", controller.SetNodeTimeout)
		nodeInstRouter.POST("neartimeout", controller.SetNodeNearTimeout)
	}
	eventRouter := s.engineRouter.Group("/event/").Use()
	{
		eventRouter.POST("senddriveevent", controller.SendDriveEvent)
		eventRouter.POST("sendexternalevent", controller.SendExternalEvent)
		eventRouter.POST("sendtriggerevent", controller.SendTriggerEvent)
		eventRouter.POST("sendcronpresetevent", controller.SendCronPresetEvent)
	}

	s.openAPIRouter.Use(controller.CallAuth())
	defOpenAPIRouter := s.openAPIRouter.Group("/def/").Use()
	{
		defOpenAPIRouter.GET("list", controller.GetDefList)
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
