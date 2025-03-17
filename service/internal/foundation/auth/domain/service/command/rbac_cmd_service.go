package command

import (
	"context"
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/constants"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/repository/repo"
	pc "github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"sort"
	"strings"

	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
)

// RbacCommandService 权限服务
type RbacCommandService struct {
	namespaceRepo      ports.NamespaceRepository
	namespaceTokenRepo ports.NamespaceTokenRepository
	rbacClient         *RbacClient
}

// NewRbacCommandService 新建服务
func NewRbacCommandService(namespaceRepo *repo.NamespaceRepo,
	namespaceTokenRepo *repo.NamespaceTokenRepo,
	rbacClient *RbacClient) *RbacCommandService {
	return &RbacCommandService{
		namespaceRepo:      namespaceRepo,
		namespaceTokenRepo: namespaceTokenRepo,
		rbacClient:         rbacClient,
	}
}

// GetPermissionsForUserInDomain 获取用户在域内的权限点
func (s *RbacCommandService) GetPermissionsForUserInDomain(ctx context.Context,
	req *dto.RbacReqDTO) ([]string, error) {
	var results []string
	roles := s.rbacClient.GetRolesForUserInDomain(req.User, req.Domain)
	for _, role := range roles {
		rolePermissionsList := s.rbacClient.GetPermissionsForUserInDomain(role, "")
		for _, permissions := range rolePermissionsList {
			results = append(results, strings.Join([]string{permissions[2], permissions[3]}, "#"))
		}
	}

	return results, nil
}

// GetRolesForUserInDomain 获取用户在域内的角色
func (s *RbacCommandService) GetRolesForUserInDomain(ctx context.Context,
	req *dto.RbacReqDTO) ([]string, error) {
	return s.rbacClient.GetRolesForUserInDomain(req.User, req.Domain), nil
}

// AddPermissionToRole 添加权限给对应角色
func (s *RbacCommandService) AddPermissionToRole(ctx context.Context, req *dto.RbacReqDTO) error {
	if req.Permission != "" {
		_, err := s.rbacClient.AddPermissionForUser(req.Role, req.Domain, req.Object, req.Permission)
		if err != nil {
			return err
		}
	}

	for _, permission := range req.Permissions {
		_, err := s.rbacClient.AddPermissionForUser(req.Role, req.Domain, req.Object, permission)

		if err != nil {
			return err
		}
	}

	return nil
}

// DeletePermissionForRole 在角色里删除对应的权限
func (s *RbacCommandService) DeletePermissionForRole(ctx context.Context, req *dto.RbacReqDTO) error {
	_, err := s.rbacClient.DeletePermissionForUser(req.Role, req.Domain, req.Object, req.Permission)
	return err
}

// AddRoleForUserInDomain 给单个用户添加对应的角色
func (s *RbacCommandService) AddRoleForUserInDomain(ctx context.Context, req *dto.RbacReqDTO) error {
	_, err := s.rbacClient.AddRoleForUserInDomain(req.User, req.Role, req.Domain)
	return err
}

// GetDomainsForUser 获取用户的 Domain
func (s *RbacCommandService) GetDomainsForUser(ctx context.Context, req *dto.RbacReqDTO) ([]string, error) {
	namespaces, err := s.rbacClient.GetDomainsForUser(req.User)
	sort.Strings(namespaces)
	return append(namespaces, constants.DefaultNamespace), err
}

// AddRoleForUsersInDomain 给多个用户添加对应的角色
func (s *RbacCommandService) AddRoleForUsersInDomain(ctx context.Context, req *dto.RbacReqDTO) error {
	var addReqs []*dto.RbacReqDTO
	for _, user := range req.Users {
		getReq := &dto.RbacReqDTO{
			User:   user,
			Domain: req.Domain,
		}

		roles, err := s.GetRolesForUserInDomain(ctx, getReq)
		if utils.StrContains(roles, req.Role) {
			continue
		}

		addReq := &dto.RbacReqDTO{
			User:   user,
			Role:   req.Role,
			Domain: req.Domain,
		}
		err = s.AddRoleForUserInDomain(ctx, addReq)
		if err != nil {
			s.rollbackAddRoles(ctx, addReqs)
			return err
		}
		addReqs = append(addReqs, addReq)
	}

	return nil
}

// rollbackAddRoles 如果没添加成功就把之前加入成功的删除掉
func (s *RbacCommandService) rollbackAddRoles(ctx context.Context, reqs []*dto.RbacReqDTO) {
	for _, req := range reqs {
		err := s.DeleteRoleForUserInDomain(ctx, req)
		if err != nil {
			log.Errorf("Failed to DeleteRoleForUserInDomain, caused by %s", err.Error())
			continue
		}
	}
}

