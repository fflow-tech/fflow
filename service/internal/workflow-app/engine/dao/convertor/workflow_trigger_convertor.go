package convertor

import (
	"encoding/json"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	TriggerConvertor = &triggerConvertorImpl{} // 转换器
)

type triggerConvertorImpl struct {
}

// ConvertCreateDTOToPO  转换
func (*triggerConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateTriggerDTO) (*po.TriggerPO, error) {
	p := &po.TriggerPO{
		Type:       d.Type.String(),
		Event:      d.Event,
		Expr:       d.Expr,
		Level:      d.Level.IntValue(),
		DefVersion: d.DefVersion,
		Status:     d.Status.IntValue(),
	}

	var err error
	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}

	p.InstID, err = utils.StrToUInt64(d.InstID)
	if err != nil {
		return nil, err
	}

	triggerBytes, err := json.Marshal(d.Trigger)
	p.Attribute = string(triggerBytes)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ConvertGetDTOToPO  转换
func (*triggerConvertorImpl) ConvertGetDTOToPO(d *dto.GetTriggerDTO) (*po.TriggerPO, error) {
	p := &po.TriggerPO{
		Event:  d.Event,
		Type:   d.Type.String(),
		Level:  d.Level.IntValue(),
		Status: d.Status.IntValue(),
	}
	var err error
	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}

	p.InstID, err = utils.StrToUInt64(d.InstID)
	if err != nil {
		return nil, err
	}

	p.ID, err = utils.StrToUInt(d.TriggerID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ConvertDeleteDTOToPO 转换
func (*triggerConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteTriggerDTO) (*po.TriggerPO, error) {
	p := &po.TriggerPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}
	p.ID, err = utils.StrToUInt(d.TriggerID)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertUpdateDTOToPO 转换
func (*triggerConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateTriggerDTO) *po.TriggerPO {
	p := &po.TriggerPO{
		Status: d.Status.IntValue(),
	}
	return p
}

// ConvertPageQueryDTOToPO 转化
func (*triggerConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryTriggerDTO) (*po.TriggerPO, error) {
	p := &po.TriggerPO{
		Type:       d.Type.String(),
		Event:      d.Event,
		Level:      d.Level.IntValue(),
		DefVersion: d.DefVersion,
		Status:     d.Status.IntValue(),
	}
	var err error
	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}

	p.InstID, err = utils.StrToUInt64(d.InstID)
	if err != nil {
		return nil, err
	}

	return p, nil
}
