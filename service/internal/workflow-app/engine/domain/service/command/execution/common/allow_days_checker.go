// Package common 提供流程执行过程中的公共能力
package common

import (
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// AllowDaysChecker 执行时间检查器
type AllowDaysChecker interface {
	Check(checkTime time.Time, allowDaysPolicy entity.AllowDaysPolicy) (bool, error) // 检查时间是否满足可执行的条件
}

// DefaultAllowDaysChecker 默认执行时间检查器
type DefaultAllowDaysChecker struct {
}

// NewDefaultAllowDaysChecker 初始化
func NewDefaultAllowDaysChecker() (*DefaultAllowDaysChecker, error) {
	return &DefaultAllowDaysChecker{}, nil
}

// checkAllowDaysFunc 判断是否在可执行的时间段内对应的检测方法
var checkAllowDaysFunc = map[entity.AllowDaysPolicy]func(checkTime time.Time) bool{
	entity.Any: func(checkTime time.Time) bool { return true },
	entity.Week: func(checkTime time.Time) bool {
		Weekend := checkTime.Weekday()
		return Weekend != time.Saturday && Weekend != time.Sunday
	},
	entity.Weekend: func(checkTime time.Time) bool {
		Weekend := checkTime.Weekday()
		return Weekend == time.Sunday || Weekend == time.Saturday
	},
}

// Check 检查时间是否满足可执行的条件
func (d *DefaultAllowDaysChecker) Check(checkTime time.Time, allowDaysPolicy entity.AllowDaysPolicy) (bool, error) {
	checkFunc, ok := checkAllowDaysFunc[allowDaysPolicy]
	// 默认使用 ANY 策略
	if !ok {
		return checkAllowDaysFunc[entity.Any](checkTime), nil
	}
	return checkFunc(checkTime), nil
}
