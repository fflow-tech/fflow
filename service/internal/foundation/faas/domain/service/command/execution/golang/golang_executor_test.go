package golang

import (
	c "context"
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	context "github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/runtimecontext"
)

func Test_golangExecutor_Execute(t *testing.T) {
	const src = `// Package p
	package p

	import	(
		"fmt"
		"github.com/fflow-tech/fflow-sdk-go/faas"
	)

	func handler(ctx faas.Context, s map[string]interface{}) (interface{}, error) {
		return map[string]interface{}{"result": s["key"].(string) + "-FOO"}, nil
	}`

	type args struct {
		ctx    context.RuntimeContext
		code   string
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		want1   []string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				ctx:    context.NewRuntimeContext(c.Background(), &entity.Function{}, nil),
				code:   src,
				params: map[string]interface{}{"key": "bar"},
			},
			want:    map[string]interface{}{"result": "bar-FOO"},
			want1:   []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &golangExecutor{}
			got, got1, err := e.Execute(tt.args.ctx, tt.args.code, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("golangExecutor.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("golangExecutor.Execute() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("golangExecutor.Execute() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
