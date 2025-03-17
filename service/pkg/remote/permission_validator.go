package remote

import (
	"context"
	"fmt"
	pb "github.com/fflow-tech/fflow/api/foundation/auth"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DefaultPermissionValidatorConfig 默认客户端配置
type DefaultPermissionValidatorConfig struct {
	AuthTarget          string `json:"faasTarget,omitempty"`
	LoadBalancingPolicy string `json:"loadBalancingPolicy,omitempty"`
}

// DefaultPermissionValidator 默认校验器
type DefaultPermissionValidator struct {
	config     *DefaultPermissionValidatorConfig
	authClient pb.AuthClient
}

// NewDefaultPermissionValidator 创建默认客户端
func NewDefaultPermissionValidator(config *DefaultPermissionValidatorConfig) (*DefaultPermissionValidator, error) {
	conn, err := grpc.Dial(config.AuthTarget,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, config.LoadBalancingPolicy)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &DefaultPermissionValidator{
		config:     config,
		authClient: pb.NewAuthClient(conn),
	}, nil
}

// ValidateToken 校验 token
func (c *DefaultPermissionValidator) ValidateToken(ctx context.Context, req *ValidateTokenReqDTO) error {
	rsp, err := c.authClient.ValidateToken(ctx, &pb.ValidateTokenReq{
		BasicReq: &pb.BasicReq{
			Namespace:   req.Namespace,
			AccessToken: req.AccessToken,
		},
	})
	if err != nil {
		return err
	}

	if rsp.BasicRsp.Code != errno.OK.Code {
		return fmt.Errorf("invalid token: %s", rsp.BasicRsp.Message)
	}

	return nil
}
