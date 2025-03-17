package execution

import (
	"reflect"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/expr"
)

// TestDefaultWorkflowDecider_BasicDecide 测试
func TestDefaultWorkflowDecider_BasicDecide(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"流程开始启动的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
			}},
			[]string{"t1"},
			false,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2"},
			false,
		},
		{
			"使用next正常配置下一个节点为end的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1", Next: "end"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "end"},
				},
			}},
			[]string{},
			false,
		},
		{
			"使用next正常配置下一个节点的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1", Next: "t2"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2"},
				},
			}},
			[]string{"t2"},
			false,
		},
		{
			"存在引用节点的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "SERVICE",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "REF",
							"ref":  "t1",
							"next": "t3",
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1", Next: "t2"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2"},
				},
			}},
			[]string{"t2"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_SwitchDecide 测试
func TestDefaultWorkflowDecider_SwitchDecide(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"当前节点是Switch,且Switch满足条件的分支下一个节点是end的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "SWITCH",
							"switch": []map[string]interface{}{
								{
									"condition": `${"true" == "true"}`,
									"next":      "t1",
								},
								{
									"condition": `${"true" == "true"}`,
									"next":      "end",
								},
							},
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.SwitchNodeDef{
						Switch: []entity.SwitchCase{
							{Condition: `${"true" == "true"}`, Next: "t1"},
							{Condition: `${"true" == "false"}`, Next: "end"},
						},
						BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "end", Type: entity.SwitchNode},
					},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "end", Type: entity.SwitchNode},
				},
			}},
			[]string{"t1"},
			false,
		},
		{
			"Switch下一个节点是end的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
						{"t2": map[string]interface{}{
							"type": "SWITCH",
							"switch": []map[string]interface{}{
								{
									"condition": `${"true" == "false"}`,
									"next":      "t1",
								},
								{
									"condition": `${"true" == "false"}`,
									"next":      "end",
								},
							},
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.SwitchNodeDef{
						Switch: []entity.SwitchCase{
							{Condition: `${"true" == "false"}`, Next: "t1"},
							{Condition: `${"true" == "true"}`, Next: "end"},
						},
						BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "end", Type: entity.SwitchNode},
					},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "end", Type: entity.SwitchNode},
				},
			}},
			[]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_ForkDecide 测试
