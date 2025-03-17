package storage

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
)

// Transaction 执行事务
type Transaction interface {
	Transaction(fun func(*mysql.Client) error) error
}

// UserDAO 存储层接口
type UserDAO interface {
	Transaction
	Create(*dto.CreateUserDTO) (*po.UserPO, error)
	Get(*dto.GetUserDTO) (*po.UserPO, error)
	Delete(*dto.DeleteUserDTO) error
	Update(*dto.UpdateUserDTO) error
	PageQuery(userDTO *dto.PageQueryUserDTO) ([]*po.UserPO, error)
	Count(*dto.PageQueryUserDTO) (int64, error)
}

// NamespaceDAO 存储层接口
type NamespaceDAO interface {
	Transaction
	Create(namespaceDTO *dto.CreateNamespaceDTO) (*po.NamespacePO, error)
	Get(*dto.GetNamespaceDTO) (*po.NamespacePO, error)
	Delete(*dto.DeleteNamespaceDTO) error
	Update(*dto.UpdateNamespaceDTO) error
	PageQuery(NamespaceDTO *dto.PageQueryNamespaceDTO) ([]*po.NamespacePO, error)
	Count(*dto.PageQueryNamespaceDTO) (int64, error)
}

// NamespaceTokenDAO 存储层接口
type NamespaceTokenDAO interface {
	Transaction
	Create(*dto.CreateNamespaceTokenDTO) (*po.NamespaceTokenPO, error)
	Get(*dto.GetNamespaceTokenDTO) (*po.NamespaceTokenPO, error)
	Delete(*dto.DeleteNamespaceTokenDTO) error
	Update(*dto.UpdateNamespaceTokenDTO) error
	PageQuery(NamespaceTokenDTO *dto.PageQueryNamespaceTokenDTO) ([]*po.NamespaceTokenPO, error)
	Count(*dto.PageQueryNamespaceTokenDTO) (int64, error)
}
