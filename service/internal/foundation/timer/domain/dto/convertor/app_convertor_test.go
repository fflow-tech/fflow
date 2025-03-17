package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

func Test_appConvertorImpl_ConvertAppEntitiesToDTOs(t *testing.T) {
	type args struct {
		e []*entity.App
	}
	tests := []struct {
		name    string
		args    args
		want    []*dto.App
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				e: []*entity.App{
					{Name: "test"},
				},
			},
			want: []*dto.App{
				{Name: "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &appConvertorImpl{}
			got, err := c.ConvertAppEntitiesToDTOs(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAppEntitiesToDTOs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertAppEntitiesToDTOs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appConvertorImpl_ConvertAppEntityToDTO(t *testing.T) {
	type args struct {
		e *entity.App
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.App
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				e: &entity.App{
					Name: "test",
				},
			},
			want: &dto.App{
				Name: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &appConvertorImpl{}
			got, err := c.ConvertAppEntityToDTO(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAppEntityToDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertAppEntityToDTO() got = %v, want %v", got, tt.want)
			}
		})
	}
}