func TestDefaultWorkflowDecider_ForkDecide(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"下一个节点是Fork的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "FORK",
							"fork": []string{"t3", "t4"},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
						{"t4": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.AssignNodeDef{
						BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2", Type: entity.AssignNode},
					},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2", Type: entity.AssignNode},
				},
			}},
			[]string{"t2"},
			false,
		},
		{
			"当前节点是Fork的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "FORK",
							"fork": []string{"t3", "t4"},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
						{"t4": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.ForkNodeDef{
						BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Type: entity.ForkNode},
						Fork:         []string{"t3", "t4"},
					},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Type: entity.ForkNode},
				},
			}},
			[]string{"t3", "t4"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_FailedPolicyDecide 测试
func TestDefaultWorkflowDecider_FailedPolicyDecide(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name       string
		args       args
		want       []string
		statusWant entity.InstStatus
		wantErr    bool
	}{
		{
			"节点失败且没有配置失败策略的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1", Next: "t2"},
					Status:       entity.NodeInstFailed,
					Reason:       &entity.NodeReason{FailedReason: "test"},
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2"},
				},
			}},
			[]string{},
			entity.InstFailed,
			false,
		},
		{
			"节点失败且配置失败策略为Terminal的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2"},
					Status:  entity.NodeInstFailed,
					Reason:  &entity.NodeReason{FailedReason: "test"},
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2", Schedule: entity.Schedule{
						FailedPolicy: entity.Terminal,
					}},
				},
			}},
			[]string{},
			entity.InstFailed,
			false,
		},
		{
			"节点失败且配置失败策略为Ignore的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t2",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2"},
					Status:  entity.NodeInstFailed,
					Reason:  &entity.NodeReason{FailedReason: "test"},
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2", Schedule: entity.Schedule{
						FailedPolicy: entity.Ignore,
					}},
				},
			}},
			[]string{"t2"},
			entity.InstRunning,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.InstStatus, tt.statusWant) {
				t.Errorf("Decide() got = %v, want %v", got.InstStatus, tt.statusWant)
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_JoinDecide 测试
func TestDefaultWorkflowDecider_JoinDecide(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name       string
		args       args
		want       []string
		statusWant entity.InstStatus
		wantErr    bool
	}{
		{
			"JOIN的父节点还没有全部成功的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "FORK",
							"fork": []string{"t2", "t3"},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t4",
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t4",
						}},
						{"t4": map[string]interface{}{
							"type": "JOIN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.JoinNodeDef{
						BasicNodeDef: entity.BasicNodeDef{RefName: "t4", Next: "end", Type: entity.JoinNode}},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t4", Next: "end", Type: entity.JoinNode},
				},
			}},
			[]string{},
			entity.InstRunning,
			false,
		},
		{
			"JOIN的父节点全部成功, 且下一个节点是end的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "FORK",
							"fork": []string{"t2", "t3"},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t4",
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t4",
						}},
						{"t4": map[string]interface{}{
							"type": "JOIN",
							"next": "end",
						}},
					},
				},
				SchedNodeInsts: []*entity.NodeInst{{
					NodeDef:      entity.BasicNodeDef{RefName: "t2", Next: "t4"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "t4"},
				}, {
					NodeDef:      entity.BasicNodeDef{RefName: "t3", Next: "t4"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t3", Next: "t4"},
				}},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t4", Next: "end"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t4", Next: "end", Type: entity.JoinNode},
				},
			}},
			[]string{},
			entity.InstSucceed,
			false,
		},
		{
			"JOIN的父节点全部成功, 且下一个节点不是end的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "FORK",
							"fork": []string{"t2", "t3"},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t4",
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "t4",
						}},
						{"t4": map[string]interface{}{
							"type": "JOIN",
							"next": "t5",
						}},
						{"t5": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				SchedNodeInsts: []*entity.NodeInst{{
					NodeDef:      entity.BasicNodeDef{RefName: "t2", Next: "t4"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "t4"},
				}, {
					NodeDef:      entity.BasicNodeDef{RefName: "t3", Next: "t4"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t3", Next: "t4"},
				}},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t2", Next: "t4"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "t4", Type: entity.JoinNode},
				},
			}},
			[]string{"t4"},
			entity.InstRunning,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.InstStatus, tt.statusWant) {
				t.Errorf("Decide() got = %v, want %v", got.InstStatus, tt.statusWant)
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_SchedulePolicyDecide_Ignore 测试调度策略决策_跳过第一次调度
func TestDefaultWorkflowDecider_SchedulePolicyDecide_Ignore(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"流程开始启动的情况-跳过第一次调度",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "IGNORE_FIRST_SCHEDULE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
			}},
			[]string{"t2"},
			false,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-跳过第一次调度",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "IGNORE_FIRST_SCHEDULE",
							},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "IGNORE_FIRST_SCHEDULE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "IGNORE_FIRST_SCHEDULE",
							},
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2"},
			false,
		},
		{
			"使用next正常配置下一个节点为end的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "IGNORE_FIRST_SCHEDULE",
							},
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_SchedulePolicyDecide_Not_Complete 测试调度策略决策_不等待
func TestDefaultWorkflowDecider_SchedulePolicyDecide_Not_Complete(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"流程开始启动的情况-直接执行下个节点",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "SCHEDULE_NEXT_IF_NOT_COMPLETE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
			}},
			[]string{"t1", "t2"},
			false,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-直接执行下个节点",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "SCHEDULE_NEXT_IF_NOT_COMPLETE",
							},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2", "t3"},
			false,
		},
		{
			"使用next正常配置下一个节点为end的情况-直接执行下个节点",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "SCHEDULE_NEXT_IF_NOT_COMPLETE",
							},
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2"},
			false,
		},
		{
			"使用next正常配置下一个节点为schedulePolicy的情况-直接执行下个节点",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "SCHEDULE_NEXT_IF_NOT_COMPLETE",
							},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "SCHEDULE_NEXT_IF_NOT_COMPLETE",
							},
						}},
						{"t4": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2", "t3", "t4"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_ExecuteTimesPolicy_EXACTLY_ONCE 测试执行次数_有且执行一次
func TestDefaultWorkflowDecider_ExecuteTimesPolicy_EXACTLY_ONCE(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"流程开始启动的情况-正常执行",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "EXACTLY_ONCE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
			}},
			[]string{"t1"},
			false,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-跳过执行",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "EXACTLY_ONCE",
							},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				SchedNodeInsts: []*entity.NodeInst{
					&entity.NodeInst{
						NodeDef:      entity.BasicNodeDef{RefName: "t2"},
						Status:       entity.NodeInstSucceed,
						BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t3"},
			false,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-正常执行",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "EXACTLY_ONCE",
							},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2"},
			false,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-跳过后的正常执行",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "EXACTLY_ONCE",
							},
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"schedulePolicy": "IGNORE_FIRST_SCHEDULE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "EXACTLY_ONCE",
							},
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultWorkflowDecider_ExecuteTimesPolicy_AtLeastOnce 测试执行次数_最少执行一次
func TestDefaultWorkflowDecider_ExecuteTimesPolicy_AtLeastOnce(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name      string
		args      args
		want      []string
		wantErr   bool
		instState entity.InstStatus
	}{
		{
			"流程开始启动的情况-正常执行",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "AT_LEAST_ONCE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
			}},
			[]string{"t1"},
			false,
			entity.InstRunning,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-正常结束",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "AT_LEAST_ONCE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				SchedNodeInsts: []*entity.NodeInst{
					{
						NodeDef:      entity.BasicNodeDef{RefName: "t1"},
						Status:       entity.NodeInstSucceed,
						BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t2", Next: "end"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "end"},
				},
			}},
			[]string{},
			false,
			entity.InstSucceed,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-不结束",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "AT_LEAST_ONCE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t2", Next: "end"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2", Next: "end"},
				},
			}},
			[]string{},
			false,
			entity.InstRunning,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-等待孤立节点的结束",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
							"schedule": map[string]interface{}{
								"executeTimesPolicy": "AT_LEAST_ONCE",
							},
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				SchedNodeInsts: []*entity.NodeInst{{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef: entity.BasicNodeDef{RefName: "t1", Schedule: entity.Schedule{
						ExecuteTimesPolicy: entity.AtLeastOnce,
					}},
					Status: entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Schedule: entity.Schedule{
						ExecuteTimesPolicy: entity.AtLeastOnce,
					}},
				},
			}},
			[]string{"t2"},
			false,
			entity.InstRunning,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
				return
			}
			if got.InstStatus != tt.instState {
				t.Errorf("Decide() instState = %v, want %v", got.InstStatus, tt.instState)
			}
		})
	}
}

