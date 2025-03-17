package common

import (
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/stretchr/testify/suite"
)

// TestTimeoutSuite 测试流程定义命令
func TestTimeoutSuite(t *testing.T) {
	suite.Run(t, new(timeoutSuite))
}

type timeoutSuite struct {
	suite.Suite
	instTimeoutChecker *DefaultInstTimeoutChecker
}

// SetupTest 执行用例执行前准备工作
func (s *timeoutSuite) SetupTest() {
}

// TestCheckWorkflowInst 流程实例超时检测测试用例
func (s *timeoutSuite) TestCheckWorkflowInst() {
	type args struct {
		inst    *entity.WorkflowInst
		addTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantGot bool
		wantErr bool
	}{
		{
			name: "5s timeout",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "5s",
						},
					},
				},
				addTime: -5 * time.Second, // 减去5s
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "5s not timeout",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "5s",
						},
					},
				},
				addTime: -1 * time.Second, // 减去1s
			},
			wantGot: false,
			wantErr: false,
		},
		{
			name: "5m timeout",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "5m",
						},
					},
				},
				addTime: -6 * time.Minute, // 减去5m
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "5m not timeout",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "5m",
						},
					},
				},
				addTime: -1 * time.Minute, // 减去1m
			},
			wantGot: false,
			wantErr: false,
		},
		{
			name: "5h timeout",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "5h",
						},
					},
				},
				addTime: -6 * time.Hour, // 减去5h
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "5h not timeout",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "5h",
						},
					},
				},
				addTime: -1 * time.Hour, // 减去1h
			},
			wantGot: false,
			wantErr: false,
		},
		{
			name: "9w format err",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "9w",
						},
					},
				},
				addTime: -9 * 7 * 24 * time.Hour, // 减去9w
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "duration format err",
			args: args{
				inst: &entity.WorkflowInst{
					Status: entity.InstRunning,
					WorkflowDef: &entity.WorkflowDef{
						Timeout: entity.Timeout{
							Duration: "9w3",
						},
					},
				},
				addTime: -9 * 7 * 24 * time.Hour, // 减去9w
			},
			wantGot: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			now := time.Now()
			tt.args.inst.StartAt = now.Add(tt.args.addTime)
			got, err := s.instTimeoutChecker.CheckWorkflowInst(tt.args.inst)
			s.Equal(tt.wantErr, err != nil, "CheckWorkflowInst() got = %v, want %v", err != nil, tt.wantErr)
			if got != tt.wantGot {
				s.T().Errorf("CheckWorkflowInst() got != wantgot got:%v,wantgot:%v", got, tt.wantGot)
			}
		})
	}
}

// TestCheckNodeInst 节点实例超时检测测试用例
func (s *timeoutSuite) TestCheckNodeInst() {
	type args struct {
		inst    *entity.NodeInst
		addTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantGot bool
		wantErr bool
	}{
		{
			name: "5s timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "5s",
						},
					},
				},
				addTime: -5 * time.Second, // 减去5s
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "5s not timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "5s",
						},
					},
				},
				addTime: -1 * time.Second, // 减去1s
			},
			wantGot: false,
			wantErr: false,
		},
		{
			name: "5d timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "5d",
						},
					},
				},
				addTime: -6 * 24 * time.Hour, // 减去5d
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "5d not timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "5d",
						},
					},
				},
				addTime: -1 * 24 * time.Hour, // 减去1d
			},
			wantGot: false,
			wantErr: false,
		},
		{
			name: "超过最大超时时间",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "5w",
						},
					},
				},
				addTime: -6 * 7 * 24 * time.Hour, // 减去5w
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "超过最大超时时间",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "5w",
						},
					},
				},
				addTime: -1 * 7 * 24 * time.Hour, // 减去1w
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "9w format err",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "9w",
						},
					},
				},
				addTime: -9 * 7 * 24 * time.Hour, // 减去9w
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "duration format err",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Duration: "9w3",
						},
					},
				},
				addTime: -9 * 7 * 24 * time.Hour, // 减去9w
			},
			wantGot: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			now := time.Now()
			tt.args.inst.ScheduledAt = now.Add(tt.args.addTime)
			got, err := s.instTimeoutChecker.CheckNodeInst(tt.args.inst)
			s.Equal(tt.wantErr, err != nil, "CheckNodeInst() got = %v, want %v", err != nil, tt.wantErr)
			if got != tt.wantGot {
				s.T().Errorf("CheckNodeInst() got != wantgot got:%v,wantgot:%v", got, tt.wantGot)
			}
		})
	}
}

