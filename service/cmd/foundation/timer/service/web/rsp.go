package web

import (
	"github.com/fflow-tech/fflow/service/pkg/errno"
)

// WebRsp Web返回
type WebRsp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewFailedWebRspWithMsg 新建失败返回
func NewFailedWebRspWithMsg(err *errno.BasicError, message string) WebRsp {
	return WebRsp{
		Code:    err.Code,
		Message: message,
	}
}

// NewFailedWebRsp 新建失败返回
func NewFailedWebRsp(err *errno.BasicError) WebRsp {
	return WebRsp{
		Code:    err.Code,
		Message: err.Message,
	}
}

// NewSucceedWebRsp 新建成功返回
func NewSucceedWebRsp(data interface{}) WebRsp {
	return WebRsp{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
		Data:    data,
	}
}
