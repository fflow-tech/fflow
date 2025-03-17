package ports

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/repository/repo"
)

// RepoProviderSet 仓储层集合
type RepoProviderSet struct {
	functionRepo FunctionRepository
}

// FunctionRepo 函数仓储层
func (r *RepoProviderSet) FunctionRepo() FunctionRepository {
	return r.functionRepo
}

// NewRepoSet 实例化
func NewRepoSet(functionRepo *repo.FunctionRepo) *RepoProviderSet {
	return &RepoProviderSet{
		functionRepo: functionRepo,
	}
}
