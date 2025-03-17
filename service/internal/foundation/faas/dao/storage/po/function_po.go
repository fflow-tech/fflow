package po

import (
	"gorm.io/gorm"
)

// FunctionPO 函数对象
type FunctionPO struct {
	gorm.Model
	Name         string `gorm:"column:name;NOT NULL"`              // 函数名
	Namespace    string `gorm:"column:namespace;NOT NULL"`         // 命名空间
	Creator      string `gorm:"column:creator;NOT NULL"`           // 创建人
	Updater      string `gorm:"column:updater;NOT NULL"`           // 更新人
	Code         string `gorm:"column:code"`                       // 代码
	InputSchema  string `gorm:"column:input_schema;default:null"`  // 函数入参格式
	OutputSchema string `gorm:"column:output_schema;default:null"` // 函数返回结果格式
	Description  string `gorm:"column:description"`                // 描述
	Language     string `gorm:"column:language;NOT NULL"`          // 所使用语言 javascript/golang
	Version      int    `gorm:"column:version;NOT NULL"`           // 版本号
	Token        string `gorm:"column:token"`                      // 函数 token
}

// TableName 表名
func (m *FunctionPO) TableName() string {
	return "function"
}
