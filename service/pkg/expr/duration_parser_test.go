package expr

import (
	"testing"
	"time"
)

// TestGetParseDuration 获取持续时间
func TestGetParseDuration(t *testing.T) {
	type args struct {
		duration string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "1s",
			args: args{
				duration: "1s",
			},
			want:    time.Second,
			wantErr: false,
		},
		{
			name: "2m",
			args: args{
				duration: "2m",
			},
			want:    2 * time.Minute,
			wantErr: false,
		},
		{
			name: "3h",
			args: args{
				duration: "3h",
			},
			want:    3 * time.Hour,
			wantErr: false,
		},
		{
			name: "4d",
			args: args{
				duration: "4d",
			},
			want:    4 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name: "1d",
			args: args{
				duration: "1d",
			},
			want:    24 * time.Hour,
			wantErr: false,
		},
		{
			name: "5w",
			args: args{
				duration: "5w",
			},
			want:    5 * 7 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name: "format error",
			args: args{
				duration: "52tw",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
