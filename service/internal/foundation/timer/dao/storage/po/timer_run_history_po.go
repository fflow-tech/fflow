package po

import "gorm.io/gorm"

// RunHistoryPO 运行流水记录
type RunHistoryPO struct {
	gorm.Model
	DefID    string `gorm:"column:def_id;NOT NULL"`        // 定义ID
	Name     string `gorm:"column:name;NOT NULL"`          // 定时器名称
	Output   string `gorm:"column:output;default:null"`    // 执行结果
	RunTimer string `gorm:"column:run_timer;default:null"` // 执行时间
	CostTime int    `gorm:"column:cost_time"`              // 执行耗时
	Status   string `gorm:"column:status;NOT NULL"`        // 当前状态
}

// TableName 表名
func (m *RunHistoryPO) TableName() string {
	return "run_history"
}
