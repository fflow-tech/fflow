/*
Package execution 用于多语言 FAAS 脚本的具体执行
根据语言类型调用不同的执行器，并获取返回结果，记录每一次调用的元数据，统一的超时控制等
JS 的执行器参考 https://github.com/grafana/k6 实现，通过 babel 编译支持 ES6+ 语法，
提供了一些内置的 modules，出于安全考虑暂不支持用户引入任意包
*/
package execution

import (
	c "context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/fflow-tech/fflow-sdk-go/faas"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/golang"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/js"
	"github.com/panjf2000/ants/v2"
)

const defaultTimeout = 60                 // 函数执行的默认超时时间暂设置为较宽泛的值，以便用于一些扫描场景，后续考虑超时时间可配置
const maxLogLength = 4096                 // 日志字段的截断长度，防止用户生成过大的日志对数据库性能造成影响
const defaultGoroutinePoolPoolSize = 2000 // 函数执行器默认协程池大小

// Executor 脚本执行器接口
type Executor interface {
	Execute(faas.Context, string, map[string]interface{}) (interface{}, []string, error)
}

// CodeExecutor 执行器
type CodeExecutor struct {
	goroutinePool *ants.Pool
	executor      map[entity.LanguageType]Executor
	functionRepo  ports.FunctionRepository
}

// NewCodeExecutor 初始化执行器
func NewCodeExecutor(repoProviderSet *ports.RepoProviderSet) (*CodeExecutor, error) {
	// 如果协程池的协程已经满了，Submit 的时候不会再阻塞，会直接返回错误
	goroutinePool, err := ants.NewPool(defaultGoroutinePoolPoolSize, ants.WithNonblocking(true))
	if err != nil {
		return nil, err
	}
	return &CodeExecutor{
		goroutinePool: goroutinePool,
		executor: map[entity.LanguageType]Executor{
			entity.Golang: golang.NewGolangExecutor(),
			entity.Js:     js.NewJavascriptExecutor(),
		},
		functionRepo: repoProviderSet.FunctionRepo(),
	}, nil
}

// Execute 函数执行器
func (e *CodeExecutor) Execute(context faas.Context, req *dto.CallFunctionReqDTO, function *entity.Function) (
	result interface{}, record *dto.UpdateRunHistoryDTO, err error) {
	startTime := time.Now().UnixNano()
	record = &dto.UpdateRunHistoryDTO{
		Status: string(entity.Succeed),
	}
	var logs []string

	defer func() {
		// 执行器统一的 recover，panic 后将 panic 的信息更新到 log 中
		if r := recover(); r != nil {
			errorInfo := fmt.Sprintf("panic in execute caused by %s", r)
			err = fmt.Errorf(errorInfo)

			logs = append(logs, errorInfo)
			record.CostTime = getCostTimeOfMillisecond(startTime, time.Now().UnixNano())
			record.Status = string(entity.Failed)
			record.Log = subStr(strings.Join(logs, `\n`), maxLogLength)
		}
	}()

	// 执行函数
	result, logs, err = e.executeWithTimeOut(context, function, req.Input)
	endTime := time.Now().UnixNano()
	record.CostTime = getCostTimeOfMillisecond(startTime, endTime)
	// 执行报错后，将返回的 Error 记录到 logs 中
	if err != nil {
		record.Status = string(entity.Failed)
		logs = append(logs, fmt.Sprintf("ERROR: %s", err.Error()))
	}
	record.Log = subStr(strings.Join(logs, `\n`), maxLogLength)
	output, _ := json.Marshal(result)
	record.Output = string(output)

	return
}

// executeWithTimeOut 带超时时间控制的执行
// REF: https://github.com/zeromicro/go-zero
func (e *CodeExecutor) executeWithTimeOut(ctx faas.Context, function *entity.Function,
	input map[string]interface{}) (interface{}, []string, error) {
	ctxWithTimeout, cancel := c.WithTimeout(ctx.Context(), time.Second*defaultTimeout)
	defer cancel()

	var result interface{}
	var logs []string
	var err error

	var lock sync.Mutex
	done := make(chan struct{})
	// create channel with buffer size 1 to avoid goroutine leak
	panicChan := make(chan interface{}, 1)
	if err := e.goroutinePool.Submit(func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		lock.Lock()
		defer lock.Unlock()
		result, logs, err = e.execute(ctx, function, input)
		close(done)
	}); err != nil {
		return "", nil, fmt.Errorf("get function executor goroutine failed: %w", err)
	}

	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		lock.Lock()
		defer lock.Unlock()
		return result, logs, err
	case <-ctxWithTimeout.Done():
		err := ctxWithTimeout.Err()
		return "", nil, fmt.Errorf("execute function timeout: %w", err)
	}
}

// execute 函数执行
func (e *CodeExecutor) execute(context faas.Context, function *entity.Function,
	input map[string]interface{}) (interface{}, []string, error) {
	// 根据 language 获取对应的执行器
	executor, ok := e.executor[function.Language]
	if !ok {
		err := fmt.Errorf("invalid language %s", function.Language.String())
		return "", nil, err
	}

	result, logs, err := executor.Execute(context, function.Code, input)
	if err != nil {
		return "", logs, err
	}
	// 对返回值进行检查，如果不能转为 json 字符串则抛出错误，防止接口层 panic
	_, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		return "", logs, fmt.Errorf("the return of the funciton is invalid: %w", marshalErr)
	}
	return result, logs, nil
}

// Debug 函数 debug，debug 不需要记录执行历史
func (e *CodeExecutor) Debug(context faas.Context, function *entity.Function, input map[string]interface{}) (
	result interface{}, logs []string,
	err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic in execute: %s ", r)
		}
	}()

	return e.executeWithTimeOut(context, function, input)
}

// getCostTime 获取时间差，单位为 ms
func getCostTimeOfMillisecond(start int64, complete int64) int64 {
	costTime := complete - start
	if costTime < 0 {
		costTime = time.Now().UnixNano() - start
	}

	return costTime / 1e6
}

// subStr 截取字符串
func subStr(s string, l int) string {
	// rune 的截断性能会好一些
	if utf8.RuneCountInString(s) > l {
		rs := []rune(s)
		return string(rs[:l])
	}

	return s
}