// DeleteRoleForUserInDomain 删除域内用户对应的角色
func (s *RbacCommandService) DeleteRoleForUserInDomain(ctx context.Context, req *dto.RbacReqDTO) error {
	_, err := s.rbacClient.DeleteRoleForUserInDomain(req.User, req.Role, req.Domain)
	return err
}

// HasPermission 判断是否有权限
func (s *RbacCommandService) HasPermission(ctx context.Context, req *dto.RbacReqDTO) (bool, error) {
	permissions, err := s.GetPermissionsForUserInDomain(ctx, req)
	if err != nil {
		return false, err
	}

	if utils.StrContains(permissions, permissionKey(req.Object, req.Permission)) {
		return true, nil
	}

	return false, nil
}

// CreateNamespace 添加命名空间
func (s *RbacCommandService) CreateNamespace(ctx context.Context, req *dto.CreateNamespaceDTO) error {
	_, err := s.namespaceRepo.Create(req)
	if err != nil {
		return err
	}
	// 添加命名空间的时候给创建者添加命名空间的管理员权限
	return s.AddRoleForUserInDomain(ctx, &dto.RbacReqDTO{
		User:   req.Creator,
		Role:   constants.NamespaceAdminRole,
		Domain: req.Namespace,
	})
}

// GetNamespaces 获取所有的命名空间
func (s *RbacCommandService) GetNamespaces(ctx context.Context, req *dto.GetNamespacesReqDTO) (
	[]*dto.NamespacePermissionDTO, int64, error) {
	if req.PageQuery == nil {
		req.PageQuery = pc.NewDefaultPageQuery()
	}
	if req.Order == nil {
		req.Order = pc.NewDefaultOrder()
	}
	namespaces, total, err := s.namespaceRepo.PageQuery(&dto.PageQueryNamespaceDTO{
		Namespace: req.Namespace,
		PageQuery: req.PageQuery,
		Order:     req.Order,
	})
	if err != nil {
		return nil, 0, err
	}

	var namespacePermissions []*dto.NamespacePermissionDTO

	for _, namespace := range namespaces {
		permissions, err := s.GetPermissionsForUserInDomain(ctx,
			&dto.RbacReqDTO{User: req.Username, Domain: namespace.Namespace},
		)
		if err != nil {
			return nil, 0, err
		}
		namespacePermissions = append(namespacePermissions,
			&dto.NamespacePermissionDTO{
				Namespace:   namespace.Namespace,
				Permissions: permissions,
				Creator:     namespace.Creator,
				CreatedAt:   namespace.CreatedAt,
			})
	}

	return namespacePermissions, total, nil
}

func permissionKey(object, permission string) string {
	return strings.Join([]string{object, permission}, "#")
}

// RegisterNamespaceAPIToken 获取命名空间的 Token
func (s *RbacCommandService) RegisterNamespaceAPIToken(ctx context.Context, req *dto.RegisterNamespaceAPITokenReqDTO) (
	string, error) {
	ok, err := s.HasPermission(ctx, &dto.RbacReqDTO{
		User:       req.Creator,
		Domain:     req.Namespace,
		Object:     "token",
		Permission: "write"},
	)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("user %s has no permission to create token for namespace %s", req.Creator, req.Namespace)
	}

	token := utils.GenerateToken()
	if len(req.Name) <= 0 {
		req.Name = "My Test Key"
	}

	_, err = s.namespaceTokenRepo.Create(&dto.CreateNamespaceTokenDTO{
		Name:      req.Name,
		Namespace: req.Namespace,
		Creator:   req.Creator,
		Token:     token,
	})
	return token, err
}

// GetNamespaceTokens 获取命名空间的 Token 列表
func (s *RbacCommandService) GetNamespaceTokens(ctx context.Context, req *dto.PageQueryNamespaceTokenDTO) (
	[]*dto.NamespaceTokenDTO, int64, error) {
	ok, err := s.HasPermission(ctx, &dto.RbacReqDTO{
		User:       req.Creator,
		Domain:     req.Namespace,
		Object:     "token",
		Permission: "read"},
	)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		return nil, 0, fmt.Errorf("user %s has no permission to read token for namespace %s", req.Creator, req.Namespace)
	}

	if req.PageQuery == nil {
		req.PageQuery = pc.NewDefaultPageQuery()
	}
	if req.Order == nil {
		req.Order = pc.NewDefaultOrder()
	}

	req.Creator = ""
	tokens, total, err := s.namespaceTokenRepo.PageQuery(req)
	if err != nil {
		return nil, 0, err
	}

	return convertor.AuthConvertor.ConvertNamespaceTokenEntitiesToDTOs(tokens), total, nil
}
