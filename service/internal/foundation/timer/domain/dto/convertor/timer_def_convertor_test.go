package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

// Test_defConvertorImpl_ConvertEntityToDTO 转换实体结构
func Test_defConvertorImpl_ConvertEntityToDTO(t *testing.T) {
	type args struct {
		e *entity.TimerDef
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.TimerDefDTO
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				e: &entity.TimerDef{
					Name:            "test",
					DefID:           "test",
					NotifyRpcParam:  "{}",
					NotifyHttpParam: "{}",
				},
			},
			want: &dto.TimerDefDTO{
				Name:  "test",
				DefID: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &defConvertorImpl{}
			got, err := c.ConvertEntityToDTO(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertEntityToDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertEntityToDTO() got = %v, want %v", got, tt.want)
			}
		})
	}
}
