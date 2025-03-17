package dto

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// NodeInstDTO 节点实例定义
type NodeInstDTO struct {
	entity.NodeInst
}

// CreateNodeInstDTO 创建
type CreateNodeInstDTO struct {
	Namespace   string    `json:"namespace,omitempty" form:"namespace"`
	DefID       string    `json:"def_id,omitempty"`
	DefVersion  int       `json:"def_version,omitempty"`
	InstID      string    `json:"inst_id,omitempty"`
	RefName     string    `json:"ref_name,omitempty"`
	Context     string    `json:"context,omitempty"`
	Status      int       `json:"status,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

// GetNodeInstDTO 查询
type GetNodeInstDTO struct {
	Namespace  string `json:"namespace,omitempty" form:"namespace"`
	Operator   string `json:"operator,omitempty" form:"operator"`
	NodeInstID string `json:"node_inst_id,omitempty" form:"node_inst_id"`
	DefID      string `json:"def_id,omitempty" form:"def_id"`
	DefVersion int    `json:"def_version,omitempty" form:"def_version"`
	InstID     string `json:"inst_id,omitempty" form:"inst_id"`
	RefName    string `json:"ref_name,omitempty" form:"ref_name"`
	Status     int    `json:"status,omitempty" form:"status"`
}

// ArchiveNodeInstsDTO 实例归档请求
type ArchiveNodeInstsDTO struct {
	DefID       string   `form:"def_id" json:"def_id,omitempty"`
	InstID      string   `json:"inst_id,omitempty" form:"inst_id"`
	NodeInstIDs []string `form:"node_inst_ids" json:"node_inst_ids,omitempty" binding:"required"`
}

// GetNodeInstsByIDsDTO 查询请求
type GetNodeInstsByIDsDTO struct {
	DefID       string   `form:"def_id" json:"def_ids,omitempty"`
	InstID      string   `json:"inst_id,omitempty" form:"inst_id"`
	NodeInstIDs []string `form:"node_inst_ids" json:"node_inst_ids,omitempty" binding:"required"`
}

// DeleteNodeInstsByIDsDTO 删除请求
type DeleteNodeInstsByIDsDTO struct {
	DefID       string   `form:"def_id" json:"def_ids,omitempty"`
	InstID      string   `json:"inst_id,omitempty" form:"inst_id"`
	NodeInstIDs []string `form:"node_inst_ids" json:"node_inst_ids,omitempty" binding:"required"`
}

// DeleteNodeInstDTO 删除
type DeleteNodeInstDTO struct {
	Namespace  string
	NodeInstID string
	DefID      string
}

// UpdateNodeInstDTO 更新
type UpdateNodeInstDTO struct {
	Namespace     string
	NodeInstID    string
	DefID         string
	Context       string
	Status        int
	WaitAt        time.Time
	ExecuteAt     time.Time
	AsynWaitResAt time.Time
	CompletedAt   time.Time
}

// PageQueryNodeInstDTO 分页查询
type PageQueryNodeInstDTO struct {
	Namespace     string                  `json:"namespace,omitempty"`
	DefID         string                  `json:"def_id,omitempty"`
	DefVersion    int                     `json:"def_version,omitempty"`
	InstID        string                  `json:"inst_id,omitempty"`
	RefName       string                  `json:"ref_name,omitempty"`
	Statuses      []entity.NodeInstStatus `json:"statuses,omitempty"`
	ReadFromSlave bool                    `json:"read_from_slave,omitempty"` // 是否从备库读取数据
	*constants.PageQuery
	*constants.Order
}

// RerunNodeDTO 重跑节点
type RerunNodeDTO struct {
	Namespace   string                 `json:"namespace,omitempty"`
	DefID       string                 `json:"def_id,omitempty"`
	InstID      string                 `json:"inst_id,omitempty"`
	NodeRefName string                 `json:"node_ref_name,omitempty"`
	Operator    string                 `json:"operator,omitempty"`
	Reason      string                 `json:"reason,omitempty"`
	Input       map[string]interface{} `json:"input,omitempty"`
}

// RunNodeDTO 运行节点
type RunNodeDTO struct {
	Namespace   string                 `json:"namespace,omitempty"`
	DefID       string                 `json:"def_id,omitempty"`
	InstID      string                 `json:"inst_id,omitempty"`
	NodeInstID  string                 `json:"node_inst_id,omitempty"` // 因为已经调度了节点，这里会有 NodeInstID
	NodeRefName string                 `json:"node_ref_name,omitempty"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Operator    string                 `json:"operator,omitempty"`
	Reason      string                 `json:"reason,omitempty"`
}

