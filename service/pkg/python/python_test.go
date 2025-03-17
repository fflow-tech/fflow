package python

import "testing"

func TestExecutor_Execute(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test", args{"testdata/test.py"}, "hello world, this is python script hello\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Executor{}
			got, err := e.Execute(tt.args.script, "hello")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
