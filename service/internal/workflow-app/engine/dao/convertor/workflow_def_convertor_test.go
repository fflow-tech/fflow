package convertor

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
)

func Test_defConvertorImpl_ConvertCreateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.CreateWorkflowDefDTO
	}
	tests := []struct {
		name string
		args args
		want *po.WorkflowDefPO
	}{
		{"正常情况", args{d: &dto.CreateWorkflowDefDTO{DefID: "100"}}, &po.WorkflowDefPO{DefID: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			if got, _ := de.ConvertCreateDTOToPO(tt.args.d); !reflect.DeepEqual(got.DefID, tt.want.DefID) {
				t.Errorf("ConvertCreateDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defConvertorImpl_ConvertDTOToPO(t *testing.T) {
	type args struct {
		d *dto.WorkflowDefDTO
	}
	tests := []struct {
		name string
		args args
		want *po.WorkflowDefPO
	}{
		{"正常情况", args{d: &dto.WorkflowDefDTO{}}, &po.WorkflowDefPO{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			if got, _ := de.ConvertDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defConvertorImpl_ConvertDeleteDTOToPO(t *testing.T) {
	type args struct {
		d *dto.DeleteWorkflowDefDTO
	}
	tests := []struct {
		name string
		args args
		want *po.WorkflowDefPO
	}{
		{"正常情况", args{d: &dto.DeleteWorkflowDefDTO{DefID: "100"}}, &po.WorkflowDefPO{DefID: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			if got, _ := de.ConvertDeleteDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertDeleteDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defConvertorImpl_ConvertGetDTOToPO(t *testing.T) {
	type args struct {
		d *dto.GetWorkflowDefDTO
	}
	tests := []struct {
		name string
		args args
		want *po.WorkflowDefPO
	}{
		{"正常情况", args{d: &dto.GetWorkflowDefDTO{DefID: "100"}}, &po.WorkflowDefPO{DefID: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			if got, _ := de.ConvertGetDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertGetDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defConvertorImpl_ConvertPageQueryDTOToPO(t *testing.T) {
	type args struct {
		d *dto.PageQueryWorkflowDefDTO
	}
	tests := []struct {
		name string
		args args
		want *po.WorkflowDefPO
	}{
		{"正常情况", args{d: &dto.PageQueryWorkflowDefDTO{DefID: "100"}}, &po.WorkflowDefPO{DefID: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			if got, _ := de.ConvertPageQueryDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertPageQueryDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defConvertorImpl_ConvertUpdateDTOToPO(t *testing.T) {
	type args struct {
		d *dto.UpdateWorkflowDefDTO
	}
	tests := []struct {
		name string
		args args
		want *po.WorkflowDefPO
	}{
		{"正常情况", args{d: &dto.UpdateWorkflowDefDTO{}}, &po.WorkflowDefPO{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de := &defConvertorImpl{}
			if got := de.ConvertUpdateDTOToPO(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertUpdateDTOToPO() = %v, want %v", got, tt.want)
			}
		})
	}
}
