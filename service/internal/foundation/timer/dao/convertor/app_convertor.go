package convertor

import (
	"github.com/jinzhu/copier"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

var (
	// AppConvertor dto->po 转换体
	AppConvertor = &appConvertor{}
)

type appConvertor struct {
}

// ConvertCreateDTOToPO 创建定时器定义 DTO->PO
func (*appConvertor) ConvertCreateDTOToPO(d *dto.CreateAppDTO) *po.App {
	p := &po.App{}
	copier.Copy(p, d)
	return p
}

// ConvertGetDTOToPO 查询定时器定义 DTO->PO
func (*appConvertor) ConvertGetDTOToPO(d *dto.GetAppDTO) *po.App {
	p := &po.App{}
	copier.Copy(p, d)
	return p
}

// ConvertDeleteDTOToPO 删除定时器定义 DTO->PO
func (*appConvertor) ConvertDeleteDTOToPO(d *dto.DeleteAppDTO) *po.App {
	p := &po.App{}
	copier.Copy(p, d)
	return p
}
