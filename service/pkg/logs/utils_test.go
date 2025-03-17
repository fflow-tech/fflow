package logs

import "testing"

// TestGetFlowTraceID 测试
func TestGetFlowTraceID(t *testing.T) {
	type args struct {
		defID  string
		instID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"正常情况", args{
			defID:  "200",
			instID: "100",
		}, "$$200$$##100##"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFlowTraceID(tt.args.defID, tt.args.instID); got != tt.want {
				t.Errorf("GetFlowTraceIDByStrValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
