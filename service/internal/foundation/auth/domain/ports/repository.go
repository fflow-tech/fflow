package ports

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/entity"
)

// UserRepository 仓储层接口
type UserRepository interface {
	Create(req *dto.CreateUserDTO) (string, error)
	CreateIfNotExists(req *dto.CreateUserDTO) (*entity.User, error)
	Get(req *dto.GetUserDTO) (*entity.User, error)
	Delete(req *dto.DeleteUserDTO) error
	Update(req *dto.UpdateUserDTO) error
	PageQuery(req *dto.PageQueryUserDTO) ([]*entity.User, error)
}

// NamespaceRepository 仓储层接口
type NamespaceRepository interface {
	Create(req *dto.CreateNamespaceDTO) (string, error)
	CreateIfNotExists(req *dto.CreateNamespaceDTO) (*entity.Namespace, error)
	Get(req *dto.GetNamespaceDTO) (*entity.Namespace, error)
	Delete(req *dto.DeleteNamespaceDTO) error
	Update(req *dto.UpdateNamespaceDTO) error
	PageQuery(req *dto.PageQueryNamespaceDTO) ([]*entity.Namespace, int64, error)
	Count(req *dto.PageQueryNamespaceDTO) (int64, error)
}

// NamespaceTokenRepository 仓储层接口
type NamespaceTokenRepository interface {
	Create(req *dto.CreateNamespaceTokenDTO) (string, error)
	CreateIfNotExists(req *dto.CreateNamespaceTokenDTO) (*entity.NamespaceToken, error)
	Get(req *dto.GetNamespaceTokenDTO) (*entity.NamespaceToken, error)
	Delete(req *dto.DeleteNamespaceTokenDTO) error
	Update(req *dto.UpdateNamespaceTokenDTO) error
	PageQuery(req *dto.PageQueryNamespaceTokenDTO) ([]*entity.NamespaceToken, int64, error)
	Count(req *dto.PageQueryNamespaceTokenDTO) (int64, error)
}
