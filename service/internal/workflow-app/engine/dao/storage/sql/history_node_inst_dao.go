package sql

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// HistoryNodeInstDAO 数据访问对象
type HistoryNodeInstDAO struct {
	db *mysql.Client
}

// NewHistoryNodeInstDAO 数据访问对象构造函数
func NewHistoryNodeInstDAO(db *mysql.Client) *HistoryNodeInstDAO {
	return &HistoryNodeInstDAO{db: db}
}

// Transaction 事务
func (dao *HistoryNodeInstDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// BatchCreate	批量创建
func (dao *HistoryNodeInstDAO) BatchCreate(req []*dto.HistoryNodeInstDTO) error {
	p := convertor.HistoryNodeInstConvertor.ConvertDTOsToPOs(req)

	if err := dao.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&p, 100).Error; err != nil {
		log.Errorf("Failed to batch create node inst, caused by %s", err)
		return err
	}

	return nil
}
