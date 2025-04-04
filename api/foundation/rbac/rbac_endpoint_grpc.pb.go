// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: rbac_endpoint.proto

package rbac

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Rbac_GetPermissionsForUserInDomain_FullMethodName = "/rbac.Rbac/GetPermissionsForUserInDomain"
	Rbac_GetRolesForUserInDomain_FullMethodName       = "/rbac.Rbac/GetRolesForUserInDomain"
	Rbac_GetDomainsForUser_FullMethodName             = "/rbac.Rbac/GetDomainsForUser"
	Rbac_AddPermissionToRole_FullMethodName           = "/rbac.Rbac/AddPermissionToRole"
	Rbac_DeletePermissionForRole_FullMethodName       = "/rbac.Rbac/DeletePermissionForRole"
	Rbac_AddRoleForUserInDomain_FullMethodName        = "/rbac.Rbac/AddRoleForUserInDomain"
	Rbac_AddRoleForUsersInDomain_FullMethodName       = "/rbac.Rbac/AddRoleForUsersInDomain"
	Rbac_DeleteRoleForUserInDomain_FullMethodName     = "/rbac.Rbac/DeleteRoleForUserInDomain"
	Rbac_HasPermission_FullMethodName                 = "/rbac.Rbac/HasPermission"
)

// RbacClient is the client API for Rbac service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RbacClient interface {
	GetPermissionsForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*GetPermissionsForUserInDomainRsp, error)
	GetRolesForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*GetRolesForUserInDomainRsp, error)
	GetDomainsForUser(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*GetDomainsForUserRsp, error)
	AddPermissionToRole(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error)
	DeletePermissionForRole(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error)
	AddRoleForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error)
	AddRoleForUsersInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error)
	DeleteRoleForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error)
	HasPermission(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error)
}

type rbacClient struct {
	cc grpc.ClientConnInterface
}

func NewRbacClient(cc grpc.ClientConnInterface) RbacClient {
	return &rbacClient{cc}
}

