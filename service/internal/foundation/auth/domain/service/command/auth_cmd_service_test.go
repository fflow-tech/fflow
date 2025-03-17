package command

import (
	"context"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"testing"
)

func TestAuthCommandService_Oauth2Callback(t *testing.T) {
	type args struct {
		ctx context.Context
		req *dto.Oauth2CallbackReqDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"case1", args{
				ctx: context.Background(),
				req: &dto.Oauth2CallbackReqDTO{Code: "3aa2268d5c7651dd2b0b"},
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AuthCommandService{}
			_, err := m.Oauth2Callback(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Oauth2Callback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_generateCaptcha(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Case1", args{n: 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateCaptcha(tt.args.n)
			if len(got) != 4 {
				t.Errorf("generateCaptcha() = %v", got)
			}
		})
	}
}
