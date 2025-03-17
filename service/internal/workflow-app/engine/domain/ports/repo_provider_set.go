package ports

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
)

// RepoProviderSet 仓储层集合
type RepoProviderSet struct {
	workflowDefRepo     WorkflowDefRepository
	workflowInstRepo    WorkflowInstRepository
	workflowArchiveRepo WorkflowArchiveRepository
	nodeInstRepo        NodeInstRepository
	eventBusRepo        EventBusRepository
	cacheRepo           CacheRepository
	remoteRepo          RemoteRepository
	triggerRepo         TriggerRepository
}

// WorkflowDefRepo 流程定义仓储层
func (r *RepoProviderSet) WorkflowDefRepo() WorkflowDefRepository {
	return r.workflowDefRepo
}

// WorkflowInstRepo 流程实例仓储层
func (r *RepoProviderSet) WorkflowInstRepo() WorkflowInstRepository {
	return r.workflowInstRepo
}

// WorkflowArchiveRepo 归档实例仓储层
func (r *RepoProviderSet) WorkflowArchiveRepo() WorkflowArchiveRepository {
	return r.workflowArchiveRepo
}

// NodeInstRepo 节点实例仓储层
func (r *RepoProviderSet) NodeInstRepo() NodeInstRepository {
	return r.nodeInstRepo
}

// EventBusRepo 事件仓储层
func (r *RepoProviderSet) EventBusRepo() EventBusRepository {
	return r.eventBusRepo
}

// CacheRepo 缓存仓储层
func (r *RepoProviderSet) CacheRepo() CacheRepository {
	return r.cacheRepo
}

// RemoteRepo 远程调用仓储层
func (r *RepoProviderSet) RemoteRepo() RemoteRepository {
	return r.remoteRepo
}

// TriggerRepo 触发器仓储层
func (r *RepoProviderSet) TriggerRepo() TriggerRepository {
	return r.triggerRepo
}

// NewRepoSet 实例化
func NewRepoSet(defRepo *repo.WorkflowDefRepo,
	instRepo *repo.WorkflowInstRepo,
	nodeInstRepo *repo.NodeInstRepo,
	eventBusRepo *repo.EventBusRepo,
	archiveRepo *repo.WorkflowArchiveRepo,
	cacheRepo *repo.CacheRepo,
	remoteRepo *repo.RemoteRepo,
	triggerRepo *repo.TriggerRepo) *RepoProviderSet {
	return &RepoProviderSet{
		workflowDefRepo:     defRepo,
		workflowInstRepo:    instRepo,
		nodeInstRepo:        nodeInstRepo,
		workflowArchiveRepo: archiveRepo,
		eventBusRepo:        eventBusRepo,
		cacheRepo:           cacheRepo,
		remoteRepo:          remoteRepo,
		triggerRepo:         triggerRepo,
	}
}
