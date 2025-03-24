package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/nodeexecutor"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// ServiceNodeExecutor 服务节点执行器实现
type ServiceNodeExecutor struct {
	workflowInstRepo ports.WorkflowInstRepository
	remoteRepository ports.RemoteRepository
	executor         map[entity.ServiceType]ServiceNodeRemoteExecutor
	exprEvaluator    expr.Evaluator
}

// ServiceNodeRemoteExecutor 服务节点执行器
type ServiceNodeRemoteExecutor interface {
	Execute(context.Context, *entity.NodeInst, interface{}) error
	Polling(context.Context, *entity.NodeInst, interface{}) error
	Cancel(context.Context, *entity.NodeInst, interface{}) error
}

// NewServiceNodeExecutor 初始化执行器
func NewServiceNodeExecutor(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet) *ServiceNodeExecutor {
	return &ServiceNodeExecutor{
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		exprEvaluator:    workflowProviderSet.ExprEvaluator(),
		executor: map[entity.ServiceType]ServiceNodeRemoteExecutor{
			entity.HTTPService:   nodeexecutor.NewServiceHTTPNodeExecutor(repoProviderSet.RemoteRepo()),
			entity.FAASService:   nodeexecutor.NewServiceFAASNodeExecutor(repoProviderSet.RemoteRepo()),
			entity.MCPService:    nodeexecutor.NewServiceMCPNodeExecutor(repoProviderSet.RemoteRepo()),
			entity.OpenAIService: nodeexecutor.NewServiceOpenAINodeExecutor(repoProviderSet.RemoteRepo()),
		},
	}
}

// AsyncComplete 是否异步完成
func (d *ServiceNodeExecutor) AsyncComplete(nodeInst *entity.NodeInst) bool {
	return d.AsyncByTrigger(nodeInst) || d.AsyncByPolling(nodeInst)
}

// AsyncByTrigger 通过触发器实现异步
func (d *ServiceNodeExecutor) AsyncByTrigger(nodeInst *entity.NodeInst) bool {
	nodeDef, err := d.getNodeDef(nodeInst)
	if err != nil {
		log.Warnf("Failed to decide async by trigger, caused by %s", err)
		return false
	}

	return nodeDef.AsyncComplete && !d.AsyncByPolling(nodeInst)
}

// AsyncByPolling 通过轮询实现异步
func (d *ServiceNodeExecutor) AsyncByPolling(nodeInst *entity.NodeInst) bool {
	// 如果没有拿到轮询的配置, 则不属于通过轮询实现异步
	return !d.argsNotExists(nodeInst, entity.PollingArgs)
}

// Execute 执行节点
func (d *ServiceNodeExecutor) Execute(ctx context.Context, nodeInst *entity.NodeInst) error {
	executor, args, err := d.getExecutorAndArgs(nodeInst, entity.NormalArgs)
	if err != nil {
		return err
	}
	if err := executor.Execute(ctx, nodeInst, args); err != nil {
		return err
	}

	return d.setNodeInstStatusIfExecuteFailed(nodeInst)
}

func (d *ServiceNodeExecutor) setNodeInstStatusIfExecuteFailed(nodeInst *entity.NodeInst) error {
	inst, err := d.workflowInstRepo.Get(&dto.GetWorkflowInstDTO{
		InstID: nodeInst.InstID,
		DefID:  nodeInst.DefID,
	})
	if err != nil {
		return err
	}

	nodeBasicArgs, err := entity.GetServiceNodeBasicArgs(nodeInst.NodeDef, entity.NormalArgs)
	if err != nil {
		return err
	}

	// 如果条件没配, 默认为执行成功
	if nodeBasicArgs.SuccessCondition == "" {
		return nil
	}

	match, err := matchConditionForCurNodeInst(d.exprEvaluator, inst, nodeInst, nodeBasicArgs.SuccessCondition)
	if err != nil {
		return err
	}

	if !match {
		nodeInst.Status = entity.NodeInstFailed
		nodeInst.Reason.FailedReason = fmt.Sprintf(successConditionNotMatchErrFormat,
			nodeBasicArgs.SuccessCondition, utils.StructToJsonStr(nodeInst.Output))
	}

	return nil
}

