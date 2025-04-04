// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: rbac_endpoint.proto

package rbac

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RbacReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User        string   `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Role        string   `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	Domain      string   `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Object      string   `protobuf:"bytes,4,opt,name=object,proto3" json:"object,omitempty"`
	Permission  string   `protobuf:"bytes,5,opt,name=permission,proto3" json:"permission,omitempty"`
	Permissions []string `protobuf:"bytes,6,rep,name=permissions,proto3" json:"permissions,omitempty"`
	Users       []string `protobuf:"bytes,7,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *RbacReq) Reset() {
	*x = RbacReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rbac_endpoint_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RbacReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RbacReq) ProtoMessage() {}

func (x *RbacReq) ProtoReflect() protoreflect.Message {
	mi := &file_rbac_endpoint_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RbacReq.ProtoReflect.Descriptor instead.
func (*RbacReq) Descriptor() ([]byte, []int) {
	return file_rbac_endpoint_proto_rawDescGZIP(), []int{0}
}

func (x *RbacReq) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *RbacReq) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *RbacReq) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *RbacReq) GetObject() string {
	if x != nil {
		return x.Object
	}
	return ""
}

func (x *RbacReq) GetPermission() string {
	if x != nil {
		return x.Permission
	}
	return ""
}

func (x *RbacReq) GetPermissions() []string {
	if x != nil {
		return x.Permissions
	}
	return nil
}

func (x *RbacReq) GetUsers() []string {
	if x != nil {
		return x.Users
	}
	return nil
}

type GetPermissionsForUserInDomainRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BasicRsp    *BasicRsp `protobuf:"bytes,1,opt,name=basicRsp,proto3" json:"basicRsp,omitempty"`
	Permissions []string  `protobuf:"bytes,2,rep,name=permissions,proto3" json:"permissions,omitempty"`
}

func (x *GetPermissionsForUserInDomainRsp) Reset() {
	*x = GetPermissionsForUserInDomainRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rbac_endpoint_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPermissionsForUserInDomainRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPermissionsForUserInDomainRsp) ProtoMessage() {}

