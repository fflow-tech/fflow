package convertor

import (
	pb "github.com/fflow-tech/fflow/api/foundation/faas"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	// FaasConvertor 转换器
	FaasConvertor = &faasConvertor{}
)

type faasConvertor struct {
}

// ConvertCallPbToDTO 转换
func (*faasConvertor) ConvertCallPbToDTO(req *pb.CallReq) (*dto.CallFunctionReqDTO, error) {
	ret := dto.CallFunctionReqDTO{
		Namespace: req.BasicReq.Namespace,
		Function:  req.Function,
	}
	var err error
	if ret.Input, err = utils.JsonStrToMap(req.Input); err != nil {
		return nil, err
	}

	return &ret, nil
}