func (d *ServiceNodeExecutor) setReqBodyToArgs(args interface{}, reqBody map[string]interface{}) {
	v := reflect.ValueOf(args).Elem()
	v.FieldByName("Body").Set(reflect.ValueOf(reqBody))
}

// setServiceNodeIsMockMode 如果当前实例处于调试模式，且当前节点被标记为 MOCK 节点，那么标记执行 MOCK 请求
func (d *ServiceNodeExecutor) setServiceNodeIsMockMode(workflowInst *entity.WorkflowInst,
	curNodeInst *entity.NodeInst, args interface{}) {
	nodeMockNode := workflowInst.InDebugMode() &&
		utils.StrContains(workflowInst.DebugMockNodes, curNodeInst.BasicNodeDef.RefName)
	reflect.ValueOf(args).Elem().FieldByName("MockMode").Set(reflect.ValueOf(nodeMockNode))
}

func (d *ServiceNodeExecutor) buildReqBody(workflowInst *entity.WorkflowInst,
	nodeInst *entity.NodeInst, argsType entity.ServiceNodeArgsType) (
	map[string]interface{}, error) {
	if len(nodeInst.Input) > 0 && argsType == entity.NormalArgs {
		return nodeInst.Input, nil
	}

	basicArgs, err := entity.GetServiceNodeBasicArgs(nodeInst.NodeDef, argsType)
	if err != nil {
		return nil, err
	}

	result, err := evaluateMapForCurNodeInst(d.exprEvaluator, workflowInst, nodeInst, basicArgs.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s]failed to evaluate map for argsType=%s: %w",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), argsType, err)
	}

	result, err = pathsAppendVariables(result, basicArgs.AppendVariables, workflowInst.Variables)
	if err != nil {
		return nil, err
	}

	result, err = pathsStringify(result, basicArgs.Stringify)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// pathsAppendVariables 将全局变量追加到配置的路径中
func pathsAppendVariables(source map[string]interface{}, pathStrs []string,
	oldVariables map[string]interface{}) (map[string]interface{}, error) {
	for _, pathStr := range pathStrs {
		paths := strings.Split(pathStr, ".")
		if len(paths) == 0 {
			continue
		}

		var err error
		newVariables := map[string]interface{}{}
		utils.CopyMap(oldVariables, newVariables)
		source, err = utils.MergeMap(source, generateNewMapByPath(paths, newVariables))
		if err != nil {
			return nil, err
		}
	}

	return source, nil
}

// pathsStringify 将配置的路径值转换成字符串
func pathsStringify(source map[string]interface{}, pathStrs []string) (map[string]interface{}, error) {
	// 动态替换之后还需要将 json 结构转换成字符串
	for _, pathStr := range pathStrs {
		paths := strings.Split(pathStr, ".")
		if len(paths) == 0 {
			continue
		}
		var err error
		source, err = stringifyPath(source, paths)
		if err != nil {
			return nil, err
		}
	}
	return source, nil
}

// stringifyPath 转换 map 指定位置的值为字符串
func stringifyPath(source map[string]interface{}, paths []string) (map[string]interface{}, error) {
	j, err := stringify(source, paths)
	if err != nil {
		return nil, err
	}

	return utils.MergeMap(source, generateNewMapByPath(paths, j))
}

func stringify(source map[string]interface{}, paths []string) (string, error) {
	if len(paths) <= 0 {
		return "", nil
	}

	m, ok := source, false
	for _, path := range paths {
		m, ok = m[path].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("failed to stringify, caused by paths %s not exists in input", path)
		}
	}
	b, err := json.Marshal(m)
	return string(b), err
}

