package ports

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
)

// CommandPorts 写入接口
type CommandPorts interface {
	AuthCommandPorts
	RbacCommandPorts
	NamespaceCommandPorts
}

// QueryPorts 读入接口
type QueryPorts interface {
	AuthQueryPorts
}

// AuthCommandPorts 函数相关接口
type AuthCommandPorts interface {
	Login(ctx context.Context, req *dto.LoginReqDTO) (*dto.LoginRspDTO, error)
	OutLogin(ctx context.Context, req *dto.OutLoginReqDTO) (*dto.OutLoginRspDTO, error)
	Oauth2Callback(ctx context.Context, req *dto.Oauth2CallbackReqDTO) (*dto.Oauth2CallbackRspDTO, error)
	GetCaptcha(ctx context.Context, req *dto.GetCaptchaReqDTO) (
		*dto.GetCaptchaRspDTO, error)
	VerifyCaptcha(ctx context.Context, req *dto.VerifyCaptchaReqDTO) (*dto.VerifyCaptchaRspDTO, error)
	ValidateToken(ctx context.Context, req *dto.ValidateTokenReqDTO) error
}

// NamespaceCommandPorts 命名空间相关接口
type NamespaceCommandPorts interface {
	CreateNamespace(ctx context.Context, req *dto.CreateNamespaceDTO) error
	GetNamespaces(ctx context.Context, req *dto.GetNamespacesReqDTO) ([]*dto.NamespacePermissionDTO, int64, error)
	RegisterNamespaceAPIToken(ctx context.Context, req *dto.RegisterNamespaceAPITokenReqDTO) (string, error)
	GetNamespaceTokens(ctx context.Context, req *dto.PageQueryNamespaceTokenDTO) ([]*dto.NamespaceTokenDTO, int64, error)
}

// RbacCommandPorts 权限操作接口
type RbacCommandPorts interface {
	GetPermissionsForUserInDomain(ctx context.Context, req *dto.RbacReqDTO) ([]string, error)
	GetRolesForUserInDomain(ctx context.Context, req *dto.RbacReqDTO) ([]string, error)
	GetDomainsForUser(ctx context.Context, req *dto.RbacReqDTO) ([]string, error)
	AddPermissionToRole(ctx context.Context, req *dto.RbacReqDTO) error
	DeletePermissionForRole(ctx context.Context, req *dto.RbacReqDTO) error
	AddRoleForUserInDomain(ctx context.Context, req *dto.RbacReqDTO) error
	AddRoleForUsersInDomain(ctx context.Context, req *dto.RbacReqDTO) error
	DeleteRoleForUserInDomain(ctx context.Context, req *dto.RbacReqDTO) error
	HasPermission(ctx context.Context, req *dto.RbacReqDTO) (bool, error)
}

// AuthQueryPorts 鉴权相关接口
type AuthQueryPorts interface {
}