func (x *GetPermissionsForUserInDomainRsp) ProtoReflect() protoreflect.Message {
	mi := &file_rbac_endpoint_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPermissionsForUserInDomainRsp.ProtoReflect.Descriptor instead.
func (*GetPermissionsForUserInDomainRsp) Descriptor() ([]byte, []int) {
	return file_rbac_endpoint_proto_rawDescGZIP(), []int{1}
}

func (x *GetPermissionsForUserInDomainRsp) GetBasicRsp() *BasicRsp {
	if x != nil {
		return x.BasicRsp
	}
	return nil
}

func (x *GetPermissionsForUserInDomainRsp) GetPermissions() []string {
	if x != nil {
		return x.Permissions
	}
	return nil
}

type GetRolesForUserInDomainRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BasicRsp *BasicRsp `protobuf:"bytes,1,opt,name=basicRsp,proto3" json:"basicRsp,omitempty"`
	Roles    []string  `protobuf:"bytes,2,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *GetRolesForUserInDomainRsp) Reset() {
	*x = GetRolesForUserInDomainRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rbac_endpoint_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRolesForUserInDomainRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRolesForUserInDomainRsp) ProtoMessage() {}

func (x *GetRolesForUserInDomainRsp) ProtoReflect() protoreflect.Message {
	mi := &file_rbac_endpoint_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRolesForUserInDomainRsp.ProtoReflect.Descriptor instead.
func (*GetRolesForUserInDomainRsp) Descriptor() ([]byte, []int) {
	return file_rbac_endpoint_proto_rawDescGZIP(), []int{2}
}

func (x *GetRolesForUserInDomainRsp) GetBasicRsp() *BasicRsp {
	if x != nil {
		return x.BasicRsp
	}
	return nil
}

func (x *GetRolesForUserInDomainRsp) GetRoles() []string {
	if x != nil {
		return x.Roles
	}
	return nil
}

type GetDomainsForUserRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BasicRsp *BasicRsp `protobuf:"bytes,1,opt,name=basicRsp,proto3" json:"basicRsp,omitempty"`
	Domains  []string  `protobuf:"bytes,2,rep,name=domains,proto3" json:"domains,omitempty"`
}

func (x *GetDomainsForUserRsp) Reset() {
	*x = GetDomainsForUserRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rbac_endpoint_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDomainsForUserRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDomainsForUserRsp) ProtoMessage() {}

func (x *GetDomainsForUserRsp) ProtoReflect() protoreflect.Message {
	mi := &file_rbac_endpoint_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDomainsForUserRsp.ProtoReflect.Descriptor instead.
func (*GetDomainsForUserRsp) Descriptor() ([]byte, []int) {
	return file_rbac_endpoint_proto_rawDescGZIP(), []int{3}
}

func (x *GetDomainsForUserRsp) GetBasicRsp() *BasicRsp {
	if x != nil {
		return x.BasicRsp
	}
	return nil
}

func (x *GetDomainsForUserRsp) GetDomains() []string {
	if x != nil {
		return x.Domains
	}
	return nil
}

type BasicRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BasicRsp) Reset() {
	*x = BasicRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rbac_endpoint_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicRsp) ProtoMessage() {}

func (x *BasicRsp) ProtoReflect() protoreflect.Message {
	mi := &file_rbac_endpoint_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicRsp.ProtoReflect.Descriptor instead.
func (*BasicRsp) Descriptor() ([]byte, []int) {
	return file_rbac_endpoint_proto_rawDescGZIP(), []int{4}
}

func (x *BasicRsp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *BasicRsp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_rbac_endpoint_proto protoreflect.FileDescriptor

var file_rbac_endpoint_proto_rawDesc = []byte{
	0x0a, 0x13, 0x72, 0x62, 0x61, 0x63, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x72, 0x62, 0x61, 0x63, 0x22, 0xb9, 0x01, 0x0a, 0x07,
	0x52, 0x62, 0x61, 0x63, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x20, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x70, 0x0a, 0x20, 0x47, 0x65, 0x74, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52, 0x73, 0x70, 0x12, 0x2a, 0x0a, 0x08, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x72, 0x62, 0x61, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x52, 0x08, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x5e, 0x0a, 0x1a, 0x47, 0x65, 0x74,
	0x52, 0x6f, 0x6c, 0x65, 0x73, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x44, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x52, 0x73, 0x70, 0x12, 0x2a, 0x0a, 0x08, 0x62, 0x61, 0x73, 0x69, 0x63,
	0x52, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72, 0x62, 0x61, 0x63,
	0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x52, 0x08, 0x62, 0x61, 0x73, 0x69, 0x63,
	0x52, 0x73, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x22, 0x5c, 0x0a, 0x14, 0x47, 0x65, 0x74,
	0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x73,
	0x70, 0x12, 0x2a, 0x0a, 0x08, 0x62, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63,
	0x52, 0x73, 0x70, 0x52, 0x08, 0x62, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x22, 0x38, 0x0a, 0x08, 0x42, 0x61, 0x73, 0x69, 0x63,
	0x52, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x32, 0xb9, 0x04, 0x0a, 0x04, 0x52, 0x62, 0x61, 0x63, 0x12, 0x56, 0x0a, 0x1d, 0x47, 0x65,
	0x74, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x46, 0x6f, 0x72, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x0d, 0x2e, 0x72, 0x62,
	0x61, 0x63, 0x2e, 0x52, 0x62, 0x61, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x26, 0x2e, 0x72, 0x62, 0x61,
	0x63, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52,
	0x73, 0x70, 0x12, 0x4a, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x46, 0x6f,
	0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x0d, 0x2e,
	0x72, 0x62, 0x61, 0x63, 0x2e, 0x52, 0x62, 0x61, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x20, 0x2e, 0x72,
	0x62, 0x61, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x46, 0x6f, 0x72, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x52, 0x73, 0x70, 0x12, 0x3e,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x46, 0x6f, 0x72, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x0d, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x52, 0x62, 0x61, 0x63, 0x52,
	0x65, 0x71, 0x1a, 0x1a, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x73, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x73, 0x70, 0x12, 0x34,
	0x0a, 0x13, 0x41, 0x64, 0x64, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54,
	0x6f, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0d, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x52, 0x62, 0x61,
	0x63, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x69,
	0x63, 0x52, 0x73, 0x70, 0x12, 0x38, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x12,
	0x0d, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x52, 0x62, 0x61, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x0e,
	0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x12, 0x37,
	0x0a, 0x16, 0x41, 0x64, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x0d, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e,
	0x52, 0x62, 0x61, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x42,
	0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x12, 0x38, 0x0a, 0x17, 0x41, 0x64, 0x64, 0x52, 0x6f,
	0x6c, 0x65, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x73, 0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x12, 0x0d, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x52, 0x62, 0x61, 0x63, 0x52, 0x65,
	0x71, 0x1a, 0x0e, 0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73,
	0x70, 0x12, 0x3a, 0x0a, 0x19, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x46,
	0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x0d,
	0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x52, 0x62, 0x61, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e,
	0x72, 0x62, 0x61, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x12, 0x2e, 0x0a,
	0x0d, 0x48, 0x61, 0x73, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0d,
	0x2e, 0x72, 0x62, 0x61, 0x63, 0x2e, 0x52, 0x62, 0x61, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e,
	0x72, 0x62, 0x61, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x73, 0x70, 0x42, 0x44, 0x0a,
	0x12, 0x69, 0x6f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x66, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x72,
	0x62, 0x61, 0x63, 0x42, 0x09, 0x72, 0x62, 0x61, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x21, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e,
	0x6f, 0x72, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x66, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x72,
	0x62, 0x61, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rbac_endpoint_proto_rawDescOnce sync.Once
	file_rbac_endpoint_proto_rawDescData = file_rbac_endpoint_proto_rawDesc
)

func file_rbac_endpoint_proto_rawDescGZIP() []byte {
	file_rbac_endpoint_proto_rawDescOnce.Do(func() {
		file_rbac_endpoint_proto_rawDescData = protoimpl.X.CompressGZIP(file_rbac_endpoint_proto_rawDescData)
	})
	return file_rbac_endpoint_proto_rawDescData
}

var file_rbac_endpoint_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_rbac_endpoint_proto_goTypes = []interface{}{
	(*RbacReq)(nil),                          // 0: rbac.RbacReq
	(*GetPermissionsForUserInDomainRsp)(nil), // 1: rbac.GetPermissionsForUserInDomainRsp
	(*GetRolesForUserInDomainRsp)(nil),       // 2: rbac.GetRolesForUserInDomainRsp
	(*GetDomainsForUserRsp)(nil),             // 3: rbac.GetDomainsForUserRsp
	(*BasicRsp)(nil),                         // 4: rbac.BasicRsp
}
var file_rbac_endpoint_proto_depIdxs = []int32{
	4,  // 0: rbac.GetPermissionsForUserInDomainRsp.basicRsp:type_name -> rbac.BasicRsp
	4,  // 1: rbac.GetRolesForUserInDomainRsp.basicRsp:type_name -> rbac.BasicRsp
	4,  // 2: rbac.GetDomainsForUserRsp.basicRsp:type_name -> rbac.BasicRsp
	0,  // 3: rbac.Rbac.GetPermissionsForUserInDomain:input_type -> rbac.RbacReq
	0,  // 4: rbac.Rbac.GetRolesForUserInDomain:input_type -> rbac.RbacReq
	0,  // 5: rbac.Rbac.GetDomainsForUser:input_type -> rbac.RbacReq
	0,  // 6: rbac.Rbac.AddPermissionToRole:input_type -> rbac.RbacReq
	0,  // 7: rbac.Rbac.DeletePermissionForRole:input_type -> rbac.RbacReq
	0,  // 8: rbac.Rbac.AddRoleForUserInDomain:input_type -> rbac.RbacReq
	0,  // 9: rbac.Rbac.AddRoleForUsersInDomain:input_type -> rbac.RbacReq
	0,  // 10: rbac.Rbac.DeleteRoleForUserInDomain:input_type -> rbac.RbacReq
	0,  // 11: rbac.Rbac.HasPermission:input_type -> rbac.RbacReq
	1,  // 12: rbac.Rbac.GetPermissionsForUserInDomain:output_type -> rbac.GetPermissionsForUserInDomainRsp
	2,  // 13: rbac.Rbac.GetRolesForUserInDomain:output_type -> rbac.GetRolesForUserInDomainRsp
	3,  // 14: rbac.Rbac.GetDomainsForUser:output_type -> rbac.GetDomainsForUserRsp
	4,  // 15: rbac.Rbac.AddPermissionToRole:output_type -> rbac.BasicRsp
	4,  // 16: rbac.Rbac.DeletePermissionForRole:output_type -> rbac.BasicRsp
	4,  // 17: rbac.Rbac.AddRoleForUserInDomain:output_type -> rbac.BasicRsp
	4,  // 18: rbac.Rbac.AddRoleForUsersInDomain:output_type -> rbac.BasicRsp
	4,  // 19: rbac.Rbac.DeleteRoleForUserInDomain:output_type -> rbac.BasicRsp
	4,  // 20: rbac.Rbac.HasPermission:output_type -> rbac.BasicRsp
	12, // [12:21] is the sub-list for method output_type
	3,  // [3:12] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_rbac_endpoint_proto_init() }
func file_rbac_endpoint_proto_init() {
	if File_rbac_endpoint_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rbac_endpoint_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RbacReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rbac_endpoint_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPermissionsForUserInDomainRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rbac_endpoint_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRolesForUserInDomainRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rbac_endpoint_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDomainsForUserRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rbac_endpoint_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rbac_endpoint_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rbac_endpoint_proto_goTypes,
		DependencyIndexes: file_rbac_endpoint_proto_depIdxs,
		MessageInfos:      file_rbac_endpoint_proto_msgTypes,
	}.Build()
	File_rbac_endpoint_proto = out.File
	file_rbac_endpoint_proto_rawDesc = nil
	file_rbac_endpoint_proto_goTypes = nil
	file_rbac_endpoint_proto_depIdxs = nil
}
