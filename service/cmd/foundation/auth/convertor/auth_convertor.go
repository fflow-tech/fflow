package convertor

import (
	pb "github.com/fflow-tech/fflow/api/foundation/auth"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
)

var (
	// AuthConvertor 转换器
	AuthConvertor = &authConvertor{}
)

type authConvertor struct {
}

// ConvertValidateTokenPbToDTO 转换
func (*authConvertor) ConvertValidateTokenPbToDTO(req *pb.ValidateTokenReq) (*dto.ValidateTokenReqDTO, error) {
	return &dto.ValidateTokenReqDTO{
		Namespace:   req.BasicReq.Namespace,
		AccessToken: req.BasicReq.AccessToken,
	}, nil
}
