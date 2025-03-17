package rpc

import (
	"context"
	"github.com/fflow-tech/fflow/service/cmd/foundation/auth/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/errno"

	pb "github.com/fflow-tech/fflow/api/foundation/rbac"
)

// RbacService Rbac服务
type RbacService struct {
	pb.UnimplementedRbacServer
	domainService *service.DomainService
}

// NewRbacService 新建
func NewRbacService(domainService *service.DomainService) *RbacService {
	return &RbacService{domainService: domainService}
}

// GetPermissionsForUserInDomain 获取用户在域内的权限
func (s *RbacService) GetPermissionsForUserInDomain(ctx context.Context, req *pb.RbacReq) (
	*pb.GetPermissionsForUserInDomainRsp, error) {
	rsp := &pb.GetPermissionsForUserInDomainRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		baseRsp := NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		*rsp = NewGetPermissionsForUserInDomainRsp(&baseRsp, []string{})
		return rsp, nil
	}

	permissions, err := s.domainService.Commands.GetPermissionsForUserInDomain(ctx, rbacReq)
	if err != nil {
		baseRsp := NewRbacFailedRsp(errno.Internal.Code, err.Error())
		*rsp = NewGetPermissionsForUserInDomainRsp(&baseRsp, []string{})
		return rsp, nil
	}

	baseRsp := NewRbacSucceedRsp()
	*rsp = NewGetPermissionsForUserInDomainRsp(&baseRsp, permissions)
	return rsp, nil
}

// GetRolesForUserInDomain 获取用户在域内的角色
func (s *RbacService) GetRolesForUserInDomain(ctx context.Context, req *pb.RbacReq) (*pb.GetRolesForUserInDomainRsp, error) {
	rsp := &pb.GetRolesForUserInDomainRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		baseRsp := NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		*rsp = NewGetRolesForUserInDomainRsp(&baseRsp, []string{})
		return rsp, nil
	}

	roles, err := s.domainService.Commands.GetRolesForUserInDomain(ctx, rbacReq)
	if err != nil {
		baseRsp := NewRbacFailedRsp(errno.Internal.Code, err.Error())
		*rsp = NewGetRolesForUserInDomainRsp(&baseRsp, []string{})
		return rsp, nil
	}

	baseRsp := NewRbacSucceedRsp()
	*rsp = NewGetRolesForUserInDomainRsp(&baseRsp, roles)
	return rsp, nil
}

// GetDomainsForUser 获取用户所有的域
func (s *RbacService) GetDomainsForUser(ctx context.Context, req *pb.RbacReq) (*pb.GetDomainsForUserRsp, error) {
	rsp := &pb.GetDomainsForUserRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		baseRsp := NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		*rsp = NewGetDomainsForUserRsp(&baseRsp, []string{})
		return rsp, nil
	}

	domains, err := s.domainService.Commands.GetDomainsForUser(ctx, rbacReq)
	if err != nil {
		baseRsp := NewRbacFailedRsp(errno.Internal.Code, err.Error())
		*rsp = NewGetDomainsForUserRsp(&baseRsp, []string{})
		return rsp, nil
	}

	baseRsp := NewRbacSucceedRsp()
	*rsp = NewGetDomainsForUserRsp(&baseRsp, domains)
	return rsp, nil
}

