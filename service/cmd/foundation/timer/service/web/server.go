// Package web 负责应用 HTTP 通讯能力。
package web

import (
	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/factory"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/ginfilter"
	"github.com/fflow-tech/fflow/service/pkg/log"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Server Web容器
type Server struct {
	addrs       []string
	rootRouter  *gin.Engine
	basicRouter *gin.RouterGroup
	timerRouter *gin.RouterGroup
}

// NewServer 新建Web容器
func NewServer(addrs ...string) *Server {
	s := &Server{}
	s.addrs = addrs
	s.rootRouter = gin.Default()
	s.rootRouter.MaxMultipartMemory = maxMultipartMemory
	s.rootRouter.Use(
		gin.Recovery(),
		ginfilter.Cors(config.GetCorsConfig()),
	)
	s.basicRouter = s.rootRouter.Group("timer")
	s.timerRouter = s.basicRouter.Group("api/v1")

	domainService, err := factory.GetDomainService()
	if err != nil {
		panic(err)
	}
	s.RegisterBasicController()
	s.RegisterAppController(NewAppController(domainService))
	s.RegisterTimerController(NewTimerController(domainService))
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

// RegisterTimerController 注册处理器
func (s *Server) RegisterTimerController(controller *TimerController) {
	s.timerRouter.Use(ginfilter.Auth(config.GetAppConfig().AuthConfig))
	defRouter := s.timerRouter.Group("/def/").Use()
	{
		defRouter.GET("get", controller.GetDefDetail)
		defRouter.POST("create", controller.CreateTimerDef)
		defRouter.POST("change", controller.ChangeDefStatus)
		defRouter.GET("list", controller.GetTimerDefList)
		defRouter.DELETE("delete", controller.DeleteTimer)
		defRouter.GET("runHistory", controller.GetTimerRunHistory)
		defRouter.GET("timerTaskList", controller.GetTimerTaskList)
		defRouter.DELETE("deleteRunHistories", controller.DeleteRunHistories)
		defRouter.POST("timerListSend", controller.TimerListSendNotify)
	}
}

// RegisterAppController 注册 app 处理器
func (s *Server) RegisterAppController(controller *AppController) {
	appRouter := s.timerRouter.Group("/app/").Use()
	{
		appRouter.GET("list", controller.GetAppList)
		appRouter.POST("create", controller.CreateApp)
		appRouter.DELETE("deleteApp", controller.DeleteApp)
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