// TestCheckNodeInstExpr 节点实例超时检测测试用例
func (s *timeoutSuite) TestCheckNodeInstExpr() {
	type args struct {
		inst    *entity.NodeInst
		addTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantGot bool
		wantErr bool
	}{
		{
			name: "short time",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Expr: "*/10 * * * * * *", //这里为了测试使用
						},
					},
				},
				addTime: -11 * time.Second, // 减去5s
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "long timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Expr: "* * * * * */1 *", //这里为了测试使用
						},
					},
				},
				addTime: -2 * time.Second, // 减去2s
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "10 min",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Expr: "0 */10 * * * * *", //这里为了测试使用
						},
					},
				},
				addTime: -20 * time.Minute, // 减去20分钟
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "expr error",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							Expr: "*/10 12 a",
						},
					},
				},
				addTime: -2 * time.Second, // 减去2s
			},
			wantGot: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			now := time.Now()
			tt.args.inst.ScheduledAt = now.Add(tt.args.addTime)
			got, err := s.instTimeoutChecker.CheckNodeInst(tt.args.inst)
			if err != nil {
				log.Errorf("%s", err)
			}
			s.Equal(tt.wantErr, err != nil, "CheckNodeInst() got = %v, want %v", err != nil, tt.wantErr)
			if got != tt.wantGot {
				s.T().Errorf("CheckNodeInst() got != wantgot got:%v,wantgot:%v", got, tt.wantGot)
			}
		})
	}
}

// TestCheckNearNodeInst 节点实例接近超时检测测试用例
func (s *timeoutSuite) TestCheckNearNodeInst() {
	type args struct {
		inst    *entity.NodeInst
		addTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantGot bool
		wantErr bool
	}{
		{
			name: "5s timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "5s",
						},
					},
				},
				addTime: -5 * time.Second, // 减去5s
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "5s not timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "5s",
						},
					},
				},
				addTime: -1 * time.Second, // 减去1s
			},
			wantGot: false,
			wantErr: false,
		},
		{
			name: "5d timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "5d",
						},
					},
				},
				addTime: -6 * 24 * time.Hour, // 减去5d
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "5d not timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "5d",
						},
					},
				},
				addTime: -1 * 24 * time.Hour, // 减去1d
			},
			wantGot: false,
			wantErr: false,
		},
		{
			name: "超过最大超时时间",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "5w",
						},
					},
				},
				addTime: -6 * 7 * 24 * time.Hour, // 减去5w
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "超过最大超时时间",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "5w",
						},
					},
				},
				addTime: -1 * 7 * 24 * time.Hour, // 减去1w
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "9w format err",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "9w",
						},
					},
				},
				addTime: -9 * 7 * 24 * time.Hour, // 减去9w
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "duration format err",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutDuration: "9w3",
						},
					},
				},
				addTime: -9 * 7 * 24 * time.Hour, // 减去9w
			},
			wantGot: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			now := time.Now()
			tt.args.inst.ScheduledAt = now.Add(tt.args.addTime)
			got, err := s.instTimeoutChecker.CheckNodeInstNearTimeout(tt.args.inst)
			s.Equal(tt.wantErr, err != nil, "CheckNearNodeInst() got = %v, want %v", err != nil, tt.wantErr)
			if got != tt.wantGot {
				s.T().Errorf("CheckNearNodeInst() got != wantgot got:%v,wantgot:%v", got, tt.wantGot)
			}
		})
	}
}

// TestCheckNearNodeInstExpr 节点实例接近超时检测测试用例
func (s *timeoutSuite) TestCheckNearNodeInstExpr() {
	type args struct {
		inst    *entity.NodeInst
		addTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantGot bool
		wantErr bool
	}{
		{
			name: "node instance is terminal",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutExpr: "*/10 * * * * * *",
						},
					},
				},
				addTime: -5 * time.Second, // 减去 5s
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "short time",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutExpr: "*/10 * * * * * *",
						},
					},
				},
				addTime: -11 * time.Second, // 减去5s
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "long timeout",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutExpr: "* * * * * */1 *",
						},
					},
				},
				addTime: -2 * time.Second, // 减去2s
			},
			wantGot: false,
			wantErr: true,
		},
		{
			name: "1 min",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutExpr: "0 */10 * * * * *", //这里为了测试使用
						},
					},
				},
				addTime: -20 * time.Minute, // 减去两分钟
			},
			wantGot: true,
			wantErr: false,
		},
		{
			name: "expr error",
			args: args{
				inst: &entity.NodeInst{
					Status: entity.NodeInstRunning,
					BasicNodeDef: entity.BasicNodeDef{
						Timeout: entity.NodeTimeout{
							NearTimeoutExpr: "*/10 12 a",
						},
					},
				},
				addTime: -2 * time.Second, // 减去2s
			},
			wantGot: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			now := time.Now()
			tt.args.inst.ScheduledAt = now.Add(tt.args.addTime)
			got, err := s.instTimeoutChecker.CheckNodeInstNearTimeout(tt.args.inst)
			if err != nil {
				log.Errorf("%s", err)
			}
			s.Equal(tt.wantErr, err != nil, "CheckNearNodeInst() got = %v, want %v", err != nil, tt.wantErr)
			if got != tt.wantGot {
				s.T().Errorf("CheckNearNodeInst() got != wantgot got:%v,wantgot:%v", got, tt.wantGot)
			}
		})
	}
}
