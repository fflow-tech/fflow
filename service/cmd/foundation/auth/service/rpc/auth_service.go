// Package rpc 提供 grpc 调用的入口
package rpc

import (
	"context"
	"fmt"

	pb "github.com/fflow-tech/fflow/api/foundation/auth"
	"github.com/fflow-tech/fflow/service/cmd/foundation/auth/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/errno"
)

// AuthService foundation/Auth 后端服务实现
type AuthService struct {
	pb.UnimplementedAuthServer
	domainService *service.DomainService
}

// NewAuthService AuthService构造函数
func NewAuthService(domainService *service.DomainService) *AuthService {
	return &AuthService{domainService: domainService}
}

// NewSucceedRsp 生成成功返回
func NewSucceedRsp() *pb.BasicRsp {
	return &pb.BasicRsp{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
	}
}

// ValidateToken 校验 Token
func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenReq) (*pb.ValidateTokenRsp, error) {
	rsp := &pb.ValidateTokenRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	validateTokenReqDTO, err := convertor.AuthConvertor.ConvertValidateTokenPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	if err := s.domainService.Commands.ValidateToken(ctx, validateTokenReqDTO); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.PermissionDenied.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// NewFailedRsp 通过自定义的错误码生成请求返回
func NewFailedRsp(code int32, message string) *pb.BasicRsp {
	return &pb.BasicRsp{
		Code:    code,
		Message: message,
	}
}

func validateBasicReq(req *pb.BasicReq) error {
	if req == nil || req.Namespace == "" {
		return fmt.Errorf("the req or req's operator must not be empty")
	}

	return nil
}
