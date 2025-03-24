package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	workflowInstIDKeys     = []string{"instID", "inst_id"}
	workflowInstNameKeys   = []string{"instName", "inst_name"}
	workflowDefIDKeys      = []string{"defID", "def_id"}
	workflowDefVersionKeys = []string{"defVersion", "def_version"}
	nodeInstIDKeys         = []string{"nodeInstID", "node_inst_id"}
	nodeRefNameKeys        = []string{"nodeRefName", "node_ref_name"}
)

// WorkflowDef 流程实体定义
// 因为和定义文件相关联, 所以 json 使用驼峰命名
type WorkflowDef struct {
	Namespace        string                   `json:"namespace,omitempty"`
	ID               string                   `json:"id,omitempty"`
	DefID            string                   `json:"defID,omitempty"`
	ParentDefID      string                   `json:"parentDefID,omitempty"`
	ParentDefVersion int                      `json:"parentDefVersion,omitempty"`
	RefName          string                   `json:"refName,omitempty"` // 子流程对应的 RefName
	Version          int                      `json:"version,omitempty"`
	Creator          string                   `json:"creator,omitempty"`
	Status           DefStatus                `json:"status,omitempty"`
	Name             string                   `json:"name,omitempty"`
	Desc             string                   `json:"desc,omitempty"`
	Timeout          Timeout                  `json:"timeout,omitempty"`
	Triggers         []map[string]TriggerDef  `json:"triggers,omitempty"`
	Input            []map[string]InputKeyDef `json:"input,omitempty"`
	Owner            Owner                    `json:"owner,omitempty"`
	Msg              WorkflowMsg              `json:"msg,omitempty"`
	Biz              map[string]interface{}   `json:"biz,omitempty"`
	Variables        map[string]interface{}   `json:"variables,omitempty"`
	Webhooks         []string                 `json:"webhooks,omitempty"`
	Nodes            []map[string]interface{} `json:"nodes,omitempty"`
	Subworkflows     []map[string]WorkflowDef `json:"subworkflows,omitempty"`
	CreatedAt        time.Time                `json:"createdAt,omitempty"`
}

// NodeMsg 节点消息
type NodeMsg struct {
	WaitMsg        string `json:"wait,omitempty"`        // 等待开始的消息
	StartMsg       string `json:"start,omitempty"`       // 执行开始的消息
	AsynWaitMsg    string `json:"asynWait,omitempty"`    // 异步等待的消息
	SuccessMsg     string `json:"success,omitempty"`     // 执行成功的消息
	FailMsg        string `json:"fail,omitempty"`        // 执行失败的消息
	CancelMsg      string `json:"cancel,omitempty"`      // 执行取消的消息
	TimeoutMsg     string `json:"timeout,omitempty"`     // 执行超时的消息
	NearTimeoutMsg string `json:"nearTimeout,omitempty"` // 执行接近超时的消息
}

// WorkflowMsg 流程消息
type WorkflowMsg struct {
	StartMsg          string `json:"start,omitempty"`          // 执行开始的消息
	SuccessMsg        string `json:"success,omitempty"`        // 执行成功的消息
	FailMsg           string `json:"fail,omitempty"`           // 执行失败的消息
	PauseMsg          string `json:"pause,omitempty"`          // 执行暂停的消息
	ResumeMsg         string `json:"resume,omitempty"`         // 执行恢复的消息
	CancelMsg         string `json:"cancel,omitempty"`         // 执行取消的消息
	SkipNodeMsg       string `json:"skipNode,omitempty"`       // 执行节点跳过的消息
	CancelSkipNodeMsg string `json:"cancelNodeSkip,omitempty"` // 执行节点取消跳过的消息
	TimeoutMsg        string `json:"timeout,omitempty"`        // 执行超时的消息
	NearTimeoutMsg    string `json:"nearTimeout,omitempty"`    // 执行接近超时的消息
}

