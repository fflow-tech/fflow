package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"

	"github.com/jinzhu/copier"
)

var (
	// RunHistoryConvertor 运行记录转换体
	RunHistoryConvertor = &functionConvertor{}
)

type functionConvertor struct {
}

// ConvertEntityToGetDTO 运行记录 entity->dto
func (c *functionConvertor) ConvertEntityToGetDTO(e *entity.RunHistory) (*dto.GetRunHistoryRspDTO, error) {
	d := &dto.GetRunHistoryRspDTO{}
	if err := copier.Copy(d, e); err != nil {
		return nil, err
	}

	return d, nil
}

// ConvertEntitiesToDTOs 运行记录 entities->dtos
func (c *functionConvertor) ConvertEntitiesToDTOs(e []*entity.RunHistory) ([]*dto.GetRunHistoryRspDTO, error) {
	historyRsp := make([]*dto.GetRunHistoryRspDTO, 0, len(e))
	for _, runHistory := range e {
		history, err := c.ConvertEntityToGetDTO(runHistory)
		if err != nil {
			return nil, err
		}

		historyRsp = append(historyRsp, history)
	}
	return historyRsp, nil
}
