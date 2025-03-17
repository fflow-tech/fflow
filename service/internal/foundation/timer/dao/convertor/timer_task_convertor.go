package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"

	"github.com/jinzhu/copier"
)

var (
	// TaskConvertor 定时器任务转换体
	TaskConvertor = &taskConvertorImpl{}
)

type taskConvertorImpl struct {
}

// ConvertCreateDTOToPO  创建运行记录 DTO->PO
func (*taskConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateRunHistoryDTO) (*po.RunHistoryPO, error) {
	p := &po.RunHistoryPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertGetDTOToPO  查询单个运行记录 DTO->PO
func (*taskConvertorImpl) ConvertGetDTOToPO(d *dto.GetRunHistoryDTO) (*po.RunHistoryPO, error) {
	p := &po.RunHistoryPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertDeleteDTOToPO  删除运行记录 DTO->PO
func (*taskConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteRunHistoryDTO) (*po.RunHistoryPO, error) {
	p := &po.RunHistoryPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertUpdateDTOToPO  更新运行记录 DTO->PO
func (*taskConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateRunHistoryDTO) (*po.RunHistoryPO, error) {
	p := &po.RunHistoryPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}

// ConvertPageQueryDTOToPO 分页查询运行记录 DTO->PO
func (*taskConvertorImpl) ConvertPageQueryDTOToPO(d *dto.PageQueryRunHistoryDTO) *po.RunHistoryPO {
	p := &po.RunHistoryPO{
		DefID:    d.DefID,
		Name:     d.Name,
		RunTimer: d.RunTimer,
	}
	return p
}
