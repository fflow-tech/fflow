package convertor

import (
	"encoding/json"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/jinzhu/copier"
)

var (
	NodeInstConvertor = &nodeInstConvertorImpl{} // 转换器
)

type nodeInstConvertorImpl struct {
}

// ConvertEntityToCreateDTO 转换
func (*nodeInstConvertorImpl) ConvertEntityToCreateDTO(e *entity.NodeInst) (*dto.CreateNodeInstDTO, error) {
	d := &dto.CreateNodeInstDTO{
		Namespace:   e.Namespace,
		DefID:       e.DefID,
		DefVersion:  e.DefVersion,
		InstID:      e.InstID,
		RefName:     e.BasicNodeDef.RefName,
		Status:      e.Status.IntValue(),
		ScheduledAt: e.ScheduledAt,
	}

	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	d.Context = string(b)
	return d, nil
}

// ConvertEntityToDTO 转换
func (*nodeInstConvertorImpl) ConvertEntityToDTO(e *entity.NodeInst) (*dto.NodeInstDTO, error) {
	return &dto.NodeInstDTO{
		NodeInst: *e,
	}, nil
}

// ConvertEntityToUpdateDTO 转换
func (*nodeInstConvertorImpl) ConvertEntityToUpdateDTO(e *entity.NodeInst) (*dto.UpdateNodeInstDTO, error) {
	d := &dto.UpdateNodeInstDTO{}

	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	if err = copier.Copy(d, e); err != nil {
		return nil, err
	}

	d.NodeInstID = e.NodeInstID
	d.Context = string(b)
	d.Status = e.Status.IntValue()
	return d, nil
}
