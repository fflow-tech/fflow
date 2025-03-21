package expr

import (
	"reflect"
	"testing"
	"time"
)

// TestExpressionEvaluator_Match 测试是否匹配
func TestExpressionEvaluator_Match(t *testing.T) {
	type args struct {
		ctx  map[string]interface{}
		expr string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"Case1", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				expr: `${w.i.ww == "vv"}`,
			}, true, false,
		},
		{
			"Case1", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				expr: "$w.i.ww == 'vv'}",
			}, false, true,
		},
		{
			"Case3", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"poll_output": map[string]interface{}{"data": map[string]interface{}{
						"status": 1,
					}},
				}},
				expr: "${w.poll_output.data.status == 1}",
			}, true, false,
		},
		{
			"Case4", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"variables": map[string]interface{}{
						"projectName":            "测试项目",
						"test":                   1,
						"workflow_definition_id": 90,
						"workflow_instance_id":   1197,
					},
				},
				},
				expr: "${w.variables.n == 1}",
			}, false, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DefaultEvaluator{}
			got, err := c.Match(tt.args.ctx, tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Match() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultExpressionEvaluator_IsExpression 测试是否为正确的表达式
func TestDefaultExpressionEvaluator_IsExpression(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Case1", args{expr: "${dddd}"}, true,
		},
		{
			"Case2", args{expr: "${dddd"}, false,
		},
		{
			"Case3", args{expr: "{dddd}"}, false,
		},
		{
			"Case4", args{expr: "$dddd}"}, false,
		},
		{
			"Case5", args{expr: "dddd"}, false,
		},
		{
			"Case6", args{expr: "${dddd}dd"}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DefaultEvaluator{}
			if got := c.IsExpression(tt.args.expr); got != tt.want {
				t.Errorf("IsExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultExpressionEvaluator_Evaluate 表达式计算
func TestDefaultExpressionEvaluator_Evaluate(t *testing.T) {
	type args struct {
		ctx  map[string]interface{}
		expr string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Case1", args{
			ctx: map[string]interface{}{
				"test": "test",
			},
			expr: `${test == "test"}`,
		}, true, false},
		{"Case2", args{
			ctx: map[string]interface{}{
				"test": "test",
			},
			expr: `${test + "test"}`,
		}, "testtest", false},
		{"Case3", args{
			ctx: map[string]interface{}{
				"test": "test",
				"k2": map[string]interface{}{
					"k3": "v3",
				},
			},
			expr: `${k2.k3 + "test"}`,
		}, "v3test", false},
		{"Case4", args{
			ctx: map[string]interface{}{
				"test": "test",
				"k2": map[string]interface{}{
					"k3": "v3",
					"k4": 4,
				},
			},
			expr: `${k2.k4 + 5}`,
		}, float64(9), false},
		{"三目表达式", args{
			ctx: map[string]interface{}{
				"test": "test",
				"k2": map[string]interface{}{
					"k3": "v3",
					"k4": 4,
				},
			},
			expr: `${(k2.k4 + 5) == 9 ? true : false }`,
		}, true, false},
		{"当前时间1", args{expr: `${curtimeformat("0102")}`}, time.Now().Format("0102"), false},
		{"字符串格式化1", args{
			ctx: map[string]interface{}{
				"test": "dd",
				"k2": map[string]interface{}{
					"k3": "v3",
					"k4": 4,
				},
			},
			expr: `${sprintf("%s_%s",test,k2.k3)}`,
		}, "dd_v3", false},
		{"字符串格式化2", args{
			ctx: map[string]interface{}{
				"test": "dd",
				"k2": map[string]interface{}{
					"k3": "v3",
					"k4": 4,
				},
			},
			expr: `${sprintf("%s_%d",test,k2.k4)}`,
		}, "dd_4", false},
		{"字符串格式化3", args{
			ctx: map[string]interface{}{
				"test": "dd",
				"k2": map[string]interface{}{
					"k3": "v3",
					"k4": 4,
				},
			},
			expr: `${sprintf("{\"flow_name\":\"%s\",\"instance_id\":\"%s\"}", test, k2.k3)}`,
		}, "{\"flow_name\":\"dd\",\"instance_id\":\"v3\"}", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DefaultEvaluator{}
			got, err := c.Evaluate(tt.args.ctx, tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExpressionEvaluator_ReplaceExpressionMap 测试
func TestExpressionEvaluator_ReplaceExpressionMap(t *testing.T) {
	type args struct {
		ctx     map[string]interface{}
		exprMap map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			"string", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				exprMap: map[string]interface{}{"w": `${w.i.ww == "vv"}`},
			}, map[string]interface{}{"w": true}, false,
		},
		{
			"[]interface", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				exprMap: map[string]interface{}{"w": map[string]interface{}{
					"i": []interface{}{`${w.i.ww == "vv"}`, `${w.i.ww == "tt"}`}}},
			}, map[string]interface{}{"w": map[string]interface{}{
				"i": []interface{}{true, false},
			}}, false,
		},
		{
			"map[string]interface", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				exprMap: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": `${w.i.ww == "vv"}`},
				}},
			}, map[string]interface{}{"w": map[string]interface{}{
				"i": map[string]interface{}{"ww": true},
			}}, false,
		},
		{
			"int", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				exprMap: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": 12, "tt": "ss"},
				}},
			}, map[string]interface{}{"w": map[string]interface{}{
				"i": map[string]interface{}{"ww": 12, "tt": "ss"},
			}}, false,
		},
		{
			"expr error", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				exprMap: map[string]interface{}{"w": `$w.i.ww == "vv"}`},
			}, map[string]interface{}{"w": `$w.i.ww == "vv"}`}, false,
		},
		{
			"not found", args{
				ctx: map[string]interface{}{"w": map[string]interface{}{
					"i": map[string]interface{}{"ww": "vv"},
				}},
				exprMap: map[string]interface{}{"w": `${w.i.tt == "vv"}`},
			}, map[string]interface{}{"w": false}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DefaultEvaluator{}
			got, err := c.EvaluateMap(tt.args.ctx, tt.args.exprMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceExpressionMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("ReplaceExpressionMap() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_curtimeformat(t *testing.T) {
	c := &DefaultEvaluator{}
	got, err := c.Evaluate(map[string]interface{}{}, `${curtimeformat("2006-01-02 15:04:05")}`)
	if err != nil {
		t.Errorf("Evaluate() error = %v", err)
		return
	}
	t.Logf("got = %v", got)
}
