package sql

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// HistoryWorkflowInstDAO 历史数据访问对象
type HistoryWorkflowInstDAO struct {
	db *mysql.Client
}

// NewHistoryWorkflowInstDAO 数据访问对象构造函数
func NewHistoryWorkflowInstDAO(db *mysql.Client) *HistoryWorkflowInstDAO {
	return &HistoryWorkflowInstDAO{db: db}
}

// Transaction 事务
func (dao *HistoryWorkflowInstDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// BatchCreate	批量创建
func (dao *HistoryWorkflowInstDAO) BatchCreate(req []*dto.HistoryWorkflowInstDTO) error {
	p, err := convertor.HistoryWorkflowInstConvertor.ConvertDTOsToPOs(req)
	if err != nil {
		return err
	}

	if err := dao.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&p, 100).Error; err != nil {
		log.Errorf("Failed to batch create workflow inst, caused by %s", err)
		return err
	}

	return nil
}