// AddPermissionToRole 给角色添加权限
func (s *RbacService) AddPermissionToRole(ctx context.Context, req *pb.RbacReq) (*pb.BasicRsp, error) {
	rsp := &pb.BasicRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		*rsp = NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := s.domainService.Commands.AddPermissionToRole(ctx, rbacReq); err != nil {
		*rsp = NewRbacFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	*rsp = NewRbacSucceedRsp()
	return rsp, nil
}

// DeletePermissionForRole 给用户删除权限
func (s *RbacService) DeletePermissionForRole(ctx context.Context, req *pb.RbacReq) (*pb.BasicRsp, error) {
	rsp := &pb.BasicRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		*rsp = NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := s.domainService.Commands.DeletePermissionForRole(ctx, rbacReq); err != nil {
		*rsp = NewRbacFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	*rsp = NewRbacSucceedRsp()
	return rsp, nil
}

// AddRoleForUserInDomain 给域内单个用户添加角色
func (s *RbacService) AddRoleForUserInDomain(ctx context.Context, req *pb.RbacReq) (*pb.BasicRsp, error) {
	rsp := &pb.BasicRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		*rsp = NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := s.domainService.Commands.AddRoleForUserInDomain(ctx, rbacReq); err != nil {
		*rsp = NewRbacFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	*rsp = NewRbacSucceedRsp()
	return rsp, nil
}

// AddRoleForUsersInDomain 给域内多个用户添加角色
func (s *RbacService) AddRoleForUsersInDomain(ctx context.Context, req *pb.RbacReq) (*pb.BasicRsp, error) {
	rsp := &pb.BasicRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		*rsp = NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := s.domainService.Commands.AddRoleForUsersInDomain(ctx, rbacReq); err != nil {
		*rsp = NewRbacFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	*rsp = NewRbacSucceedRsp()
	return rsp, nil
}

// DeleteRoleForUserInDomain 给域内用户删除角色
func (s *RbacService) DeleteRoleForUserInDomain(ctx context.Context, req *pb.RbacReq) (*pb.BasicRsp, error) {
	rsp := &pb.BasicRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		*rsp = NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := s.domainService.Commands.DeleteRoleForUserInDomain(ctx, rbacReq); err != nil {
		*rsp = NewRbacFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	*rsp = NewRbacSucceedRsp()
	return rsp, nil
}

// HasPermission 判断用户是否有权限
func (s *RbacService) HasPermission(ctx context.Context, req *pb.RbacReq) (*pb.BasicRsp, error) {
	rsp := &pb.BasicRsp{}
	rbacReq, err := convertor.RbacConvertor.ConvertPbToDTO(req)
	if err != nil {
		*rsp = NewRbacFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	hasPermission, err := s.domainService.Commands.HasPermission(ctx, rbacReq)
	if err != nil {
		*rsp = NewRbacFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	if !hasPermission {
		*rsp = NewRbacFailedRsp(errno.Unauthenticated.Code, errno.Unauthenticated.Message)
		return rsp, nil
	}

	*rsp = NewRbacSucceedRsp()
	return rsp, nil
}

// NewRbacSucceedRsp 生成成功返回
func NewRbacSucceedRsp() pb.BasicRsp {
	return pb.BasicRsp{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
	}
}

// NewRbacFailedRsp 通过自定义的错误码生成请求返回
func NewRbacFailedRsp(code int32, message string) pb.BasicRsp {
	return pb.BasicRsp{
		Code:    code,
		Message: message,
	}
}

// NewGetPermissionsForUserInDomainRsp 初始化用户权限的返回
func NewGetPermissionsForUserInDomainRsp(basicRsp *pb.BasicRsp,
	permissions []string) pb.GetPermissionsForUserInDomainRsp {
	return pb.GetPermissionsForUserInDomainRsp{BasicRsp: basicRsp, Permissions: permissions}
}

// NewGetRolesForUserInDomainRsp 初始化用户角色的返回
func NewGetRolesForUserInDomainRsp(baseRsp *pb.BasicRsp, roles []string) pb.GetRolesForUserInDomainRsp {
	return pb.GetRolesForUserInDomainRsp{BasicRsp: baseRsp, Roles: roles}
}

// NewGetDomainsForUserRsp 返回用户所有的域
func NewGetDomainsForUserRsp(baseRsp *pb.BasicRsp, domains []string) pb.GetDomainsForUserRsp {
	return pb.GetDomainsForUserRsp{BasicRsp: baseRsp, Domains: domains}
}
