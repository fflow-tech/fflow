package modules

import (
	"github.com/fflow-tech/fflow-sdk-go/faas"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/js/modules/console"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution/js/modules/storage"
)

// InitModules 初始化 golang 的方法功 js 调用，注意调用时方法名都是小写字母开头
func InitModules(ctx faas.Context) map[string]interface{} {
	var Symbols = map[string]interface{}{
		"console": console.New(ctx),
		"storage": storage.New(ctx),
	}

	return Symbols
}
