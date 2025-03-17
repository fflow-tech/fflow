// Package query 领域 query 服务，对外提供查询能力。
package query

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
)

// TimerDefQueryService 定时器定义查询服务
type TimerDefQueryService struct {
	timerDefRepo ports.TimerDefRepository
}

// NewTimerDefQueryService 新建查询服务
func NewTimerDefQueryService(r *ports.RepoProviderSet) *TimerDefQueryService {
	return &TimerDefQueryService{timerDefRepo: r.TimerDefRepo()}
}

// GetTimerDef 获取定时器定义
func (m *TimerDefQueryService) GetTimerDef(d *dto.GetTimerDefDTO) (*dto.TimerDefDTO, error) {
	timerDefEntity, err := m.timerDefRepo.GetTimerDef(d)
	if err != nil {
		return nil, err
	}
	return convertor.DefConvertor.ConvertEntityToDTO(timerDefEntity)
}

// GetTimerDefList 获取定时器定义列表
func (m *TimerDefQueryService) GetTimerDefList(d *dto.PageQueryTimeDefDTO) ([]*dto.TimerDefDTO, int64, error) {
	timerDefList, total, err := m.timerDefRepo.GetTimerDefList(d)
	if err != nil {
		return nil, 0, err
	}

	timerDTOs := make([]*dto.TimerDefDTO, 0, len(timerDefList))
	for _, timerEntity := range timerDefList {
		timerDefDTO, err := convertor.DefConvertor.ConvertEntityToDTO(timerEntity)
		if err != nil {
			return nil, 0, err
		}
		timerDTOs = append(timerDTOs, timerDefDTO)
	}
	return timerDTOs, total, nil
}

// CountTimersByStatus 根据状态统计定时器数量.
func (m *TimerDefQueryService) CountTimersByStatus(status entity.TimerDefStatus) (int64, error) {
	return m.timerDefRepo.CountTimersByStatus(status)
}
