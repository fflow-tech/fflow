package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// NodeInst 节点实例定义
type NodeInst struct {
	NodeDef           interface{}            `json:"node_def,omitempty"`
	BasicNodeDef      BasicNodeDef           `json:"basic_node_def,omitempty"`
	Namespace         string                 `json:"namespace,omitempty"`
	DefID             string                 `json:"def_id,omitempty"`
	DefVersion        int                    `json:"def_version,omitempty"`
	SubworkflowInstID string                 `json:"subworkflow_inst_id,omitempty"`
	SubworkflowDefID  string                 `json:"subworkflow_def_id,omitempty"`
	InstID            string                 `json:"inst_id,omitempty"`
	NodeInstID        string                 `json:"node_inst_id,omitempty"`
	RetryCount        int                    `json:"retry_count,omitempty"` // 当前重试次数
	PreStatus         NodeInstStatus         `json:"pre_status"`
	Status            NodeInstStatus         `json:"status"`
	Owner             *Owner                 `json:"owner,omitempty"`
	Biz               map[string]interface{} `json:"biz,omitempty"`
	Input             map[string]interface{} `json:"input,omitempty"`
	Output            map[string]interface{} `json:"output,omitempty"`
	PollInput         map[string]interface{} `json:"poll_input,omitempty"`
	PollOutput        map[string]interface{} `json:"poll_output,omitempty"`
	CancelInput       map[string]interface{} `json:"cancel_input,omitempty"`
	CancelOutput      map[string]interface{} `json:"cancel_output,omitempty"`
	ScheduledAt       time.Time              `json:"scheduled_at"`
	WaitAt            time.Time              `json:"wait_at"`
	ExecuteAt         time.Time              `json:"execute_at"`
	AsynWaitResAt     time.Time              `json:"asyn_wait_res_at"`
	CompletedAt       time.Time              `json:"completed_at"`
	Operator          *NodeOperator          `json:"operator,omitempty"`
	Reason            *NodeReason            `json:"reason,omitempty"`
	PollCount         int                    `json:"poll_count,omitempty"`        // 当前轮询次数
	PollFailedCount   int                    `json:"poll_failed_count,omitempty"` // 轮询失败的次数
	FromRerun         bool                   `json:"from_rerun,omitempty"`        // 节点实例是否通过 RerunNode 创建
	Retrying          bool                   `json:"retrying,omitempty"`          // 节点实例重试中
	Nexts             []string               `json:"nexts,omitempty"`             // 节点所有可能的下一个节点
	Parents           []string               `json:"parents,omitempty"`           // 节点所有的父节点
	WaitForDebug      bool                   `json:"wait_for_debug,omitempty"`    // 因为调试阻塞
}

// NodeReason 原因
type NodeReason struct {
	RunReason        string `json:"run_reason,omitempty"`
	RerunReason      string `json:"rerun_reason,omitempty"`
	SucceedReason    string `json:"succeed_reason,omitempty"`
	CancelledReason  string `json:"cancelled_reason,omitempty"`
	FailedReason     string `json:"failed_reason,omitempty"`
	TimeoutReason    string `json:"timeout_reason,omitempty"`
	PollFailedReason string `json:"poll_failed_reason,omitempty"`
}

// NodeOperator 操作人
type NodeOperator struct {
	RunOperator       string `json:"run_operator,omitempty"`
	RerunOperator     string `json:"rerun_operator,omitempty"`
	SucceedOperator   string `json:"succeed_operator,omitempty"`
	CancelledOperator string `json:"cancelled_operator,omitempty"`
	FailedOperator    string `json:"failed_operator,omitempty"`
}

// newNodeInst 根据流程实例和真正的节点定义初始化节点实例
func newNodeInst(inst WorkflowInst, nodeDef interface{}) (*NodeInst, error) {
	basicDef, err := GetBasicNodeDefFromNodeDef(nodeDef)
	if err != nil {
		return nil, err
	}
	r := &NodeInst{
		BasicNodeDef: *basicDef,
		NodeDef:      nodeDef,
		Namespace:    inst.WorkflowDef.Namespace,
		DefID:        inst.WorkflowDef.DefID,
		DefVersion:   inst.WorkflowDef.Version,
		InstID:       inst.InstID,
		Status:       NodeInstScheduled,
		ScheduledAt:  time.Now(),
		Operator:     &NodeOperator{},
		Reason:       &NodeReason{},
		Owner: &Owner{
			Wechat:    basicDef.Owner.Wechat,
			ChatGroup: basicDef.Owner.ChatGroup,
		},
	}
	if r.BasicNodeDef.Name == "" {
		r.BasicNodeDef.Name = r.BasicNodeDef.RefName
	}

	// 节点实例写入节点下一个节点和父节点信息
	nexts, err := getAllConnectNextNodes(inst.WorkflowDef, basicDef)
	if err != nil {
		return nil, err
	}
	r.Nexts = nexts
	parents, err := GetNodeParents(inst.WorkflowDef, basicDef.RefName)
	if err != nil {
		return nil, err
	}
	r.Parents = parents

	return r, nil
}

