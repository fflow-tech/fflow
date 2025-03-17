package convertor

import (
	pb "github.com/fflow-tech/fflow/api/foundation/rbac"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/jinzhu/copier"
)

var (
	// RbacConvertor 转换器
	RbacConvertor = &rbacConvertor{}
)

type rbacConvertor struct {
}

// ConvertPbToDTO 转换
func (*rbacConvertor) ConvertPbToDTO(req *pb.RbacReq) (*dto.RbacReqDTO, error) {
	rbacDTO := &dto.RbacReqDTO{}
	if err := copier.Copy(rbacDTO, req); err != nil {
		return nil, err
	}
	return rbacDTO, nil
}
