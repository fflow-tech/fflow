package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/entity"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/utils"
)

var (
	AuthConvertor = &authConvertorImpl{}
)

type authConvertorImpl struct{}

// ConvertPOToEntity 转换
func (c *authConvertorImpl) ConvertPOToEntity(p *po.UserPO) (*entity.User, error) {
	u := &entity.User{}
	if err := copier.Copy(u, p); err != nil {
		return nil, err
	}
	u.ID = utils.ToString(p.ID)
	u.Phone = p.Phone.String
	u.Email = p.Email.String
	return u, nil
}

// ConvertPOsToEntities 转换
func (c *authConvertorImpl) ConvertPOsToEntities(p []*po.UserPO) ([]*entity.User, error) {
	users := make([]*entity.User, 0, len(p))
	for _, userPO := range p {
		user, err := c.ConvertPOToEntity(userPO)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

var (
	NamespaceConvertor = &namespaceConvertorImpl{}
)

type namespaceConvertorImpl struct{}

// ConvertPOToEntity 转换
func (c *namespaceConvertorImpl) ConvertPOToEntity(p *po.NamespacePO) (*entity.Namespace, error) {
	u := &entity.Namespace{}
	if err := copier.Copy(u, p); err != nil {
		return nil, err
	}
	return u, nil
}

// ConvertPOsToEntities 转换
func (c *namespaceConvertorImpl) ConvertPOsToEntities(p []*po.NamespacePO) ([]*entity.Namespace, error) {
	namespaces := make([]*entity.Namespace, 0, len(p))
	for _, namespacePO := range p {
		namespace, err := c.ConvertPOToEntity(namespacePO)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, namespace)
	}
	return namespaces, nil
}

var (
	NamespaceTokenConvertor = &namespaceTokenConvertorImpl{}
)

type namespaceTokenConvertorImpl struct{}

// ConvertPOToEntity 转换
func (c *namespaceTokenConvertorImpl) ConvertPOToEntity(p *po.NamespaceTokenPO) (*entity.NamespaceToken, error) {
	u := &entity.NamespaceToken{}
	if err := copier.Copy(u, p); err != nil {
		return nil, err
	}
	return u, nil
}

// ConvertPOsToEntities 转换
func (c *namespaceTokenConvertorImpl) ConvertPOsToEntities(p []*po.NamespaceTokenPO) ([]*entity.NamespaceToken, error) {
	namespaceTokens := make([]*entity.NamespaceToken, 0, len(p))
	for _, namespacePO := range p {
		namespaceToken, err := c.ConvertPOToEntity(namespacePO)
		if err != nil {
			return nil, err
		}
		namespaceTokens = append(namespaceTokens, namespaceToken)
	}
	return namespaceTokens, nil
}
