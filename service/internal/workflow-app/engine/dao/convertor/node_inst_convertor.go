package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/jinzhu/copier"
)

var (
	NodeInstConvertor = &nodeInstConvertorImpl{} // 转换器
)

type nodeInstConvertorImpl struct {
}

// ConvertCreateDTOToPO  转换
func (*nodeInstConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateNodeInstDTO) (*po.NodeInstPO, error) {
	p := &po.NodeInstPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
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

// ConvertGetDTOToPO  转换
func (*nodeInstConvertorImpl) ConvertGetDTOToPO(d *dto.GetNodeInstDTO) (*po.NodeInstPO, error) {
	p := &po.NodeInstPO{}
	p.Namespace = d.Namespace
	var err error
	p.ID, err = utils.StrToUInt(d.NodeInstID)
	if err != nil {
		return nil, err
	}

	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}

	p.DefVersion = d.DefVersion

	p.InstID, err = utils.StrToUInt64(d.InstID)
	if err != nil {
		return nil, err
	}

	p.RefName = d.RefName
	p.Status = d.Status
	return p, nil
}

// ConvertDeleteDTOToPO 转换
func (*nodeInstConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteNodeInstDTO) (*po.NodeInstPO, error) {
	p := &po.NodeInstPO{}

	var err error
	p.ID, err = utils.StrToUInt(d.NodeInstID)
	if err != nil {
		return nil, err
	}
	p.DefID, err = utils.StrToUInt64(d.DefID)
	return p, err
}

// ConvertUpdateDTOToPO 转换
func (*nodeInstConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateNodeInstDTO) *po.NodeInstPO {
	return &po.NodeInstPO{
		Context:       d.Context,
		Status:        d.Status,
		WaitAt:        d.WaitAt,
		ExecuteAt:     d.ExecuteAt,
		AsynWaitResAt: d.AsynWaitResAt,
		CompletedAt:   d.CompletedAt,
	}
}

// ConvertPageQueryDTOToPO 转换
func (*nodeInstConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryNodeInstDTO) (*po.NodeInstPO, error) {
	p := &po.NodeInstPO{
		DefVersion: d.DefVersion,
		RefName:    d.RefName,
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
