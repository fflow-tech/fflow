package redis

import (
	"os"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	redisclient "github.com/fflow-tech/fflow/service/pkg/redis"
)

var (
	pollingTask  *PollingTask
	timerDefDAO  *TimerDefDAO
	timerTaskDAO *TimerTaskDAO
	timeDuration = "2021-01-01 01-01"
	timerDefID   = "timerDefID"
)

// TestMain 测试初始化方法
func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		os.Exit(1)
		return
	}
	exitCode := m.Run()
	tearDown()
	os.Exit(exitCode)
}
func setup() error {
	if err := createDAOs(); err != nil {
		return err
	}
	if err := createData(); err != nil {
		return err
	}
	return nil
}
func createData() error {
	if err := createDefData(); err != nil {
		return err
	}
	return nil
}
func createDefData() error {
	timerDef := &dto.CreateTimerDefDTO{
		DefID: timerDefID,
	}
	return timerDefDAO.AddTimerDef(timerDef)
}
func tearDown() {
	deleteDefData()
}
func deleteDefData() {
	timerDef := &dto.DeleteTimerDefDTO{
		DefID: timerDefID,
	}
	timerDefDAO.DelTimerDef(timerDef)
	return
}

// createDAOs 创建 DAO 对象
func createDAOs() error {
	redisConfig := config.GetRedisConfig()
	client := redisclient.GetClient(redisConfig)

	pollingTask = NewPollingTaskClient(client)
	timerDefDAO = NewTimerDefClient(client)
	timerTaskDAO = NewTimerTaskClient(client, concurrency.GetDefaultWorkerPool())
	return nil
}

// TestPollingTask_GetBucketNum 获取桶数量
func TestPollingTask_GetBucketNum(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "test",
			want:    20,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pollingTask.GetBucketNum(); got != tt.want {
				t.Errorf("GetBucketNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPollingTask_GetTaskBucketID 获取桶ID
func TestPollingTask_GetTaskBucketID(t *testing.T) {
	type args struct {
		defID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				defID: "200",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pollingTask.GetTaskBucketID(tt.args.defID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskBucketID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestPollingTask_SetBucketNum  设置桶数
func TestPollingTask_SetBucketNum(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				num: 3,
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				num: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pollingTask.SetBucketNum(tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("SetBucketNum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPollingTask_SetTimeSlice 设置时间切片
func TestPollingTask_SetTimeSlice(t *testing.T) {
	type args struct {
		timeDuration string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				timeDuration: timeDuration,
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				timeDuration: timeDuration,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pollingTask.SetTimeSlice(tt.args.timeDuration); (err != nil) != tt.wantErr {
				t.Errorf("SetTimeSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err := pollingTask.deleteTimeSlice(timeDuration); err != nil {
		t.Errorf("deleteTimeSlice() error = %v", err)
	}
}

// TestPollingTask_GetTimeSlice 获取时间分测试
func TestPollingTask_GetTimeSlice(t *testing.T) {
	type args struct {
		timeDuration string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				timeDuration: timeDuration,
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				timeDuration: timeDuration,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pollingTask.GetTimeSlice(tt.args.timeDuration); (err != nil) != tt.wantErr {
				t.Errorf("GetTimeSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err := pollingTask.deleteTimeSlice(timeDuration); err != nil {
		t.Errorf("deleteTimeSlice() error = %v", err)
	}
}

// TestPollingTask_SuccessTimeSlice 时间分片测试
func TestPollingTask_SuccessTimeSlice(t *testing.T) {
	type args struct {
		timeDuration string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				timeDuration: timeDuration,
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				timeDuration: timeDuration,
			},
			wantErr: true,
		},
	}
	if err := pollingTask.SetTimeSlice(timeDuration); err != nil {
		t.Errorf("SetTimeSlice() error = %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pollingTask.SuccessTimeSlice(tt.args.timeDuration); (err != nil) != tt.wantErr {
				t.Errorf("SuccessTimeSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err := pollingTask.deleteTimeSlice(timeDuration); err != nil {
		t.Errorf("deleteTimeSlice() error = %v", err)
	}
}
