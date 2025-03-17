package entity

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowInst 流程实例定义
type WorkflowInst struct {
	WorkflowDef                             *WorkflowDef           `json:"workflow_def,omitempty"`
	InstID                                  string                 `json:"inst_id,omitempty"`
	ParentInstID                            string                 `json:"parent_inst_id,omitempty"`
	ParentNodeInstID                        string                 `json:"parent_node_inst_id,omitempty"`
	Name                                    string                 `json:"name,omitempty"`
	Creator                                 string                 `json:"creator,omitempty"`
	PreStatus                               InstStatus             `json:"pre_status"`
	Status                                  InstStatus             `json:"status"`
	StartAt                                 time.Time              `json:"start_at"`
	LastRestartAt                           time.Time              `json:"last_restart_at"`                                // 最后一次重启时间
	LastRestartNode                         string                 `json:"last_restart_node,omitempty"`                    // 最后一次重启开始的节点
	BeforeLastRestartMaxNodeInstID          string                 `json:"before_last_restart_max_node_inst_id,omitempty"` // 最后一次重启最大的节点的实例ID
	CompletedAt                             time.Time              `json:"completed_at"`
	SchedNodeInsts                          []*NodeInst            `json:"sched_node_insts,omitempty"` // 已经被调度过的节点实例, 只取最新的
	CurNodeInst                             *NodeInst              `json:"cur_node_inst,omitempty"`    // 根据实际的情况实时生成
	Input                                   map[string]interface{} `json:"input,omitempty"`
	Output                                  map[string]interface{} `json:"output,omitempty"`
	Variables                               map[string]interface{} `json:"variables,omitempty"`
	Biz                                     map[string]interface{} `json:"biz,omitempty"`
	ExecutePath                             [][]string             `json:"execute_path,omitempty"` // 流程执行路径
	Owner                                   *Owner                 `json:"owner,omitempty"`
	FailedNodeRefNames                      []string               `json:"failed_node_ref_names,omitempty"`
	Reason                                  *InstReason            `json:"reason"`
	Operator                                *InstOperator          `json:"operator,omitempty"`
	IgnoreFirstScheduleNodes                []string               `json:"ignore_first_schedule_nodes,omitempty"`                  // 跳过了第一次调度的节点
	WaitSomeNodesCompleteBeforeInstComplete bool                   `json:"wait_some_nodes_execute_before_inst_complete,omitempty"` // 等待节点执行中
	SkipNodes                               []string               `json:"skip_nodes,omitempty"`                                   // 标记需要跳过的节点
	RunCompletedNodeInstIDsAfterPaused      []string               `json:"run_completed_node_inst_ids_after_paused,omitempty"`     // 在暂停后完成的节点实例ID
	WaitCompletedNodeInstIDsAfterPaused     []string               `json:"wait_completed_node_inst_ids_after_paused,omitempty"`    // 在暂停后完成等待的节点实例ID
	NodeInstsCount                          int                    `json:"node_insts_count"`                                       // 所有节点实例数量
	Breakpoints                             []string               `json:"breakpoints,omitempty"`                                  // 所有断点名称(调试模式下)
	CurBlockedBreakpoint                    string                 `json:"cur_blocked_breakpoint,omitempty"`                       // 当前被阻塞的断点(调试模式下)
	DebugMockNodes                          []string               `json:"debug_mock_nodes,omitempty"`                             // 调试模式下需要 MOCK 的节点
	CurDebugMode                            DebugMode              `json:"cur_debug_mode,omitempty"`                               // 当前调试模式
}

// InDebugMode 是否处于调试模式
func (w *WorkflowInst) InDebugMode() bool {
	return w.CurDebugMode != ""
}

// DebugMode 调试类型
type DebugMode string

var (
	SingleStepMode DebugMode = "SingleStepMode" // 单步调试模式
	BreakpointMode DebugMode = "BreakpointMode" // 断点调试模式
)

// InstReason 原因
type InstReason struct {
	StartReason     string              `json:"start_reason,omitempty"`
	RestartReason   string              `json:"restart_reason,omitempty"`
	PauseReason     string              `json:"pause_reason,omitempty"`
	ResumeReason    string              `json:"resume_reason,omitempty"`
	SucceedReason   string              `json:"succeed_reason,omitempty"`
	CancelledReason string              `json:"cancelled_reason,omitempty"`
	TimeoutReason   string              `json:"timeout_reason,omitempty"`
	FailedRootCause InstFailedRootCause `json:"failed_root_cause"`
}

// InstOperator 操作人
type InstOperator struct {
	RestartOperator   string `json:"restart_operator,omitempty"`
	PauseOperator     string `json:"pause_operator,omitempty"`
	ResumeOperator    string `json:"resume_operator,omitempty"`
	SucceedOperator   string `json:"succeed_operator,omitempty"`
	CancelledOperator string `json:"cancelled_operator,omitempty"`
	FailedOperator    string `json:"failed_operator,omitempty"`
}

// InstStatus 流程实例状态枚举
type InstStatus struct {
	intValue    int
	strValue    string
	isTerminal  bool
	isCompleted bool
}

// 流程状态枚举类型
var (
	InstRunning   = InstStatus{1, "running", false, false}  // 运行中
	InstPaused    = InstStatus{2, "paused", false, false}   // 暂停
	InstSucceed   = InstStatus{3, "succeed", true, true}    // 成功
	InstFailed    = InstStatus{4, "failed", true, true}     // 失败
	InstCancelled = InstStatus{5, "cancelled", true, false} // 取消
	InstTimeout   = InstStatus{6, "timeout", true, false}   // 超时
)

