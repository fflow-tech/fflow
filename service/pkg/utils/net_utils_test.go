package utils

import "testing"

// TestIsValidURL 测试是否是一个合法的 URL
func TestIsValidURL(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"正常情况", args{urlStr: "https://www.baidu.com"}, true},
		{"非http开头的情况", args{urlStr: "www.baidu.com"}, false},
		{"非正常字符结尾的情况", args{urlStr: "https://www.baidu.com\\"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidURL(tt.args.urlStr); got != tt.want {
				t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
