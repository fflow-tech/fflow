// Package convertor 负责 entity 到 dto 的转换。
package convertor

import (
	"encoding/json"

	"github.com/jinzhu/copier"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

var (
	// DefConvertor 定时器转换体
	DefConvertor = &defConvertorImpl{}
)

type defConvertorImpl struct {
}

// ConvertEntityToDTO 定时器领域模型 entity -> dto
func (c *defConvertorImpl) ConvertEntityToDTO(e *entity.TimerDef) (*dto.TimerDefDTO, error) {
	notifyRpcParam := dto.NotifyRpcParam{}
	if err := json.Unmarshal([]byte(e.NotifyRpcParam), &notifyRpcParam); err != nil {
		return nil, err
	}

	notifyHttpParam := dto.NotifyHttpParam{}
	if err := json.Unmarshal([]byte(e.NotifyHttpParam), &notifyHttpParam); err != nil {
		return nil, err
	}

	d := &dto.TimerDefDTO{}
	if err := copier.Copy(d, e); err != nil {
		return nil, err
	}

	d.NotifyRpcParam = notifyRpcParam
	d.NotifyHttpParam = notifyHttpParam
	return d, nil
}
