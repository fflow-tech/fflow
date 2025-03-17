package utils

import (
	"encoding/gob"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	gob.Register(testData)
	// 退出
	os.Exit(exitCode)
}

// TestJsonStrToMap 测试 json 字符串转换成 map
func TestJsonStrToMap(t *testing.T) {
	type args struct {
		j string
	}

	w1 := map[string]interface{}{}
	w1["hello"] = "world"
	w1["test"] = map[string]interface{}{"lala": "haha", "test": "1", "test2": "2"}

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"Case1", args{j: `{"hello":"world","test":{"lala":"haha", "test": "1", "test2":"2"}}`}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := JsonStrToMap(tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonStrToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCopyMap 测试拷贝 map
func TestCopyMap(t *testing.T) {
	type args struct {
		m1 map[string]interface{}
		m2 map[string]interface{}
	}
	m1 := map[string]interface{}{
		"Test": map[string]interface{}{
			"Test2": map[string]interface{}{
				"Test3": "test4",
			},
		}}

	m2 := map[string]interface{}{}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"Case1", args{
			m1: m1,
			m2: m2,
		}, m1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CopyMap(tt.args.m1, tt.args.m2)
			if !reflect.DeepEqual(tt.args.m2, tt.want) {
				t.Errorf("CopyMap() = %v, want %v", tt.args.m2, tt.want)
			}
		})
	}
}

