package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	RunHistoryConvertor = &runHistoryConvertorImpl{} // 转换器
)

type runHistoryConvertorImpl struct {
}

// ConvertCreateDTOToPO  转换
func (*runHistoryConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateRunHistoryDTO) *po.RunHistoryPO {
	p := &po.RunHistoryPO{}
	copier.Copy(p, d)
	p.Name = d.FunctionName
	return p
}

// ConvertGetDTOToPO  转换
func (*runHistoryConvertorImpl) ConvertGetDTOToPO(d *dto.GetRunHistoryDTO) *po.RunHistoryPO {
	p := &po.RunHistoryPO{}
	copier.Copy(p, d)
	return p
}

// ConvertDeleteDTOToPO  转换
func (*runHistoryConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteRunHistoryDTO) *po.RunHistoryPO {
	p := &po.RunHistoryPO{}
	copier.Copy(p, d)
	return p
}

// ConvertBatchDeleteDTOToPO  转换
func (*runHistoryConvertorImpl) ConvertBatchDeleteDTOToPO(d *dto.BatchDeleteRunHistoryDTO) *po.RunHistoryPO {
	p := po.RunHistoryPO{}
	copier.Copy(&p, d)
	return &p
}

// ConvertUpdateDTOToPO  转换
func (*runHistoryConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateRunHistoryDTO) *po.RunHistoryPO {
	p := po.RunHistoryPO{}
	copier.Copy(&p, d)
	return &p
}

// ConvertPageQueryDTOToPO 转换
func (*runHistoryConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryRunHistoryDTO) *po.RunHistoryPO {
	return &po.RunHistoryPO{
		Namespace: d.Namespace,
		Version:   d.Version,
		Name:      d.FunctionName,
	}
}
