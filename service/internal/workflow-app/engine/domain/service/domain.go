// Package service 领域服务
package service

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/query"
)

// DomainService 领域服务
type DomainService struct {
	Commands ports.CommandPorts
	Queries  ports.QueryPorts
}

// NewDomainService 初始化领域服务
func NewDomainService(commands *command.Adapters, queries *query.Adapters) *DomainService {
	return &DomainService{
		Commands: commands,
		Queries:  queries,
	}
}
