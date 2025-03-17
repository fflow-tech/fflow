package redis

import (
	"reflect"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
)

// TestTimerTaskDAO_GetTimerTask 测试获取定时器任务
func TestTimerTaskDAO_GetTimerTask(t *testing.T) {
	type args struct {
		d *dto.GetTimerTaskDTO
	}
	startTime, err := time.ParseInLocation(dto.TimerTriggerTimeFormat, "2021-01-01 01:01:01", time.Local)
	if err != nil {
		t.Errorf("ParseInLocation() got = %v", err)
	}
	endTime, err := time.ParseInLocation(dto.TimerTriggerTimeFormat, "2021-01-01 01:02:01", time.Local)
	if err != nil {
		t.Errorf("ParseInLocation() got = %v", err)
	}
	tests := []struct {
		name    string
		want    []string
		wantErr bool
		args    args
	}{
		{
			name:    "success",
			want:    []string{"123"},
			wantErr: false,
			args: args{
				d: &dto.GetTimerTaskDTO{
					BucketTime: "1_2021-01-01 01:01",
					StartTime:  startTime,
					EndTime:    endTime,
				},
			},
		},
	}
	addDTO := &dto.AddTimerTaskDTO{
		BucketID:  "1",
		HashID:    "123",
		TimerTime: startTime,
	}
	if err := timerTaskDAO.AddTimerTask(addDTO); err != nil {
		t.Errorf("AddTimerTask() error = %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timerTaskDAO.GetTimerTasks(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimerTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTimerTask() got = %v, want %v", got, tt.want)
			}
		})
	}
	if err := timerTaskDAO.DelTimerTask(&dto.DelTimerTaskDTO{BucketTime: "2_2021-01-01 01-01", HashID: "123"}); err != nil {
		t.Errorf("AddTimerTask() error = %v", err)
	}
}

// TestTimerTaskDAO_GetNotTriggeredTimers 获取没有触发的定时器
func TestTimerTaskDAO_GetNotTriggeredTimers(t *testing.T) {
	type args struct {
		bucketTime string
	}
	startTime, err := time.ParseInLocation(dto.TimerTriggerTimeFormat, "2021-01-01 01:01:01", time.Local)
	if err != nil {
		t.Errorf("ParseInLocation() got = %v", err)
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "success",
			want:    []string{"123"},
			wantErr: false,
			args: args{
				bucketTime: "1_2021-01-01 01:01",
			},
		},
	}
	addDTO := &dto.AddTimerTaskDTO{
		BucketID:  "1",
		HashID:    "123",
		TimerTime: startTime,
	}
	if err := timerTaskDAO.AddTimerTask(addDTO); err != nil {
		t.Errorf("AddTimerTask() error = %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timerTaskDAO.GetNotTriggeredTimers(tt.args.bucketTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotTriggeredTimers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotTriggeredTimers() got = %v, want %v", got, tt.want)
			}
		})
	}
	if err := timerTaskDAO.DelTimerTask(&dto.DelTimerTaskDTO{BucketTime: "2_2021-01-01 01-01", HashID: "123"}); err != nil {
		t.Errorf("AddTimerTask() error = %v", err)
	}
}
