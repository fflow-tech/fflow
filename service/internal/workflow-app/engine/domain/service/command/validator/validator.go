// Package validator 校验器实现
package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/xeipuuv/gojsonschema"
)

const (
	maxSubWorkflowCount = 10 // 最大子流程数量
)

// ValidateDefJson 检查传入内容格式
func ValidateDefJson(defJson string) error {
	if err := ValidateWorkflowDefJson(defJson); err != nil {
		return err
	}

	valueJson, err := simplejson.NewJson([]byte(defJson))
	if err != nil {
		return err
	}
	subWorkflows := valueJson.GetPath("subworkflows").MustArray()
	if len(subWorkflows) == 0 {
		return nil
	}

	// 如果定义中有子流程，则对子流程也进行校验
	return ValidateSubworkflowDefJson(subWorkflows)
}

// ValidateSubworkflowDefJson 检查子流程定义格式
func ValidateSubworkflowDefJson(subWorkflows []interface{}) error {
	for _, subWorkflow := range subWorkflows {
		subWorkflowMap := subWorkflow.(map[string]interface{})
		for _, def := range subWorkflowMap {
			defJson, err := json.Marshal(def)
			if err != nil {
				return fmt.Errorf("subworkflow def is not a valid json: %v", def)
			}
			if err := ValidateWorkflowDefJson(string(defJson)); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateWorkflowDefJson 检查流程定义格式
func ValidateWorkflowDefJson(defJson string) error {
	if err := ValidateJsonSchema(defJson); err != nil {
		return err
	}
	if err := ValidateDuplicateName(defJson); err != nil {
		return err
	}
	if err := ValidateInfiniteLoop(defJson); err != nil {
		return err
	}
	if err := ValidateDefJsonSize(defJson); err != nil {
		return err
	}

	return validateNodeConfig(defJson)
}

// ValidateJsonSchema 检查传入内容格式的schema格式
func ValidateJsonSchema(defJson string) error {
	schemaLoader := gojsonschema.NewStringLoader(config.GetSchemaConfig())
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("failed to new schema for validate def json: %w", err)
	}

	checkContentLoader := gojsonschema.NewStringLoader(defJson)
	result, err := schema.Validate(checkContentLoader)
	if err != nil {
		return fmt.Errorf("failed to new schema for validate def json: %w", err)
	}

	if result.Valid() {
		return nil
	}
	resultErrors := result.Errors()
	validateErrors := fmt.Errorf("illegal def json：")
	for _, resultError := range resultErrors {
		validateErrors = fmt.Errorf("%s[%s]", validateErrors, resultError.String())
	}

	return validateErrors
}

// ValidateDuplicateName 检测重复名称 RefName不能重复
func ValidateDuplicateName(defJson string) error {
	workflowDefEntity := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(defJson), workflowDefEntity); err != nil {
		return err
	}
	return validateWorkflowDuplicateName(workflowDefEntity)
}

func validateWorkflowDuplicateName(def *entity.WorkflowDef) error {
	refNameMap := make(map[string]struct{})
	if err := validateNodeCount(def); err != nil {
		return err
	}
	for _, node := range def.Nodes {
		for nodeName := range node {
			lowNodeName := strings.ToLower(nodeName)
			if _, ok := refNameMap[lowNodeName]; ok {
				return fmt.Errorf("duplicate ref name=[%s]", nodeName)
			}
			refNameMap[lowNodeName] = struct{}{}
		}
	}
	for _, subworkflow := range def.Subworkflows {
		for subworkflowName, workflowDef := range subworkflow {
			lowSubworkflowName := strings.ToLower(subworkflowName)
			if _, ok := refNameMap[lowSubworkflowName]; ok {
				return fmt.Errorf("duplicate ref name=[%s]", subworkflowName)
			}
			if err := validateWorkflowDuplicateName(&workflowDef); err != nil {
				return err
			}
			refNameMap[lowSubworkflowName] = struct{}{}
		}
	}
	return nil
}

// ValidateStartInstInput 检测开始流程实例的输入
func ValidateStartInstInput(def *entity.WorkflowDef, input map[string]interface{}) error {
	// 校验输入参数值大小
	if err := validateInputArgsSize(input); err != nil {
		return err
	}

	for _, inputKeyMap := range def.Input {
		for optionName, inputKeyDef := range inputKeyMap {
			if err := validateOptionInput(optionName, inputKeyDef, input); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateOptionInput 检查单个 option 的输入校验
func validateOptionInput(optionName string, inputKeyDef entity.InputKeyDef, input map[string]interface{}) error {
	inputOption, ok := input[optionName]
	if !ok {
		// 如果没有填写配置
		if inputKeyDef.Required {
			// 当必须要填时 则直接报错
			return fmt.Errorf("input field `%s` is required", optionName)
		}
		return nil
	}

	// 当填写了配置时，如果对应的 options 为空则可以填入任意值
	if inputKeyDef.Options == nil {
		return nil
	}

	// 如果对应的 options 不为空则必须是里面填入的值
	if !hasInputOption(inputKeyDef.Options, inputOption) {
		// 没有找到对应的选项则返回错误
		return fmt.Errorf("input field `%s` value %v must in options:%v",
			optionName, inputOption, inputKeyDef.Options)
	}
	return nil
}

// validateInputArgsSize 校验输入参数值大小
func validateInputArgsSize(input map[string]interface{}) error {
	inputOptionBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	maxInputArgsSize := config.GetValidationRulesConfig().MaxInputArgsSize
	inputArgLen := len(inputOptionBytes)
	if inputArgLen > maxInputArgsSize {
		return fmt.Errorf("input args length %d must < %dbytes", inputArgLen, maxInputArgsSize)
	}

	return nil
}

// hasInputOption 检测是否包含输入选项
func hasInputOption(options []interface{}, inputOption interface{}) bool {
	for _, option := range options {
		if option == inputOption {
			return true
		}
	}
	return false
}

// ValidateInfiniteLoop 死循环校验
func ValidateInfiniteLoop(defJson string) error {
	workflowDef := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(defJson), workflowDef); err != nil {
		return err
	}

	startNodeName, err := getStartNodeName(workflowDef)
	if err != nil {
		return err
	}

	// 没有节点时也默认成功
	if startNodeName == "" {
		return nil
	}

	ok, err := simulateRunWorkflowToEnd(startNodeName, workflowDef, nil)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("illegal workflow, caused by has endless loop")
	}

	return nil
}

// getStartNodeName 获取开始节点名称
func getStartNodeName(workflowDef *entity.WorkflowDef) (string, error) {
	nodeIndexMap, err := entity.GetNodeIndexDefMap(workflowDef)
	if err != nil {
		return "", err
	}

	node, ok := nodeIndexMap[0]
	if !ok {
		return "", nil
	}
	basicNodeDef, err := entity.GetBasicNodeDefFromNode(workflowDef, node)
	if err != nil {
		return "", err
	}
	return basicNodeDef.RefName, nil
}

// simulateRunWorkflowToEnd 模拟执行工作流到结束
func simulateRunWorkflowToEnd(node string, workflowDef *entity.WorkflowDef, runNodes []string) (bool, error) {
	runNodes = append(runNodes, node)
	nodeDef, err := entity.GetBasicNodeDefByRefName(workflowDef, node)
	if err != nil {
		return false, err
	}
	nextNodes, err := getNextNodesName(nodeDef, workflowDef, runNodes)
	if err != nil {
		return false, err
	}
	for _, nextNode := range nextNodes {
		if nextNode == entity.EndNode {
			return true, nil
		}
		end, err := simulateRunWorkflowToEnd(nextNode, workflowDef, runNodes)
		if err != nil {
			return false, err
		}
		if end {
			return true, nil
		}
	}
	return false, nil
}

// isSchedNode 是已经调度过的节点
func isSchedNode(curNodeName string, schedNodes []string) bool {
	for _, schedNode := range schedNodes {
		if schedNode == curNodeName {
			return true
		}
	}
	return false
}

// getNextNodesName 获取后续可能执行的节点名称，已经执行过了节点这里不会再返回
func getNextNodesName(curNode *entity.BasicNodeDef, workflowDef *entity.WorkflowDef, schedNodes []string) (
	[]string, error) {
	switch curNode.Type {
	case entity.ForkNode:
		{
			return getForkNodeTypeNextNodes(workflowDef, schedNodes, curNode.RefName)
		}
	case entity.SwitchNode:
		{
			return getSwitchNodeTypeNextNodes(workflowDef, schedNodes, curNode.RefName)
		}
	default:
		return getCommonTypeNextNodes(workflowDef, schedNodes, curNode)
	}
}

// getForkNodeTypeNextNodes 获取 fork 类型的后续执行节点
func getForkNodeTypeNextNodes(workflowDef *entity.WorkflowDef, schedNodes []string, curNodeName string) (
	[]string, error) {
	var nextNodes []string
	// 遍历所有的路径
	nodeDef, err := entity.GetNodeDefByRefName(workflowDef, curNodeName)
	if err != nil {
		return nil, err
	}
	forkNodes := nodeDef.(entity.ForkNodeDef)
	for _, forkNode := range forkNodes.Fork {
		if !isSchedNode(forkNode, schedNodes) {
			nextNodes = append(nextNodes, forkNode)
		}
	}
	return nextNodes, nil
}

// getSwitchNodeTypeNextNodes 获取 switch 类型的后续执行节点
func getSwitchNodeTypeNextNodes(workflowDef *entity.WorkflowDef, schedNodes []string, curNodeName string) (
	[]string, error) {
	var nextNodes []string
	// switch 遍历所有的路径
	nodeDef, err := entity.GetNodeDefByRefName(workflowDef, curNodeName)
	if err != nil {
		return nil, err
	}
	forkNodes := nodeDef.(entity.SwitchNodeDef)
	for _, forkNode := range forkNodes.Switch {
		if !isSchedNode(forkNode.Next, schedNodes) {
			nextNodes = append(nextNodes, forkNode.Next)
		}
	}
	return nextNodes, nil
}

// getCommonTypeNextNodes 获取普通类型的后续执行节点
func getCommonTypeNextNodes(workflowDef *entity.WorkflowDef, schedNodes []string, curNode *entity.BasicNodeDef) (
	[]string, error) {
	var nextNodes []string
	nextNodeName, err := entity.GetNextNode(workflowDef, curNode)
	if err != nil {
		return nil, err
	}
	if !isSchedNode(nextNodeName, schedNodes) {
		nextNodes = append(nextNodes, nextNodeName)
	}
	return nextNodes, nil
}

// ValidateSubworkflow 子流程配置校验
func ValidateSubworkflow(defJson string) error {
	workflowDef := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(defJson), workflowDef); err != nil {
		return err
	}
	if workflowDef.Subworkflows == nil {
		return nil
	}
	if err := validateSubworkflowCount(workflowDef); err != nil {
		return err
	}
	return validateSubworkflowDef(workflowDef)
}

// validateSubworkflowCount 检查子流程数量
func validateSubworkflowCount(workflowDef *entity.WorkflowDef) error {
	if len(workflowDef.Subworkflows) > maxSubWorkflowCount {
		return fmt.Errorf("subworkflow count %d > %d", len(workflowDef.Subworkflows), maxSubWorkflowCount)
	}
	return nil
}

// validateSubworkflowDef 检查子流程定义
func validateSubworkflowDef(workflowDef *entity.WorkflowDef) error {
	for _, subWorkflow := range workflowDef.Subworkflows {
		for _, subWorkflowDef := range subWorkflow {
			if subWorkflowDef.Subworkflows != nil {
				// 子流程不能包含子流程
				return fmt.Errorf("illegal subworkflow, caused by subworkflow must not contains subworkflow")
			}
		}
	}
	return nil
}

// ValidateDefJsonSize 检测定义内容的大小
func ValidateDefJsonSize(defJson string) error {
	defJsonSize := config.GetValidationRulesConfig().DefJsonSize
	jsonSize := len([]byte(defJson))
	if jsonSize > defJsonSize {
		return fmt.Errorf("workflow json size must < %dbytes", defJsonSize)
	}
	return nil
}

// validateNodeCount
func validateNodeCount(def *entity.WorkflowDef) error {
	maxNodesCount := config.GetValidationRulesConfig().MaxNodesCount
	minNodesCount := config.GetValidationRulesConfig().MinNodesCount
	if len(def.Nodes) < minNodesCount || len(def.Nodes) > maxNodesCount {
		return fmt.Errorf("total nodes count [%d] must < %d, > %d", len(def.Nodes), maxNodesCount, minNodesCount)
	}
	return nil
}

// validateNodeConfig 校验节点配置规则
func validateNodeConfig(defJson string) error {
	workflowDefEntity := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(defJson), workflowDefEntity); err != nil {
		return err
	}
	nodes, err := entity.GetNodeRefNameDefMap(workflowDefEntity)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		basicNodeDef := &entity.BasicNodeDef{}
		for _, nodeDef := range node {
			if err := utils.ToOtherInterfaceValue(&basicNodeDef, nodeDef); err != nil {
				return err
			}
			if err := validateRefNodeMustConfigNextField(basicNodeDef); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateRefNodeMustConfigNextField 检测引用类型节点是否具有下一个节点的配置
func validateRefNodeMustConfigNextField(nodeDef *entity.BasicNodeDef) error {
	if nodeDef.Type == entity.RefNode && nodeDef.Next == "" && len(nodeDef.Return) <= 0 {
		return fmt.Errorf("ref node [%s] must configure next node or return value", nodeDef.Name)
	}
	return nil
}
