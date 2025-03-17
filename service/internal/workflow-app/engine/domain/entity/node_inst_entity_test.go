package entity

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestGetNodeParents 测试获取父节点是否正确
func TestGetNodeParents(t *testing.T) {
	type args struct {
		workflowDef *WorkflowDef
		refName     string
	}
	def := &WorkflowDef{}
	json.Unmarshal([]byte(defJson), def)

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"节点t1", args{
			workflowDef: def,
			refName:     "t1",
		}, []string{}, false},
		{"节点fork1", args{
			workflowDef: def,
			refName:     "fork1",
		}, []string{"t1"}, false},
		{"节点a1", args{
			workflowDef: def,
			refName:     "a1",
		}, []string{"fork1"}, false},
		{"节点a2", args{
			workflowDef: def,
			refName:     "a2",
		}, []string{"fork1"}, false},
		{"节点join1", args{
			workflowDef: def,
			refName:     "join1",
		}, []string{"a1"}, false},
		{"节点t2", args{
			workflowDef: def,
			refName:     "t2",
		}, []string{"a2", "join1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeParents(tt.args.workflowDef, tt.args.refName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeParents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNodeParents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetNextNode 测试获取下一个节点
func TestGetNextNode(t *testing.T) {
	type args struct {
		workflowDef     *WorkflowDef
		curBasicNodeDef *BasicNodeDef
	}
	def := &WorkflowDef{}
	json.Unmarshal([]byte(defJson), def)

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"正常情况", args{
			workflowDef:     def,
			curBasicNodeDef: &BasicNodeDef{RefName: "join1"},
		}, "fork1", false},
		{"通过 return 结束的情况", args{
			workflowDef:     def,
			curBasicNodeDef: &BasicNodeDef{RefName: "t2", Return: map[string]interface{}{"k1": "v1"}},
		}, "end", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNextNode(tt.args.workflowDef, tt.args.curBasicNodeDef)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNextNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNextNode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var defJson = `
{
    "name": "测试join",
    "desc": "test",
    "timeout": {},
    "input": [
      {
        "operator": {
          "options": [
            "workflow",
            "conductor"
          ],
          "default": "workflow"
        }
      }
    ],
    "owner": {
      "wechat": "beihaizheng",
      "chatGroup": "wx2021"
    },
    "biz": {
      "example1": "测试1",
      "example2": "测试2"
    },
    "variables": {
      "projectName": "测试项目"
    },
    "nodes": [
      {
        "t1": {
          "assign": [
            {
              "owner": {
                "wechat": "beihaizheng"
              }
            },
            {
              "variables": {
                "test": 0
              }
            }
          ],
          "next": "fork1",
          "type": "ASSIGN"
        }
      },
      {
        "fork1": {
          "fork": [
            "a1",
            "a2"
          ],
          "type": "FORK"
        }
      },
      {
        "a1": {
          "assign": [
            {
              "owner": {
                "wechat": "beihaizheng"
              }
            },
            {
              "variables": {
                "test": 1
              }
            }
          ],
          "next": "join1",
          "type": "ASSIGN"
        }
      },
      {
        "a2": {
          "next": "t2",
          "type": "ASSIGN"
        }
      },
      {
        "join1": {
          "type": "JOIN"
        }
      },
      {
        "t2": {
          "assign": [
            {
              "owner": {
                "wechat": "beihaizheng"
              }
            },
            {
              "variables": {
                "test": 1
              }
            }
          ],
          "return": {"k1":"v1"},
          "type": "ASSIGN",
          "wait": {
            "allowDays": "WORKDAY",
            "expr": "0 25 16  ? * WED *"
          }
        }
      }
    ]
}
`
