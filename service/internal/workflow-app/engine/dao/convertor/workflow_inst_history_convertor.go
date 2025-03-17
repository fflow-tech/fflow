package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	HistoryWorkflowInstConvertor = &historyWorkflowInstConvertorImpl{} // 转换器
)

type historyWorkflowInstConvertorImpl struct {
}

// ConvertDTOsToPOs 转换
func (*historyWorkflowInstConvertorImpl) ConvertDTOsToPOs(
	req []*dto.HistoryWorkflowInstDTO) ([]*po.HistoryWorkflowInstPO, error) {
	hs := make([]*po.HistoryWorkflowInstPO, 0)
	for _, s := range req {
		d := &po.HistoryWorkflowInstPO{}
		var err error
		if err = copier.Copy(d, s); err != nil {
			return nil, err
		}

		hs = append(hs, d)
	}
	return hs, nil
}