func generateNewMapByPath(paths []string, v interface{}) map[string]interface{} {
	if len(paths) == 0 {
		return map[string]interface{}{}
	}

	if len(paths) == 1 {
		return map[string]interface{}{paths[0]: v}
	}

	r := map[string]interface{}{}
	c := map[string]interface{}{}
	r[paths[0]] = c
	for i := 1; i < len(paths); i++ {
		path := paths[i]
		if i < len(paths)-1 {
			c[path] = map[string]interface{}{}
			c = c[path].(map[string]interface{})
		}

		if i == len(paths)-1 {
			c[path] = v
		}
	}

	return r
}

// Polling 轮询节点
func (d *ServiceNodeExecutor) Polling(ctx context.Context, nodeInst *entity.NodeInst) error {
	executor, args, err := d.getExecutorAndArgs(nodeInst, entity.PollingArgs)
	if err != nil {
		return err
	}
	return executor.Polling(ctx, nodeInst, args)
}

// Cancel 取消执行节点
func (d *ServiceNodeExecutor) Cancel(ctx context.Context, nodeInst *entity.NodeInst) error {
	if d.argsNotExists(nodeInst, entity.CancelArgs) {
		log.Warnf("[%s]Node %s not exist cancel args, skip cancel node",
			logs.GetFlowTraceID(nodeInst.DefID, nodeInst.InstID), nodeInst.BasicNodeDef.RefName)
		nodeInst.Status = entity.NodeInstCancelled
		return nil
	}

	executor, args, err := d.getExecutorAndArgs(nodeInst, entity.CancelArgs)
	if err != nil {
		return err
	}
	return executor.Cancel(ctx, nodeInst, args)
}

func (d *ServiceNodeExecutor) argsNotExists(nodeInst *entity.NodeInst, argsType entity.ServiceNodeArgsType) bool {
	nodeDef, err := entity.GetServiceNodeBasicArgs(nodeInst.NodeDef, argsType)
	if err != nil {
		return true
	}

	return nodeDef.Protocol == ""
}

func (d *ServiceNodeExecutor) getExecutorAndArgs(nodeInst *entity.NodeInst, argsType entity.ServiceNodeArgsType) (
	ServiceNodeRemoteExecutor, interface{}, error) {
	args, err := entity.GetServiceNodeArgs(nodeInst.NodeDef, argsType)
	if err != nil {
		return nil, nil, err
	}
	executor, err := d.getExecutor(nodeInst.NodeDef, argsType)
	if err != nil {
		return nil, nil, err
	}

	workflowInst, err := d.workflowInstRepo.Get(&dto.GetWorkflowInstDTO{
		InstID: nodeInst.InstID,
		DefID:  nodeInst.DefID,
	})
	if err != nil {
		return nil, nil, err
	}

	reqBody, err := d.buildReqBody(workflowInst, nodeInst, argsType)
	if err != nil {
		return nil, nil, err
	}

	d.setServiceNodeIsMockMode(workflowInst, nodeInst, args)

	d.setReqBodyToArgs(args, reqBody)
	return executor, args, nil
}

// Type 获取是哪种节点类型的处理器
func (d *ServiceNodeExecutor) Type() entity.NodeType {
	return entity.ServiceNode
}

// getNodeDef 获取节点定义
func (d *ServiceNodeExecutor) getNodeDef(nodeInst *entity.NodeInst) (entity.ServiceNodeDef, error) {
	nodeDef, err := entity.ToActualNodeDef(nodeInst.BasicNodeDef.Type, nodeInst.NodeDef)
	if err != nil {
		return entity.ServiceNodeDef{}, err
	}
	return nodeDef.(entity.ServiceNodeDef), nil
}

// getExecutor 根据 args.protocol 获取对应的执行器
func (d *ServiceNodeExecutor) getExecutor(originNodeDef interface{},
	argsType entity.ServiceNodeArgsType) (ServiceNodeRemoteExecutor, error) {
	basicArgs, err := entity.GetServiceNodeBasicArgs(originNodeDef, argsType)
	if err != nil {
		return nil, err
	}
	protocol := strings.ToUpper(basicArgs.Protocol)
	executor, hasKey := d.executor[entity.ServiceType(protocol)]
	if !hasKey {
		return nil, fmt.Errorf("No executor for this protocol: %s", protocol)
	}
	return executor, nil
}
