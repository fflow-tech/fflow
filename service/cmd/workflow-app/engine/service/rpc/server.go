// Package rpc 提供 grpc 调用的入口
package rpc

import (
	"net"
	"sync"

	pb "github.com/fflow-tech/fflow/api/workflow-app/engine"
	"github.com/fflow-tech/fflow/service/cmd/workflow-app/engine/factory"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/provider"
	"google.golang.org/grpc"
)

// Server Web容器
type Server struct {
	grpcAddr   string
	grpcServer *grpc.Server
}

// NewServer 新建 Web 容器
func NewServer(grpcAddr string) *Server {
	return &Server{grpcAddr: grpcAddr}
}

// Serve 启动监听
func (s *Server) Serve() error {
	log.Infof("Start to bind grpc")
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		panic(err)
	}

	s.grpcServer = grpc.NewServer()
	domainService, err := factory.GetDomainService()
	if err != nil {
		panic(err)
	}

	pb.RegisterWorkflowServer(s.grpcServer, NewWorkflowEngineService(domainService))

	if err := provider.GetRegistryProvider().Register(constants.ServiceName, s.grpcAddr); err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}

var (
	mutex sync.Mutex
)

// Close 关闭 RPC 服务
func (s *Server) Close(ch chan struct{}) error {
	log.Infof("Shutdown RPC Server...")
	defer log.Infof("RPC Server exit")
	defer func() {
		if ch != nil {
			ch <- struct{}{}
		}
	}()

	mutex.Lock()
	defer mutex.Unlock()

	s.grpcServer.GracefulStop()
	return nil
}
