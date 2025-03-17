package execution

import (
	"testing"
	"time"
)

// Test_getPollMaxDuration 测试
func Test_getPollMaxDuration(t *testing.T) {
	type args struct {
		durationStr string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"正常情况", args{durationStr: "30s"}, 30 * time.Second},
		{"大于最大值的情况", args{durationStr: "2d"}, 30 * time.Minute},
		{"配置非法的情况", args{durationStr: "**"}, pollDefaultMaxDuration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPollMaxDuration(tt.args.durationStr); got != tt.want {
				t.Errorf("getPollMaxDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_getPollInitialDuration 测试
func Test_getPollInitialDuration(t *testing.T) {
	type args struct {
		initialDurationStr string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"正常情况", args{initialDurationStr: "30s"}, 30 * time.Second},
		{"大于最大值的情况", args{initialDurationStr: "2d"}, 30 * time.Minute},
		{"配置非法的情况", args{initialDurationStr: "**"}, pollDefaultInitialDuration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPollInitialDuration(tt.args.initialDurationStr); got != tt.want {
				t.Errorf("getPollInitialDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
