package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

// Test_taskConvertorImpl_ConvertPageQueryDTOToPO 结构体转换
func Test_taskConvertorImpl_ConvertPageQueryDTOToPO(t *testing.T) {
	type args struct {
		d *dto.PageQueryRunHistoryDTO
	}
	tests := []struct {
		name string
		args args
		want *po.RunHistoryPO
	}{
		{
			name: "success",
			args: args{
				d: &dto.PageQueryRunHistoryDTO{
					DefID: "test",
				},
			},
			want: &po.RunHistoryPO{
				DefID: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &taskConvertorImpl{}
			if got := ta.ConvertPageQueryDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertPageQueryDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskConvertorImpl_ConvertCreateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.CreateRunHistoryDTO
	}
	tests := []struct {
		name string
		args args
		want *po.RunHistoryPO
	}{
		{
			name: "success",
			args: args{
				d: &dto.CreateRunHistoryDTO{
					DefID: "test",
				},
			},
			want: &po.RunHistoryPO{
				DefID: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &taskConvertorImpl{}
			if got, _ := ta.ConvertCreateDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertCreateDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskConvertorImpl_ConvertGetDTOToPO(t *testing.T) {
	type args struct {
		d *dto.GetRunHistoryDTO
	}
	tests := []struct {
		name string
		args args
		want *po.RunHistoryPO
	}{
		{
			name: "success",
			args: args{
				d: &dto.GetRunHistoryDTO{
					DefID: "test",
				},
			},
			want: &po.RunHistoryPO{
				DefID: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &taskConvertorImpl{}
			if got, _ := ta.ConvertGetDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertGetDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskConvertorImpl_ConvertDeleteDTOToPO(t *testing.T) {
	type args struct {
		d *dto.DeleteRunHistoryDTO
	}
	tests := []struct {
		name string
		args args
		want *po.RunHistoryPO
	}{
		{
			name: "success",
			args: args{
				d: &dto.DeleteRunHistoryDTO{
					DefID: "test",
				},
			},
			want: &po.RunHistoryPO{
				DefID: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &taskConvertorImpl{}
			if got, _ := ta.ConvertDeleteDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertDeleteDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskConvertorImpl_ConvertUpdateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.UpdateRunHistoryDTO
	}
	tests := []struct {
		name string
		args args
		want *po.RunHistoryPO
	}{
		{
			name: "success",
			args: args{
				d: &dto.UpdateRunHistoryDTO{
					DefID: "test",
				},
			},
			want: &po.RunHistoryPO{
				DefID: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &taskConvertorImpl{}
			if got, _ := ta.ConvertUpdateDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertUpdateDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}