// TestDefaultWorkflowDecider_SkipNode 跳过节点
func TestDefaultWorkflowDecider_SkipNode(t *testing.T) {
	type args struct {
		inst *entity.WorkflowInst
	}
	tests := []struct {
		name      string
		args      args
		want      []string
		wantErr   bool
		instState entity.InstStatus
	}{
		{
			"流程开始启动的情况-正常执行",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
			}},
			[]string{"t1"},
			false,
			entity.InstRunning,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-正常结束",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1", Next: "t2"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1", Next: "t2"},
				},
				SkipNodes: []string{"t2"},
			}},
			[]string{},
			false,
			entity.InstSucceed,
		},
		{
			"使用数组中后面节点作为默认下一个节点的情况-不结束",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t2"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				},
				SkipNodes: []string{"t2"},
			}},
			[]string{"t3"},
			false,
			entity.InstRunning,
		},
		{
			"配置了节点执行条件且条件不满足的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"condition": "${false}",
							"type":      "ASSIGN",
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t3"},
			false,
			entity.InstRunning,
		},
		{
			"配置了节点执行条件且条件满足的情况",
			args{inst: &entity.WorkflowInst{
				WorkflowDef: &entity.WorkflowDef{
					Nodes: []map[string]interface{}{
						{"t1": map[string]interface{}{
							"type": "ASSIGN",
						}},
						{"t2": map[string]interface{}{
							"condition": "${true}",
							"type":      "ASSIGN",
						}},
						{"t3": map[string]interface{}{
							"type": "ASSIGN",
							"next": "end",
						}},
					},
				},
				CurNodeInst: &entity.NodeInst{
					NodeDef:      entity.BasicNodeDef{RefName: "t1"},
					Status:       entity.NodeInstSucceed,
					BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				},
			}},
			[]string{"t2"},
			false,
			entity.InstRunning,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultWorkflowDecider(expr.NewDefaultEvaluator())
			got, err := d.Decide(tt.args.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(entity.GetNodeRefNames(got.NodesToBeScheduled), tt.want) {
				t.Errorf("Decide() got = %v, want %v", got, tt.want)
				return
			}
			if got.InstStatus != tt.instState {
				t.Errorf("Decide() instState = %v, want %v", got.InstStatus, tt.instState)
			}
		})
	}
}
