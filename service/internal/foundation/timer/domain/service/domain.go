// Package service 领域服务层，提供 app 层增删改查能力。
package service

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service/query"
)

// DomainService 领域服务
type DomainService struct {
	Commands ports.CommandPorts
	Queries  ports.QueryPorts
}

// NewDomainService 创建领域服务
func NewDomainService(commands *command.Adapters, queries *query.Adapters) *DomainService {
	return &DomainService{
		Commands: commands,
		Queries:  queries,
	}
}
