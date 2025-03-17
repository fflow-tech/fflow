package dto

import "github.com/fflow-tech/fflow/service/pkg/constants"

// App 应用定义
type App struct {
	ID      int    // 应用ID
	Name    string // 应用名称
	Creator string // 创建人
}

// CountAppDTO 查询 APP 总数参数体
type CountAppDTO struct {
	Name    string `json:"name,omitempty"`    // 应用名称
	Creator string `json:"creator,omitempty"` // 创建人
}

// PageQueryAppDTO 分页查询 APP 参数体
type PageQueryAppDTO struct {
	Name    string `json:"name,omitempty" form:"name,omitempty"`       // 应用名称
	Creator string `json:"creator,omitempty" form:"creator,omitempty"` // 创建人
	*constants.PageQuery
	*constants.Order
}

// CreateAppDTO 创建 APP 参数体
type CreateAppDTO struct {
	Name    string `json:"name,omitempty" form:"name,omitempty"`       // 应用名称
	Creator string `json:"creator,omitempty" form:"creator,omitempty"` // 创建人
}

// GetAppDTO 获取 APP 定义参数体
type GetAppDTO struct {
	Name string `json:"name,omitempty" form:"name,omitempty"` // 应用名称
}

// DeleteAppDTO 删除 APP 定义参数体
type DeleteAppDTO struct {
	Name string `json:"name,omitempty" form:"name,omitempty"` // 应用名称
}
