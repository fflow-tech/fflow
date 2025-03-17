package event

import "testing"

// TestGetEventInstID 测试
func TestGetEventInstID(t *testing.T) {
	type args struct {
		msg []byte
	}
	payload := `{"event_type":"NodeCompleteDriveEvent","router_value":"233",
    "def_id":233,"def_version":0,"inst_id":238,"node_inst_id":250}`

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Normal Case", args{msg: []byte(payload)}, "238", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEventInstID(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventInstID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEventInstID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
