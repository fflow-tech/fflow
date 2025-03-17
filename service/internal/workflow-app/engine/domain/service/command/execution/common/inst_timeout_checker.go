package common

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/expr"
)

var (
	maxWorkflowInstDurationTime = 60 * 24 * time.Hour // 流程最长持续时间为两个月
	maxNodeInstDurationTime     = 30 * 24 * time.Hour // 流程最长持续时间为一个月
	maxExprDurationTime         = 30 * 24 * time.Hour // 表达式最长持续一个月
	minExprDurationTime         = 30 * time.Second    // 表达式类型最小 30s
)

// InstTimeoutChecker 实例超时检查器
type InstTimeoutChecker interface {
	CheckWorkflowInst(inst *entity.WorkflowInst) (bool, error)
	CheckNodeInst(nodeInst *entity.NodeInst) (bool, error)
	CheckNodeInstNearTimeout(nodeInst *entity.NodeInst) (bool, error)
}

// DefaultInstTimeoutChecker 默认超时检查器
type DefaultInstTimeoutChecker struct {
}

// NewDefaultInstTimeoutChecker 初始化默认超时检查器
func NewDefaultInstTimeoutChecker() *DefaultInstTimeoutChecker {
	return &DefaultInstTimeoutChecker{}
}

// CheckWorkflowInst 检查实例是否超时
func (d *DefaultInstTimeoutChecker) CheckWorkflowInst(inst *entity.WorkflowInst) (bool, error) {
	// 先判断流程状态
	if inst.Status.IsTerminal() {
		return false, fmt.Errorf("inst status=%s is already terminal", inst.Status.String())
	}
	// 流程的默认超时时间是 4 周
	if inst.WorkflowDef.Timeout.Duration == "" {
		inst.WorkflowDef.Timeout.Duration = "4w"
	}

	return d.checkDurationTimeout(inst.WorkflowDef.Timeout.Duration, maxWorkflowInstDurationTime, inst.StartAt)
}

// CheckNodeInst 检查节点是否超时
func (d *DefaultInstTimeoutChecker) CheckNodeInst(nodeInst *entity.NodeInst) (bool, error) {
	// 0. 判断节点状态
	if nodeInst.Status.IsTerminal() {
		return false, fmt.Errorf("nodeInst status=%s is already terminal", nodeInst.Status.String())
	}
	// 1. 先检测时间间隔的超时，超时则直接返回
	durationTimeout, err := d.checkDurationTimeout(nodeInst.BasicNodeDef.Timeout.Duration,
		maxNodeInstDurationTime, nodeInst.ScheduledAt)
	if err != nil {
		return false, err
	}
	if durationTimeout {
		return true, err
	}

	// 2. 如果没有配置事件间隔的超时，则检查表达式的超时
	return d.checkExprTimeout(nodeInst.BasicNodeDef.Timeout.Expr, nodeInst.ScheduledAt)
}

// CheckNodeInstNearTimeout 检查节点是否接近超时
func (d *DefaultInstTimeoutChecker) CheckNodeInstNearTimeout(nodeInst *entity.NodeInst) (bool, error) {
	// 节点状态是结束时不进行超时判断，返回错误，其他节点状态均要进接近行超时判断
	if nodeInst.Status.IsTerminal() {
		return false, fmt.Errorf("nodeInst status=%s is already terminal", nodeInst.Status)
	}
	// 检查是否接近 duration 字段设置的超时时间
	durationTimeout, err := d.checkDurationTimeout(nodeInst.BasicNodeDef.Timeout.NearTimeoutDuration,
		maxNodeInstDurationTime, nodeInst.ScheduledAt)
	if err != nil {
		return false, err
	}
	if durationTimeout {
		return true, err
	}
	return d.checkExprTimeout(nodeInst.BasicNodeDef.Timeout.NearTimeoutExpr, nodeInst.ScheduledAt)
}

// checkDurationTimeout 检测时间间隔类型超时
func (d *DefaultInstTimeoutChecker) checkDurationTimeout(durationConfig string, maxDuration time.Duration,
	startAt time.Time) (bool, error) {
	// 当duration没有配置时 不返回错误 算未超时
	if durationConfig == "" {
		return false, nil
	}
	duration, err := expr.ParseDuration(durationConfig)
	if err != nil {
		return false, err
	}
	// 判断配置是否超过最大时间
	if duration > maxDuration {
		return false, fmt.Errorf("duration %d exceed max duration %d", duration, maxDuration)
	}
	return startAt.Add(duration).Before(time.Now()), nil
}

// checkExprTimeout 检测表达式类型超时
func (d *DefaultInstTimeoutChecker) checkExprTimeout(exprConfig string, startAt time.Time) (bool, error) {
	// 当表达式没有配置时 则认为没有超时 不算错误
	if exprConfig == "" {
		return false, nil
	}
	expression, err := cronexpr.Parse(exprConfig)
	if err != nil {
		return false, fmt.Errorf("illegal expr format:%v", err)
	}
	nextTimeout := expression.Next(startAt)
	if nextTimeout.After(startAt.Add(maxExprDurationTime)) || nextTimeout.Before(startAt.Add(minExprDurationTime)) {
		return false, fmt.Errorf("expr next timeout:%v must < %v > %v",
			nextTimeout, startAt.Add(maxExprDurationTime), startAt.Add(minExprDurationTime))
	}
	return nextTimeout.Before(time.Now()), nil
}
