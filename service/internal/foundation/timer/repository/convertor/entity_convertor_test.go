package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
)

// Test_defConvertorImpl_ConvertPOToEntity def po转换实体测试
func Test_defConvertorImpl_ConvertPOToEntity(t *testing.T) {
	type args struct {
		p *po.TimerDefPO
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.TimerDef
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				p: &po.TimerDefPO{
					DefID: "test",
				},
			},
			want: &entity.TimerDef{
				DefID: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			got, err := de.ConvertPOToEntity(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertPOToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertPOToEntity() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_taskConvertorImpl_ConvertPOToEntity po转换实体测试
func Test_taskConvertorImpl_ConvertPOToEntity(t *testing.T) {
	type args struct {
		p *po.RunHistoryPO
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.RunHistory
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				p: &po.RunHistoryPO{
					DefID: "test",
				},
			},
			want: &entity.RunHistory{
				DefID: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &taskConvertorImpl{}
			got, err := ta.ConvertPOToEntity(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertPOToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertPOToEntity() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskConvertorImpl_ConvertHistoryPOsToEntities(t *testing.T) {
	type args struct {
		p []*po.RunHistoryPO
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.RunHistory
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				p: []*po.RunHistoryPO{
					{DefID: "test"},
				},
			},
			want: []*entity.RunHistory{
				{DefID: "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &taskConvertorImpl{}
			got, err := ta.ConvertHistoryPOsToEntities(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertHistoryPOsToEntities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertHistoryPOsToEntities() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appConvertorImpl_ConvertAppPOsToEntities(t *testing.T) {
	type args struct {
		p []*po.App
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.App
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				p: []*po.App{
					{Name: "test"},
				},
			},
			want: []*entity.App{
				{Name: "test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &appConvertorImpl{}
			got, err := ta.ConvertAppPOsToEntities(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAppPOsToEntities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertAppPOsToEntities() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appConvertorImpl_ConvertAppPOToEntity(t *testing.T) {
	type args struct {
		p *po.App
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.App
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				p: &po.App{
					Name: "test",
				},
			},
			want: &entity.App{
				Name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &appConvertorImpl{}
			got, err := ta.ConvertAppPOToEntity(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAppPOToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertAppPOToEntity() got = %v, want %v", got, tt.want)
			}
		})
	}
}
