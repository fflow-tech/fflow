package runtimecontext

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// RuntimeLogger 运行时 logger
type RuntimeLogger struct {
	function *entity.Function
	logs     []string
	ctx      context.Context
}

// 获取一个 logger
func newRuntimeLogger(ctx context.Context, function *entity.Function) *RuntimeLogger {
	return &RuntimeLogger{
		logs:     []string{},
		function: function,
		ctx:      ctx,
	}
}

// Logf 基础的 log 方法
func (r *RuntimeLogger) Logf(message string, args ...interface{}) {
	r.Infof(message, args...)
}

type logType string

var (
	debugLog logType = "DEBUG"
	infoLog  logType = "INFO"
	warnLog  logType = "WARN"
	errorLog logType = "ERROR"
)

func (r *RuntimeLogger) buildLogStr(t logType, message string, args ...interface{}) string {
	return fmt.Sprintf("%s %s %s", utils.GetCurrentLogTimestamp(), t, r.buildFileLogStr(message, args...))
}

func (r *RuntimeLogger) buildFileLogStr(message string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", r.buildFunctionKey(), fmt.Sprintf(message, args...))
}

func (r *RuntimeLogger) buildFunctionKey() string {
	debugMode := r.ctx.Value("debugMode")
	function := r.function

	// debug 模式下，或是缺失函数基础信息时，直接返回空字符串
	if debugMode.(bool) || function == nil || function.Namespace == "" {
		return ""
	}

	return fmt.Sprintf("[%s.%s]", function.Namespace, function.Name)
}

// Infof 基础的 log 方法
func (r *RuntimeLogger) Infof(message string, args ...any) {
	r.logs = append(r.logs, r.buildLogStr(infoLog, message, args...))
	log.Infof(r.buildFileLogStr(message, args...))
}

// Errorf 基础的 log 方法
func (r *RuntimeLogger) Errorf(message string, args ...any) {
	r.logs = append(r.logs, r.buildLogStr(errorLog, message, args...))
	log.Errorf(r.buildFileLogStr(message, args...))
}

// Debugf 基础的 log 方法
func (r *RuntimeLogger) Debugf(message string, args ...any) {
	r.logs = append(r.logs, r.buildLogStr(debugLog, message, args...))
	log.Debugf(r.buildFileLogStr(message, args...))
}

// Warnf 基础的 log 方法
func (r *RuntimeLogger) Warnf(message string, args ...any) {
	r.logs = append(r.logs, r.buildLogStr(warnLog, message, args...))
	log.Warnf(r.buildFileLogStr(message, args...))
}
