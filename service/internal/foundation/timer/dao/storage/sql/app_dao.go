package sql

import (
	"fmt"
	"gorm.io/gorm"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// AppDAO App 数据访问对象
type AppDAO struct {
	db *mysql.Client
}

// NewAppDAO NewAppDAO 数据访问对象构造函数
func NewAppDAO(db *mysql.Client) *AppDAO {
	return &AppDAO{db: db}
}

// Transaction 事务
func (dao *AppDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建 App 定义
func (dao *AppDAO) Create(d *dto.CreateAppDTO) (*po.App, error) {
	if utils.IsZero(d.Name) || utils.IsZero(d.Creator) {
		return nil, fmt.Errorf("create app `GlobalConfigName`、`Creator` must not be zero, GlobalConfigName:[%s] Creator:[%s]",
			d.Name, d.Creator)
	}
	p := convertor.AppConvertor.ConvertCreateDTOToPO(d)
	if err := dao.db.Create(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

// Get 查询 App 定义
func (dao *AppDAO) Get(d *dto.GetAppDTO) (*po.App, error) {
	if utils.IsZero(d.Name) {
		return nil, fmt.Errorf("get app `GlobalConfigName` must not be zero, GlobalConfigName:[%s]", d.Name)
	}

	app := &po.App{}
	p := convertor.AppConvertor.ConvertGetDTOToPO(d)
	if err := dao.db.Where(p).Take(app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

// PageQuery 分页查询 App 定义列表
func (dao *AppDAO) PageQuery(d *dto.PageQueryAppDTO) ([]*po.App, error) {
	var apps []*po.App
	db := dao.db.Model(&po.App{})
	if !utils.IsZero(d.Name) {
		db = db.Where("name like ?", "%"+d.Name+"%")
	}

	if !utils.IsZero(d.Creator) {
		db = db.Where("creator like ?", "%"+d.Creator+"%")
	}

	if err := db.Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

// Delete 删除 app 定义
func (dao *AppDAO) Delete(d *dto.DeleteAppDTO) error {
	if utils.IsZero(d.Name) {
		return fmt.Errorf("Delete app `GlobalConfigName` must not be zero, GlobalConfigName:[%s]", d.Name)
	}

	p := convertor.AppConvertor.ConvertDeleteDTOToPO(d)
	if err := dao.db.Unscoped().Where(p).Delete(p).Error; err != nil {
		return err
	}

	return nil
}

// Count 查询 app 总数
func (dao *AppDAO) Count(d *dto.CountAppDTO) (int64, error) {
	var total int64
	db := dao.db.Model(&po.App{})
	if !utils.IsZero(d.Name) {
		db = db.Where("name like ?", "%"+d.Name+"%")
	}

	if !utils.IsZero(d.Creator) {
		db = db.Where("creator like ?", "%"+d.Creator+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