func (c *rbacClient) GetPermissionsForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*GetPermissionsForUserInDomainRsp, error) {
	out := new(GetPermissionsForUserInDomainRsp)
	err := c.cc.Invoke(ctx, Rbac_GetPermissionsForUserInDomain_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) GetRolesForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*GetRolesForUserInDomainRsp, error) {
	out := new(GetRolesForUserInDomainRsp)
	err := c.cc.Invoke(ctx, Rbac_GetRolesForUserInDomain_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) GetDomainsForUser(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*GetDomainsForUserRsp, error) {
	out := new(GetDomainsForUserRsp)
	err := c.cc.Invoke(ctx, Rbac_GetDomainsForUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) AddPermissionToRole(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error) {
	out := new(BasicRsp)
	err := c.cc.Invoke(ctx, Rbac_AddPermissionToRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) DeletePermissionForRole(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error) {
	out := new(BasicRsp)
	err := c.cc.Invoke(ctx, Rbac_DeletePermissionForRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) AddRoleForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error) {
	out := new(BasicRsp)
	err := c.cc.Invoke(ctx, Rbac_AddRoleForUserInDomain_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) AddRoleForUsersInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error) {
	out := new(BasicRsp)
	err := c.cc.Invoke(ctx, Rbac_AddRoleForUsersInDomain_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) DeleteRoleForUserInDomain(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error) {
	out := new(BasicRsp)
	err := c.cc.Invoke(ctx, Rbac_DeleteRoleForUserInDomain_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rbacClient) HasPermission(ctx context.Context, in *RbacReq, opts ...grpc.CallOption) (*BasicRsp, error) {
	out := new(BasicRsp)
	err := c.cc.Invoke(ctx, Rbac_HasPermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RbacServer is the server API for Rbac service.
// All implementations must embed UnimplementedRbacServer
// for forward compatibility
type RbacServer interface {
	GetPermissionsForUserInDomain(context.Context, *RbacReq) (*GetPermissionsForUserInDomainRsp, error)
	GetRolesForUserInDomain(context.Context, *RbacReq) (*GetRolesForUserInDomainRsp, error)
	GetDomainsForUser(context.Context, *RbacReq) (*GetDomainsForUserRsp, error)
	AddPermissionToRole(context.Context, *RbacReq) (*BasicRsp, error)
	DeletePermissionForRole(context.Context, *RbacReq) (*BasicRsp, error)
	AddRoleForUserInDomain(context.Context, *RbacReq) (*BasicRsp, error)
	AddRoleForUsersInDomain(context.Context, *RbacReq) (*BasicRsp, error)
	DeleteRoleForUserInDomain(context.Context, *RbacReq) (*BasicRsp, error)
	HasPermission(context.Context, *RbacReq) (*BasicRsp, error)
	mustEmbedUnimplementedRbacServer()
}

// UnimplementedRbacServer must be embedded to have forward compatible implementations.
type UnimplementedRbacServer struct {
}

func (UnimplementedRbacServer) GetPermissionsForUserInDomain(context.Context, *RbacReq) (*GetPermissionsForUserInDomainRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPermissionsForUserInDomain not implemented")
}
func (UnimplementedRbacServer) GetRolesForUserInDomain(context.Context, *RbacReq) (*GetRolesForUserInDomainRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRolesForUserInDomain not implemented")
}
func (UnimplementedRbacServer) GetDomainsForUser(context.Context, *RbacReq) (*GetDomainsForUserRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDomainsForUser not implemented")
}
func (UnimplementedRbacServer) AddPermissionToRole(context.Context, *RbacReq) (*BasicRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPermissionToRole not implemented")
}
func (UnimplementedRbacServer) DeletePermissionForRole(context.Context, *RbacReq) (*BasicRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePermissionForRole not implemented")
}
func (UnimplementedRbacServer) AddRoleForUserInDomain(context.Context, *RbacReq) (*BasicRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRoleForUserInDomain not implemented")
}
func (UnimplementedRbacServer) AddRoleForUsersInDomain(context.Context, *RbacReq) (*BasicRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRoleForUsersInDomain not implemented")
}
func (UnimplementedRbacServer) DeleteRoleForUserInDomain(context.Context, *RbacReq) (*BasicRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRoleForUserInDomain not implemented")
}
func (UnimplementedRbacServer) HasPermission(context.Context, *RbacReq) (*BasicRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasPermission not implemented")
}
func (UnimplementedRbacServer) mustEmbedUnimplementedRbacServer() {}

// UnsafeRbacServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RbacServer will
// result in compilation errors.
type UnsafeRbacServer interface {
	mustEmbedUnimplementedRbacServer()
}

func RegisterRbacServer(s grpc.ServiceRegistrar, srv RbacServer) {
	s.RegisterService(&Rbac_ServiceDesc, srv)
}

func _Rbac_GetPermissionsForUserInDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).GetPermissionsForUserInDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_GetPermissionsForUserInDomain_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).GetPermissionsForUserInDomain(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_GetRolesForUserInDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).GetRolesForUserInDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_GetRolesForUserInDomain_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).GetRolesForUserInDomain(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_GetDomainsForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).GetDomainsForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_GetDomainsForUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).GetDomainsForUser(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_AddPermissionToRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).AddPermissionToRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_AddPermissionToRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).AddPermissionToRole(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_DeletePermissionForRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).DeletePermissionForRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_DeletePermissionForRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).DeletePermissionForRole(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_AddRoleForUserInDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).AddRoleForUserInDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_AddRoleForUserInDomain_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).AddRoleForUserInDomain(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_AddRoleForUsersInDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).AddRoleForUsersInDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_AddRoleForUsersInDomain_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).AddRoleForUsersInDomain(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_DeleteRoleForUserInDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).DeleteRoleForUserInDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_DeleteRoleForUserInDomain_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).DeleteRoleForUserInDomain(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rbac_HasPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RbacReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RbacServer).HasPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rbac_HasPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RbacServer).HasPermission(ctx, req.(*RbacReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Rbac_ServiceDesc is the grpc.ServiceDesc for Rbac service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rbac_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rbac.Rbac",
	HandlerType: (*RbacServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPermissionsForUserInDomain",
			Handler:    _Rbac_GetPermissionsForUserInDomain_Handler,
		},
		{
			MethodName: "GetRolesForUserInDomain",
			Handler:    _Rbac_GetRolesForUserInDomain_Handler,
		},
		{
			MethodName: "GetDomainsForUser",
			Handler:    _Rbac_GetDomainsForUser_Handler,
		},
		{
			MethodName: "AddPermissionToRole",
			Handler:    _Rbac_AddPermissionToRole_Handler,
		},
		{
			MethodName: "DeletePermissionForRole",
			Handler:    _Rbac_DeletePermissionForRole_Handler,
		},
		{
			MethodName: "AddRoleForUserInDomain",
			Handler:    _Rbac_AddRoleForUserInDomain_Handler,
		},
		{
			MethodName: "AddRoleForUsersInDomain",
			Handler:    _Rbac_AddRoleForUsersInDomain_Handler,
		},
		{
			MethodName: "DeleteRoleForUserInDomain",
			Handler:    _Rbac_DeleteRoleForUserInDomain_Handler,
		},
		{
			MethodName: "HasPermission",
			Handler:    _Rbac_HasPermission_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rbac_endpoint.proto",
}
