// Package redis 负责与 redis 交互，结合业务完成增删改查。
package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/redis"
)

const (
	tickerTableName = "ticker_table" // 定时器的存放表名称
)

// TimerDefDAO 定时器定义
type TimerDefDAO struct {
	redisClient *redis.Client
}

// NewTimerDefClient 新建定时器定义客户端
func NewTimerDefClient(db *redis.Client) *TimerDefDAO {
	return &TimerDefDAO{
		redisClient: db,
	}
}

// AddTimerDef 增加定时器定义
func (t *TimerDefDAO) AddTimerDef(def *dto.CreateTimerDefDTO) error {
	v, err := t.redisClient.Hexists(context.Background(), tickerTableName, def.DefID)
	if err != nil {
		return err
	}
	if v == 1 {
		return fmt.Errorf("failed to AddTimerDef, caused by def has been exists")
	}

	p, err := convertor.DefConvertor.ConvertCreateDTOToPO(def)
	if err != nil {
		return err
	}

	jsonValue, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return t.redisClient.HSet(context.Background(), tickerTableName, def.DefID, string(jsonValue))
}

// GetTimerDef 获取定时器定义
func (t *TimerDefDAO) GetTimerDef(d *dto.GetTimerDefDTO) (*po.TimerDefPO, error) {
	v, err := t.redisClient.HGet(context.Background(), tickerTableName, d.DefID)
	if err != nil {
		return nil, err
	}

	timerDef := &po.TimerDefPO{}
	err = json.Unmarshal([]byte(v), timerDef)
	if err != nil {
		return nil, fmt.Errorf("failed to GetTimerDef, caused by Unmarshal err:%v", err)
	}
	return timerDef, nil
}

// DelTimerDef 删除定时器定义
func (t *TimerDefDAO) DelTimerDef(d *dto.DeleteTimerDefDTO) error {
	return t.redisClient.HDel(context.Background(), tickerTableName, d.DefID)
}

// ChangeTimerStatus 更改定时器定义状态
func (t *TimerDefDAO) ChangeTimerStatus(d *dto.ChangeTimerStatusDTO) error {
	timerDef, err := t.GetTimerDef(&dto.GetTimerDefDTO{DefID: d.DefID})
	if err != nil {
		return err
	}

	timerDef.Status = d.Status
	jsonValue, err := json.Marshal(timerDef)
	if err != nil {
		return err
	}
	return t.redisClient.HSet(context.Background(), tickerTableName, timerDef.DefID, string(jsonValue))
}
