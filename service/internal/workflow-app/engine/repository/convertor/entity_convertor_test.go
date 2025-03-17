package convertor

import (
	"reflect"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// Test_filterSchedNodeInsts 针对重启的节点还没有执行到或者没有重启的情况
func Test_filterSchedNodeInsts(t *testing.T) {
	type args struct {
		inst      *entity.WorkflowInst
		nodeInsts []*entity.NodeInst
	}
	tests := []struct {
		name string
		args args
		want []*entity.NodeInst
	}{
		{"没有重启过的情况", args{
			inst: &entity.WorkflowInst{
				LastRestartAt:   time.Time{},
				LastRestartNode: "",
			},
			nodeInsts: []*entity.NodeInst{{
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				NodeInstID:   "1",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "3",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t3"},
				NodeInstID:   "2",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "4",
			}}},
			[]*entity.NodeInst{{
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "4",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t3"},
				NodeInstID:   "2",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				NodeInstID:   "1",
			}},
		},
		{"重启的节点还没执行到的情况", args{
			inst: &entity.WorkflowInst{
				BeforeLastRestartMaxNodeInstID: "3",
				LastRestartNode:                "t4",
			},
			nodeInsts: []*entity.NodeInst{{
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				NodeInstID:   "1",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "3",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t3"},
				NodeInstID:   "2",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "4",
			}}},
			[]*entity.NodeInst{{
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "4",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t3"},
				NodeInstID:   "2",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
				NodeInstID:   "1",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterSchedNodeInsts(tt.args.inst, tt.args.nodeInsts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterSchedNodeInsts() = %s, want %s", utils.StructToJsonStr(got), utils.StructToJsonStr(tt.want))
			}
		})
	}
}

// Test_filterSchedNodeInstsForAlreadyRunToNode 针对重启的节点已执行到的情况
func Test_filterSchedNodeInstsForAlreadyRunToNode(t *testing.T) {
	type args struct {
		inst      *entity.WorkflowInst
		nodeInsts []*entity.NodeInst
	}
	tests := []struct {
		name string
		args args
		want []*entity.NodeInst
	}{
		{"重启的节点已执行到且执行过一次的情况", args{
			inst: &entity.WorkflowInst{
				BeforeLastRestartMaxNodeInstID: "4",
				LastRestartNode:                "t3",
			},
			nodeInsts: []*entity.NodeInst{{
				NodeInstID:   "1",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
			}, {
				NodeInstID:   "3",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
			}, {
				NodeInstID:   "2",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t3"},
			}, {
				NodeInstID:   "4",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
			}}},
			[]*entity.NodeInst{{
				NodeInstID:   "1",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
			}},
		},
		{"重启的节点已执行到且执行过多次的情况", args{
			inst: &entity.WorkflowInst{
				BeforeLastRestartMaxNodeInstID: "3",
				LastRestartNode:                "t2",
			},
			nodeInsts: []*entity.NodeInst{{
				NodeInstID:   "1",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
			}, {
				NodeInstID:   "2",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
			}, {
				NodeInstID:   "3",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t3"},
			}, {
				NodeInstID:   "4",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
			}, {
				NodeInstID:   "5",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t4"},
			}}},
			[]*entity.NodeInst{{
				NodeInstID:   "5",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t4"},
			}, {
				NodeInstID:   "4",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
			}, {
				NodeInstID:   "1",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
			}},
		},
		{"重启的节点已执行到且前面没有任何节点的的情况", args{
			inst: &entity.WorkflowInst{
				BeforeLastRestartMaxNodeInstID: "4",
				LastRestartNode:                "t1",
			},
			nodeInsts: []*entity.NodeInst{{
				NodeInstID:   "1",
				BasicNodeDef: entity.BasicNodeDef{RefName: "t1"},
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "3",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t3"},
				NodeInstID:   "2",
			}, {
				BasicNodeDef: entity.BasicNodeDef{RefName: "t2"},
				NodeInstID:   "4",
			}}},
			[]*entity.NodeInst{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterSchedNodeInsts(tt.args.inst, tt.args.nodeInsts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterSchedNodeInsts() = %s, want %s", utils.StructToJsonStr(got), utils.StructToJsonStr(tt.want))
			}
		})
	}
}
