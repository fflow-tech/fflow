package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/jinzhu/copier"
)

var (
	InstConvertor = &instConvertorImpl{} // 转换器
)

type instConvertorImpl struct {
}

// ConvertCreateDTOToPO  转换
func (*instConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateWorkflowInstRepoDTO) (*po.WorkflowInstPO, error) {
	p := &po.WorkflowInstPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertGetDTOToPO  转换
func (*instConvertorImpl) ConvertGetDTOToPO(d *dto.GetWorkflowInstDTO) (*po.WorkflowInstPO, error) {
	p := &po.WorkflowInstPO{}
	instID, err := utils.StrToUInt(d.InstID)
	if err != nil {
		return nil, err
	}
	p.ID = instID
	p.Creator = d.Operator

	defID, err := utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}
	p.DefID = defID

	return p, nil
}

// ConvertDeleteDTOToPO 转化
func (*instConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteWorkflowInstsDTO) (*po.WorkflowInstPO, error) {
	p := &po.WorkflowInstPO{}
	defID, err := utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}
	p.DefID = defID

	instID, err := utils.StrToUInt(d.InstID)
	if err != nil {
		return nil, err
	}
	p.ID = instID
	return p, nil
}

// ConvertUpdateDTOToPO 转化
func (*instConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateWorkflowInstDTO) *po.WorkflowInstPO {
	p := &po.WorkflowInstPO{
		Context:     d.Context,
		Status:      d.Status,
		CompletedAt: d.CompletedAt,
	}
	return p
}

// ConvertPageQueryDTOToPO 转化
func (*instConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryWorkflowInstDTO) (*po.WorkflowInstPO, error) {
	p := &po.WorkflowInstPO{}
	defID, err := utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}
	p.Namespace = d.Namespace
	p.DefID = defID
	p.Status = d.Status.IntValue()
	p.Creator = d.Creator

	return p, nil
}