// NewNodeInstByNodeRefName 通过节点引用名称初始化实例
func NewNodeInstByNodeRefName(inst WorkflowInst, refName string) (*NodeInst, error) {
	nodeRefNameMap, err := GetNodeRefNameDefMap(inst.WorkflowDef)
	if err != nil {
		return nil, err
	}
	node, ok := nodeRefNameMap[refName]
	if !ok {
		return nil, fmt.Errorf("[%d]not found refName=[%s] node", inst.InstID, refName)
	}

	index, err := GetIndexByRefName(inst.WorkflowDef, refName)
	if err != nil {
		return nil, err
	}

	nodeDef, err := NewNodeDef(node, index)
	if err != nil {
		return nil, err
	}
	switch nodeDef.(type) {
	case RefNodeDef:
		refNodeDef := nodeDef.(RefNodeDef)
		return newRefNodeInst(inst, nodeRefNameMap, refNodeDef.RefName, refNodeDef.Ref)
	default:
		return newNodeInst(inst, nodeDef)
	}
}

func newRefNodeInst(inst WorkflowInst, nodeRefNameMap map[string]map[string]interface{},
	refNodeName string, refedNodeName string) (*NodeInst, error) {
	realNodeInnerMap, err := getRealNodeInnerMap(nodeRefNameMap, refNodeName, refedNodeName)
	if err != nil {
		return nil, err
	}

	realNodeMap := map[string]interface{}{refNodeName: realNodeInnerMap}

	index, err := GetIndexByRefName(inst.WorkflowDef, refNodeName)
	if err != nil {
		return nil, err
	}

	nodeDef, err := NewNodeDef(realNodeMap, index)
	if err != nil {
		return nil, err
	}
	return newNodeInst(inst, nodeDef)
}

func getRealNodeInnerMap(nodeRefNameMap map[string]map[string]interface{},
	refNodeName string, refedNodeName string) (map[string]interface{}, error) {
	refNode := nodeRefNameMap[refNodeName]
	refedNode, ok := nodeRefNameMap[refedNodeName]
	if !ok {
		return nil, fmt.Errorf("inst not found ref=[%s] node", refedNodeName)
	}
	refNodeInnerMap := getInnerDefMap(refNode)
	if _, ok := refNodeInnerMap["return"]; !ok {
		if _, ok := refNodeInnerMap["next"]; !ok {
			return nil, fmt.Errorf("ref node [%s] must configure next node or return value", refNodeName)
		}
	}

	refedNodeInnerMap := getInnerDefMap(refedNode)

	refedNodeType, ok := refedNodeInnerMap["type"].(string)
	// 不允许引用非服务节点
	if !ok || !strings.EqualFold(refedNodeType, ServiceNode.String()) {
		return nil, fmt.Errorf("cannot ref not service node=[%s]", refedNodeName)
	}

	// 将引用节点的配置和被引用节点的配置做聚合
	realNodeInnerMap, err := utils.MergeMap(refedNodeInnerMap, refNodeInnerMap)
	if err != nil {
		return nil, err
	}

	// 保留被引用节点的节点类型
	realNodeInnerMap["type"] = refedNodeType
	return realNodeInnerMap, nil
}

func getInnerDefMap(node map[string]interface{}) map[string]interface{} {
	for _, m := range node {
		return m.(map[string]interface{})
	}

	return map[string]interface{}{}
}

// NodeInstStatus 流程节点状态枚举
type NodeInstStatus struct {
	intValue    int
	strValue    string
	isTerminal  bool
	isCompleted bool
}

// IntValue 整数值
func (s NodeInstStatus) IntValue() int {
	return s.intValue
}

// IsTerminal 是否是终态
func (s NodeInstStatus) IsTerminal() bool {
	return s.isTerminal
}

// IntValue 整数值
func (s NodeInstStatus) String() string {
	return s.strValue
}

// IsCompleted 是否是完成的状态
func (s NodeInstStatus) IsCompleted() bool {
	return s.isCompleted
}

var strNodeInstStatusMap = map[string]NodeInstStatus{
	NodeInstScheduled.String(): NodeInstScheduled,
	NodeInstWaiting.String():   NodeInstWaiting,
	NodeInstPaused.String():    NodeInstPaused,
	NodeInstRunning.String():   NodeInstRunning,
	NodeInstSucceed.String():   NodeInstSucceed,
	NodeInstFailed.String():    NodeInstFailed,
	NodeInstCancelled.String(): NodeInstCancelled,
	NodeInstTimeout.String():   NodeInstTimeout,
}

