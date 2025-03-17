package utils

import "testing"

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Case1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(); len(got) <= 0 {
				t.Errorf("GenerateToken() = %v, want length > 0", got)
			}
		})
	}
}
