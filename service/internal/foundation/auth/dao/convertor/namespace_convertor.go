package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	NamespaceConvertor = &namespaceConvertorImpl{} // 转换器
)

type namespaceConvertorImpl struct {
}

// ConvertDTOToPO 转换
func (*namespaceConvertorImpl) ConvertDTOToPO(req *dto.CreateNamespaceDTO) (*po.NamespacePO, error) {
	p := &po.NamespacePO{}
	if err := copier.Copy(p, req); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertGetDTOToPO 转换
func (c *namespaceConvertorImpl) ConvertGetDTOToPO(req *dto.GetNamespaceDTO) (*po.NamespacePO, error) {
	p := &po.NamespacePO{}
	if err := copier.Copy(p, req); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertDeleteDTOToPO 转换
func (c *namespaceConvertorImpl) ConvertDeleteDTOToPO(req *dto.DeleteNamespaceDTO) (*po.NamespacePO, error) {
	p := &po.NamespacePO{}
	if err := copier.Copy(p, req); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertUpdateDTOToPO 转换
func (c *namespaceConvertorImpl) ConvertUpdateDTOToPO(req *dto.UpdateNamespaceDTO) (*po.NamespacePO, error) {
	p := &po.NamespacePO{}
	if err := copier.Copy(p, req); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertPageQueryDTOToPO 转换
func (c *namespaceConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryNamespaceDTO) (*po.NamespacePO, error) {
	p := &po.NamespacePO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	p.Namespace = ""
	return p, nil
}
