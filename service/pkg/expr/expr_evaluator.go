package expr

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"github.com/spf13/cast"
)

// Evaluator 表达式计算器
type Evaluator interface {
	// IsExpression 是否是合法的表达式
	IsExpression(expr string) bool
	// Evaluate 根据上下文计算值
	Evaluate(ctx map[string]interface{}, expr string) (interface{}, error)
	// Match 根据上下文计算是否匹配
	Match(ctx map[string]interface{}, expr string) (bool, error)
	// EvaluateMap 根据上下文计算 map
	EvaluateMap(ctx map[string]interface{}, exprMap map[string]interface{}) (map[string]interface{}, error)
}

// DefaultEvaluator 计算器
type DefaultEvaluator struct {
}

// NewDefaultEvaluator 初始化计算器
func NewDefaultEvaluator() *DefaultEvaluator {
	return &DefaultEvaluator{}
}

// IsExpression 是否是表达式
func (c *DefaultEvaluator) IsExpression(expr string) bool {
	return strings.HasPrefix(expr, "${") && strings.HasSuffix(expr, "}")
}

var (
	curTimeFormatFunc = gval.Function("curtimeformat", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("curtimeformat() expects exactly one string argument")
		}

		return time.Now().Format(args[0].(string)), nil
	})
	sprintfFunc = gval.Function("sprintf", func(args ...interface{}) (interface{}, error) {
		if len(args) <= 1 {
			return nil, fmt.Errorf("sprintf() expects > 1 argument")
		}

		return fmt.Sprintf(args[0].(string), args[1:]...), nil
	})
)

// Evaluate 计算
func (c *DefaultEvaluator) Evaluate(ctx map[string]interface{}, expr string) (interface{}, error) {
	realExpr, err := c.getRealExpr(expr)
	if err != nil {
		return false, err
	}

	return gval.Evaluate(realExpr, ctx, jsonpath.Language(), curTimeFormatFunc, sprintfFunc)
}

// Match 是否匹配
func (c *DefaultEvaluator) Match(ctx map[string]interface{}, expr string) (bool, error) {
	realExpr, err := c.getRealExpr(expr)
	if err != nil {
		return false, err
	}

	result, err := gval.Evaluate(realExpr, ctx, jsonpath.Language())
	if err != nil {
		return false, err
	}

	return cast.ToBool(result), nil
}

// EvaluateMap 替换map中的表达式
func (c *DefaultEvaluator) EvaluateMap(ctx map[string]interface{},
	exprMap map[string]interface{}) (map[string]interface{}, error) {
	rspExpressionMap := make(map[string]interface{})
	for k, v := range exprMap {
		var err error
		rspExpressionMap[k], err = c.evaluateInterface(ctx, v)
		if err != nil {
			return nil, err
		}
	}
	return rspExpressionMap, nil
}

// evaluateInterface 计算 interface 中的表达式
func (c *DefaultEvaluator) evaluateInterface(ctx map[string]interface{},
	expr interface{}) (interface{}, error) {
	switch exprValue := expr.(type) {
	case string:
		// 如果不是表达式 直接返回 是的话就做替换
		if !c.IsExpression(exprValue) {
			return exprValue, nil
		}
		return c.Evaluate(ctx, exprValue)
	case map[string]interface{}:
		var err error
		rsp := make(map[string]interface{})
		for k, v := range exprValue {
			rsp[k], err = c.evaluateInterface(ctx, v)
			if err != nil {
				return nil, err
			}
		}
		return rsp, nil
	case []interface{}:
		rsp := make([]interface{}, 0, len(exprValue))
		for _, v := range exprValue {
			interfaceValue, err := c.evaluateInterface(ctx, v)
			if err != nil {
				return nil, err
			}
			rsp = append(rsp, interfaceValue)
		}
		return rsp, nil
	default:
		// 其它类型直接返回值
		return exprValue, nil
	}
}

// getRealExpr 获取实际的表达式
func (c *DefaultEvaluator) getRealExpr(expr string) (string, error) {
	if !c.IsExpression(expr) {
		return "", fmt.Errorf("illegal expr=[%s]", expr)
	}

	return expr[2 : len(expr)-1], nil
}
