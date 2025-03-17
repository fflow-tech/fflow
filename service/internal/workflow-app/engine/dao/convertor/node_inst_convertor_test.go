package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
)

func Test_nodeInstConvertorImpl_ConvertCreateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.CreateNodeInstDTO
	}
	tests := []struct {
		name string
		args args
		want *po.NodeInstPO
	}{
		{"正常情况", args{d: &dto.CreateNodeInstDTO{}}, &po.NodeInstPO{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &nodeInstConvertorImpl{}
			if got, _ := no.ConvertCreateDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertCreateDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeInstConvertorImpl_ConvertDeleteDTOToPO(t *testing.T) {
	type args struct {
		d *dto.DeleteNodeInstDTO
	}
	tests := []struct {
		name string
		args args
		want *po.NodeInstPO
	}{
		{"正常情况", args{d: &dto.DeleteNodeInstDTO{}}, &po.NodeInstPO{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &nodeInstConvertorImpl{}
			if got, _ := no.ConvertDeleteDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertDeleteDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeInstConvertorImpl_ConvertGetDTOToPO(t *testing.T) {
	type args struct {
		d *dto.GetNodeInstDTO
	}
	tests := []struct {
		name string
		args args
		want *po.NodeInstPO
	}{
		{"正常情况", args{d: &dto.GetNodeInstDTO{}}, &po.NodeInstPO{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &nodeInstConvertorImpl{}
			if got, _ := no.ConvertGetDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertGetDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeInstConvertorImpl_ConvertPageQueryDTOToPO(t *testing.T) {
	type args struct {
		d *dto.PageQueryNodeInstDTO
	}
	tests := []struct {
		name string
		args args
		want *po.NodeInstPO
	}{
		{"正常情况", args{d: &dto.PageQueryNodeInstDTO{DefID: "100"}}, &po.NodeInstPO{DefID: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &nodeInstConvertorImpl{}
			if got, _ := no.ConvertPageQueryDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertPageQueryDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeInstConvertorImpl_ConvertUpdateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.UpdateNodeInstDTO
	}
	tests := []struct {
		name string
		args args
		want *po.NodeInstPO
	}{
		{"正常情况", args{d: &dto.UpdateNodeInstDTO{}}, &po.NodeInstPO{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &nodeInstConvertorImpl{}
			if got := no.ConvertUpdateDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertUpdateDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}
