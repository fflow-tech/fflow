package dto

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// WorkflowInstDTO 流程定义
type WorkflowInstDTO struct {
	entity.WorkflowInst
}

// StartWorkflowInstDTO 创建实例请求
type StartWorkflowInstDTO struct {
	Namespace        string                 `json:"namespace,omitempty" binding:"required"`
	DefID            string                 `json:"def_id,omitempty" binding:"required"`
	ParentInstID     string                 `json:"parent_inst_id,omitempty"`
	ParentNodeInstID string                 `json:"parent_node_inst_id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Creator          string                 `json:"creator,omitempty"`
	Input            map[string]interface{} `json:"input,omitempty"`
	Reason           string                 `json:"reason,omitempty"`
	DebugMode        bool                   `json:"debug_mode,omitempty"`
}

// RestartWorkflowInstDTO 重启实例请求
type RestartWorkflowInstDTO struct {
	Namespace   string                 `json:"namespace,omitempty"`
	DefID       string                 `json:"def_id,omitempty"`
	InstID      string                 `json:"inst_id,omitempty" binding:"required"`
	NodeRefName string                 `json:"node_ref_name,omitempty"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Operator    string                 `json:"operator,omitempty"`
	Reason      string                 `json:"reason,omitempty"`
}

// CancelWorkflowInstDTO 取消实例请求
type CancelWorkflowInstDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
	InstID    string `json:"inst_id,omitempty" binding:"required"`
	Operator  string `json:"operator,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// PauseWorkflowInstDTO 暂停实例请求
type PauseWorkflowInstDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
	InstID    string `json:"inst_id,omitempty" binding:"required"`
	Operator  string `json:"operator,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// ResumeWorkflowInstDTO 恢复实例执行请求
type ResumeWorkflowInstDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
	InstID    string `json:"inst_id,omitempty" binding:"required"`
	Operator  string `json:"operator,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// UpdateWorkflowInstCtxDTO 更新流程实例上下文请求
type UpdateWorkflowInstCtxDTO struct {
	Namespace string                 `json:"namespace,omitempty"`
	DefID     string                 `json:"def_id,omitempty"`
	InstID    string                 `json:"inst_id,omitempty" binding:"required"`
	Operator  string                 `json:"operator,omitempty"`
	Reason    string                 `json:"reason,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// DebugWorkflowInstDTO 流程实例调试信息请求
type DebugWorkflowInstDTO struct {
	Namespace         string           `json:"namespace,omitempty"`
	DefID             string           `json:"def_id,omitempty"`
	InstID            string           `json:"inst_id,omitempty" binding:"required"`
	Operator          string           `json:"operator,omitempty"`
	Reason            string           `json:"reason,omitempty"`
	AddBreakpoints    []string         `json:"add_breakpoints,omitempty"`
	DeleteBreakpoints []string         `json:"delete_breakpoints,omitempty"`
	AddMockNodes      []string         `json:"add_mock_nodes,omitempty"`
	DeleteMockNodes   []string         `json:"delete_mock_nodes,omitempty"`
	DebugMode         entity.DebugMode `json:"debug_mode,omitempty"`
}

// SetWorkflowInstTimeoutDTO 标记实例超时请求
type SetWorkflowInstTimeoutDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
	InstID    string `json:"inst_id,omitempty" binding:"required"`
	Reason    string `json:"reason,omitempty"`
}

// DriveWorkflowInstNextNodesDTO 驱动下一个节点请求
type DriveWorkflowInstNextNodesDTO struct {
	Namespace      string `json:"namespace,omitempty"`
	InstID         string `json:"inst_id,omitempty" binding:"required"`
	DefID          string `json:"def_id,omitempty"`
	FromResumeInst bool   `json:"from_resume_inst"`           // 来自恢复实例
	CurNodeInstID  string `json:"cur_node_inst_id,omitempty"` // 流程启动的时候不填, 其他时候都需要填
	Reason         string `json:"reason,omitempty"`
}

// CompleteWorkflowInstDTO 标记流程完成
type CompleteWorkflowInstDTO struct {
	Namespace string                 `json:"namespace,omitempty"`
	DefID     string                 `json:"def_id,omitempty"`
	InstID    string                 `json:"inst_id,omitempty" binding:"required"`
	Status    entity.InstStatus      `json:"status,omitempty"`
	Output    map[string]interface{} `json:"output,omitempty"`
	Operator  string                 `json:"operator,omitempty"`
	Reason    string                 `json:"reason,omitempty"`
}

// CreateWorkflowInstRepoDTO 创建实例内部请求
type CreateWorkflowInstRepoDTO struct {
	Namespace  string    `json:"namespace,omitempty"`
	DefID      string    `json:"def_id,omitempty" binding:"required"`
	DefVersion int       `json:"def_version,omitempty"`
	Name       string    `json:"name,omitempty"`
	Creator    string    `json:"creator,omitempty"`
	Context    string    `json:"context,omitempty"`
	Status     int       `json:"status,omitempty"`
	StartAt    time.Time `json:"start_at"`
}

