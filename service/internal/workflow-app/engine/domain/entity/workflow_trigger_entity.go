package entity

import (
	"fmt"
	"strconv"

	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// Trigger 触发器
type Trigger struct {
	BasicTriggerDef
	Action
	TriggerID  string        `json:"triggerID,omitempty"`
	Level      TriggerLevel  `json:"level,omitempty"`
	Namespace  string        `json:"namespace,omitempty"`
	DefID      string        `json:"defID,omitempty"`
	DefVersion int           `json:"defVersion,omitempty"`
	InstID     string        `json:"instID,omitempty"`
	Status     TriggerStatus `json:"status,omitempty"`
}

// TriggerDef 触发器配置
type TriggerDef struct {
	BasicTriggerDef
	Actions []map[string]Action `json:"actions,omitempty"`
}

// BasicTriggerDef 基础触发器配置
type BasicTriggerDef struct {
	RefName string      `json:"refName,omitempty"` // 引用名称
	Type    TriggerType `json:"type,omitempty"`
	Event   string      `json:"event,omitempty"` // 事件名称
	Expr    string      `json:"expr,omitempty"`
}

// Action 节点动作
type Action struct {
	Name       string          `json:"name,omitempty"`      // 动作的名称
	Condition  string          `json:"condition,omitempty"` // 事件执行条件表达式
	AllowDays  AllowDaysPolicy `json:"allowDays,omitempty"`
	ActionType ActionType      `json:"action,omitempty"`
	Args       interface{}     `json:"args,omitempty"`
}

// BasicActionArgs 基础定时事件
type BasicActionArgs struct {
	Operator string                 `json:"operator,omitempty"` // 操作人
	Input    map[string]interface{} `json:"input,omitempty"`    // 参数
}

// StartWorkflowActionArgs 启动流程消息
type StartWorkflowActionArgs struct {
	BasicActionArgs
	DefID      string `json:"defID,omitempty"` // 流程定义ID
	DefVersion int    `json:"defVersion"`      // 流程定义版本号
	Name       string `json:"name,omitempty"`  // 流程实例名称
}

// RerunNodeActionArgs 重跑节点
type RerunNodeActionArgs struct {
	BasicActionArgs
	DefID  string `json:"defID,omitempty"`
	InstID string `json:"instID,omitempty"`
	Node   string `json:"node,omitempty"`
}

// ResumeNodeActionArgs 恢复节点
type ResumeNodeActionArgs struct {
	BasicActionArgs
	DefID  string `json:"defID,omitempty"`
	InstID string `json:"instID,omitempty"`
	Node   string `json:"node,omitempty"`
}

// CompleteNodeActionArgs 完成节点
type CompleteNodeActionArgs struct {
	BasicActionArgs
	DefID  string                 `json:"def_id,omitempty"`
	InstID string                 `json:"inst_id,omitempty"`
	Node   string                 `json:"node,omitempty"`
	Status string                 `json:"status,omitempty"`
	Output map[string]interface{} `json:"output,omitempty"`
}

// TriggerLevel 触发器级别
type TriggerLevel int

const (
	DefTrigger  TriggerLevel = 1 // 流程定义级别触发器
	InstTrigger TriggerLevel = 2 // 流程实例级别触发器
)

var (
	levelMap = map[int]string{1: "DefTrigger", 2: "InstTrigger"}
)

// IntValue 返回触发器级别对应整数值
func (t TriggerLevel) IntValue() int {
	return int(t)
}

// String 字符串
func (t TriggerLevel) String() string {
	return levelMap[t.IntValue()]
}

// TriggerStatus 触发器状态
type TriggerStatus int

const (
	DisabledTrigger TriggerStatus = 1 // 未激活
	EnabledTrigger  TriggerStatus = 2 // 已激活
)

// IntValue 获取对应整数值
func (t TriggerStatus) IntValue() int {
	return int(t)
}

// DefConvertToCtx 转换成上下文
func DefConvertToCtx(def *WorkflowDef) (map[string]interface{}, error) {
	workflowDefCtxMap, err := utils.StructToMap(def)
	if err != nil {
		return nil, err
	}
	if len(def.Variables) == 0 {
		def.Variables = map[string]interface{}{}
	}
	ownerMap, err := utils.StructToMap(def.Owner)
	if err != nil {
		return nil, err
	}
	appendDefDefaultVariables(def)
	workflowDefCtxMap["w"] = map[string]interface{}{
		"i":         def.Input,
		"input":     def.Input,
		"b":         def.Biz,
		"biz":       def.Biz,
		"v":         def.Variables,
		"variables": def.Variables,
		"owner":     ownerMap,
	}
	log.Infof("workflow def ctx map:%s", utils.StructToJsonStr(workflowDefCtxMap))
	return workflowDefCtxMap, nil
}

func appendDefDefaultVariables(def *WorkflowDef) {
	for _, k := range workflowDefIDKeys {
		def.Variables[k] = def.DefID
	}
	for _, k := range workflowDefVersionKeys {
		def.Variables[k] = strconv.Itoa(def.Version)
	}
}

// GetActionArgs 获取动作
func GetActionArgs(originActionArgs interface{}, actionType ActionType) (interface{}, error) {
	switch actionType {
	case StartWorkflow:
		action := &StartWorkflowActionArgs{}
		err := utils.ToOtherInterfaceValue(&action, originActionArgs)
		return action, err
	case RerunNode:
		action := &RerunNodeActionArgs{}
		err := utils.ToOtherInterfaceValue(&action, originActionArgs)
		return action, err
	case ResumeNode:
		action := &ResumeNodeActionArgs{}
		err := utils.ToOtherInterfaceValue(&action, originActionArgs)
		return action, err
	case CompleteNode:
		action := &CompleteNodeActionArgs{}
		err := utils.ToOtherInterfaceValue(&action, originActionArgs)
		return action, err
	default:
		return nil, fmt.Errorf("unsupported actionType=%s", actionType)
	}
}

var (
	triggerLevelMap = map[ActionType]TriggerLevel{
		StartWorkflow:      DefTrigger,
		RestartWorkflow:    InstTrigger,
		CancelWorkflow:     InstTrigger,
		PauseWorkflow:      InstTrigger,
		ResumeWorkflow:     InstTrigger,
		CompleteWorkflow:   InstTrigger,
		SetWorkflowTimeout: InstTrigger,
		RerunNode:          InstTrigger,
		ResumeNode:         InstTrigger,
		SkipNode:           InstTrigger,
		CancelNode:         InstTrigger,
		CompleteNode:       InstTrigger,
		SetNodeTimeout:     InstTrigger,
	}
)

// GetTriggerLevel 获取触发器级别
func GetTriggerLevel(actionType ActionType) TriggerLevel {
	return triggerLevelMap[actionType]
}

// GetAllAction 从[]map[string]Action中提取所有Action
func GetAllAction(actionMaps []map[string]Action) []Action {
	actions := make([]Action, 0, len(actionMaps))
	for _, actionMap := range actionMaps {
		for k, v := range actionMap {
			v.Name = k
			actions = append(actions, v)
		}
	}

	return actions
}
