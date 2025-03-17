package dto

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
)

// HistoryWorkflowInstDTO 历史流程定义
type HistoryWorkflowInstDTO struct {
	po.HistoryWorkflowInstPO
}

// ArchiveHistoryWorkflowInstsDTO 归档历史流程请求
type ArchiveHistoryWorkflowInstsDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
}
