// Package rpc 提供 grpc 调用的入口
package rpc

import (
	"context"
	"fmt"

	pb "github.com/fflow-tech/fflow/api/foundation/faas"
	"github.com/fflow-tech/fflow/service/cmd/foundation/faas/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// FAASService foundation/faas 后端服务实现
type FAASService struct {
	pb.UnimplementedFaasServer
	domainService *service.DomainService
}

// NewFAASService FAASService构造函数
func NewFAASService(domainService *service.DomainService) *FAASService {
	return &FAASService{domainService: domainService}
}

// NewSucceedRsp 生成成功返回
func NewSucceedRsp() *pb.BasicRsp {
	return &pb.BasicRsp{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
	}
}

// Call 调用函数
func (s *FAASService) Call(ctx context.Context, req *pb.CallReq) (*pb.CallRsp, error) {
	rsp := &pb.CallRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	callReq, err := convertor.FaasConvertor.ConvertCallPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}
	data, err := s.domainService.Commands.CallFunction(ctx, callReq)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	switch data.(type) {
	case string:
		rsp.Output = data.(string)
	default:
		rsp.Output = utils.StructToJsonStr(data)
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
