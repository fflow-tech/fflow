package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/jinzhu/copier"
)

var (
	DefConvertor = &defConvertorImpl{} // 转换器
)

type defConvertorImpl struct {
}

// ConvertCreateDTOToPO  转换
func (*defConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	p := &po.WorkflowDefPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}
	p.ParentDefID, err = utils.StrToUInt64(d.ParentDefID)
	if err != nil {
		return nil, err
	}
	p.Status = entity.Disabled.IntValue()
	return p, nil
}

// ConvertGetDTOToPO  转换
func (*defConvertorImpl) ConvertGetDTOToPO(d *dto.GetWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	p := &po.WorkflowDefPO{}
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

// ConvertDeleteDTOToPO 转化
func (*defConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	p := &po.WorkflowDefPO{}
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

// ConvertUpdateDTOToPO 转化
func (*defConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateWorkflowDefDTO) *po.WorkflowDefPO {
	p := &po.WorkflowDefPO{
		Status: d.Status.IntValue(),
	}
	return p
}

// ConvertPageQueryDTOToPO 转化
func (*defConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	p := &po.WorkflowDefPO{
		Namespace: d.Namespace,
		Name:      d.Name,
		Version:   d.Version,
		Status:    d.Status,
	}
	var err error
	if p.DefID, err = utils.StrToUInt64(d.DefID); err != nil {
		return nil, err
	}

	return p, nil
}

// ConvertDTOToPO 转化
func (*defConvertorImpl) ConvertDTOToPO(d *dto.WorkflowDefDTO) (*po.WorkflowDefPO, error) {
	p := &po.WorkflowDefPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	p.DefID, err = utils.StrToUInt64(d.DefID)
	if err != nil {
		return nil, err
	}
	p.ParentDefID, err = utils.StrToUInt64(d.ParentDefID)
	if err != nil {
		return nil, err
	}
	return p, nil
}
