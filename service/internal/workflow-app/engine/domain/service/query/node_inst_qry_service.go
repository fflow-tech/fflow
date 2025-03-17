package query

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// NodeInstQueryService 查询服务
type NodeInstQueryService struct {
	nodeInstRepo ports.NodeInstRepository
}

// NewNodeInstQueryService 新建查询服务
func NewNodeInstQueryService(r *repo.NodeInstRepo) *NodeInstQueryService {
	return &NodeInstQueryService{nodeInstRepo: r}
}

// GetNodeInstDetail 查询节点详情
func (m *NodeInstQueryService) GetNodeInstDetail(ctx context.Context,
	req *dto.GetNodeInstDTO) (*dto.NodeInstDTO, error) {
	e, err := m.nodeInstRepo.Get(req)
	if err != nil {
		log.Errorf("Failed to get node inst detail, caused by %s", err)
		return nil, err
	}

	return convertor.NodeInstConvertor.ConvertEntityToDTO(e)
}
