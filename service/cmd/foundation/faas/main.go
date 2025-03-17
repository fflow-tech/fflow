package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/fflow-tech/fflow/service/cmd/foundation/faas/docs"
	_ "github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"

	"github.com/fflow-tech/fflow/service/cmd/foundation/faas/factory"
	"github.com/fflow-tech/fflow/service/cmd/foundation/faas/service/rpc"
	"github.com/fflow-tech/fflow/service/cmd/foundation/faas/service/web"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/k8s"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/registry"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	httpAddr         = flag.String("http.addr", ":50031", "The http server address")
	grpcAddr         = flag.String("grpc.addr", ":50032", "The grpc server address")
	globalConfigName = flag.String("config.name", "app", "The global config name")
	globalConfigType = flag.String("config.type", "yaml", "The global config type")
	globalConfigPath = flag.String("config.path", "./config/", "The global config path")
)

func main() {
	flag.Parse()
	// 先初始化工厂才能进行后面的操作
	if err := factory.New(factory.WithRegistryClientType(registry.Kubernetes),
		factory.WithConfigClientType(config.Kubernetes),
		factory.WithK8sConfig(k8s.Config{
			GlobalConfigName: *globalConfigName,
			GlobalConfigType: *globalConfigType,
			GlobalConfigPath: *globalConfigPath,
		}),
	); err != nil {
		log.Fatalf("Factory init failed: %v", err)
		panic(err)
	}
	// 初始化服务
	httpServer := web.NewServer(*httpAddr)
	go func() {
		if err := httpServer.Serve(); err != nil {
			log.Fatalf("Http server not serve: %v", err)
		}
		log.Infof("Http server closed")
	}()

	grpcServer := rpc.NewServer(*grpcAddr)
	go func() {
		if err := grpcServer.Serve(); err != nil {
			log.Fatalf("Grpc server not serve: %v", err)
		}
		log.Infof("Grpc server closed")
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个 5 秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 Ctrl+C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify 把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给 quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV) // 此处不会阻塞
	<-quit                                                                // 阻塞在此，当接收到上述两种信号时才会往下执行
	utils.ShutdownGraceful(grpcServer.Close, httpServer.Close)
}
