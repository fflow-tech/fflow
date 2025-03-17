package ports

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/repository/repo"
)

// RepoProviderSet 仓储层集合
type RepoProviderSet struct {
	userRepo UserRepository
}

// UserRepository 仓储层
func (r *RepoProviderSet) UserRepository() UserRepository {
	return r.userRepo
}

// NewRepoSet 实例化
func NewRepoSet(userRepo *repo.UserRepo) *RepoProviderSet {
	return &RepoProviderSet{
		userRepo: userRepo,
	}
}
