// Package entity 领域实体定义
package entity

import "fmt"

// DecideResult 决策结果
type DecideResult struct {
	NodesToBeScheduled  []*NodeInst          `json:"nodes_to_be_scheduled,omitempty"`
	NodesToBeUpdated    []*NodeInst          `json:"nodes_to_be_updated,omitempty"`
	InstStatus          InstStatus           `json:"inst_status"`
	InstFailedRootCause *InstFailedRootCause `json:"inst_failed_root_cause"`
}

// InstFailedRootCause 实例失败的根因
type InstFailedRootCause struct {
	FailedNodeRefNames []string `json:"failed_node_ref_names"`   // 失败的节点
	FailedReason       string   `json:"failed_reason,omitempty"` // 失败原因
}

// NewDecideResult 实例化
func NewDecideResult() *DecideResult {
	return &DecideResult{
		NodesToBeScheduled:  []*NodeInst{},
		NodesToBeUpdated:    []*NodeInst{},
		InstStatus:          InstRunning,
		InstFailedRootCause: &InstFailedRootCause{},
	}
}

// String 字符串
func (r DecideResult) String() string {
	return fmt.Sprintf("NodesToBeScheduled:%s, NodesToBeUpdated:%s, InstStatus:%s",
		getNodeRefNames(r.NodesToBeScheduled),
		getNodeRefNames(r.NodesToBeUpdated),
		r.InstStatus.String())
}

func getNodeRefNames(nodeInsts []*NodeInst) []string {
	r := []string{}
	for _, nodeInst := range nodeInsts {
		r = append(r, nodeInst.BasicNodeDef.RefName)
	}
	return r
}
