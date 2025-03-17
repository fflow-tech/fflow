package convertor

import (
	"database/sql"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	UserConvertor = &userConvertorImpl{} // 转换器
)

type userConvertorImpl struct {
}

// ConvertDTOToPO 转换
func (*userConvertorImpl) ConvertDTOToPO(req *dto.CreateUserDTO) (*po.UserPO, error) {
	p := &po.UserPO{}
	var err error
	if err = copier.Copy(p, req); err != nil {
		return nil, err
	}
	if len(req.Phone) > 0 {
		p.Phone = sql.NullString{String: req.Phone, Valid: true}
	} else {
		p.Phone = sql.NullString{String: req.Phone, Valid: false}
	}
	if len(req.Email) > 0 {
		p.Email = sql.NullString{String: req.Email, Valid: true}
	} else {
		p.Email = sql.NullString{String: req.Email, Valid: false}
	}
	return p, nil
}

// ConvertGetDTOToPO 转换
func (c *userConvertorImpl) ConvertGetDTOToPO(d *dto.GetUserDTO) (*po.UserPO, error) {
	p := &po.UserPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	if len(d.Email) > 0 {
		p.Email = sql.NullString{String: d.Email, Valid: true}
	} else {
		p.Email = sql.NullString{String: d.Email, Valid: false}
	}

	if len(d.Phone) > 0 {
		p.Phone = sql.NullString{String: d.Phone, Valid: true}
	} else {
		p.Phone = sql.NullString{String: d.Phone, Valid: false}
	}

	return p, nil
}

// ConvertDeleteDTOToPO 转换
func (c *userConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteUserDTO) (*po.UserPO, error) {
	p := &po.UserPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertUpdateDTOToPO 转换
func (c *userConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateUserDTO) (*po.UserPO, error) {
	p := &po.UserPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertPageQueryDTOToPO 转换
func (c *userConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryUserDTO) (*po.UserPO, error) {
	p := &po.UserPO{}
	var err error
	if err = copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}