// TestUintToStr 测试 uint 转换为字符串
func TestUintToStr(t *testing.T) {
	type args struct {
		i uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Case1", args{i: 100}, "100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UintToStr(tt.args.i); got != tt.want {
				t.Errorf("UintToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testType struct {
	Name string `structs:"name"`
	Age  string `structs:"Age"`
	// Score 忽略空字符串, 也就是说在转出的 map 里面不会有名为 score 的 key
	Score string `structs:"score,omitempty"`
}

// TestStructToMap 测试结构体转换成 map
func TestStructToMap(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"Case1", args{s: testType{
			Name:  "hunter",
			Age:   "18",
			Score: "",
		}}, map[string]interface{}{"Age": "18", "GlobalConfigName": "hunter", "Score": ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := StructToMap(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasAnyPrefix 测试是否包含前缀方法
func TestHasAnyPrefix(t *testing.T) {
	type args struct {
		str     string
		prefixs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Case1", args{
			str:     "test112",
			prefixs: []string{"ddd", "test11"},
		}, true},
		{"Case2", args{
			str:     "test112",
			prefixs: []string{"ddd", "test113"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAnyPrefix(tt.args.str, tt.args.prefixs); got != tt.want {
				t.Errorf("HasAnyPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseTemplate 测试解析模板
func TestParseTemplate(t *testing.T) {
	type args struct {
		text string
		m    map[string]interface{}
	}
	m1 := map[string]interface{}{}
	m1["Event"] = "haiden"
	m1["Age"] = 12
	m1["haha"] = "ddd"
	m1["Lala"] = true
	m1["test"] = map[string]interface{}{"Event": "hunter", "Age": 21}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Case1", args{text: "Event: {{.test.Event}}, Age: {{.test.Age}}, haha:{{.haha}}", m: m1}, "Event: hunter, Age: 21, haha:ddd"},
		{"Case2", args{text: `{{if eq $.test.Event "hunter"}}true{{else if ne $.test.Event "hunter"}}false{{end}}`, m: m1}, "true"},
		{"Case3", args{text: `{{if eq $.Lala true}}true{{end}}`, m: m1}, "true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTemplate(tt.args.text, tt.args.m); got != tt.want {
				t.Errorf("ParseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInterfaceToSimpleJson 测试接口转换成 simpleJson
func TestInterfaceToSimpleJson(t *testing.T) {
	type args struct {
		i interface{}
	}
	exampleOut, _ := simplejson.NewJson([]byte("123"))
	tests := []struct {
		name    string
		args    args
		want    *simplejson.Json
		wantErr bool
	}{
		{"happy path", args{"123"}, exampleOut, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InterfaceToSimpleJson(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("InterfaceToSimpleJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && *got == *tt.want {
				t.Errorf("InterfaceToSimpleJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBytesToInt 测试字节数组转换成整数
func TestBytesToInt(t *testing.T) {
	type args struct {
		bys []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"TestBytesToInt 传入合法 8Bytes 的参数", args{[]byte{0, 0, 0, 0, 0, 3, 18, 98}}, 201314, false},
		{"TestBytesToInt 传入合法非标准长度的参数", args{[]byte{3, 18, 98}}, 0, true},
		{"TestBytesToInt 传入过长的参数", args{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 3, 18, 98}}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BytesToInt(tt.args.bys)
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesToInt() err = %v, wantErr %v", got, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("BytesToInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIntToBytes 测试整数转换成字节数组
func TestIntToBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"case1", args{201314}, []byte{0, 0, 0, 0, 0, 3, 18, 98}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToBytes(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

type TestDataStruct struct {
	A string
	B int
	C bool
}

var testData = TestDataStruct{
	A: "1234567",
	B: 20210806,
	C: true,
}

// TestStructToBytes 测试结构体转换成字节数组
func TestStructToBytes(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"case1", args{testData}, false},
		{"case1", args{nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StructToBytes(tt.args.value)
			newValue, err := BytesToStruct(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.value, newValue) {
				t.Errorf("StructToBytes() got = %v, want %v", got, newValue)
			}
		})
	}
}

// TestGetCurrentGoroutineID 测试获取当前的协程 ID
func TestGetCurrentGoroutineID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				id1 := GetCurrentGoroutineID()
				id2 := GetCurrentGoroutineID()
				t.Logf("id1=%v, id2=%v", id1, id2)
				if id1 != id2 {
					t.Errorf("GetCurrentGoroutineID() = %v, want %v", id1, id2)
				}
			}
		})
	}
}

// TestGetNextTimeByExpr 测试
func TestGetNextTimeByExpr(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"正常情况", args{expr: "0 22 ? ? *"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNextTimeByExpr(tt.args.expr, time.Now())
			t.Logf("GetNextTimeByExpr got=%s", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNextTimeByExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetUInt64FromJson(t *testing.T) {
	type args struct {
		jsonBytes []byte
		branch    []string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"正常情况", args{
			jsonBytes: []byte(`{"event_type":"NodeFailEvent","router_value":"1516_1410","reason":"[14]type:framework, code:141, msg:tcp client transport ReadFrame: read tcp-: read: connection reset by peer, cost:5.002162441s","def_id":"1516","def_version":"1","inst_id":"1410","node":"t7","node_inst_id":"1772"}`),
			branch:    []string{"inst_id"},
		}, 1410, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUInt64FromJson(tt.args.jsonBytes, tt.args.branch...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUInt64FromJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUInt64FromJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStructToJsonStr 测试
func TestStructToJsonStr(t *testing.T) {
	type args struct {
		s interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"正常情况1", args{s: map[string]interface{}{"hello": "world"}}, `{"hello":"world"}`},
		{"正常情况2", args{s: &map[string]interface{}{"hello": "world"}}, `{"hello":"world"}`},
		{"正常情况3", args{s: map[string]interface{}{}}, `{}`},
		{"正常情况3", args{s: map[interface{}]interface{}{}}, `{}`},
		{"为nil的情况", args{s: nil}, `{}`},
		{"为空的情况", args{s: ""}, `{}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StructToJsonStr(tt.args.s); got != tt.want {
				t.Errorf("StructToJsonStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetIntervalTimeByCronExpr 测试获取定时间隔时间
func TestGetIntervalTimeByCronExpr(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{" 每小时的10分30秒触发任务", args{expr: "30 10 * * * ?"}, false},
		{" 每天1点10分30秒触发任务", args{expr: "30 10 1 * * ?"}, false},
		{"每月20号1点10分30秒触发任务", args{expr: "30 10 1 20 * ?"}, false},
		{"星期一到星期五的10点15分0秒触发任务", args{expr: "0 15 10 ? * MON-FRI"}, false},
		{"每个月最后一天的10点15分0秒触发任务", args{expr: "0 15 10 L * ?"}, false},
		{"每周二、四、六下午五点", args{expr: "0 0 17 ? * TUE,THU,SAT"}, false},
		{" 每天5-15点整点触发", args{expr: "0 0 5-15 * * ?"}, false},
		{" 每天上午10点，下午2点，4点", args{expr: "0 0 10,14,16 * * ?"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIntervalTimeByCronExpr(tt.args.expr)
			t.Logf("GetIntervalTimeByCronExpr got=%d", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIntervalTimeByCronExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAddElementsToSliceIfNotExists(t *testing.T) {
	type args struct {
		slice    []string
		elements []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "success",
			args: args{
				slice:    []string{"s1", "s2"},
				elements: []string{"s2", "s3"},
			},
			want: []string{"s1", "s2", "s3"},
		},
		{
			name: "success1",
			args: args{
				slice:    []string{"s1", "s2"},
				elements: []string{"s4", "s3"},
			},
			want: []string{"s1", "s2", "s4", "s3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddElementsToSliceIfNotExists(tt.args.slice, tt.args.elements...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddElementsToSliceIfNotExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveElementsFromSlice(t *testing.T) {
	type args struct {
		slice    []string
		elements []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "success",
			args: args{
				slice:    []string{"s1", "s2", "s3"},
				elements: []string{"s2", "s4"},
			},
			want: []string{"s1", "s3"},
		},
		{
			name: "success1",
			args: args{
				slice:    []string{"s1", "s2", "s3"},
				elements: []string{"s2"},
			},
			want: []string{"s1", "s3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteElementsFromSlice(tt.args.slice, tt.args.elements...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteElementsFromSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
