package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"

	"github.com/jinzhu/copier"
)

var (
	FunctionConvertor = &functionConvertor{} // 转换器
)

type functionConvertor struct {
}

// ConvertEntityToGetDTO 转换
func (c *functionConvertor) ConvertEntityToGetDTO(e *entity.Function) (*dto.GetFunctionRspDTO, error) {
	d := &dto.GetFunctionRspDTO{}
	if err := copier.Copy(d, e); err != nil {
		return nil, err
	}
	d.Function = e.Name
	d.Language = e.Language.String()
	return d, nil
}

// ConvertEntitiesToDTOs 转换
func (c *functionConvertor) ConvertEntitiesToDTOs(es []*entity.Function) ([]*dto.GetFunctionRspDTO, error) {
	functions := make([]*dto.GetFunctionRspDTO, 0, len(es))
	for _, e := range es {
		function, err := c.ConvertEntityToGetDTO(e)
		if err != nil {
			return nil, err
		}

		functions = append(functions, function)
	}
	return functions, nil
}

// ConvertHistoryEntityToGetDTO 转换
func (c *functionConvertor) ConvertHistoryEntityToGetDTO(e *entity.RunHistory) (*dto.GetRunHistoryRspDTO, error) {
	d := &dto.GetRunHistoryRspDTO{}
	if err := copier.Copy(d, e); err != nil {
		return nil, err
	}
	d.FunctionName = e.Name
	return d, nil
}

// ConvertHistoryEntitiesToDTOs 转换
func (c *functionConvertor) ConvertHistoryEntitiesToDTOs(es []*entity.RunHistory) ([]*dto.GetRunHistoryRspDTO, error) {
	histories := make([]*dto.GetRunHistoryRspDTO, 0, len(es))
	for _, e := range es {
		history, err := c.ConvertHistoryEntityToGetDTO(e)
		if err != nil {
			return nil, err
		}

		histories = append(histories, history)
	}
	return histories, nil
}
