package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/entity"
	"github.com/jinzhu/copier"
)

var (
	AuthConvertor = &authConvertorImpl{} // 转换器
)

type authConvertorImpl struct{}

// ConvertEntityToVerifyCaptchaDTO 转换
func (*authConvertorImpl) ConvertEntityToVerifyCaptchaDTO(e *entity.User) (*dto.VerifyCaptchaRspDTO, error) {
	d := &dto.VerifyCaptchaRspDTO{}
	if err := copier.Copy(d, e); err != nil {
		return nil, err
	}
	return d, nil
}

// ConvertEntityToOauth2CallbackRspDTO 转换
func (*authConvertorImpl) ConvertEntityToOauth2CallbackRspDTO(e *entity.User) (*dto.Oauth2CallbackRspDTO, error) {
	d := &dto.Oauth2CallbackRspDTO{}
	if err := copier.Copy(d, e); err != nil {
		return nil, err
	}
	return d, nil
}

// ConvertNamespaceTokenEntityToDTO 转换
func (*authConvertorImpl) ConvertNamespaceTokenEntityToDTO(e *entity.NamespaceToken) *dto.NamespaceTokenDTO {
	return &dto.NamespaceTokenDTO{
		Namespace: e.Namespace,
		Name:      e.Name,
		Token:     e.Token,
		Creator:   e.Creator,
		CreatedAt: e.CreatedAt,
	}
}

// ConvertNamespaceTokenEntitiesToDTOs 转换
func (*authConvertorImpl) ConvertNamespaceTokenEntitiesToDTOs(es []*entity.NamespaceToken) []*dto.NamespaceTokenDTO {
	var dtos []*dto.NamespaceTokenDTO
	for _, e := range es {
		dtos = append(dtos, AuthConvertor.ConvertNamespaceTokenEntityToDTO(e))
	}
	return dtos
}
