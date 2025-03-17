package constants

import "github.com/fflow-tech/fflow-sdk-go/faas"

// FunctionType 函数的类型：第一个参数为 context，编写函数时可以获取到 context 中的值
type FunctionType = func(faas.Context, map[string]interface{}) (interface{}, error)

const (
	DefaultPackageName       = "p"                                                               // 包名
	DefaultEntryFunctionName = "handler"                                                         // 入口函数名
	FunctionTypeStr          = "func(faas.Context, map[string]interface{}) (interface{}, error)" // 函数类型
)
