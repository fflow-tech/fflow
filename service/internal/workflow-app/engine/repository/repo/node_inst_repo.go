package repo

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/sql"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/repository/convertor"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// NodeInstRepo 节点实例仓储层
type NodeInstRepo struct {
	nodeInstDAO storage.NodeInstDAO
}

// NewNodeInstRepo 实体构造函数
func NewNodeInstRepo(d *sql.NodeInstDAO) *NodeInstRepo {
	return &NodeInstRepo{nodeInstDAO: d}
}

// Create 创建
func (t *NodeInstRepo) Create(d *dto.CreateNodeInstDTO) (string, error) {
	if err := ValidateNodeInstCtxSize(d.Context); err != nil {
		return "", err
	}
	r, err := t.nodeInstDAO.Create(d)
	if err != nil {
		return "", err
	}

	return utils.UintToStr(r.ID), err
}

// UpdateWithDefID 更新
func (t *NodeInstRepo) UpdateWithDefID(d *dto.UpdateNodeInstDTO) error {
	if err := ValidateNodeInstCtxSize(d.Context); err != nil {
		return err
	}
	return t.nodeInstDAO.Update(d)
}

// Get 查询
func (t *NodeInstRepo) Get(d *dto.GetNodeInstDTO) (*entity.NodeInst, error) {
	p, err := t.nodeInstDAO.Get(d)
	if err != nil {
		return &entity.NodeInst{}, err
	}
	return convertor.NodeInstConvertor.ConvertPOToEntity(p)
}

// PageQuery 分页查询
func (t *NodeInstRepo) PageQuery(req *dto.PageQueryNodeInstDTO) ([]*entity.NodeInst, error) {
	if req.PageQuery == nil {
		req.PageQuery = constants.NewDefaultPageQuery()
	}

	nodeInstPOs, err := t.nodeInstDAO.PageQuery(req)
	if err != nil {
		return nil, err
	}
	var nodeInsts []*entity.NodeInst
	for _, nodeInstPO := range nodeInstPOs {
		nodeInst, err := convertor.NodeInstConvertor.ConvertPOToEntity(nodeInstPO)
		if err != nil {
			return nil, err
		}
		nodeInsts = append(nodeInsts, nodeInst)
	}

	return nodeInsts, nil
}
