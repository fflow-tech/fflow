package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	HistoryNodeInstConvertor = &historyNodeInstConvertorImpl{} // 转换器
)

type historyNodeInstConvertorImpl struct {
}

// ConvertDTOsToPOs 转换
func (*historyNodeInstConvertorImpl) ConvertDTOsToPOs(req []*dto.HistoryNodeInstDTO) []*po.HistoryNodeInstPO {
	historyPOs := make([]*po.HistoryNodeInstPO, 0)
	for _, historyDTO := range req {
		historyPO := po.HistoryNodeInstPO{}
		copier.Copy(&historyPO, historyDTO)
		historyPOs = append(historyPOs, &historyPO)
	}
	return historyPOs
}
