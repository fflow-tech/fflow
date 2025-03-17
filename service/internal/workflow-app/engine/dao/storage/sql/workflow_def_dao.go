package sql

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/seq"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"gorm.io/gorm"
)

// WorkflowDefDAO WorkFlowDef数据访问对象
type WorkflowDefDAO struct {
	db *mysql.Client
}

// NewWorkflowDefDAO WorkFlowDef数据访问对象构造函数
func NewWorkflowDefDAO(db *mysql.Client) *WorkflowDefDAO {
	return &WorkflowDefDAO{db: db}
}

// Transaction 事务
func (dao *WorkflowDefDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建流程定义
// DefID通过请求返回
func (dao *WorkflowDefDAO) Create(def *dto.CreateWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	// 创建时当不带defID时 获取defID
	if def.DefID == "" {
		var err error
		def.DefID, err = seq.NewString()
		if err != nil {
			return nil, err
		}
	}

	p, err := convertor.DefConvertor.ConvertCreateDTOToPO(def)
	if err != nil {
		return nil, err
	}
	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create workflow def, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// BatchCreate 批量创建流程定义
// DefID通过原始的创建请求返回
func (dao *WorkflowDefDAO) BatchCreate(defs []*dto.CreateWorkflowDefDTO) error {
	var defPOs []*po.WorkflowDefPO
	// 创建时当不带defID时 获取defID
	for _, def := range defs {
		if def.DefID == "" {
			var err error
			def.DefID, err = seq.NewString()
			if err != nil {
				return err
			}
		}
		p, err := convertor.DefConvertor.ConvertCreateDTOToPO(def)
		if err != nil {
			return err
		}
		defPOs = append(defPOs, p)
	}

	if err := dao.db.Create(&defPOs).Error; err != nil {
		log.Errorf("Failed to create workflow defs, caused by %s", err)
		return err
	}

	return nil
}

// Get 获取流程定义信息
func (dao *WorkflowDefDAO) Get(def *dto.GetWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	if utils.IsZero(def.DefID) || utils.IsZero(def.Version) {
		return nil, fmt.Errorf("get def `DefID` and `DefVersion`  must not be zero,"+
			" DefID:[%s] DefVersion:[%d]", def.DefID, def.Version)
	}

	p, err := convertor.DefConvertor.ConvertGetDTOToPO(def)
	if err != nil {
		return nil, err
	}
	r := &po.WorkflowDefPO{}
	if err := dao.db.ReadFromSlave(def.ReadFromSlave).Where(p).Take(r).Error; err != nil {
		log.Errorf("Failed to get workflow def, caused by %s", err)
		return nil, err
	}
	return r, nil
}

// Delete 删除流程
func (dao *WorkflowDefDAO) Delete(d *dto.DeleteWorkflowDefDTO) error {
	if utils.IsZero(d.DefID) || utils.IsZero(d.Version) {
		return fmt.Errorf("delete def `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	p, err := convertor.DefConvertor.ConvertDeleteDTOToPO(d)
	if err != nil {
		return err
	}
	if err := dao.db.Where(p).Delete(p).Error; err != nil {
		log.Errorf("Failed to delete workflow def, caused by %s", err)
		return err
	}

	return nil
}

// Update 更新流程
func (dao *WorkflowDefDAO) Update(d *dto.UpdateWorkflowDefDTO) error {
	if utils.IsZero(d.DefID) || utils.IsZero(d.Version) {
		return fmt.Errorf("Update def `DefID` must not be zero, DefID:[%d]", d.DefID)
	}

	p := convertor.DefConvertor.ConvertUpdateDTOToPO(d)
	if err := dao.db.Where("def_id = ? and version = ?", d.DefID, d.Version).
		Updates(p).Error; err != nil {
		log.Errorf("Failed to update workflow def, caused by %s", err)
		return err
	}

	return nil
}

// PageQueryLastVersion 分页查询流程定义的最新版本
func (dao *WorkflowDefDAO) PageQueryLastVersion(d *dto.PageQueryWorkflowDefDTO) (
	[]*po.WorkflowDefPO, error) {
	ids, err := dao.getIDs(d)
	if err != nil {
		return nil, err
	}

	var workflowDefs []*po.WorkflowDefPO
	db := dao.db.ReadFromSlave(d.ReadFromSlave).Model(&po.WorkflowDefPO{})
	if err := db.Where("id in (?)", ids).Order(d.OrderStr()).Find(&workflowDefs).Error; err != nil {
		log.Errorf("Failed to page query workflow def, caused by %s", err)
		return nil, err
	}

	return workflowDefs, nil
}

func (dao *WorkflowDefDAO) getIDs(d *dto.PageQueryWorkflowDefDTO) ([]uint64, error) {
	if d.PageQuery == nil {
		d.PageQuery = constants.NewDefaultPageQuery()
	}

	var workflowDefs []*po.WorkflowDefPO
	db := dao.db.Model(&po.WorkflowDefPO{})
	// 流程名称支持模糊查询
	if d.Name != "" {
		db = db.Where("name like ?", "%"+d.Name+"%")
		d.Name = ""
	}
	if d.Operator != "" {
		db = db.Where("creator = ?", d.Operator)
	}

	p, err := convertor.DefConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return nil, err
	}
	db = db.Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).
		Group("def_id").Select("max(id) as id")
	if err := db.Find(&workflowDefs).Error; err != nil {
		log.Errorf("Failed to page query workflow def, caused by %s", err)
		return nil, err
	}

	var r []uint64
	for _, d := range workflowDefs {
		r = append(r, uint64(d.ID))
	}
	return r, nil
}

// Count 根据条件获取总数
func (dao *WorkflowDefDAO) Count(d *dto.PageQueryWorkflowDefDTO) (int64, error) {
	if d.PageQuery == nil {
		d.PageQuery = constants.NewDefaultPageQuery()
	}

	var totalCount int64
	db := dao.db.ReadFromSlave(d.ReadFromSlave).Model(&po.WorkflowDefPO{})
	// 流程名称支持模糊查询
	if d.Name != "" {
		db = db.Where("name like ?", "%"+d.Name+"%")
		d.Name = ""
	}
	if d.Operator != "" {
		db = db.Where("creator = ?", d.Operator)
	}

	p, err := convertor.DefConvertor.ConvertPageQueryDTOToPO(d)
	if err != nil {
		return 0, err
	}
	if err := db.Model(po.WorkflowDefPO{}).Where(p).Group("def_id").Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to get workflow def count, caused by %s", err)
		return 0, err
	}

	return totalCount, nil
}

// GetLastVersion 获取定义的最新版本
func (dao *WorkflowDefDAO) GetLastVersion(d *dto.GetWorkflowDefDTO) (*po.WorkflowDefPO, error) {
	if utils.IsZero(d.DefID) {
		return nil, fmt.Errorf("getLastVersion `DefID` must not be zero, DefID:[%s]", d.DefID)
	}

	r := &po.WorkflowDefPO{}
	p, err := convertor.DefConvertor.ConvertGetDTOToPO(d)
	if err != nil {
		return nil, err
	}
	if err := dao.db.ReadFromSlave(d.ReadFromSlave).Where(p).Last(r).Error; err != nil {
		log.Errorf("Failed to get workflow def last version, caused by %s", err)
		return nil, err
	}

	return r, nil
}

// GetSubWorkflowLastVersion 获取 subworkflow 最新版本
func (dao *WorkflowDefDAO) GetSubWorkflowLastVersion(d *dto.GetSubworkflowDefDTO) (*po.WorkflowDefPO, error) {
	if utils.IsZero(d.ParentDefID) || d.RefName == "" {
		return nil, fmt.Errorf("failed to get subWorkflow last version `ParentDefID` must not be zero And `RefName`"+
			" must not be empty, ParentDefID:[%s], RefName:[%s]", d.ParentDefID, d.RefName)
	}

	r := &po.WorkflowDefPO{}
	if err := dao.db.ReadFromSlave(false).
		Where("parent_def_id = ?", d.ParentDefID).
		Where("attribute -> '$.ref_name' = (?)", d.RefName).
		Where("attribute -> '$.parent_def_version' = (?)", d.ParentDefVersion).
		Last(r).Error; err != nil {
		log.Errorf("Failed to get subworkflow last version for workflow [%d][%d][%s], caused by %s",
			d.ParentDefID, d.ParentDefVersion, d.RefName, err)
		return nil, err
	}

	return r, nil
}