// BasicNodeDef 节点定义
type BasicNodeDef struct {
	RefName       string                 `json:"refName,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Owner         Owner                  `json:"owner,omitempty"`
	Biz           map[string]interface{} `json:"biz,omitempty"`
	Type          NodeType               `json:"type,omitempty"`
	Condition     string                 `json:"condition,omitempty"`
	Retry         Retry                  `json:"retry,omitempty"`
	Timeout       NodeTimeout            `json:"timeout,omitempty"`
	Wait          Wait                   `json:"wait,omitempty"`
	Schedule      Schedule               `json:"schedule,omitempty"`
	Args          interface{}            `json:"args,omitempty"`
	Next          string                 `json:"next,omitempty"`
	Index         int                    `json:"index,omitempty"`         // 在整体中的序列号
	AsyncComplete bool                   `json:"asyncComplete,omitempty"` // 是否异步完成
	Return        map[string]interface{} `json:"return,omitempty"`        // 流程的返回
	Webhooks      []string               `json:"webhooks,omitempty"`
	Msg           NodeMsg                `json:"msg,omitempty"`
}

// Poll 轮询配置
type Poll struct {
	TimeoutDuration string     `json:"timeoutDuration,omitempty"`
	PollCondition   string     `json:"pollCondition,omitempty"`
	CancelCondition string     `json:"cancelCondition,omitempty"`
	InitialDuration string     `json:"initialDuration,omitempty"`
	MaxDuration     string     `json:"maxDuration,omitempty"`
	Policy          PollPolicy `json:"policy,omitempty"`
}

// Retry 重试配置
type Retry struct {
	Count    int         `json:"count,omitempty"`
	Duration string      `json:"duration,omitempty"`
	Policy   RetryPolicy `json:"policy,omitempty"`
}

// NodeTimeout 节点超时配置
type NodeTimeout struct {
	Duration            string        `json:"duration,omitempty"`
	Expr                string        `json:"expr,omitempty"`
	Policy              TimeoutPolicy `json:"policy,omitempty"`
	NearTimeoutDuration string        `json:"nearTimeoutDuration,omitempty"`
	NearTimeoutExpr     string        `json:"nearTimeoutExpr,omitempty"`
	NearTimeoutPolicy   TimeoutPolicy `json:"nearTimeoutPolicy,omitempty"`
}

// Wait 节点等待配置
type Wait struct {
	Duration  string          `json:"duration,omitempty"`
	Expr      string          `json:"expr,omitempty"`
	AllowDays AllowDaysPolicy `json:"allowDays,omitempty"`
}

// Schedule 调度配置
type Schedule struct {
	FailedPolicy       FailedPolicy       `json:"failedPolicy,omitempty"`
	SchedulePolicy     SchedulePolicy     `json:"schedulePolicy,omitempty"`
	ExecuteTimesPolicy ExecuteTimesPolicy `json:"executeTimesPolicy,omitempty"`
}

// ServiceNodeDef 服务节点定义
type ServiceNodeDef struct {
	BasicNodeDef
	PollArgs   interface{} `json:"pollArgs,omitempty"`
	CancelArgs interface{} `json:"cancelArgs,omitempty"`
}

// SwitchNodeDef Switch节点定义
type SwitchNodeDef struct {
	BasicNodeDef
	Switch []SwitchCase `json:"switch,omitempty"`
}

// SwitchCase Switch分支廷议
type SwitchCase struct {
	Condition string `json:"condition,omitempty"`
	Next      string `json:"next,omitempty"`
}

// ExclusiveJoinNodeDef 互斥JOIN节点定义
type ExclusiveJoinNodeDef struct {
	BasicNodeDef
}

// JoinNodeDef JOIN节点的定义
type JoinNodeDef struct {
	BasicNodeDef
}

// ForkNodeDef Fork节点定义
type ForkNodeDef struct {
	BasicNodeDef
	Fork []string `json:"fork,omitempty"`
}

// SubworkflowNodeDef 子流程节点定义
type SubworkflowNodeDef struct {
	BasicNodeDef
	Subworkflow string          `json:"subworkflow,omitempty"` // 内部定义的子流程的 refname，和 ID 必填一个
	ID          string          `json:"id,omitempty"`          // 流程定义 ID
	Version     int             `json:"version,omitempty"`     // 流程定义版本
	Args        SubworkflowArgs `json:"args,omitempty"`
}

// SubworkflowArgs 子流程节点参数
type SubworkflowArgs struct {
	Operator string                 `json:"operator,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
}

// AssignNodeDef Assign节点定义
type AssignNodeDef struct {
	BasicNodeDef
	Assign []AssignKey `json:"assign,omitempty"`
}

// EventNodeDef 事件节点定义
type EventNodeDef struct {
	BasicNodeDef
	Args EventArgs `json:"args,omitempty"`
}

// EventArgs 事件节点参数
type EventArgs struct {
	Target string   `json:"target,omitempty"`
	Event  EventDef `json:"event,omitempty"`
}

