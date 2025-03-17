package convertor

import (
	"encoding/json"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/jinzhu/copier"
)

var (
	DefConvertor = &defConvertorImpl{} // 流程定义转换器
)

type defConvertorImpl struct {
}

// ConvertEntityToDTO 转换
func (c *defConvertorImpl) ConvertEntityToDTO(e *entity.WorkflowDef) (*dto.WorkflowDefDTO, error) {
	d := &dto.WorkflowDefDTO{}
	if err := copier.Copy(d, e); err != nil {
		return nil, err
	}
	d.DefID = e.DefID
	d.Description = e.Desc
	d.Attribute.ParentDefVersion = e.ParentDefVersion
	d.Attribute.RefName = e.RefName
	defJson, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	d.DefJson = string(defJson)
	return d, nil
}

// ConvertEntitiesToDTOs 批量转换
func (c *defConvertorImpl) ConvertEntitiesToDTOs(e []*entity.WorkflowDef) ([]*dto.WorkflowDefDTO, error) {
	defs := make([]*dto.WorkflowDefDTO, 0, len(e))
	for _, defEntity := range e {
		def, err := c.ConvertEntityToDTO(defEntity)
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}

	return defs, nil
}
