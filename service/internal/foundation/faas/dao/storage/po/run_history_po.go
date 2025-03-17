package po

import "gorm.io/gorm"

// RunHistoryPO 运行流水记录
type RunHistoryPO struct {
	gorm.Model
	Namespace string `gorm:"column:namespace;NOT NULL"`  // 命名空间
	Name      string `gorm:"column:name;NOT NULL"`       // 函数名称
	Operator  string `gorm:"column:operator;NOT NULL"`   // 执行者
	Input     string `gorm:"column:input;default:null"`  // 函数入参
	Output    string `gorm:"column:output;default:null"` // 执行结果
	Log       string `gorm:"column:log"`                 // 执行日志
	CostTime  int    `gorm:"column:cost_time"`           // 执行耗时
	Version   int    `gorm:"column:version;NOT NULL"`    // 版本号
	Status    string `gorm:"column:status;NOT NULL"`     // 当前状态
}

// TableName 表名
func (m *RunHistoryPO) TableName() string {
	return "run_history"
}
