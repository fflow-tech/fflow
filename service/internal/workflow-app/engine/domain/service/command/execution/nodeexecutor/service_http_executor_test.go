package nodeexecutor

import (
	"testing"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

func TestServiceHTTPNodeExecutor_validateArgs(t *testing.T) {
	type args struct {
		args *entity.HTTPArgs
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"URL 和 Method 没有的情况", args{args: &entity.HTTPArgs{}}, true},
		{"URL 没有的情况", args{args: &entity.HTTPArgs{Method: "test"}}, true},
		{"都有的情况", args{args: &entity.HTTPArgs{Method: "test", URL: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ServiceHTTPNodeExecutor{}
			if err := d.validateArgs(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("validateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
