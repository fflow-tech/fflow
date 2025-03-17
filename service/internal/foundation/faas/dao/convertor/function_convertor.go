package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	FunctionConvertor = &functionConvertorImpl{} // 转换器
)

type functionConvertorImpl struct {
}

// ConvertCreateDTOToPO  转换
func (*functionConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateFunctionDTO) *po.FunctionPO {
	p := &po.FunctionPO{}
	copier.Copy(p, d)
	p.Name = d.Function
	p.Language = d.Language.String()
	return p
}

// ConvertGetDTOToPO  转换
func (*functionConvertorImpl) ConvertGetDTOToPO(d *dto.GetFunctionReqDTO) *po.FunctionPO {
	p := &po.FunctionPO{}
	copier.Copy(p, d)
	p.Name = d.Function
	return p
}

// ConvertDeleteDTOToPO  转换
func (*functionConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteFunctionDTO) *po.FunctionPO {
	p := &po.FunctionPO{}
	copier.Copy(p, d)
	p.Name = d.Function
	return p
}

// ConvertPageQueryDTOToPO 转换
func (*functionConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryFunctionDTO) *po.FunctionPO {
	p := &po.FunctionPO{
		Namespace: d.Namespace,
		Creator:   d.Creator,
		Version:   d.Version,
		Name:      d.Function,
		Language:  d.Language,
	}
	return p
}
