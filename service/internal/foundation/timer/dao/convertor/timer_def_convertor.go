// Package convertor 负责 DTO 到 PO 转换。
package convertor

import (
	"encoding/json"
	"strconv"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

var (
	// DefConvertor  定时器定义转换体
	DefConvertor = &defConvertorImpl{}
)

type defConvertorImpl struct {
}

// ConvertCreateDTOToPO  创建定时器DTO转换PO
func (*defConvertorImpl) ConvertCreateDTOToPO(d *dto.CreateTimerDefDTO) (*po.TimerDefPO, error) {
	notifyRpcParam, err := json.Marshal(d.NotifyRpcParam)
	if err != nil {
		return nil, err
	}

	notifyHttpParam, err := json.Marshal(d.NotifyHttpParam)
	if err != nil {
		return nil, err
	}

	p := &po.TimerDefPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}

	p.NotifyRpcParam = string(notifyRpcParam)
	p.NotifyHttpParam = string(notifyHttpParam)
	p.ExecuteTimeLimit = d.ExecuteTimeLimit
	return p, nil
}

// ConvertDeleteDTOToPO  删除定时器定义 DTO->PO
func (*defConvertorImpl) ConvertDeleteDTOToPO(d *dto.DeleteTimerDefDTO) (*po.TimerDefPO, error) {
	id, err := strconv.Atoi(d.DefID)
	if err != nil {
		return nil, err
	}

	p := &po.TimerDefPO{Model: gorm.Model{ID: uint(id)}}
	return p, nil
}

// ConvertUpdateDTOToPO 更新定时器定义 DTO->PO
func (*defConvertorImpl) ConvertUpdateDTOToPO(d *dto.UpdateTimerDefDTO) (*po.TimerDefPO, error) {
	p := &po.TimerDefPO{}
	if err := copier.Copy(p, d); err != nil {
		return nil, err
	}
	return p, nil
}
