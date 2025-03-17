package nodeexecutor

import (
	"testing"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

func TestServiceFAASNodeExecutor_validateArgs(t *testing.T) {
	type args struct {
		args *entity.FAASArgs
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Namespace 没有的情况", args{args: &entity.FAASArgs{}}, true},
		{"Func 没有的情况", args{args: &entity.FAASArgs{Namespace: "test"}}, true},
		{"都有的情况",
			args{args: &entity.FAASArgs{Namespace: "test", Func: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ServiceFAASNodeExecutor{}
			err := d.validateArgs(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
