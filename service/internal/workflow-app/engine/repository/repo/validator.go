package repo

import (
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
)

// ValidateWorkflowInstCtxSize 检测流程实例上下文的大小
func ValidateWorkflowInstCtxSize(ctxInfo string) error {
	workflowInstCtxSize := config.GetValidationRulesConfig().WorkflowInstCtxSize
	ctxSize := len([]byte(ctxInfo))
	if ctxSize > workflowInstCtxSize {
		return fmt.Errorf("failed to ValidateWorkflowInstCtxSize,caused by: must not more than %d",
			workflowInstCtxSize)
	}
	return nil
}

// ValidateNodeInstCtxSize 检测节点实例上下文大小
func ValidateNodeInstCtxSize(ctxInfo string) error {
	nodeInstCtxSize := config.GetValidationRulesConfig().NodeInstCtxSize
	ctxSize := len([]byte(ctxInfo))
	if ctxSize > nodeInstCtxSize {
		return fmt.Errorf("failed to ValidateNodeInstCtxSize,caused by: must not more than %d", nodeInstCtxSize)
	}
	return nil
}
