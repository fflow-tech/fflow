package dto

import (
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// TriggerDTO 触发器定义
type TriggerDTO struct {
	Namespace  string `json:"namespace,omitempty"`
	Name       string `json:"name,omitempty"`        // 事件名称
	Condition  string `json:"condition,omitempty"`   // 事件触发条件
	Action     string `json:"action,omitempty"`      // 事件响应
	Type       int    `json:"type,omitempty"`        // 触发器类型 1:流程级别触发器, 2:节点级别触发器
	DefID      string `json:"def_id,omitempty"`      // 主键ID
	DefVersion int    `json:"def_version,omitempty"` // 流程的版本号
	InstID     string `json:"inst_id,omitempty"`     // 流程实例ID
	Status     int    `json:"status,omitempty"`      // 触发器状态，1:未激活, 2:已激活
}

// CreateTriggerDTO 创建触发器请求
type CreateTriggerDTO struct {
	entity.Trigger
	Status entity.TriggerStatus `json:"status"` // 触发器状态
}

// GetTriggerDTO 获取触发器请求
type GetTriggerDTO struct {
	Namespace string               `json:"namespace,omitempty"`
	TriggerID string               `json:"trigger_id,omitempty"`
	DefID     string               `json:"def_id,omitempty"`  // 主键ID
	InstID    string               `json:"inst_id,omitempty"` // 实例ID
	Event     string               `json:"name,omitempty"`    // 事件名称
	Level     entity.TriggerLevel  `json:"level,omitempty"`   // 触发器级别
	Type      entity.TriggerType   `json:"type,omitempty"`    // 触发器类型
	Status    entity.TriggerStatus `json:"status,omitempty"`  // 触发器状态，1:未激活, 2:已激活
}

// DeleteTriggerDTO 删除触发器请求
type DeleteTriggerDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
	TriggerID string `json:"trigger_id,omitempty"`
}

// UpdateTriggerDTO 更新触发器请求
type UpdateTriggerDTO struct {
	Namespace  string               `json:"namespace,omitempty"`
	DefID      string               `json:"def_id,omitempty"`
	DefVersion int                  `json:"def_version,omitempty"`
	InstID     string               `json:"inst_id,omitempty"`
	TriggerID  string               `json:"trigger_id,omitempty"`
	Level      entity.TriggerLevel  `json:"level,omitempty"`
	Status     entity.TriggerStatus `json:"status,omitempty"`
}

// PageQueryTriggerDTO 分页查询
type PageQueryTriggerDTO struct {
	Namespace  string               `json:"namespace,omitempty"`
	Event      string               `json:"event,omitempty"`       // 事件名称
	Level      entity.TriggerLevel  `json:"level,omitempty"`       // 触发器级别
	Type       entity.TriggerType   `json:"type,omitempty"`        // 触发器类型
	DefID      string               `json:"def_id,omitempty"`      // 主键ID
	DefVersion int                  `json:"def_version,omitempty"` // 版本号
	InstID     string               `json:"inst_id,omitempty"`     // 实例ID
	Status     entity.TriggerStatus `json:"status,omitempty"`      // 触发器状态
	*constants.PageQuery
	*constants.Order
}

// QueryTriggerDTO 批量查询
type QueryTriggerDTO struct {
	Namespace string `json:"namespace,omitempty"`
	Event     string `json:"event,omitempty"`  // 事件名称
	Type      int    `json:"type,omitempty"`   // 触发器类型 1:流程级别触发器, 2:节点级别触发器
	DefID     string `json:"def_id,omitempty"` // 主键 ID
	Status    int    `json:"status,omitempty"` // 触发器状态，1:未激活, 2:已激活
}

// RegisterTriggersDTO 注册多个触发器请求
type RegisterTriggersDTO struct {
	Namespace  string                         `json:"namespace,omitempty"`
	DefID      string                         `json:"def_id,omitempty"`
	DefVersion int                            `json:"def_version,omitempty"`
	InstID     string                         `json:"inst_id,omitempty"`
	Operator   string                         `json:"operator,omitempty"` // 操作人
	Level      entity.TriggerLevel            `json:"level,omitempty"`    // 触发器级别
	Triggers   []map[string]entity.TriggerDef `json:"triggers,omitempty"` // 触发器配置
}

// UnRegisterTriggersDTO 去注册多个触发器请求
type UnRegisterTriggersDTO struct {
	Namespace  string              `json:"namespace,omitempty"`
	DefID      string              `json:"def_id,omitempty"`
	DefVersion int                 `json:"def_version,omitempty"`
	InstID     string              `json:"inst_id,omitempty"`
	Level      entity.TriggerLevel `json:"level,omitempty"` // 触发器级别
}

// RegisterTriggerDTO 注册单个触发器请求
type RegisterTriggerDTO struct {
	Namespace  string              `json:"namespace,omitempty"`
	DefID      string              `json:"def_id,omitempty"`
	DefVersion int                 `json:"def_version,omitempty"`
	InstID     string              `json:"inst_id,omitempty"`
	Level      entity.TriggerLevel `json:"level,omitempty"`
	entity.TriggerDef
}

// UnregisterTriggerDTO 去注册单个触发器请求
type UnregisterTriggerDTO struct {
	Namespace  string              `json:"namespace,omitempty"`
	DefID      string              `json:"def_id,omitempty"`
	DefVersion int                 `json:"def_version,omitempty"`
	InstID     string              `json:"inst_id,omitempty"` // 实例 ID
	Level      entity.TriggerLevel `json:"level,omitempty"`   // 触发器级别 流程/定义
}
