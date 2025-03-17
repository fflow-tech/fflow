package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"

	pb "github.com/fflow-tech/fflow/api/foundation/timer"
)

// Test_timerDefConvertor_ConvertCreatePbToDTO 测试创建定时器PB转DTO
func Test_timerDefConvertor_ConvertCreatePbToDTO(t *testing.T) {
	type args struct {
		req *pb.CreateTimerReq
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.CreateTimerDefDTO
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				req: &pb.CreateTimerReq{
					Name: "test",
				},
			},
			want: &dto.CreateTimerDefDTO{
				Name: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := &timerDefConvertor{}
			got, err := ti.ConvertCreatePbToDTO(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertCreatePbToDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertCreatePbToDTO() got = %v, want %v", got, tt.want)
			}
		})
	}
}
