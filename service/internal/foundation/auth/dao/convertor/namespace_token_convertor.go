package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	NamespaceTokenConvertor = &namespaceTokenConvertorImpl{} // 转换器
)

type namespaceTokenConvertorImpl struct {
}

// ConvertDTOToPO 转换
func (*namespaceTokenConvertorImpl) ConvertDTOToPO(req *dto.CreateNamespaceTokenDTO) (*po.NamespaceTokenPO, error) {
	p := &po.NamespaceTokenPO{}
	var err error
	if err = copier.Copy(p, req); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertGetDTOToPO 转换
func (c *namespaceTokenConvertorImpl) ConvertGetDTOToPO(d *dto.GetNamespaceTokenDTO) (*po.NamespaceTokenPO, error) {
	p := &po.NamespaceTokenPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}

	return p, nil
}

// ConvertDeleteDTOToPO 转换
func (c *namespaceTokenConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteNamespaceTokenDTO) (*po.NamespaceTokenPO, error) {
	p := &po.NamespaceTokenPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertUpdateDTOToPO 转换
func (c *namespaceTokenConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateNamespaceTokenDTO) (*po.NamespaceTokenPO, error) {
	p := &po.NamespaceTokenPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertPageQueryDTOToPO 转换
func (c *namespaceTokenConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryNamespaceTokenDTO) (*po.NamespaceTokenPO, error) {
	p := &po.NamespaceTokenPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}
