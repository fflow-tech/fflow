package monitor

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/redis"

	"github.com/agiledragon/gomonkey/v2"
)

type mockReporter struct{}

func (m *mockReporter) ReportTimerEnabledTotalNum(total float64) {}

func (m *mockReporter) ReportTimerFailedNum(total float64) {}

type mockTimerCounter struct{}

func (m *mockTimerCounter) CountTimersByStatus(status entity.TimerDefStatus) (int64, error) {
	return 0, nil
}

func (m *mockTimerCounter) CountPendingTimers(curTime time.Time) (int, error) {
	return 0, nil
}

type mockWorkerPool struct{}

func (m *mockWorkerPool) Submit(f func()) error {
	f()
	return nil
}

type mockLockProvider struct{}

func (m *mockLockProvider) GetDistributeLock(name string, expireTime time.Duration) *redis.DefaultDistributeLock {
	return &redis.DefaultDistributeLock{}
}

func Test_Monitor_ReportRecord(t *testing.T) {
	patch := gomonkey.ApplyMethod(reflect.TypeOf(&redis.DefaultDistributeLock{}), "Lock", func(*redis.DefaultDistributeLock) error {
		return nil
	})
	defer patch.Reset()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "success",
		},
	}

	mockMonitor := &Monitor{
		redisClient: &mockLockProvider{},
		pool:        &mockWorkerPool{},
		counter:     &mockTimerCounter{},
		reporter:    &mockReporter{},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockMonitor.ReportRecord(ctx); (err != nil) != tt.wantErr {
				t.Errorf("ReportRecord() err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}