// EventDef 事件定义，遵循 Cloud Event 协议
type EventDef struct {
	Source          string                 `json:"source,omitempty"`
	Type            string                 `json:"type,omitempty"`
	DataContentType string                 `json:"dataContentType,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
}

// TransformNodeDef Transform节点定义
type TransformNodeDef struct {
	BasicNodeDef
	Output map[string]interface{} `json:"output,omitempty"`
}

// AssignKey 设置的变量配置
type AssignKey struct {
	Biz       map[string]interface{} `json:"biz,omitempty"`
	Owner     map[string]interface{} `json:"owner,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// RefNodeDef Ref节点定义
type RefNodeDef struct {
	ServiceNodeDef
	Ref string `json:"ref,omitempty"`
}

// WaitNodeDef Wait节点定义
type WaitNodeDef struct {
	BasicNodeDef
}

// ServiceNodeBasicArgs 服务节点基础参数
type ServiceNodeBasicArgs struct {
	Poll
	SuccessCondition string                 `json:"successCondition,omitempty"`
	Protocol         string                 `json:"protocol,omitempty"`
	Body             map[string]interface{} `json:"body,omitempty"`
	Stringify        []string               `json:"stringify,omitempty"`
	AppendVariables  []string               `json:"appendVariables,omitempty"`
	MockMode         bool                   `json:"mockMode"`
}

// HTTPArgs HTTP服务节点参数
type HTTPArgs struct {
	ServiceNodeBasicArgs
	Method     string            `json:"method,omitempty"`
	URL        string            `json:"url,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// FAASArgs Faas服务节点参数
type FAASArgs struct {
	ServiceNodeBasicArgs
	Namespace string `json:"namespace,omitempty"`
	Func      string `json:"func,omitempty"`
}

// MCPArgs MCP服务节点参数
type MCPArgs struct {
	ServiceNodeBasicArgs
	URL        string            `json:"url,omitempty"`
	Tool       string            `json:"tool,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// OpenAIArgs 定义了 OpenAI 节点执行器所需的参数
type OpenAIArgs struct {
	ServiceNodeBasicArgs
	Prompt      string              `json:"prompt"`                                              // 提示词
	Model       string              `json:"model"`                                               // 模型名称，如 "gpt-3.5-turbo"
	Messages    []map[string]string `json:"messages"`                                            // 消息列表，包含 role 和 content
	Temperature float64             `json:"temperature,default:0.5"`                             // 温度参数，控制随机性
	MaxTokens   int                 `json:"maxTokens,default:1000"`                              // 最大 token 数
	Stream      bool                `json:"stream,default:false"`                                // 是否使用流式响应
	APIKey      string              `json:"apiKey"`                                              // OpenAI API Key
	BaseURL     string              `json:"baseURL,default:https://api.openai.com/v1,omitempty"` // OpenAI API Base URL
}

// GetServiceNodeArgs 获取服务节点参数
func GetServiceNodeArgs(originNodeDef interface{}, argsType ServiceNodeArgsType) (interface{}, error) {
	args, err := getServiceNodeOriginArgs(originNodeDef, argsType)
	if err != nil {
		return nil, err
	}

	// 根据 protocol 字段来区分不同的 service
	protocol := getServiceProtocol(originNodeDef, argsType)
	switch protocol {
	case HTTPService:
		protocolArgs := &HTTPArgs{}
		err := utils.ToOtherInterfaceValue(&protocolArgs, args)
		return protocolArgs, err
	case FAASService:
		protocolArgs := &FAASArgs{}
		err := utils.ToOtherInterfaceValue(&protocolArgs, args)
		return protocolArgs, err
	case OpenAIService:
		protocolArgs := &OpenAIArgs{}
		err := utils.ToOtherInterfaceValue(&protocolArgs, args)
		return protocolArgs, err
	case MCPService:
		protocolArgs := &MCPArgs{}
		err := utils.ToOtherInterfaceValue(&protocolArgs, args)
		return protocolArgs, err
	default:
		return nil, fmt.Errorf("illegal service node protocol:[%v]", protocol)
	}
}

// getServiceProtocol 获取当前 service 节点的 protocol
func getServiceProtocol(originNodeDef interface{}, argsType ServiceNodeArgsType) ServiceType {
	r, _ := GetServiceNodeBasicArgs(originNodeDef, argsType)
	return ServiceType(strings.ToUpper(r.Protocol))
}

// ServiceNodeArgsType 服务节点参数类型
type ServiceNodeArgsType string

// ServiceNodeArgsType 服务节点参数类型
const (
	NormalArgs  ServiceNodeArgsType = "NORMAL"
	PollingArgs ServiceNodeArgsType = "POLLING"
	CancelArgs  ServiceNodeArgsType = "CANCEL"
)

// GetServiceNodeBasicArgs 从服务节点里面拿出基础的参数
func GetServiceNodeBasicArgs(originNodeDef interface{}, argsType ServiceNodeArgsType) (*ServiceNodeBasicArgs, error) {
	nodeDef, err := ToActualNodeDef(ServiceNode, originNodeDef)
	if err != nil {
		return nil, err
	}
	serviceNodeDef := nodeDef.(ServiceNodeDef)
	switch argsType {
	case NormalArgs:
		return getServiceNodeBasicArgs(serviceNodeDef.Args)
	case PollingArgs:
		return getServiceNodeBasicArgs(serviceNodeDef.PollArgs)
	case CancelArgs:
		return getServiceNodeBasicArgs(serviceNodeDef.CancelArgs)
	default:
		return nil, fmt.Errorf("illegal args type:%s", argsType)
	}
}

func getServiceNodeBasicArgs(args interface{}) (*ServiceNodeBasicArgs, error) {
	r := &ServiceNodeBasicArgs{}
	if err := utils.ToOtherInterfaceValue(r, args); err != nil {
		return nil, err
	}
	return r, nil
}

// getServiceNodeOriginArgs 从服务节点里面拿出的原始参数
func getServiceNodeOriginArgs(originNodeDef interface{}, argsType ServiceNodeArgsType) (interface{}, error) {
	nodeDef, err := ToActualNodeDef(ServiceNode, originNodeDef)
	if err != nil {
		return nil, err
	}
	serviceNodeDef := nodeDef.(ServiceNodeDef)

	switch argsType {
	case NormalArgs:
		return serviceNodeDef.Args, nil
	case PollingArgs:
		return serviceNodeDef.PollArgs, nil
	case CancelArgs:
		return serviceNodeDef.CancelArgs, nil
	default:
		return nil, fmt.Errorf("illegal args type:%s", argsType)
	}
}

// Owner 流程Owner
type Owner struct {
	Wechat    string `json:"wechat,omitempty"`
	ChatGroup string `json:"chatGroup,omitempty"`
}

// InputKeyDef 输入key定义
type InputKeyDef struct {
	Options  []interface{} `json:"options,omitempty"`
	Default  interface{}   `json:"default,omitempty"`
	Required bool          `json:"required,omitempty"`
}

// Timeout 流程超时设置
type Timeout struct {
	Duration string        `json:"duration,omitempty"`
	Expr     string        `json:"expr,omitempty"`
	Policy   TimeoutPolicy `json:"policy,omitempty"`
}

// DefStatus 流程定义状态枚举
type DefStatus struct {
	intValue int
	strValue string
}

// 流程状态枚举类型
var (
	Disabled = DefStatus{1, "disabled"} // 未激活
	Enabled  = DefStatus{2, "enabled"}  // 已激活
)

// IntValue 整数值
func (s DefStatus) IntValue() int {
	return s.intValue
}

// String 整数值
func (s DefStatus) String() string {
	return s.strValue
}

var (
	intDefStatusMap = map[int]DefStatus{
		Disabled.IntValue(): Disabled,
		Enabled.IntValue():  Enabled,
	}
	strDefStatusMap = map[string]DefStatus{
		Disabled.String(): Disabled,
		Enabled.String():  Enabled,
	}
)

// GetDefStatusByIntValue 通过整数值返回状态枚举
func GetDefStatusByIntValue(i int) DefStatus {
	return intDefStatusMap[i]
}

// GetDefStatusByStrValue 通过字符串值返回状态枚举
func GetDefStatusByStrValue(s string) DefStatus {
	return strDefStatusMap[s]
}

// MarshalJSON 重写序列化方法
func (s DefStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON 重写反序列化方法
func (s *DefStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = strDefStatusMap[j]
	return nil
}

// TimeoutPolicy 超时策略
type TimeoutPolicy string

const (
	TimeoutWf TimeoutPolicy = "TIME_OUT_WF" // 超时则终止流程, 默认值
	AlertOnly TimeoutPolicy = "ALERT_ONLY"  // 仅发送服务号告警不终止继续执行
)

// UnmarshalJSON 重写反序列化方法
func (s *TimeoutPolicy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = TimeoutPolicy(strings.ToUpper(j))
	return nil
}

// AllowDaysPolicy 允许的执行时间策略
type AllowDaysPolicy string

const (
	Any     AllowDaysPolicy = "ANY"     // 任意时间, 默认值
	Week    AllowDaysPolicy = "WEEK"    // 工作日+周末 非腾讯节假日
	Weekend AllowDaysPolicy = "WEEKEND" // 周六周日
)

// UnmarshalJSON 重写反序列化方法
func (s *AllowDaysPolicy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = AllowDaysPolicy(strings.ToUpper(j))
	return nil
}

// String 字符串值
func (p AllowDaysPolicy) String() string {
	return string(p)
}

// FailedPolicy 失败策略
type FailedPolicy string

const (
	Terminal FailedPolicy = "TERMINAL" // 终止流程
	Ignore   FailedPolicy = "IGNORE"   // 忽略失败
)

// UnmarshalJSON 重写反序列化方法
func (s *FailedPolicy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = FailedPolicy(strings.ToUpper(j))
	return nil
}

// String 字符串值
func (s FailedPolicy) String() string {
	return string(s)
}

// SchedulePolicy 调度策略
type SchedulePolicy string

const (
	ScheduleNextIfNotComplete SchedulePolicy = "SCHEDULE_NEXT_IF_NOT_COMPLETE" // 没有完成就可以调度下一个节点
	ScheduleNextUntilComplete SchedulePolicy = "SCHEDULE_NEXT_UNTIL_COMPLETE"  // 调度下一个直到节点执行完成
	IgnoreFirstSchedule       SchedulePolicy = "IGNORE_FIRST_SCHEDULE"         // 忽略第一次被调度，针对当前实例
)

// UnmarshalJSON 重写反序列化方法
func (s *SchedulePolicy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = SchedulePolicy(strings.ToUpper(j))
	return nil
}

// String 字符串值
func (s SchedulePolicy) String() string {
	return string(s)
}

// ExecuteTimesPolicy 执行次数策略
type ExecuteTimesPolicy string

const (
	AtLeastOnce ExecuteTimesPolicy = "AT_LEAST_ONCE" // 最少执行一次
	ExactlyOnce ExecuteTimesPolicy = "EXACTLY_ONCE"  // 有且执行一次
	AnyTimes    ExecuteTimesPolicy = "ANY"           // 可以执行任意次
)

// UnmarshalJSON 重写反序列化方法
func (s *ExecuteTimesPolicy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = ExecuteTimesPolicy(strings.ToUpper(j))
	return nil
}

// String 字符串值
func (s ExecuteTimesPolicy) String() string {
	return string(s)
}

// TriggerType 触发器类型
type TriggerType string

const (
	Timer TriggerType = "timer" // 定时器
	Event TriggerType = "event" // 事件
)

// UnmarshalJSON 重写反序列化方法
func (s *TriggerType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = TriggerType(strings.ToLower(j))
	return nil
}

// String 字符串值
func (s TriggerType) String() string {
	return string(s)
}

const (
	EndNode = "end" // 结束节点
)

// ActionType action类型
type ActionType string

const (
	StartWorkflow      ActionType = "START_WORKFLOW"
	RestartWorkflow    ActionType = "RESTART_WORKFLOW"
	CancelWorkflow     ActionType = "CANCEL_WORKFLOW"
	PauseWorkflow      ActionType = "PAUSE_WORKFLOW"
	ResumeWorkflow     ActionType = "RESUME_WORKFLOW"
	CompleteWorkflow   ActionType = "COMPLETE_WORKFLOW"
	SetWorkflowTimeout ActionType = "SET_WORKFLOW_TIMEOUT"
	RerunNode          ActionType = "RERUN_NODE"
	ResumeNode         ActionType = "RESUME_NODE"
	SkipNode           ActionType = "SKIP_NODE"
	CancelNode         ActionType = "CANCEL_SKIP_NODE"
	CompleteNode       ActionType = "COMPLETE_NODE"
	SetNodeTimeout     ActionType = "SET_NODE_TIMEOUT"
)

// UnmarshalJSON 重写反序列化方法
func (s *ActionType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = ActionType(strings.ToUpper(j))
	return nil
}

// NodeType 节点类型
type NodeType string

const (
	ServiceNode       NodeType = "SERVICE"        // 服务接口
	SwitchNode        NodeType = "SWITCH"         // Switch节点
	ExclusiveJoinNode NodeType = "EXCLUSIVE_JOIN" // 互斥Join节点
	TransformNode     NodeType = "TRANSFORM"      // 转换节点
	ForkNode          NodeType = "FORK"           // Fork节点
	JoinNode          NodeType = "JOIN"           // Join节点
	SubWorkflowNode   NodeType = "SUB_WORKFLOW"   // 子流程节点
	AssignNode        NodeType = "ASSIGN"         // 设置参数节点
	RefNode           NodeType = "REF"            // 引用节点, 功能节点支持引用, 非功能节点不支持引用
	WaitNode          NodeType = "WAIT"           // 等待节点
	EventNode         NodeType = "EVENT"          // 事件节点
)

// UnmarshalJSON 重写反序列化方法
func (s *NodeType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = NodeType(strings.ToUpper(j))
	return nil
}

// ServiceType 服务节点类型
type ServiceType string

// HTTP、FAAS 和 MCP 服务类型常量
const (
	HTTPService   ServiceType = "HTTP"   // HTTP 服务节点
	FAASService   ServiceType = "FAAS"   // FAAS 服务节点
	MCPService    ServiceType = "MCP"    // MCP 服务节点
	OpenAIService ServiceType = "OPENAI" // OPENAI 服务节点
)

// UnmarshalJSON 重写反序列化方法
func (s *ServiceType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = ServiceType(strings.ToUpper(j))
	return nil
}

// String 转换
func (s ServiceType) String() string {
	return string(s)
}

// String 转换
func (c NodeType) String() string {
	return string(c)
}

// RetryPolicy 重试策略
type RetryPolicy string

const (
	ExponentialBackoff RetryPolicy = "EXPONENTIAL_BACKOFF" // 指数退避
	Fixed              RetryPolicy = "FIXED"               // 固定间隔
)

// UnmarshalJSON 重写反序列化方法
func (s *RetryPolicy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = RetryPolicy(strings.ToUpper(j))
	return nil
}

// PollPolicy 轮询策略
type PollPolicy = RetryPolicy

// String 转换
func (p RetryPolicy) String() string {
	return string(p)
}

// NewNodeDef 实例化
func NewNodeDef(node map[string]interface{}, nodeIndex int) (interface{}, error) {
	if node == nil {
		return nil, fmt.Errorf("Node must not nil ")
	}

	refName, nodeDefMap, err := getNodeRefNameAndDefMap(node)
	if err != nil {
		return nil, err
	}

	t := nodeDefMap["type"]
	if !utils.IsString(t) {
		return nil, fmt.Errorf("illegal NodeType=%+v ", t)
	}
	nodeDefMap["index"] = nodeIndex
	nodeType := NodeType(strings.ToUpper(t.(string)))
	switch nodeType {
	case ServiceNode:
		nodeDef := ServiceNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case SwitchNode:
		nodeDef := SwitchNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case ExclusiveJoinNode:
		nodeDef := ExclusiveJoinNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case TransformNode:
		nodeDef := TransformNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case ForkNode:
		nodeDef := ForkNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case JoinNode:
		nodeDef := JoinNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case SubWorkflowNode:
		nodeDef := SubworkflowNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case RefNode:
		nodeDef := RefNodeDef{ServiceNodeDef: ServiceNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case AssignNode:
		nodeDef := AssignNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case WaitNode:
		nodeDef := WaitNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	case EventNode:
		nodeDef := EventNodeDef{BasicNodeDef: BasicNodeDef{RefName: refName}}
		err := utils.ToOtherInterfaceValue(&nodeDef, nodeDefMap)
		return nodeDef, err
	default:
		return nil, fmt.Errorf("Unsupport NodeType=%s ", nodeType)
	}
}

// ToActualNodeDef 转换成实际的节点定义
func ToActualNodeDef(nodeType NodeType, oldDef interface{}) (interface{}, error) {
	switch nodeType {
	case ServiceNode:
		newDef := ServiceNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case SwitchNode:
		newDef := SwitchNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case ExclusiveJoinNode:
		newDef := ExclusiveJoinNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case TransformNode:
		newDef := TransformNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case ForkNode:
		newDef := ForkNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case JoinNode:
		newDef := JoinNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case SubWorkflowNode:
		newDef := SubworkflowNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case RefNode:
		newDef := RefNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case AssignNode:
		newDef := AssignNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case WaitNode:
		newDef := WaitNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	case EventNode:
		newDef := EventNodeDef{}
		err := utils.ToOtherInterfaceValue(&newDef, oldDef)
		return newDef, err
	default:
		return nil, fmt.Errorf("Unsupport NodeType=%s ", nodeType)
	}
}

// getNodeRefNameAndDefMap 获取节点的引用名称和定义MAP
func getNodeRefNameAndDefMap(node map[string]interface{}) (string, map[string]interface{}, error) {
	for refName, nodeDefMap := range node {
		return refName, nodeDefMap.(map[string]interface{}), nil
	}

	return "", nil, fmt.Errorf("illegal BasicNodeDef ")
}

// GetNodeRefNameDefMap 获取节点引用名称和节点定义的映射
func GetNodeRefNameDefMap(def *WorkflowDef) (map[string]map[string]interface{}, error) {
	r := map[string]map[string]interface{}{}
	for _, node := range def.Nodes {
		refName, _, err := getNodeRefNameAndDefMap(node)
		if err != nil {
			return nil, err
		}
		r[refName] = node
	}

	return r, nil
}

// GetNodeIndexDefMap 获取节点序列号和节点定义的映射
func GetNodeIndexDefMap(def *WorkflowDef) (map[int]map[string]interface{}, error) {
	r := map[int]map[string]interface{}{}
	for i, node := range def.Nodes {
		r[i] = node
	}
	return r, nil
}

// GetBasicNodeDefFromNodeDef 获取已经解析完成的节点定义获取基础的节点定义
func GetBasicNodeDefFromNodeDef(nodeDef interface{}) (*BasicNodeDef, error) {
	basicNodeDef := &BasicNodeDef{}
	err := utils.ToOtherInterfaceValue(basicNodeDef, nodeDef)
	if err != nil {
		return nil, err
	}
	return basicNodeDef, err
}

// GetIndexByRefName 根据节点引用名称获取序号
func GetIndexByRefName(workflowDef *WorkflowDef, refName string) (int, error) {
	for i, node := range workflowDef.Nodes {
		curRefName, _, err := getNodeRefNameAndDefMap(node)
		if err != nil {
			return 0, err
		}
		if strings.EqualFold(refName, curRefName) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("[%d]Not exists refName=[%s] index,", workflowDef.DefID, refName)
}

// GetBasicNodeDefByRefName 根据节点引用名称获取节点定义
func GetBasicNodeDefByRefName(workflowDef *WorkflowDef, refName string) (*BasicNodeDef, error) {
	for _, node := range workflowDef.Nodes {
		curRefName, _, err := getNodeRefNameAndDefMap(node)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(refName, curRefName) {
			return GetBasicNodeDefFromNode(workflowDef, node)
		}
	}

	return nil, fmt.Errorf("[%d]Not exists refName=[%s] NodeDef,", workflowDef.DefID, refName)
}

// GetNodeDefByRefName 根据节点引用名称获取节点真实的定义
func GetNodeDefByRefName(workflowDef *WorkflowDef, refName string) (interface{}, error) {
	for _, node := range workflowDef.Nodes {
		curRefName, _, err := getNodeRefNameAndDefMap(node)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(refName, curRefName) {
			index, err := GetIndexByRefName(workflowDef, refName)
			if err != nil {
				return nil, err
			}

			return NewNodeDef(node, index)
		}
	}

	return nil, fmt.Errorf("[%d]Not exists refName=[%s] NodeDef,", workflowDef.DefID, refName)
}

// GetBasicNodeDefFromNode 根据原始的节点map获取基础的节点定义
func GetBasicNodeDefFromNode(workflowDef *WorkflowDef, node map[string]interface{}) (*BasicNodeDef, error) {
	refName, _, err := getNodeRefNameAndDefMap(node)
	if err != nil {
		return nil, err
	}

	index, err := GetIndexByRefName(workflowDef, refName)
	if err != nil {
		return nil, err
	}

	nodeDef, err := NewNodeDef(node, index)
	if err != nil {
		return nil, err
	}
	basicNodeDef := &BasicNodeDef{}
	if err = utils.ToOtherInterfaceValue(&basicNodeDef, nodeDef); err != nil {
		return nil, err
	}
	return basicNodeDef, err
}
