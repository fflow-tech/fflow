package convertor

import (
	"strings"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/remote"
)

var (
	AbilityArgsConvertor = &abilityArgsConvertorImpl{} // 转换器
)

type abilityArgsConvertorImpl struct {
}

// getProtocol 获取 protocol
func (a *abilityArgsConvertorImpl) getProtocol(p string) string {
	protocol := entity.HTTPService
	return strings.ToLower(string(protocol))
}

// ConvertEntityToCallHTTPDTO 转换，参数未对齐，手动转换
func (*abilityArgsConvertorImpl) ConvertEntityToCallHTTPDTO(e *entity.HTTPArgs) *remote.CallHTTPReqDTO {
	return &remote.CallHTTPReqDTO{
		Body:     e.Body,
		URL:      e.URL,
		Method:   e.Method,
		Header:   e.Headers,
		Query:    e.Parameters,
		MockMode: e.MockMode,
	}
}

// ConvertEntityToCallFAASDTO 转换，参数未对齐，手动转换
func (*abilityArgsConvertorImpl) ConvertEntityToCallFAASDTO(e *entity.FAASArgs) *remote.CallFAASReqDTO {
	return &remote.CallFAASReqDTO{
		Namespace: e.Namespace,
		Function:  e.Func,
		Body:      e.Body,
		MockMode:  e.MockMode,
	}
}
