package sql

import (
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"gorm.io/gorm"
)

// FunctionDAO 数据访问对象
type FunctionDAO struct {
	db *mysql.Client
}

// NewFunctionDAO 构造函数
func NewFunctionDAO(db *mysql.Client) *FunctionDAO {
	return &FunctionDAO{db: db}
}

// Transaction 事务
func (dao *FunctionDAO) Transaction(f func(*mysql.Client) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		return f(mysql.NewClient(tx))
	})
}

// Create 创建函数
func (dao *FunctionDAO) Create(function *dto.CreateFunctionDTO) (*po.FunctionPO, error) {
	if utils.IsZero(function.Namespace) || utils.IsZero(function.Function) {
		return nil, fmt.Errorf("create func `namespace` + `function` must not be empty, "+
			"namespace:[%s] function:[%s]",
			function.Namespace, function.Function)
	}
	p := convertor.FunctionConvertor.ConvertCreateDTOToPO(function)
	if err := dao.db.Create(&p).Error; err != nil {
		log.Errorf("Failed to create function, caused by %s", err)
		return nil, err
	}

	return p, nil
}

// Get 获取函数信息
func (dao *FunctionDAO) Get(function *dto.GetFunctionReqDTO) (*po.FunctionPO, error) {
	if function.ID == 0 && (utils.IsZero(function.Namespace) || utils.IsZero(function.Function)) {
		return nil, fmt.Errorf("get func `namespace` + `function` must not be empty, "+
			"namespace:[%s] function:[%s]",
			function.Namespace, function.Function)
	}
	p := convertor.FunctionConvertor.ConvertGetDTOToPO(function)

	r := &po.FunctionPO{}
	// 当版本号为空时查询最新的版本
	if function.Version == 0 {
		if err := dao.db.Where(p).Last(r).Error; err != nil {
			log.Errorf("Failed to get appServer, caused by %s", err)
			return nil, err
		}
	} else {
		if err := dao.db.Where(p).Take(r).Error; err != nil {
			log.Errorf("Failed to get appServer, caused by %s", err)
			return nil, err
		}
	}

	return r, nil
}

// Delete 根据删除函数
func (dao *FunctionDAO) Delete(function *dto.DeleteFunctionDTO) error {
	if utils.IsZero(function.Namespace) || utils.IsZero(function.Function) {
		return fmt.Errorf("delete func `namespace` + `function` must not be empty, "+
			"namespace:[%s] function:[%s]",
			function.Namespace, function.Function)
	}

	p := convertor.FunctionConvertor.ConvertDeleteDTOToPO(function)
	if err := dao.db.Unscoped().Where(p).Delete(p).Error; err != nil {
		log.Errorf("Failed to delete function, caused by %s", err)
		return err
	}
	log.Infof("the function:[%s] of [namespace:%s]is deleted by [%s]", function.Function,
		function.Namespace, function.Operator)

	return nil
}

// Update 更新函数
func (dao *FunctionDAO) Update(function *dto.CreateFunctionDTO) error {
	if utils.IsZero(function.Namespace) ||
		utils.IsZero(function.Function) {
		return fmt.Errorf("update func `namespace` + `function` must not be empty, "+
			"namespace:[%s] function:[%s]",
			function.Namespace, function.Function)
	}

	p := convertor.FunctionConvertor.ConvertCreateDTOToPO(function)
	if err := dao.db.Where("app=? and server=? and service=? and name=?", function.Namespace,
		function.Function).Updates(p).Error; err != nil {
		log.Errorf("Failed to update function, caused by %s", err)
		return err
	}

	return nil
}

// PageQueryLastVersion 分页查询函数的最新版本
func (dao *FunctionDAO) PageQueryLastVersion(d *dto.PageQueryFunctionDTO) ([]*po.FunctionPO, error) {
	ids, err := dao.getIDs(d)
	if err != nil {
		return nil, err
	}

	var functions []*po.FunctionPO
	db := dao.db.Model(&po.FunctionPO{})
	if err := db.Where("id in (?)", ids).Order(d.OrderStr()).Find(&functions).Error; err != nil {
		log.Errorf("Failed to page query function, caused by %s", err)
		return nil, err
	}

	return functions, nil
}

func (dao *FunctionDAO) getIDs(d *dto.PageQueryFunctionDTO) ([]uint64, error) {
	if d.PageQuery == nil {
		d.PageQuery = constants.NewDefaultPageQuery()
	}

	var functionPOS []*po.FunctionPO
	db := dao.db.Model(&po.FunctionPO{})
	p := convertor.FunctionConvertor.ConvertPageQueryDTOToPO(d)
	db = db.Where(p).Order(d.OrderStr()).Offset(d.GetOffset()).Limit(d.GetLimit()).
		Group("name").Select("max(id) as id")
	if err := db.Find(&functionPOS).Error; err != nil {
		log.Errorf("Failed to page query function, caused by %s", err)
		return nil, err
	}

	r := []uint64{}
	for _, d := range functionPOS {
		r = append(r, uint64(d.ID))
	}
	return r, nil
}

// Count 根据条件获取总数
func (dao *FunctionDAO) Count(d *dto.PageQueryFunctionDTO) (int64, error) {
	if d.PageQuery == nil {
		d.PageQuery = constants.NewDefaultPageQuery()
	}

	var totalCount int64
	db := dao.db.Model(&po.FunctionPO{})
	p := convertor.FunctionConvertor.ConvertPageQueryDTOToPO(d)
	if err := db.Model(po.FunctionPO{}).Where(p).Group("name").Count(&totalCount).Error; err != nil {
		log.Errorf("Failed to get function count, caused by %s", err)
		return 0, err
	}

	return totalCount, nil
}
