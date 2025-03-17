package service

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/query"
)

// DomainService 领域服务
type DomainService struct {
	Commands ports.CommandPorts
	Queries  ports.QueryPorts
}

// NewDomainService 新建领域服务
func NewDomainService(commands *command.Adapters, queries *query.Adapters) *DomainService {
	return &DomainService{
		Commands: commands,
		Queries:  queries,
	}
}
