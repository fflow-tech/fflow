package golang

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/fflow-tech/fflow-sdk-go/faas"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/constants"
	"github.com/pkg/errors"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// TODO(@kekuangao): 部分语法屏蔽如 go 关键字
// 禁用包
var symbolBlackList = []string{
	"os/",
	"net/",
	"runtime/",
	"github.com/traefik/yaegi",
}

// golangExecutor golang 语言的脚本执行器
type golangExecutor struct {
}

// Symbols 允许引入的包
var Symbols = map[string]map[string]reflect.Value{}

// NewGolangExecutor 新建
func NewGolangExecutor() *golangExecutor {
	// 配置可以使用的包
	setAllowedSymbols()
	return &golangExecutor{}
}

// Execute 执行脚本
func (e *golangExecutor) Execute(ctx faas.Context, code string, params map[string]interface{}) (
	result interface{}, logs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("panic in execute caused by %s", r))
			// 即使函数 panic 也返回之前部分的日志，便于调试
			logs = ctx.Logs()
		}
	}()
	// 初始化执行器
	interpret := initGolangInterpret()
	err = e.beforeExecute()
	if err != nil {
		return nil, nil, err
	}

	// 解析当前脚本
	if _, err := interpret.Eval(code); err != nil {
		return nil, nil, err
	}

	// 执行入口函数
	f, err := interpret.Eval(fmt.Sprintf("%s.%s", constants.DefaultPackageName, constants.DefaultEntryFunctionName))
	if err != nil {
		return nil, nil, err
	}
	handler, ok := f.Interface().(constants.FunctionType)
	if !ok {
		return nil, nil, errors.New(fmt.Sprintf("the function is not %v, which is illegal", constants.FunctionTypeStr))
	}

	// 执行函数并拿到返回值
	result, err = handler(ctx, params)
	logs = ctx.Logs()
	return
}

// beforeExecute 执行前校验，如果要对要执行的代码做一些限制可以写在这里
func (e *golangExecutor) beforeExecute() error {
	return nil
}

// initInterpret 初始化执行器
func initGolangInterpret() *interp.Interpreter {
	i := interp.New(interp.Options{})
	i.Use(Symbols)
	return i
}

// setAllowedSymbols 获取所有可以使用的包（从 stdlib.Symbols 删除一些高危的包）
func setAllowedSymbols() {
	for k, v := range stdlib.Symbols {
		if isAllowedSymbol(k) {
			Symbols[k] = v
		}
	}
}

// isAllowedSymbol 判断是否是允许用的包
func isAllowedSymbol(symbol string) bool {
	p := path.Dir(symbol)
	for _, s := range symbolBlackList {
		// 包名相同
		if p == s {
			return false
		}
		// 黑名单下的所有包和其子包不能用
		if strings.HasSuffix(s, "/") && (strings.HasPrefix(p, s) || p == strings.TrimSuffix(s, "/")) {
			return false
		}
	}
	return true
}
