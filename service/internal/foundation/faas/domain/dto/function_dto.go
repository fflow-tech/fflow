package dto

import (
	"net/http"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"
)

// CallFunctionReqDTO 运行函数请求
type CallFunctionReqDTO struct {
	Namespace string                 `form:"namespace,omitempty" json:"namespace,omitempty"` // 命名空间
	Function  string                 `form:"function,omitempty" json:"function,omitempty"`   // 函数
	Input     map[string]interface{} `form:"input,omitempty" json:"input,omitempty"`         // 函数的输入
	Operator  string                 `form:"operator,omitempty" json:"operator,omitempty"`   // 操作人
	Request   *http.Request          `form:"-" json:"-"`                                     // 请求的基础信息
}

// RequestInfo 请求基础信息
type RequestInfo struct {
	Headers map[string]interface{} `form:"headers,omitempty" json:"headers,omitempty"` // 请求头
	Method  string                 `form:"method,omitempty" json:"method,omitempty"`   // 请求方法
	Body    string                 `form:"body,omitempty" json:"body,omitempty"`       // 请求的源参数
}

// CallFunctionRspDTO 运行函数结果
type CallFunctionRspDTO struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// CreateFunctionDTO 创建函数请求
type CreateFunctionDTO struct {
	Namespace    string              `form:"namespace,omitempty" json:"namespace,omitempty"`         // 命名空间
	Function     string              `form:"function,omitempty" json:"function,omitempty"`           // 函数名
	Language     entity.LanguageType `form:"language,omitempty" json:"language,omitempty"`           // 所使用的语言
	Code         string              `form:"code,omitempty" json:"code,omitempty"`                   // 代码
	InputSchema  string              `form:"input_schema,omitempty" json:"input_schema,omitempty"`   // 函数入参格式
	OutputSchema string              `form:"output_schema,omitempty" json:"output_schema,omitempty"` // 函数返回结果格式
	Description  string              `form:"description,omitempty" json:"description,omitempty"`     // 描述
	Creator      string              `form:"creator,omitempty" json:"creator,omitempty"`             // 创建人
	Updater      string              `form:"Updater,omitempty" json:"Updater,omitempty"`             // 更新人
	Version      int                 `form:"version,omitempty" json:"version,omitempty"`             // 版本号
	Token        string              `form:"token,omitempty" json:"token,omitempty"`                 // token
}

// CreateFunctionReqDTO 创建函数请求
type CreateFunctionReqDTO struct {
	Namespace    string `form:"namespace,omitempty" json:"namespace,omitempty"`         // 命名空间
	Function     string `form:"function,omitempty" json:"function,omitempty"`           // 函数名
	Language     string `form:"language,omitempty" json:"language,omitempty"`           // 所使用的语言
	Code         string `form:"code,omitempty" json:"code,omitempty"`                   // 代码
	InputSchema  string `form:"input_schema,omitempty" json:"input_schema,omitempty"`   // 函数入参格式
	OutputSchema string `form:"output_schema,omitempty" json:"output_schema,omitempty"` // 函数返回结果格式
	Description  string `form:"description,omitempty" json:"description,omitempty"`     // 描述
	Creator      string `form:"creator,omitempty" json:"creator,omitempty"`             // 创建人
}

// GetFunctionReqDTO 获取函数请求
type GetFunctionReqDTO struct {
	ID        int    `form:"id,omitempty" json:"id,omitempty"`               // 函数 id
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"` // 命名空间
	Function  string `form:"function,omitempty" json:"function,omitempty"`   // 函数名
	Version   int    `form:"version,omitempty" json:"version,omitempty"`     // 版本号
	Creator   string `form:"creator,omitempty" json:"creator,omitempty"`     // 创建人
}

// PageQueryFunctionDTO 分页查询函数的请求
type PageQueryFunctionDTO struct {
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"` // 命名空间
	Function  string `form:"function,omitempty" json:"function,omitempty"`   // 函数名
	Language  string `form:"language,omitempty" json:"language,omitempty"`   // 语言
	Version   int    `form:"version,omitempty" json:"version,omitempty"`     // 版本号
	Creator   string `form:"creator,omitempty" json:"creator,omitempty"`     // 创建人
	*constants.PageQuery
	*constants.Order
}

// GetFunctionRspDTO 函数信息返回
type GetFunctionRspDTO struct {
	ID           int       `form:"id,omitempty" json:"id,omitempty"`
	Namespace    string    `form:"namespace,omitempty" json:"namespace,omitempty"`         // 命名空间
	Creator      string    `form:"creator,omitempty" json:"creator,omitempty"`             // 创建人
	Language     string    `form:"language,omitempty" json:"language,omitempty"`           // 所使用的语言 javascript/golang/starlark
	Code         string    `form:"code,omitempty" json:"code,omitempty"`                   // 代码
	InputSchema  string    `form:"input_schema,omitempty" json:"input_schema,omitempty"`   // 函数入参格式
	OutputSchema string    `form:"output_schema,omitempty" json:"output_schema,omitempty"` // 函数返回结果格式
	Description  string    `form:"description,omitempty" json:"description,omitempty"`     // 描述
	Function     string    `form:"function,omitempty" json:"function,omitempty"`           // 函数名
	Updater      string    `form:"updater,omitempty" json:"updater,omitempty"`             // 更新人
	Version      int       `form:"version,omitempty" json:"version,omitempty"`             // 版本号
	UpdatedAt    time.Time `form:"updated_at,omitempty" json:"updated_at,omitempty"`       // 更新时间
	CreatedAt    time.Time `form:"created_at,omitempty" json:"created_at,omitempty"`       // 创建时间
}

// DeleteFunctionDTO 删除函数DTO
type DeleteFunctionDTO struct {
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"` // 命名空间
	Function  string `form:"function,omitempty" json:"function,omitempty"`   // 函数名
	Operator  string `form:"operator,omitempty" json:"operator,omitempty"`   // 操作人
}

// UpdateFunctionDTO 修改函数定义
type UpdateFunctionDTO struct {
	Namespace    string `form:"namespace,omitempty" json:"namespace,omitempty"`         // 命名空间
	Function     string `form:"function,omitempty" json:"function,omitempty"`           // 函数名
	Updater      string `form:"updater,omitempty" json:"updater,omitempty"`             // 修改人
	Code         string `form:"code,omitempty" json:"code,omitempty"`                   // 代码
	InputSchema  string `form:"input_schema,omitempty" json:"input_schema,omitempty"`   // 函数入参格式
	OutputSchema string `form:"output_schema,omitempty" json:"output_schema,omitempty"` // 函数返回结果格式
	Description  string `form:"description,omitempty" json:"description,omitempty"`     // 描述
}

// DebugFunctionDTO 调试函数请求
type DebugFunctionDTO struct {
	Namespace string                 `form:"namespace,omitempty" json:"namespace,omitempty"` // 命名空间
	Code      string                 `form:"code,omitempty" json:"code,omitempty"`           // 函数
	Language  string                 `form:"language,omitempty" json:"language,omitempty"`   // 所使用的语言
	Input     map[string]interface{} `form:"input,omitempty" json:"input,omitempty"`         // 函数的输入
	Operator  string                 `form:"operator,omitempty" json:"operator,omitempty"`   // 操作人
}

// DebugFunctionRspDTO 运行函数结果
type DebugFunctionRspDTO struct {
	Result interface{} `json:"result,omitempty"`
	Logs   []string    `json:"logs,omitempty"`
	Error  string      `json:"error,omitempty"`
}
