package js

import (
	"fmt"

	"github.com/fflow-tech/fflow-sdk-go/faas"

	"github.com/dop251/goja"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/constants"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/js/common"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/js/compiler"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/js/modules"
)

// javascriptExecutor javascirpt 语言的脚本执行器
type javascriptExecutor struct {
}

// NewJavascriptExecutor 新建
func NewJavascriptExecutor() *javascriptExecutor {
	return &javascriptExecutor{}
}

// Execute 执行脚本
func (e *javascriptExecutor) Execute(ctx faas.Context, code string, params map[string]interface{}) (
	result interface{}, logs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in execute: %s", r)
			logs = ctx.Logs()
		}
	}()

	vm := goja.New()
	pgm, err := e.beforeExecute(ctx, vm, code)
	if err != nil {
		return
	}

	// 解析当前脚本
	if _, err = vm.RunProgram(pgm); err != nil {
		return
	}
	var handler constants.FunctionType
	err = vm.ExportTo(vm.Get(constants.DefaultEntryFunctionName), &handler)
	if err != nil {
		return
	}

	// 执行脚本
	result, err = handler(ctx, params)
	logs = ctx.Logs()
	return
}

// beforeExecute 执行前引入内置包和编译
func (e *javascriptExecutor) beforeExecute(ctx faas.Context, vm *goja.Runtime, code string) (
	*goja.Program, error) {
	funcCtx := ctx.Context()
	symbols := modules.InitModules(ctx)
	for k, v := range symbols {
		vm.Set(k, common.Bind(vm, v, &funcCtx))
	}

	return e.getProgram(ctx, code)
}

// getProgram 获取编译后的执行程序
func (e *javascriptExecutor) getProgram(ctx faas.Context, code string) (*goja.Program, error) {
	c := compiler.New(ctx)
	function := ctx.Metadata()
	filename := fmt.Sprintf("[%s.%s]", function.Namespace(), function.Name())
	// ES6+ 编译为 ES5 并返回编译后的结果; 编译耗时较长，默认不转换
	pgm, _, err := c.Compile(code, filename, "", "", false)
	return pgm, err
}
