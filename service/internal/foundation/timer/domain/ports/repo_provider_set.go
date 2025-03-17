package ports

import "github.com/fflow-tech/fflow/service/internal/foundation/timer/repository/repo"

// RepoProviderSet 仓储层集合
type RepoProviderSet struct {
	timerDefDefRepo TimerDefRepository
	timerTaskRepo   TimerTaskRepository
	pollingTaskRepo PollingTaskRepository
	eventBusRepo    EventBusRepository
	remoteRepo      RemoteRepository
	appRepo         AppRepository
}

// TimerDefRepo 定时器定义repo
func (r *RepoProviderSet) TimerDefRepo() TimerDefRepository {
	return r.timerDefDefRepo
}

// TimerTaskRepo 定时器任务repo
func (r *RepoProviderSet) TimerTaskRepo() TimerTaskRepository {
	return r.timerTaskRepo
}

// PollingTaskRepo 轮询任务repo
func (r *RepoProviderSet) PollingTaskRepo() PollingTaskRepository {
	return r.pollingTaskRepo
}

// EventBusRepo 消息中间repo
func (r *RepoProviderSet) EventBusRepo() EventBusRepository {
	return r.eventBusRepo
}

// RemoteRepo 远程 repo
func (r *RepoProviderSet) RemoteRepo() RemoteRepository {
	return r.remoteRepo
}

// AppRepo 应用仓储实体
func (r *RepoProviderSet) AppRepo() AppRepository {
	return r.appRepo
}

// NewRepoSet 实例化
func NewRepoSet(defRepo *repo.TimerDefRepo,
	taskRepo *repo.TimerTaskRepo,
	pollingTaskRepo *repo.PollingTaskRepo,
	eventBusRepo *repo.EventBusRepo,
	remoteRepo *repo.RemoteRepo,
	appRepo *repo.AppRepo) *RepoProviderSet {
	return &RepoProviderSet{
		timerDefDefRepo: defRepo,
		timerTaskRepo:   taskRepo,
		pollingTaskRepo: pollingTaskRepo,
		eventBusRepo:    eventBusRepo,
		remoteRepo:      remoteRepo,
		appRepo:         appRepo,
	}
}
