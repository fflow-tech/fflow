package execution

import (
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

// TestGetRetryDelay 测试获取重试延迟时间
func TestGetRetryDelay(t *testing.T) {
	type args struct {
		nodeInst *entity.NodeInst
		err      error
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{"固定5s的延时", args{nodeInst: &entity.NodeInst{
			BasicNodeDef: entity.BasicNodeDef{
				Retry: entity.Retry{
					Count:    1,
					Duration: "5s",
					Policy:   entity.Fixed,
				},
			},
			RetryCount: 1,
		}, err: nil}, 5 * time.Second, false},
		{"退避 5s的延时 第一次", args{nodeInst: &entity.NodeInst{
			BasicNodeDef: entity.BasicNodeDef{
				Retry: entity.Retry{
					Count:    2,
					Duration: "5s",
					Policy:   entity.ExponentialBackoff,
				},
			},
			RetryCount: 0,
		}, err: nil}, 5 * time.Second, false},
		{"退避 5s的延时 第二次", args{nodeInst: &entity.NodeInst{
			BasicNodeDef: entity.BasicNodeDef{
				Retry: entity.Retry{
					Count:    2,
					Duration: "5s",
					Policy:   entity.ExponentialBackoff,
				},
			},
			RetryCount: 1,
		}, err: nil}, 10 * time.Second, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRetryDelay(tt.args.nodeInst)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRetryDelay() err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getRetryDelay() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNeedRetry 测试是否需要重试
func TestNeedRetry(t *testing.T) {
	type args struct {
		nodeInst *entity.NodeInst
		err      error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"次数判断成功", args{nodeInst: &entity.NodeInst{
			Status: entity.NodeInstFailed,
			BasicNodeDef: entity.BasicNodeDef{
				Retry: entity.Retry{
					Count:    1,
					Duration: "5s",
					Policy:   entity.Fixed,
				},
			},
			RetryCount: 0,
		}, err: nil}, true},
		{"次数判断失败", args{nodeInst: &entity.NodeInst{
			Status: entity.NodeInstFailed,
			BasicNodeDef: entity.BasicNodeDef{
				Retry: entity.Retry{
					Count:    1,
					Duration: "5s",
					Policy:   entity.Fixed,
				},
			},
			RetryCount: 1,
		}, err: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needRetry(tt.args.nodeInst)
			if got != tt.want {
				t.Errorf("checkNodeRetry() got = %v, want %v", got, tt.want)
			}
		})
	}
}
