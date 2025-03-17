package redis

import (
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

var TestDefID = "test_def"

func TestTimerDefDAO_AddTimerDef(t1 *testing.T) {
	type args struct {
		def *dto.CreateTimerDefDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				def: &dto.CreateTimerDefDTO{
					DefID: TestDefID,
				},
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				def: &dto.CreateTimerDefDTO{
					DefID: TestDefID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if err := timerDefDAO.AddTimerDef(tt.args.def); (err != nil) != tt.wantErr {
				t1.Errorf("AddTimerDef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err := timerDefDAO.DelTimerDef(&dto.DeleteTimerDefDTO{DefID: TestDefID}); err != nil {
		t1.Errorf("DelTimerDef() error = %v", err)
	}
}

func TestTimerDefDAO_GetTimerDef(t1 *testing.T) {
	type args struct {
		d *dto.GetTimerDefDTO
	}
	tests := []struct {
		name    string
		args    args
		want    *po.TimerDefPO
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				d: &dto.GetTimerDefDTO{
					DefID: timerDefID,
				},
			},
			want: &po.TimerDefPO{
				DefID: timerDefID,
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				d: &dto.GetTimerDefDTO{
					DefID: "fail",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			_, err := timerDefDAO.GetTimerDef(tt.args.d)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetTimerDef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
