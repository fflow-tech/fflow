package constants

import (
	"fmt"

	"github.com/fflow-tech/fflow/service/pkg/errno"
)

// WebRsp Web 返回
type WebRsp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
}

// NewFailedWebRspWithMsg 初始化失败返回
func NewFailedWebRspWithMsg(err *errno.BasicError, message string) WebRsp {
	return WebRsp{
		Code:    err.Code,
		Message: message,
	}
}

// NewFailedWebRsp 初始化失败返回
func NewFailedWebRsp(err *errno.BasicError) WebRsp {
	return WebRsp{
		Code:    err.Code,
		Message: err.Message,
	}
}

// NewSucceedWebRsp 初始化成功返回
func NewSucceedWebRsp(data interface{}) WebRsp {
	switch data.(type) {
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
		return WebRsp{
			Code:    errno.OK.Code,
			Message: errno.OK.Message,
			Data:    fmt.Sprintf("%d", data),
		}
	default:
		return WebRsp{
			Code:    errno.OK.Code,
			Message: errno.OK.Message,
			Data:    data,
		}
	}
}

// NewSucceedWebRspWithTotal 初始化成功返回
func NewSucceedWebRspWithTotal(data interface{}, total int64) WebRsp {
	return WebRsp{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
		Data:    data,
		Total:   total,
	}
}
