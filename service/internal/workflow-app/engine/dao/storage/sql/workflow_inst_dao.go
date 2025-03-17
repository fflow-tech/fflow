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

// WorkflowInstDAO WorkflowInst数据访问对象
type WorkflowInstDAO struct {
	db *mysql.Client
}

// NewWorkflowInstDAO WorkflowInst数据访问对象构造函数
func NewWorkflowInstDAO(db *mysql.Client) *WorkflowInstDAO {
	return &WorkflowInstDAO{db: db}
}

// Transaction 事务
func (dao *WorkflowInstDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建流程实例
func (dao *WorkflowInstDAO) Create(def *dto.CreateWorkflowInstRepoDTO) (*po.WorkflowInstPO, error) {
	id, err := seq.NewUint()
	if err != nil {
		return nil, err
	}

	p, err := convertor.InstConvertor.ConvertCreateDTOToPO(def)
	if err != nil {
		return nil, err
	}
	p.ID = id
	if err := dao.db.Create(p).Error; err != nil {
		log.Errorf("Failed to create workflow inst, caused by %s", err)
		return nil, err
	}
	return p, nil
}

// Get 获取流程实例信息
func (dao *WorkflowInstDAO) Get(d *dto.GetWorkflowInstDTO) (*po.WorkflowInstPO, error) {
	if utils.IsZero(d.InstID) {
		return nil,
			fmt.Errorf("Get inst `InstID` must not be zero, DefID:[%d] InstID:[%d] ", d.DefID, d.InstID)
	}

	p, err := convertor.InstConvertor.ConvertGetDTOToPO(d)
	if err != nil {
		return nil, err
	}
	r := &po.WorkflowInstPO{}
	if err := dao.db.ReadFromSlave(false).Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get workflow inst, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// GetWorkflowInstsByIDs 根据 ID 列表获取流程实例信息
func (dao *WorkflowInstDAO) GetWorkflowInstsByIDs(d *dto.GetWorkflowInstsByIDsDTO) ([]*po.WorkflowInstPO, error) {
	if utils.IsZero(d.DefID) {
		return nil, fmt.Errorf("get inst `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	var r []*po.WorkflowInstPO
	if err := dao.db.ReadFromSlave(false).Where("def_id = ? and id in (?)", d.DefID, d.InstIDs).
		Find(&r).Error; err != nil {
		log.Errorf("Failed to get workflow insts by ids, caused by %v", err)
		return nil, err
	}
	return r, nil
}

// Delete 删除流程实例信息
func (dao *WorkflowInstDAO) Delete(d *dto.DeleteWorkflowInstsDTO) error {
	if utils.IsZero(d.InstID) || utils.IsZero(d.DefID) {
		return fmt.Errorf("delete inst `DefID` or `InstID` must not be zero, DefID:[%s] InstID=[%s]",
			d.DefID, d.InstID)
	}

	p, err := convertor.InstConvertor.ConvertDeleteDTOToPO(d)
	if err != nil {
		return err
	}
	if err := dao.db.Where(p).Delete(&po.WorkflowInstPO{}).Error; err != nil {
		log.Errorf("Failed to delete workflow inst, caused by %s", err)
		return err
	}

	return nil
}

// DeleteWorkflowInstsByIDs 删除流程实例信息，硬删除
func (dao *WorkflowInstDAO) DeleteWorkflowInstsByIDs(d *dto.DeleteWorkflowInstsByIDsDTO) error {
	if utils.IsZero(d.DefID) {
		return fmt.Errorf("delete inst `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	if err := dao.db.Where("def_id = ? and id in (?)", d.DefID, d.InstIDs).
		Unscoped().Delete(&po.WorkflowInstPO{}).Error; err != nil {
		log.Errorf("Failed to get workflow insts by ids, caused by %v", err)
		return err
	}
	return nil
}

// Update 更新流程实例信息
func (dao *WorkflowInstDAO) Update(d *dto.UpdateWorkflowInstDTO) error {
	if utils.IsZero(d.InstID) {
		return fmt.Errorf("update inst `DefID` or `InstID` must not be zero, DefID:[%s] InstID=[%s]",
			d.DefID, d.InstID)
	}

	p := convertor.InstConvertor.ConvertUpdateDTOToPO(d)
	db := dao.db.Debug().Where("id=?", d.InstID)
	if d.DefID != "" {
		db.Where("def_id=?", d.DefID)
	}

	if err := db.Updates(p).Error; err != nil {
		log.Errorf("Failed to update workflow inst, caused by %s", err)
		return err
	}

	return nil
}

// PageQuery 流程实例分页查询
func (dao *WorkflowInstDAO) PageQuery(d *dto.PageQueryWorkflowInstDTO) ([]*po.WorkflowInstPO, error) {
	var workflowInsts []*po.WorkflowInstPO
	p, err := convertor.InstConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return nil, err
	}
	db := dao.db.ReadFromSlave(d.ReadFromSlave).Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit())
	if d.Name != "" {
		db.Where("name like ?", "%"+d.Name+"%")
	}
	if !d.CreatedBefore.IsZero() {
		db.Where("created_at < ?", d.CreatedBefore)
	}
	if err := db.Find(&workflowInsts).Error; err != nil {
		log.Errorf("Failed to page query workflow inst, caused by %s", err)
		return nil, err
	}

	return workflowInsts, nil
}

// Count 根据条件获取总数
func (dao *WorkflowInstDAO) Count(d *dto.PageQueryWorkflowInstDTO) (int64, error) {
	var totalCount int64
	p, err := convertor.InstConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return 0, err
	}
	db := dao.db.ReadFromSlave(d.ReadFromSlave).Model(po.WorkflowInstPO{}).Where(p)
	if d.Name != "" {
		db.Where("name like ?", "%"+d.Name+"%")
	}
	if !d.CreatedBefore.IsZero() {
		db.Where("created_at < ?", d.CreatedBefore)
	}

	if err := db.Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to get workflow inst count, caused by %s", err)
		return 0, err
	}

	return totalCount, nil
}

// UpdateFailed 流程流程实例为失败
func (dao *WorkflowInstDAO) UpdateFailed(d *dto.UpdateWorkflowInstFailedDTO) error {
	if utils.IsZero(d.InstID) {
		return fmt.Errorf("update inst `InstID` must not be zero, InstID=[%s]", d.InstID)
	}

	if err := dao.db.Exec("UPDATE workflow_inst "+
		"SET `context` = JSON_SET(`context`,'$.reason.failed_root_cause.failed_reason', ?), "+
		"status = ?, completed_at = now() WHERE id = ?", d.Reason, entity.InstFailed.IntValue(), d.InstID).
		Error; err != nil {
		log.Errorf("Failed to update workflow inst failed, caused by %s", err)
		return err
	}

	return nil
}
