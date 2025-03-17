package sql

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/seq"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"gorm.io/gorm"
)

// NodeInstDAO NodeInst数据访问对象
type NodeInstDAO struct {
	db *mysql.Client
}

// NewNodeInstDAO NodeInst数据访问对象构造函数
func NewNodeInstDAO(db *mysql.Client) *NodeInstDAO {
	return &NodeInstDAO{db: db}
}

// Transaction 事务
func (dao *NodeInstDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建节点实例
func (dao *NodeInstDAO) Create(def *dto.CreateNodeInstDTO) (*po.NodeInstPO, error) {
	var err error
	id, err := seq.NewUint()
	if err != nil {
		return nil, err
	}

	p, err := convertor.NodeInstConvertor.ConvertCreateDTOToPO(def)
	if err != nil {
		return nil, err
	}
	p.ID = id
	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create node inst, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Get 获取节点实例信息
func (dao *NodeInstDAO) Get(d *dto.GetNodeInstDTO) (*po.NodeInstPO, error) {
	if utils.IsZero(d.NodeInstID) {
		return nil, fmt.Errorf("get node inst `NodeInstID` must not be zero, NodeInstID=[%s]", d.NodeInstID)
	}

	p, err := convertor.NodeInstConvertor.ConvertGetDTOToPO(d)
	if err != nil {
		return nil, err
	}
	r := &po.NodeInstPO{}
	if err := dao.db.ReadFromSlave(false).Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get node inst, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// GetNodeInstsByIDs 根据 ID 列表获取节点实例信息
func (dao *NodeInstDAO) GetNodeInstsByIDs(req *dto.GetNodeInstsByIDsDTO) ([]*po.NodeInstPO, error) {
	if utils.IsZero(req.DefID) || utils.IsZero(req.InstID) {
		return nil, fmt.Errorf("get inst `DefID` and `InstID` must not be zero, DefID:[%s], InstID:[%s]",
			req.DefID, req.InstID)
	}

	var r []*po.NodeInstPO
	if err := dao.db.ReadFromSlave(false).
		Where("def_id = ? and inst_id = ? and id in (?)", req.DefID, req.InstID, req.NodeInstIDs).
		Find(&r).Error; err != nil {
		log.Errorf("Failed to get node insts by ids, caused by %v", err)
		return nil, err
	}
	return r, nil
}

// DeleteNodeInstsByIDs 删除节点实例信息，硬删除
func (dao *NodeInstDAO) DeleteNodeInstsByIDs(req *dto.DeleteNodeInstsByIDsDTO) error {
	if utils.IsZero(req.DefID) || utils.IsZero(req.InstID) {
		return fmt.Errorf("delete inst `DefID` and `InstID` must not be zero, DefID:[%s], InstID:[%s]",
			req.DefID, req.InstID)
	}

	if err := dao.db.Where("def_id = ? and inst_id = ? and id in (?)", req.DefID, req.InstID, req.NodeInstIDs).
		Unscoped().Delete(&po.NodeInstPO{}).Error; err != nil {
		log.Errorf("Failed to delete node insts by ids, caused by %v", err)
		return err
	}

	return nil
}

// Delete 根据删除节点实例信息
func (dao *NodeInstDAO) Delete(d *dto.DeleteNodeInstDTO) error {
	if utils.IsZero(d.NodeInstID) || utils.IsZero(d.DefID) {
		return fmt.Errorf("delete node inst `DefID` or `NodeInstID` must not be zero, DefID:[%s] NodeInstID=[%s]",
			d.DefID, d.NodeInstID)
	}

	p, err := convertor.NodeInstConvertor.ConvertDeleteDTOToPO(d)
	if err != nil {
		return err
	}
	if err := dao.db.Where(p).Delete(&po.NodeInstPO{}).Error; err != nil {
		log.Errorf("Failed to delete node inst, caused by %s", err)
		return err
	}

	return nil
}

// Update 更新节点实例信息
func (dao *NodeInstDAO) Update(d *dto.UpdateNodeInstDTO) error {
	if utils.IsZero(d.NodeInstID) || utils.IsZero(d.DefID) {
		return fmt.Errorf("update node inst `DefID` or `NodeInstID` must not be zero,"+
			" DefID:[%s] NodeInstID=[%s]",
			d.DefID, d.NodeInstID)
	}

	p := convertor.NodeInstConvertor.ConvertUpdateDTOToPO(d)
	if err := dao.db.Where("id=? and def_id=?", d.NodeInstID, d.DefID).Updates(p).Error; err != nil {
		log.Errorf("Failed to update node inst, caused by %v", err)
		return err
	}

	return nil
}

// PageQuery 分页查询节点实例的信息
func (dao *NodeInstDAO) PageQuery(d *dto.PageQueryNodeInstDTO) ([]*po.NodeInstPO, error) {
	if utils.IsZero(d.DefID) {
		return nil, fmt.Errorf("PageQueryLastVersion node inst `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	var nodeInsts []*po.NodeInstPO
	p, err := convertor.NodeInstConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return nil, err
	}
	db := dao.db.ReadFromSlave(d.ReadFromSlave).Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit())

	if len(d.Statuses) > 0 {
		db.Where("status in (?)", getStatusesIntValue(d.Statuses))
	}
	if err := db.Find(&nodeInsts).Error; err != nil {
		log.Errorf("Failed to page query node inst, caused by %s", err)
		return nil, err
	}

	return nodeInsts, nil
}

func getStatusesIntValue(statuses []entity.NodeInstStatus) []int {
	var r []int
	for _, status := range statuses {
		r = append(r, status.IntValue())
	}
	return r
}

// Count 根据条件获取总数
func (dao *NodeInstDAO) Count(d *dto.PageQueryNodeInstDTO) (int64, error) {
	if utils.IsZero(d.DefID) {
		return 0, fmt.Errorf("count node inst `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	var totalCount int64
	p, err := convertor.NodeInstConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return 0, err
	}
	db := dao.db.ReadFromSlave(d.ReadFromSlave).Model(po.NodeInstPO{}).Where(p)
	if len(d.Statuses) > 0 {
		db.Where("status in (?)", getStatusesIntValue(d.Statuses))
	}
	if err := db.Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to get node inst count, caused by %s", err)
		return 0, err
	}

	return totalCount, nil
}
