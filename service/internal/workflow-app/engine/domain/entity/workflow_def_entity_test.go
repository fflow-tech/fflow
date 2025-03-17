package entity

import (
	"encoding/json"
	"testing"

	"github.com/fflow-tech/fflow/service/pkg/log"
)

// TestDefToJson 测试
func TestDefToJson(t *testing.T) {
	def := &WorkflowDef{Status: Enabled}
	bytes, err := json.Marshal(def)
	if err != nil {
		t.Error(err)
	}
	log.Infof("%s", bytes)
	newDef := &WorkflowDef{}
	err = json.Unmarshal(bytes, newDef)
	if err != nil {
		t.Error(err)
	}
	log.Infof("%+v", newDef)
}

// TestNodeInstToJson 测试
func TestNodeInstToJson(t *testing.T) {
	nodeInst := &NodeInst{Status: NodeInstScheduled}
	bytes, err := json.Marshal(nodeInst)
	if err != nil {
		t.Error(err)
	}
	log.Infof("%s", bytes)
	newNodeInst := &NodeInst{}
	err = json.Unmarshal(bytes, newNodeInst)
	if err != nil {
		t.Error(err)
	}
	log.Infof("%+v", newNodeInst)
}

// TestInstToJson 测试
func TestInstToJson(t *testing.T) {
	inst := &WorkflowInst{Status: InstRunning}
	bytes, err := json.Marshal(inst)
	if err != nil {
		t.Error(err)
	}
	log.Infof("%s", bytes)
	newInst := &WorkflowInst{}
	err = json.Unmarshal(bytes, newInst)
	if err != nil {
		t.Error(err)
	}
	log.Infof("%+v", newInst)
}

type compareString struct {
	source string
	want   string
}

// TestDefJsonUnmarshal defJsonUnmarshal
func TestDefJsonUnmarshal(t *testing.T) {
	def := &WorkflowDef{}
	err := json.Unmarshal([]byte(testDefJson), def)
	if err != nil {
		t.Error(err)
	}
	compare := []compareString{
		{
			source: string(def.Timeout.Policy),
			want:   string(TimeoutWf),
		},
	}
	for _, trigger := range def.Triggers {
		compare = append(compare, compareString{string(trigger["t1"].Type), Timer.String()})

		for _, action := range trigger["t1"].Actions {
			compare = append(compare, compareString{string(action["a1"].ActionType), string(StartWorkflow)})
			compare = append(compare, compareString{string(action["a1"].AllowDays), string(Weekend)})
		}
	}

	for _, node := range def.Nodes {
		basicNodeDef, err := GetBasicNodeDefFromNode(def, node)
		if err != nil {
			t.Error(err)
		}
		compare = append(compare, compareString{string(basicNodeDef.Type), AssignNode.String()})
		compare = append(compare, compareString{string(basicNodeDef.Retry.Policy), Fixed.String()})
		compare = append(compare, compareString{string(basicNodeDef.Timeout.Policy), string(TimeoutWf)})
		compare = append(compare, compareString{string(basicNodeDef.Timeout.NearTimeoutPolicy), string(TimeoutWf)})
		compare = append(compare, compareString{string(basicNodeDef.Wait.AllowDays), Any.String()})
		compare = append(compare, compareString{string(basicNodeDef.Schedule.FailedPolicy), Ignore.String()})
		compare = append(compare, compareString{string(basicNodeDef.Schedule.ExecuteTimesPolicy),
			AtLeastOnce.String()})
		compare = append(compare, compareString{string(basicNodeDef.Schedule.SchedulePolicy),
			IgnoreFirstSchedule.String()})
	}
	for _, test := range compare {
		if test.source != test.want {
			t.Errorf("TestDefJsonUnmarshal  error source:%s want:%s", test.source, test.want)
		}
	}
}

var testDefJson = `{
    "biz": {
      "example1": "测试1",
      "example2": "测试2"
    },
    "timeout": {
        "duration": "1m",
        "policy": "TiME_OUT_WF"
    },
    "triggers": [
      {
        "t1": {
          "type": "tiMer",
          "expr": "0 */10 * * * * *",
          "actions": [
            {
              "a1": {
                "action": "STArT_WORKFLOW",
                "allowDays": "weekend",
                "args": {
                  "operator": "timer",
                  "name": "test",
                  "input": {}
                }
              }
            }
          ]
        }
      }
    ],
    "desc": "这是一个流程定义示例",
    "name": "流程定义示例",
    "nodes": [
      {
        "t1": {
          "next": "end",
          "type": "ASsIGN",
          "retry": {
            "count": 2,
            "duration": "5s",
            "policy": "FIxED"
          },
          "timeout": {
            "policy": "TIME_OUT_Wf",
            "nearTimeoutPolicy": "TIME_OUT_wf"
          },
          "wait": {
            "allowDays": "AnY"
          },
          "schedule": {
            "failedPolicy": "IGNORe",
            "schedulePolicy": "IGNORE_FIRST_SCHEdULE",
            "executeTimesPolicy": "aT_LEAST_ONCE"
          },
          "assign": [
            {
              "biz": {
                "projectName": "流程定义示例1"
              }
            },
            {
              "owner": {
                "wechat": "braumzhu"
              }
            },
            {
              "variables": {
                "test": "hh"
              }
            }
          ]
        }
      }
    ],
    "owner": {
      "wechat": "braumzhu"
    },
    "variables": {
      "projectName": "测试项目"
    }
  }`
