package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/entity"
)

func Test_authConvertorImpl_ConvertEntityToVerifyCaptchaDTO(t *testing.T) {
	type args struct {
		e *entity.User
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.VerifyCaptchaRspDTO
		wantErr bool
	}{
		{"case1", args{e: &entity.User{
			NickName: "example",
			Username: "example",
			Email:    "admin@example.com",
			Avatar:   "test",
			Status:   entity.UserStatus{},
		}}, &dto.VerifyCaptchaRspDTO{
			NickName: "example",
			Username: "example",
			Email:    "admin@example.com",
			Avatar:   "test",
		}, false},
		{"case2", args{e: &entity.User{
			NickName: "example",
			Username: "example",
			Email:    "admin@example.com",
			Avatar:   "test",
			Status:   entity.UserStatus{},
		}}, &dto.VerifyCaptchaRspDTO{
			NickName: "example",
			Username: "example",
			Email:    "admin@example.com",
			Avatar:   "test",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := &authConvertorImpl{}
			got, err := au.ConvertEntityToVerifyCaptchaDTO(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertEntityToVerifyCaptchaDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertEntityToVerifyCaptchaDTO() got = %v, want %v", got, tt.want)
			}
		})
	}
}