// GetNodeInstStatusByStrValue 通过字符串获取实例状态
func GetNodeInstStatusByStrValue(str string) NodeInstStatus {
	return strNodeInstStatusMap[str]
}

// MarshalJSON 重写序列化方法
func (s NodeInstStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON 重写反序列化方法
func (s *NodeInstStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = strNodeInstStatusMap[j]
	return nil
}

// 节点状态枚举类型
var (
	NodeInstScheduled = NodeInstStatus{1, "scheduled", false, false} // 已调度
	NodeInstWaiting   = NodeInstStatus{2, "waiting", false, false}   // 等待中
	NodeInstPaused    = NodeInstStatus{3, "paused", false, false}    // 暂停, 这个状态是预留的, 暂时不存在
	NodeInstRunning   = NodeInstStatus{4, "running", false, false}   // 运行
	NodeInstSucceed   = NodeInstStatus{5, "succeed", true, true}     // 成功
	NodeInstFailed    = NodeInstStatus{6, "failed", true, true}      // 失败
	NodeInstCancelled = NodeInstStatus{7, "cancelled", true, false}  // 取消
	NodeInstTimeout   = NodeInstStatus{8, "timeout", true, false}    // 超时
)

var (
	nodeInstStatusIntMap = map[int]NodeInstStatus{
		NodeInstScheduled.intValue: NodeInstScheduled,
		NodeInstWaiting.intValue:   NodeInstWaiting,
		NodeInstPaused.intValue:    NodeInstPaused,
		NodeInstRunning.intValue:   NodeInstRunning,
		NodeInstSucceed.intValue:   NodeInstSucceed,
		NodeInstFailed.intValue:    NodeInstFailed,
		NodeInstCancelled.intValue: NodeInstCancelled,
		NodeInstTimeout.intValue:   NodeInstTimeout,
	}
)

// GetNodeInstStatus 通过整数值返回状态枚举
func GetNodeInstStatus(i int) NodeInstStatus {
	return nodeInstStatusIntMap[i]
}

// GetNodeRefNameLatestInstMap 获取节点引用名称和节点实例的映射
// 如果一个节点被重复执行了多次取最新的
func GetNodeRefNameLatestInstMap(nodeInsts []*NodeInst) map[string]*NodeInst {
	SortByNodeInstIDDesc(nodeInsts)
	r := map[string]*NodeInst{}
	for _, inst := range nodeInsts {
		if _, exists := r[inst.BasicNodeDef.RefName]; !exists {
			r[inst.BasicNodeDef.RefName] = inst
		}
	}

	return r
}

// SortByNodeInstIDDesc 根据调度时间倒排
func SortByNodeInstIDDesc(nodeInsts []*NodeInst) {
	sort.SliceStable(nodeInsts, func(i, j int) bool {
		return nodeInsts[i].NodeInstID > nodeInsts[j].NodeInstID
	})
}

// SortByNodeInstIDAsc 根据调度时间培训
func SortByNodeInstIDAsc(nodeInsts []*NodeInst) {
	sort.SliceStable(nodeInsts, func(i, j int) bool {
		return nodeInsts[i].NodeInstID < nodeInsts[j].NodeInstID
	})
}

// GetNodeRefNames 获取节点引用名称
func GetNodeRefNames(nodeInsts []*NodeInst) []string {
	r := []string{}
	for _, nodeInst := range nodeInsts {
		r = append(r, nodeInst.BasicNodeDef.RefName)
	}

	return r
}

// GetNodeParents 获取节点引用名称
func GetNodeParents(workflowDef *WorkflowDef, refName string) ([]string, error) {
	parents := []string{}
	for _, node := range workflowDef.Nodes {
		basicNodeDef, err := GetBasicNodeDefFromNode(workflowDef, node)
		if err != nil {
			return nil, err
		}
		curRefName := basicNodeDef.RefName
		nextRefNames, err := getAllConnectNextNodes(workflowDef, basicNodeDef)
		if err != nil {
			return nil, err
		}
		for _, name := range nextRefNames {
			// 保证 parents 中元素不重复
			if name == refName && !utils.StrContains(parents, curRefName) {
				parents = append(parents, curRefName)
			}
		}
	}

	return parents, nil
}

// GetNextNode 获取一下个节点
// Switch 和 Fork 节点的下一个节点可能有多个, 不在现在这个方法考虑内
func GetNextNode(workflowDef *WorkflowDef, curBasicNodeDef *BasicNodeDef) (string, error) {
	// 如果包含 return, 下一个节点也同样为 end
	if len(curBasicNodeDef.Return) > 0 {
		return EndNode, nil
	}

	if curBasicNodeDef.Next != "" {
		return curBasicNodeDef.Next, nil
	}

	nodeIndexMap, err := GetNodeIndexDefMap(workflowDef)
	if err != nil {
		return "", err
	}

	if v, ok := nodeIndexMap[curBasicNodeDef.Index+1]; ok {
		nextBasicNodeDef, err := GetBasicNodeDefFromNode(workflowDef, v)
		if err != nil {
			return "", err
		}
		return nextBasicNodeDef.RefName, nil
	}

	return "", fmt.Errorf("[%d] Not found next node, curRefName=%s",
		workflowDef.DefID, curBasicNodeDef.RefName)
}

// getAllConnectNextNodes 获取有可能连接的下一个节点
func getAllConnectNextNodes(workflowDef *WorkflowDef, curBasicNodeDef *BasicNodeDef) ([]string, error) {
	if curBasicNodeDef.Type == ForkNode {
		return getNextsIfNodeIsForkNode(workflowDef, curBasicNodeDef)
	}

	if curBasicNodeDef.Type == SwitchNode {
		return getNextsIfNodeIsSwitchNode(workflowDef, curBasicNodeDef)
	}

	next, err := GetNextNode(workflowDef, curBasicNodeDef)
	if err != nil {
		return nil, err
	}
	return []string{next}, nil
}

func getNextsIfNodeIsForkNode(workflowDef *WorkflowDef, curBasicNodeDef *BasicNodeDef) ([]string, error) {
	if curBasicNodeDef.Type != ForkNode {
		next, err := GetNextNode(workflowDef, curBasicNodeDef)
		if err != nil {
			return nil, err
		}

		return []string{next}, nil
	}

	nodeDef, err := GetNodeDefByRefName(workflowDef, curBasicNodeDef.RefName)
	if err != nil {
		return nil, err
	}

	forkNodeDef := nodeDef.(ForkNodeDef)
	return forkNodeDef.Fork, nil
}

func getNextsIfNodeIsSwitchNode(workflowDef *WorkflowDef, curBasicNodeDef *BasicNodeDef) ([]string, error) {
	next, err := GetNextNode(workflowDef, curBasicNodeDef)
	if err != nil {
		return nil, err
	}

	if curBasicNodeDef.Type != SwitchNode {
		return []string{next}, nil
	}
	nodeDef, err := GetNodeDefByRefName(workflowDef, curBasicNodeDef.RefName)
	if err != nil {
		return nil, err
	}

	switchNodeDef := nodeDef.(SwitchNodeDef)

	r := []string{next}
	for _, c := range switchNodeDef.Switch {
		// 保证 nexts 中元素不重复
		if !utils.StrContains(r, c.Next) {
			r = append(r, c.Next)
		}
	}
	return r, nil
}

// GetRefNameNodeInstMap 获取节点引用名称和节点实例的映射
func GetRefNameNodeInstMap(nodeInsts []*NodeInst) map[string][]*NodeInst {
	r := map[string][]*NodeInst{}
	for _, nodeInst := range nodeInsts {
		if v, exists := r[nodeInst.BasicNodeDef.RefName]; exists {
			r[nodeInst.BasicNodeDef.RefName] = append(v, nodeInst)
			continue
		}
		r[nodeInst.BasicNodeDef.RefName] = []*NodeInst{nodeInst}
	}
	return r
}

// GetLastNodeInsts 拿出最新的节点执行实例, 每个节点只拿最后一次执行的实例
func GetLastNodeInsts(nodeInsts []*NodeInst) []*NodeInst {
	r := []*NodeInst{}
	refNameNodeInstMap := GetRefNameNodeInstMap(nodeInsts)
	for _, v := range refNameNodeInstMap {
		SortByNodeInstIDDesc(v)
		r = append(r, v[0])
	}

	SortByNodeInstIDDesc(r)
	return r
}

// GetOldestNodeInsts 拿出最老的节点执行实例, 每个节点只拿最先一次执行的实例
func GetOldestNodeInsts(nodeInsts []*NodeInst) []*NodeInst {
	r := []*NodeInst{}
	refNameNodeInstMap := GetRefNameNodeInstMap(nodeInsts)
	for _, v := range refNameNodeInstMap {
		SortByNodeInstIDAsc(v)
		r = append(r, v[0])
	}

	SortByNodeInstIDDesc(r)
	return r
}

// GetNodeInstByRefName 根据节点引用名称获取节点实例
func GetNodeInstByRefName(nodeInsts []*NodeInst, refName string) *NodeInst {
	for _, nodeInst := range nodeInsts {
		if nodeInst.BasicNodeDef.RefName == refName {
			return nodeInst
		}
	}

	return nil
}
