package convertor

import (
	"encoding/json"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/jinzhu/copier"
)

var (
	InstConvertor = &instConvertorImpl{} // 流程转换器
)

type instConvertorImpl struct {
}

// ConvertEntityToCreateRepoDTO 转换
func (*instConvertorImpl) ConvertEntityToCreateRepoDTO(e *entity.WorkflowInst) (*dto.CreateWorkflowInstRepoDTO, error) {
	d := &dto.CreateWorkflowInstRepoDTO{
		Namespace:  e.WorkflowDef.Namespace,
		DefID:      e.WorkflowDef.DefID,
		DefVersion: e.WorkflowDef.Version,
		Name:       e.Name,
		Creator:    e.Creator,
		Status:     e.Status.IntValue(),
		StartAt:    time.Now(),
	}
	e.CurNodeInst = nil
	e.SchedNodeInsts = nil
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	d.Context = string(b)
	return d, nil
}

// ConvertEntityToUpdateDTO 转换
func (*instConvertorImpl) ConvertEntityToUpdateDTO(oldInst *entity.WorkflowInst) (*dto.UpdateWorkflowInstDTO, error) {
	updateDTO := &dto.UpdateWorkflowInstDTO{
		Namespace:   oldInst.WorkflowDef.Namespace,
		DefID:       oldInst.WorkflowDef.DefID,
		InstID:      oldInst.InstID,
		Status:      oldInst.Status.IntValue(),
		CompletedAt: oldInst.CompletedAt,
	}
	newInst := &entity.WorkflowInst{}
	err := copier.Copy(newInst, oldInst)
	if err != nil {
		return nil, err
	}
	// 这两个字段太大了不更新到workflow_inst表里面, Get的时候从node_inst表去拿
	newInst.CurNodeInst = nil
	newInst.SchedNodeInsts = nil

	if newInst.Status != entity.InstFailed {
		newInst.Reason.FailedRootCause.FailedReason = ""
		newInst.Reason.FailedRootCause.FailedNodeRefNames = []string{}
	}

	if oldInst.Status.IsTerminal() {
		updateDTO.CompletedAt = time.Now()
	}

	b, err := json.Marshal(newInst)
	if err != nil {
		return nil, err
	}
	updateDTO.Context = string(b)
	return updateDTO, nil
}

// ConvertEntityToDTO 转换
func (*instConvertorImpl) ConvertEntityToDTO(e *entity.WorkflowInst) (*dto.WorkflowInstDTO, error) {
	d := &dto.WorkflowInstDTO{}
	err := copier.Copy(d, e)
	if err != nil {
		return nil, err
	}
	return d, nil
}
