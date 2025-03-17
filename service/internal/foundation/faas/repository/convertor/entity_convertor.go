package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"

	"github.com/jinzhu/copier"
)

var (
	FunctionConvertor = &functionConvertor{}
)

type functionConvertor struct{}

// ConvertPOToEntity 将 function 的 po 转为 entity
func (c *functionConvertor) ConvertPOToEntity(p *po.FunctionPO) (*entity.Function, error) {
	f := &entity.Function{}
	if err := copier.Copy(f, p); err != nil {
		return nil, err
	}
	f.Language = entity.GetLanguageTypeByStrValue(p.Language)
	return f, nil
}

// ConvertPOsToEntities 将 function 的 po list 转为 entity list
func (c *functionConvertor) ConvertPOsToEntities(p []*po.FunctionPO) ([]*entity.Function, error) {
	functions := make([]*entity.Function, 0, len(p))
	for _, functionPO := range p {
		function, err := c.ConvertPOToEntity(functionPO)
		if err != nil {
			return nil, err
		}
		functions = append(functions, function)
	}
	return functions, nil
}

// ConvertHistoryPOToEntity 将 function history 的 po 转为 entity
func (c *functionConvertor) ConvertHistoryPOToEntity(p *po.RunHistoryPO) (*entity.RunHistory, error) {
	f := &entity.RunHistory{}
	if err := copier.Copy(f, p); err != nil {
		return nil, err
	}
	return f, nil
}

// ConvertHistoryPOsToEntities 将 history 的 po list 转为 entity list
func (c *functionConvertor) ConvertHistoryPOsToEntities(p []*po.RunHistoryPO) ([]*entity.RunHistory, error) {
	histories := make([]*entity.RunHistory, 0, len(p))
	for _, historyPO := range p {
		history, err := c.ConvertHistoryPOToEntity(historyPO)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}
