package convertor

import (
	"reflect"
	"testing"

	"gorm.io/gorm"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

// Test_defConvertorImpl_ConvertCreateDTOToPO 测试转换方法
func Test_defConvertorImpl_ConvertCreateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.CreateTimerDefDTO
	}
	tests := []struct {
		name    string
		args    args
		want    *po.TimerDefPO
		wantErr bool
	}{
		{
			name: "success",
			args: args{d: &dto.CreateTimerDefDTO{
				DefID: "test",
				Name:  "test",
				NotifyRpcParam: dto.NotifyRpcParam{
					Service: "test",
				},
			}},
			wantErr: false,
			want: &po.TimerDefPO{
				DefID:           "test",
				Name:            "test",
				NotifyRpcParam:  "{\"service\":\"test\"}",
				NotifyHttpParam: "{}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			got, err := de.ConvertCreateDTOToPO(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertCreateDTOToPO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertCreateDTOToPO() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_defConvertorImpl_ConvertDeleteDTOToPO 测试转换方法
func Test_defConvertorImpl_ConvertDeleteDTOToPO(t *testing.T) {
	type args struct {
		d *dto.DeleteTimerDefDTO
	}
	tests := []struct {
		name    string
		args    args
		want    *po.TimerDefPO
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				d: &dto.DeleteTimerDefDTO{
					DefID: "1",
				},
			},
			want: &po.TimerDefPO{
				Model: gorm.Model{ID: 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			got, err := de.ConvertDeleteDTOToPO(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertDeleteDTOToPO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertDeleteDTOToPO() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defConvertorImpl_ConvertUpdateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.UpdateTimerDefDTO
	}
	tests := []struct {
		name    string
		args    args
		want    *po.TimerDefPO
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				d: &dto.UpdateTimerDefDTO{
					DefID: "1",
				},
			},
			want: &po.TimerDefPO{
				DefID: "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			got, err := de.ConvertUpdateDTOToPO(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertUpdateDTOToPO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertUpdateDTOToPO() got = %v, want %v", got, tt.want)
			}
		})
	}
}
