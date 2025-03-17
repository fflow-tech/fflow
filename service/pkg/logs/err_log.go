// Package logs 日志相关组件包
package logs

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/fflow-tech/fflow/service/pkg/log"
)

// DumpPanicStack 进程崩溃重启，并捕获 panic 输出
func DumpPanicStack(bizKey string, err error) {
	if e := recover(); e != nil {
		errLog().Infof("[Panic:%s] err:%s, stack:%s", bizKey, err.Error(), GetStackInfo())
	}
}

// GetStackInfo 获取调用栈信息
func GetStackInfo() string {
	return strings.Replace(fmt.Sprintf("%s", debug.Stack()), "\t", "    ", -1)
}

// GetErrLogName 错误日志
func GetErrLogName() string {
	return "err_log"
}

func errLog() log.Logger {
	return log.GetDefaultLogger()
}
