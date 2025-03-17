package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

// Test_functionConvertor_ConvertEntityToGetDTO 获取执行历史转换
func Test_functionConvertor_ConvertEntityToGetDTO(t *testing.T) {
	type args struct {
		e *entity.RunHistory
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.GetRunHistoryRspDTO
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				e: &entity.RunHistory{
					DefID: "test",
				},
			},
			want: &dto.GetRunHistoryRspDTO{
				DefID: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &functionConvertor{}
			got, err := c.ConvertEntityToGetDTO(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertEntityToGetDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertEntityToGetDTO() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_functionConvertor_ConvertEntitiesToDTOs(t *testing.T) {
	type args struct {
		e []*entity.RunHistory
	}
	tests := []struct {
		name    string
		args    args
		want    []*dto.GetRunHistoryRspDTO
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				e: []*entity.RunHistory{
					{DefID: "test"},
				},
			},
			want: []*dto.GetRunHistoryRspDTO{
				{DefID: "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &functionConvertor{}
			got, err := c.ConvertEntitiesToDTOs(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertEntitiesToDTOs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertEntitiesToDTOs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
