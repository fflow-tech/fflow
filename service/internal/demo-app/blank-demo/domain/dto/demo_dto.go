package dto

import "github.com/fflow-tech/fflow/service/pkg/constants"

// StartCollectReqDTO 开请
type StartCollectReqDTO struct {
	URLs []string `json:"urls,omitempty"`
}

// StartCollectRspDTO 登录结果
type StartCollectRspDTO struct {
	constants.WebRsp
}
