package dto

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"

	"mime/multipart"
)

// WorkflowDefDTO 流程定义
type WorkflowDefDTO struct {
	DefID       string             `gorm:"column:def_id;NOT NULL" json:"def_id,omitempty"`                 // 主键ID
	Version     int                `gorm:"column:version;NOT NULL" json:"version,omitempty"`               // 流程的版本号
	ParentDefID string             `gorm:"column:parent_def_id;NOT NULL" json:"parent_def_id,omitempty"`   // 父流程定义 ID
	Attribute   po.WorkflowDefAttr `gorm:"column:attribute;type:json;NOT NULL" json:"attribute,omitempty"` // 额外属性
	Name        string             `gorm:"column:name;NOT NULL" json:"name,omitempty"`                     // 流程定义名称
	DefJson     string             `gorm:"column:def_json;NOT NULL" json:"def_json,omitempty"`             // 流程定义的内容
	Creator     string             `gorm:"column:creator;NOT NULL" json:"creator,omitempty"`               // 创建人
	Status      entity.DefStatus   `gorm:"column:status;NOT NULL" json:"status,omitempty"`                 // 流程定义状态
	Description string             `gorm:"column:description" json:"description,omitempty"`                // 流程定义描述
	CreatedAt   time.Time          `gorm:"column:created_at" json:"created_at,omitempty"`                  // 流程定义创建时间
}

// CreateWorkflowDefDTO 创建流程请求
type CreateWorkflowDefDTO struct {
	Namespace   string             `json:"namespace,omitempty"`
	DefID       string             `json:"def_id,omitempty"`        // 主键ID
	Version     int                `json:"version,omitempty"`       // 流程的版本号
	Name        string             `json:"name,omitempty"`          // 流程定义名称
	ParentDefID string             `json:"parent_def_id,omitempty"` // 父流程定义 ID
	Attribute   po.WorkflowDefAttr `json:"attribute,omitempty"`     // 额外属性
	DefJson     string             `json:"def_json,omitempty"`      // 流程定义的内容
	Creator     string             `json:"creator,omitempty"`       // 创建人
	Description string             `json:"description,omitempty"`   // 流程定义描述
}

// GetWorkflowDefDTO 流程定义查询
type GetWorkflowDefDTO struct {
	Namespace     string `json:"namespace,omitempty"`
	Creator       string `form:"creator" json:"creator,omitempty"`
	DefID         string `form:"def_id" json:"def_id,omitempty" binding:"required"`
	Status        int    `form:"status" json:"status,omitempty"`
	Version       int    `form:"version" json:"version,omitempty"`
	ReadFromSlave bool   `json:"read_from_slave,omitempty"` // 是否从备库读取数据
}

// GetSubworkflowDefDTO 子流程定义查询
type GetSubworkflowDefDTO struct {
	Namespace        string `json:"namespace,omitempty"`
	ParentDefID      string `json:"parent_def_id,omitempty"`
	ParentDefVersion int    `json:"parent_def_version,omitempty"`
	RefName          string `json:"ref_name,omitempty"`
}

// CreateSubworkflowDefDTO 创建子流程
type CreateSubworkflowDefDTO struct {
	Namespace        string `json:"namespace,omitempty"`
	ParentDefID      string `json:"parent_def_id,omitempty"`
	ParentDefVersion int    `json:"parent_def_version,omitempty"`
	DefJson          string `json:"def_json,omitempty"`
}

// GetAllSubworkflowDefsDTO 所有子流程定义查询
type GetAllSubworkflowDefsDTO struct {
	Namespace        string `json:"namespace,omitempty"`
	ParentDefID      string `json:"parent_def_id,omitempty"`
	ParentDefVersion int    `json:"parent_def_version,omitempty"`
	DefJson          string `json:"def_json,omitempty"`
}

// DeleteWorkflowDefDTO 删除流程请求
type DeleteWorkflowDefDTO struct {
	Namespace string `json:"namespace,omitempty"`
	DefID     string `json:"def_id,omitempty"`
	Version   int    `json:"version,omitempty"`
}

// UpdateWorkflowDefDTO 更新流程请求
type UpdateWorkflowDefDTO struct {
	Namespace string           `json:"namespace,omitempty"`
	DefID     string           `json:"def_id,omitempty"`
	Version   int              `json:"version,omitempty"`
	Status    entity.DefStatus `json:"status"`
}

// PageQueryWorkflowDefDTO 分页查询流程请求
type PageQueryWorkflowDefDTO struct {
	Namespace     string `form:"namespace" json:"namespace,omitempty"`
	Operator      string `form:"operator" json:"operator,omitempty"`
	DefID         string `form:"def_id" json:"def_id,omitempty"`
	Name          string `form:"name" json:"name,omitempty"`
	Status        int    `form:"status" json:"status,omitempty"`
	Version       int    `form:"version" json:"version,omitempty"`
	GroupBy       string `form:"group_by" json:"group_by,omitempty"`
	ReadFromSlave bool   `json:"read_from_slave,omitempty"` // 是否从备库读取数据
	*constants.PageQuery
	*constants.Order
}

// PageQueryWebWorkflowDefDTO 提供给前端的分页查询流程请求
type PageQueryWebWorkflowDefDTO struct {
	Namespace string `form:"namespace" json:"namespace,omitempty"`
	Operator  string `form:"operator" json:"operator,omitempty"`
	DefID     string `form:"def_id" json:"def_id,omitempty"`
	Name      string `form:"name" json:"name,omitempty"`
	Status    string `form:"status" json:"status,omitempty"`
	Version   int    `form:"version" json:"version,omitempty"`
	GroupBy   string `form:"group_by" json:"group_by,omitempty"`
	*constants.PageQuery
	*constants.Order
}

// DisableWorkflowDefDTO 去激活流程请求
type DisableWorkflowDefDTO struct {
	Namespace string `form:"namespace" json:"namespace,omitempty"`
	Operator  string `form:"operator" json:"operator,omitempty"`
	DefID     string `form:"def_id" json:"def_id,omitempty"`
}

// EnableWorkflowDefDTO 激活流程请求
type EnableWorkflowDefDTO struct {
	Namespace string `form:"namespace" json:"namespace,omitempty"`
	Operator  string `form:"operator" json:"operator,omitempty"`
	DefID     string `form:"def_id" json:"def_id,omitempty" binding:"required"`
}

// UploadWorkflowDefDTO 上传流程请求
type UploadWorkflowDefDTO struct {
	Namespace    string                `form:"namespace" json:"namespace,omitempty"`
	DefID        string                `form:"def_id" json:"def_id,omitempty"`                         // 主键 ID
	Creator      string                `form:"creator" json:"creator,omitempty"`                       // 操作人
	Name         string                `form:"name" json:"name,omitempty"`                             // 流程定义名称
	DefJson      string                `form:"def_json" json:"def_json,omitempty"`                     // 流程定义的内容
	WorkflowFile *multipart.FileHeader `form:"workflow_file"  json:"workflow_file" binding:"required"` // [必填] 流程文件
}
