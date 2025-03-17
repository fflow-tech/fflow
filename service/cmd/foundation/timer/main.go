package main

import (
	"context"
	"flag"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service"
	"github.com/fflow-tech/fflow/service/pkg/k8s"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/fflow-tech/fflow/service/cmd/foundation/timer/docs"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/factory"
	_ "github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	_ "github.com/mbobakov/grpc-consul-resolver"

	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service/rpc"
	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/service/web"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/registry"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	httpAddr         = flag.String("http.addr", ":50021", "The http server address")
	grpcAddr         = flag.String("grpc.addr", ":50022", "The grpc server address")
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
		// 启动内部服务
		service.InitTimerTaskServer()
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
	shutdownGraceful(service.CloseEventServer)
}

// shutdownGraceful 优雅关闭
// 这里主要是用来关闭自身的资源
func shutdownGraceful(fs ...func(chan struct{}) error) {
	log.Infof("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f func(chan struct{}) error) {
			defer wg.Done()

			c := make(chan struct{}, 1)
			go f(c)
			select {
			case <-c:
			case <-ctx.Done():
			}
		}(f)
	}

	wg.Wait()
	log.Infof("Server exit")
}