// CancelNodeDTO 取消节点运行
type CancelNodeDTO struct {
	Namespace   string `json:"namespace,omitempty"`
	DefID       string `json:"def_id,omitempty"`
	InstID      string `json:"inst_id,omitempty"`
	NodeInstID  string `json:"node_inst_id,omitempty"`
	NodeRefName string `json:"node_ref_name,omitempty"`
	Operator    string `json:"operator,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// ResumeNodeDTO 恢复运行节点
type ResumeNodeDTO struct {
	Namespace   string                 `json:"namespace,omitempty"`
	DefID       string                 `json:"def_id,omitempty"`
	InstID      string                 `json:"inst_id,omitempty"`
	NodeInstID  string                 `json:"node_inst_id,omitempty"`
	NodeRefName string                 `json:"node_ref_name,omitempty"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Operator    string                 `json:"operator,omitempty"`
	Reason      string                 `json:"reason,omitempty"`
}

// CompleteNodeDTO 标记节点完成
type CompleteNodeDTO struct {
	Namespace   string                 `json:"namespace,omitempty"`
	DefID       string                 `json:"def_id,omitempty"`
	InstID      string                 `json:"inst_id,omitempty"`
	NodeInstID  string                 `json:"node_inst_id,omitempty"`
	NodeRefName string                 `json:"node_ref_name,omitempty"`
	Status      entity.NodeInstStatus  `json:"status,omitempty"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Operator    string                 `json:"operator,omitempty"`
	Reason      string                 `json:"reason,omitempty"` // 操作原因
}

// PollingNodeDTO 轮询节点
type PollingNodeDTO struct {
	Namespace  string `json:"namespace,omitempty"`
	DefID      string `json:"def_id,omitempty"`
	DefVersion int    `json:"def_version,omitempty"`
	InstID     string `json:"inst_id,omitempty"`
	NodeInstID string `json:"node_inst_id,omitempty"`
}

// ScheduleNodeDTO 调度节点
type ScheduleNodeDTO struct {
	Namespace  string `json:"namespace,omitempty"`
	DefID      string `json:"def_id,omitempty"`
	DefVersion int    `json:"def_version,omitempty"`
	InstID     string `json:"inst_id,omitempty"`
	NodeInstID string `json:"node_inst_id,omitempty"`
}

// SkipNodeDTO 跳过节点
type SkipNodeDTO struct {
	Namespace   string `json:"namespace,omitempty"`
	DefID       string `json:"def_id,omitempty"`
	InstID      string `json:"inst_id,omitempty"`
	NodeRefName string `json:"node_ref_name,omitempty"`
	Operator    string `json:"operator,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// CancelSkipNodeDTO 取消跳过节点
type CancelSkipNodeDTO struct {
	Namespace   string `json:"namespace,omitempty"`
	DefID       string `json:"def_id,omitempty"`
	InstID      string `json:"inst_id,omitempty"`
	NodeRefName string `json:"node_ref_name,omitempty"`
	Operator    string `json:"operator,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// SetNodeTimeoutDTO 标记节点超时
type SetNodeTimeoutDTO struct {
	Namespace   string `json:"namespace,omitempty"`
	DefID       string `json:"def_id,omitempty"`
	InstID      string `json:"inst_id,omitempty"`
	NodeInstID  string `json:"node_inst_id,omitempty"`
	NodeRefName string `json:"node_ref_name,omitempty"`
	Operator    string `json:"operator,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// SetNodeNearTimeoutDTO 标记节点节点超时
type SetNodeNearTimeoutDTO struct {
	Namespace   string `json:"namespace,omitempty"`
	DefID       string `json:"def_id,omitempty"`
	InstID      string `json:"inst_id,omitempty"`
	NodeInstID  string `json:"node_inst_id,omitempty"`
	NodeRefName string `json:"node_ref_name,omitempty"`
	Operator    string `json:"operator,omitempty"`
	Reason      string `json:"reason,omitempty"`
}