// GetWorkflowInstDTO 查询请求
type GetWorkflowInstDTO struct {
	Namespace     string `json:"namespace,omitempty"`
	Operator      string `form:"operator" json:"operator,omitempty"`
	InstID        string `form:"inst_id" json:"inst_id,omitempty" binding:"required"`
	DefID         string `form:"def_id" json:"def_id,omitempty"`
	CurNodeInstID string `form:"cur_node_inst_id" json:"cur_node_inst_id,omitempty"`
}

// GetWorkflowInstsByIDsDTO 查询请求
type GetWorkflowInstsByIDsDTO struct {
	Namespace string   `json:"namespace,omitempty"`
	DefID     string   `form:"def_id" json:"def_ids,omitempty"`
	InstIDs   []string `form:"inst_ids" json:"inst_ids,omitempty"`
}

// ArchiveWorkflowInstsDTO 实例归档请求
type ArchiveWorkflowInstsDTO struct {
	Namespace string   `json:"namespace,omitempty"`
	DefID     string   `form:"def_id" json:"def_id,omitempty"`
	InstIDs   []string `form:"inst_ids" json:"inst_ids,omitempty" binding:"required"`
}

// GetWorkflowInstListDTO 分页查询请求
type GetWorkflowInstListDTO struct {
	*constants.PageQuery
	Namespace     string              `form:"namespace" json:"namespace,omitempty"`
	Operator      string              `form:"operator" json:"operator,omitempty"`
	DefID         string              `form:"def_id" json:"def_id,omitempty"`
	Name          string              `form:"string" json:"string,omitempty"`
	Status        entity.InstStatus   `form:"status" json:"status,omitempty"`       // 工作流状态
	Statuses      []entity.InstStatus `form:"statuses" json:"statuses,omitempty"`   // 工作流状态列表
	CreatedBefore time.Time           `json:"created_before,omitempty"`             // 在指定时间之后创建
	AscOrder      bool                `form:"asc_order" json:"asc_order,omitempty"` // 是否按时间升序排列
	ReadFromSlave bool                `json:"read_from_slave,omitempty"`            // 是否从备库读取数据
}

// GetWebWorkflowInstListDTO 给页面提供的查询
type GetWebWorkflowInstListDTO struct {
	*constants.PageQuery
	Namespace string   `form:"namespace" json:"namespace,omitempty"`
	Operator  string   `form:"operator" json:"operator,omitempty"`
	DefID     string   `form:"def_id" json:"def_id,omitempty"`
	Name      string   `form:"name" json:"name,omitempty"`
	Status    string   `form:"status" json:"status,omitempty"`       // 工作流状态
	Statuses  []string `form:"statuses" json:"statuses,omitempty"`   // 工作流状态列表
	AscOrder  bool     `form:"asc_order" json:"asc_order,omitempty"` // 是否按时间升序排列
}

// NewGetWorkflowInstDTO 新建
func NewGetWorkflowInstDTO(instID string, defID string, curNodeInstID string) *GetWorkflowInstDTO {
	return &GetWorkflowInstDTO{InstID: instID, DefID: defID, CurNodeInstID: curNodeInstID}
}

// UpdateWorkflowInstDTO 更新请求
type UpdateWorkflowInstDTO struct {
	Namespace   string    `json:"namespace,omitempty"`
	DefID       string    `json:"def_id,omitempty"`                     // 流程定义 ID
	InstID      string    `json:"inst_id,omitempty" binding:"required"` // 流程实例 ID
	Context     string    `json:"context,omitempty"`                    // 流程实例上下文
	Status      int       `json:"status,omitempty"`                     // 流程实例状态 1:running,2:completed,3:failed,4:cancelled,5:timeout
	CompletedAt time.Time `json:"completed_at" json:"completed_at"`     // 流程执行结束时间
}

// UpdateWorkflowInstFailedDTO 流程流程实例为失败
type UpdateWorkflowInstFailedDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`  // 流程定义ID
	InstID    string `json:"inst_id,omitempty"` // 流程实例ID
	Reason    string `json:"reason"`            // 原因
}

// DeleteWorkflowInstsDTO 删除流程实例请求
type DeleteWorkflowInstsDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
	InstID    string `json:"inst_id,omitempty"`
}

// DeleteWorkflowInstsByIDsDTO 删除流程实例请求
type DeleteWorkflowInstsByIDsDTO struct {
	Namespace string   `json:"namespace,omitempty"`
	InstIDs   []string `json:"inst_ids,omitempty"`
	DefID     string   `json:"def_id,omitempty"`
}

// PageQueryWorkflowInstDTO 分页查询对象
type PageQueryWorkflowInstDTO struct {
	Namespace     string            `json:"namespace,omitempty"`
	DefID         string            `json:"def_id,omitempty"`
	DefVersion    int               `json:"def_version,omitempty"`
	Name          string            `json:"name,omitempty"`
	Creator       string            `json:"creator,omitempty"`
	Status        entity.InstStatus `json:"status,omitempty"`
	CreatedBefore time.Time         `json:"created_before,omitempty"`  // 在指定时间之后创建
	ReadFromSlave bool              `json:"read_from_slave,omitempty"` // 是否从备库读取数据
	*constants.PageQuery
	*constants.Order
}