// IntValue 整数值
func (s InstStatus) String() string {
	return s.strValue
}

// IntValue 整数值
func (s InstStatus) IntValue() int {
	return s.intValue
}

// IsTerminal 是否是终态
func (s InstStatus) IsTerminal() bool {
	return s.isTerminal
}

// IsCompleted 是否是已完成的状态
func (s InstStatus) IsCompleted() bool {
	return s.isCompleted
}

var (
	intInstStatusMap = map[int]InstStatus{
		InstRunning.IntValue():   InstRunning,
		InstPaused.IntValue():    InstPaused,
		InstSucceed.IntValue():   InstSucceed,
		InstFailed.IntValue():    InstFailed,
		InstCancelled.IntValue(): InstCancelled,
		InstTimeout.IntValue():   InstTimeout,
	}
	strInstStatusMap = map[string]InstStatus{
		InstRunning.String():   InstRunning,
		InstPaused.String():    InstPaused,
		InstSucceed.String():   InstSucceed,
		InstFailed.String():    InstFailed,
		InstCancelled.String(): InstCancelled,
		InstTimeout.String():   InstTimeout,
	}
)

// GetInstStatusByIntValue 通过整数值返回状态枚举
func GetInstStatusByIntValue(i int) InstStatus {
	return intInstStatusMap[i]
}

// GetInstStatusByStrValue 通过字符串返回状态枚举
func GetInstStatusByStrValue(s string) InstStatus {
	return strInstStatusMap[s]
}

// MarshalJSON 重写序列化方法
func (s InstStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON 重写反序列化方法
func (s *InstStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = strInstStatusMap[j]
	return nil
}

// ConvertToCtx 转换成上下文
func ConvertToCtx(inst *WorkflowInst) (map[string]interface{}, error) {
	m, err := utils.StructToMap(inst)
	if err != nil {
		return nil, err
	}
	if inst.Variables == nil || len(inst.Variables) == 0 {
		inst.Variables = map[string]interface{}{}
	}
	if inst.WorkflowDef == nil {
		inst.WorkflowDef = &WorkflowDef{}
	}

	err = appendNodesInfoToCtx(inst, m)
	if err != nil {
		return nil, err
	}

	appendDefaultVariables(inst)
	ownerMap, err := utils.StructToMap(inst.Owner)
	if err != nil {
		return nil, err
	}
	m["w"] = map[string]interface{}{
		"i":         inst.Input,
		"input":     inst.Input,
		"b":         inst.Biz,
		"biz":       inst.Biz,
		"v":         inst.Variables,
		"variables": inst.Variables,
		"owner":     ownerMap,
		"o":         ownerMap,
	}

	return m, nil
}

func appendDefaultVariables(inst *WorkflowInst) {
	for _, k := range workflowInstIDKeys {
		inst.Variables[k] = inst.InstID
	}
	for _, k := range workflowInstNameKeys {
		inst.Variables[k] = inst.Name
	}
	for _, k := range workflowDefIDKeys {
		inst.Variables[k] = inst.WorkflowDef.DefID
	}
	for _, k := range workflowDefVersionKeys {
		inst.Variables[k] = strconv.Itoa(inst.WorkflowDef.Version)
	}
}

func appendNodesInfoToCtx(inst *WorkflowInst, ctx map[string]interface{}) error {
	for _, nodeInst := range inst.SchedNodeInsts {
		if err := AppendNodeInfoToCtx(ctx, nodeInst); err != nil {
			if err != nil {
				return err
			}

			return err
		}
	}
	return nil
}

// AppendNodeInfoToCtx 向上下文中添加节点信息, 其中 key 为要追加到的字段名称
func AppendNodeInfoToCtx(ctx map[string]interface{}, nodeInst *NodeInst) error {
	return AppendNodeInfoToCtxKey(ctx, nodeInst, nodeInst.BasicNodeDef.RefName)
}

// AppendNodeInfoToCtxKey 向上下文中添加节点信息, 其中 key 为要追加到的字段名称
func AppendNodeInfoToCtxKey(ctx map[string]interface{}, nodeInst *NodeInst, key string) error {
	// 如果还没有节点实例则直接返回
	if nodeInst == nil {
		return nil
	}

	basicNodeInfoMap := map[string]interface{}{
		"output":      nodeInst.Output,
		"o":           nodeInst.Output,
		"poll_output": nodeInst.PollOutput,
		"po":          nodeInst.PollOutput,
		"owner":       nodeInst.Owner,
	}

	appendNodeDefaultVariables(basicNodeInfoMap, nodeInst)

	operatorMap, err := utils.StructToMap(nodeInst.Operator)
	if err != nil {
		return err
	}
	reasonMap, err := utils.StructToMap(nodeInst.Reason)
	if err != nil {
		return err
	}

	operatorAndReasonMap, err := utils.MergeMap(operatorMap, reasonMap)
	if err != nil {
		return err
	}
	ctx[key], err = utils.MergeMap(basicNodeInfoMap, operatorAndReasonMap)
	return err
}

func appendNodeDefaultVariables(basicNodeInfoMap map[string]interface{}, nodeInst *NodeInst) {
	for _, k := range nodeInstIDKeys {
		basicNodeInfoMap[k] = nodeInst.NodeInstID
	}
	for _, k := range nodeRefNameKeys {
		basicNodeInfoMap[k] = nodeInst.BasicNodeDef.RefName
	}
}
