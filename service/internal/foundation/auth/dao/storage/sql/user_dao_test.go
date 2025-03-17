package sql

import (
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
)

var (
	Dsn = "root:123456@tcp(127.0.0.1:3306)/auth_test?charset=utf8mb4&parseTime=True&loc=Local"
)

func TestUserDAO_Create(t *testing.T) {
	type args struct {
		req *dto.CreateUserDTO
	}

	client, _ := mysql.GetClient(config.MySQLConfig{
		Dsn: Dsn,
	})
	dao := NewUserDAO(client)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"case1", args{req: &dto.CreateUserDTO{
				Username: "admin",
				NickName: "fflow-admin",
				Email:    "fflow@admin.com",
				Phone:    "13888888888",
				AuthType: entity.Password.String(),
				Avatar:   constants.DefaultAvatar,
				Status:   entity.Enabled.IntValue(),
			}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dao.Create(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserDAO_Get(t *testing.T) {
	type args struct {
		userDTO *dto.GetUserDTO
	}
	client, _ := mysql.GetClient(config.MySQLConfig{
		Dsn:                       Dsn,
		SlowThreshold:             200,
		IgnoreRecordNotFoundError: true,
	})
	dao := NewUserDAO(client)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"case1", args{userDTO: &dto.GetUserDTO{Username: "admin", Email: "fflow@admin.com"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.Get(tt.args.userDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Logf("Get() got = %v", got)
			}
		})
	}
}
