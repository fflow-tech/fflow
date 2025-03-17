package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

func Test_appConvertor_ConvertCreateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.CreateAppDTO
	}
	tests := []struct {
		name string
		args args
		want *po.App
	}{
		{
			name: "success",
			args: args{
				d: &dto.CreateAppDTO{
					Name:    "test",
					Creator: "weixxxu",
				},
			},
			want: &po.App{
				Name:    "test",
				Creator: "weixxxu",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &appConvertor{}
			if got := de.ConvertCreateDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertCreateDTOToPO() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appConvertor_ConvertGetDTOToPO(t *testing.T) {
	type args struct {
		d *dto.GetAppDTO
	}
	tests := []struct {
		name string
		args args
		want *po.App
	}{
		{
			name: "success",
			args: args{
				d: &dto.GetAppDTO{
					Name: "test",
				},
			},
			want: &po.App{
				Name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &appConvertor{}
			if got := de.ConvertGetDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertGetDTOToPO() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appConvertor_ConvertDeleteDTOToPO(t *testing.T) {
	type args struct {
		d *dto.DeleteAppDTO
	}
	tests := []struct {
		name string
		args args
		want *po.App
	}{
		{
			name: "success",
			args: args{
				d: &dto.DeleteAppDTO{
					Name: "test",
				},
			},
			want: &po.App{
				Name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &appConvertor{}
			if got := de.ConvertDeleteDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertDeleteDTOToPO() got = %v, want %v", got, tt.want)
			}
		})
	}
}
