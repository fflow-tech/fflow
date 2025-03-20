package service

import (
	"fmt"
	"os"
	"sync"

	"context"

	"github.com/fflow-tech/fflow/service/cmd/workflow-cli/factory"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// WorkflowService 提供工作流管理和执行的服务
type WorkflowService struct {
	definitionPath string
	instancePath   string
	mutex          sync.Mutex
}

// NewWorkflowService 创建工作流服务实例
func NewWorkflowService(definitionPath, instancePath string) (*WorkflowService, error) {
	// 确保目录存在
	if err := ensureDir(definitionPath); err != nil {
		return nil, fmt.Errorf("确保定义目录存在失败: %w", err)
	}
	if err := ensureDir(instancePath); err != nil {
		return nil, fmt.Errorf("确保实例目录存在失败: %w", err)
	}

	return &WorkflowService{
		definitionPath: definitionPath,
		instancePath:   instancePath,
	}, nil
}

// ExecuteWorkflow 执行指定ID的工作流
func (s *WorkflowService) ExecuteWorkflow(defJson string, input map[string]interface{}) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 1. 获取领域服务
	domainService, err := factory.GetDomainService()
	if err != nil {
		return "", fmt.Errorf("failed to get domain service: %w", err)
	}

	// 2. 保存工作流定义到数据库
	defId, err := domainService.Commands.CreateWorkflowDef(context.Background(), &dto.CreateWorkflowDefDTO{
		DefJson: defJson,
	})
	if err != nil {
		return "", fmt.Errorf("Failed to save workflow definition: %w, defJson: %s", err, defJson)
	}
	log.Infof("Workflow definition saved: %s", defId)

	err = domainService.Commands.EnableWorkflowDef(context.Background(), &dto.EnableWorkflowDefDTO{
		DefID: defId,
	})
	if err != nil {
		return "", fmt.Errorf("Failed to enable workflow definition: %w", err)
	}
	log.Infof("Workflow definition enabled: %s", defId)

	// 4. 创建工作流实例
	instId, err := domainService.Commands.StartWorkflowInst(context.Background(), &dto.StartWorkflowInstDTO{
		DefID: defId,
		Input: input,
	})
	if err != nil {
		return "", fmt.Errorf("Failed to create workflow instance: %w", err)
	}
	log.Infof("Workflow instance started: %s", instId)

	return instId, nil
}

// GetWorkflowStatus 获取工作流状态
func (s *WorkflowService) GetWorkflowStatus(id string) (*dto.WorkflowInstDTO, error) {
	// 获取领域服务
	domainService, err := factory.GetDomainService()
	if err != nil {
		return nil, fmt.Errorf("Failed to get domain service: %w", err)
	}

	// 查询工作流实例
	return domainService.Queries.GetWorkflowInst(context.Background(), &dto.GetWorkflowInstDTO{
		InstID: id,
	})
}

// 确保目录存在，如果不存在则创建
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
