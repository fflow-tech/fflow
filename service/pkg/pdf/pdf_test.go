package pdf

import (
	"bufio"
	"context"
	"io"
	"os"
	"testing"
)

func TestExtractor_Extract(t *testing.T) {
	type fields struct {
		url string
	}
	type args struct {
		ctx   context.Context
		input io.Reader
	}
	file, _ := os.Open("testdata/test.pdf")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Case 1", fields{"https://www.fflow.link/tika"}, args{context.Background(),
			bufio.NewReader(file)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Extractor{
				url: tt.fields.url,
			}
			got, err := e.Extract(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
