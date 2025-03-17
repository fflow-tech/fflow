package js

import (
	c "context"
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	context "github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/runtimecontext"
)

func TestJavascriptExecutor_Execute(t *testing.T) {
	const wrongEntryFuncSrc = `/**
	* Get the square number of x
	*
	* @param {context.Context: ctx} request context.
	* @param {map: input} the input map of the the function.
	*/
	function handler2(ctx, input) {
		return {"result": input["x"] * input["x"]};
	}
	`

	const src = `/**
	* Get the square number of x
	*
	* @param {context.Context: ctx} request context.
	* @param {map: input} the input map of the the function.
	*/
	const handler = function (ctx, input) {
		return {"result": input["x"] * input["x"]};
	}
	`

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
			name: "wrongEntryFuncSrc path",
			args: args{
				ctx:    context.NewRuntimeContext(c.Background(), &entity.Function{}, nil),
				code:   wrongEntryFuncSrc,
				params: map[string]interface{}{"x": 2},
			},
			want:    nil,
			want1:   []string{},
			wantErr: true,
		},
		{
			name: "happy path",
			args: args{
				ctx:    context.NewRuntimeContext(c.Background(), &entity.Function{}, nil),
				code:   src,
				params: map[string]interface{}{"x": 2},
			},
			want:    map[string]interface{}{"result": int64(4)},
			want1:   []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &javascriptExecutor{}
			got, got1, err := e.Execute(tt.args.ctx, tt.args.code, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("javascriptExecutor.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("javascriptExecutor.Execute() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("javascriptExecutor.Execute() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
