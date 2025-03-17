package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/remote"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// WorkflowTriggerCommandService 流程触发器服务
type WorkflowTriggerCommandService struct {
	eventBusRepo     ports.EventBusRepository
	workflowDefRepo  ports.WorkflowDefRepository
	workflowInstRepo ports.WorkflowInstRepository
	triggerRepo      ports.TriggerRepository
	remoteRepo       ports.RemoteRepository
	exprEvaluator    expr.Evaluator
	triggerActorMap  map[entity.ActionType]func(ctx context.Context,
		trigger *entity.Trigger, actionArgs interface{}) error
}

// NewWorkflowTriggerCommandService 初始化
func NewWorkflowTriggerCommandService(workflowProviderSet *execution.WorkflowProviderSet,
	repoProviderSet *ports.RepoProviderSet,
	actor *execution.DefaultTriggerActor) *WorkflowTriggerCommandService {
	return &WorkflowTriggerCommandService{
		eventBusRepo:     repoProviderSet.EventBusRepo(),
		workflowDefRepo:  repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo: repoProviderSet.WorkflowInstRepo(),
		triggerRepo:      repoProviderSet.TriggerRepo(),
		remoteRepo:       repoProviderSet.RemoteRepo(),
		exprEvaluator:    workflowProviderSet.ExprEvaluator(),
		triggerActorMap: map[entity.ActionType]func(ctx context.Context,
			trigger *entity.Trigger, actionArgs interface{}) error{
			entity.StartWorkflow: actor.OnStartWorkflow,
			entity.RerunNode:     actor.OnRerunNode,
			entity.ResumeNode:    actor.OnResumeNode,
			entity.CompleteNode:  actor.OnCompleteNode,
		},
	}
}

// ConsumeCronTriggerEvent 消费定时触发器事件
func (m *WorkflowTriggerCommandService) ConsumeCronTriggerEvent(ctx context.Context,
	req *dto.CronTriggerEventDTO) error {
	pulsarMsg := req.Message
	cronTriggerEvent := event.CronTriggerEvent{}
	if err := json.Unmarshal(pulsarMsg.Payload(), &cronTriggerEvent); err != nil {
		return err
	}

	getDTO := &dto.GetTriggerDTO{
		DefID:     cronTriggerEvent.DefID,
		TriggerID: cronTriggerEvent.TriggerID,
	}
	cronTrigger, err := m.triggerRepo.Get(getDTO)
	if err != nil {
		return err
	}
	if cronTrigger.Status != entity.EnabledTrigger {
		return nil
	}

	return m.fireCronTriggerIfMatchCondition(ctx, err, cronTrigger)
}

func (m *WorkflowTriggerCommandService) fireCronTriggerIfMatchCondition(ctx context.Context,
	err error, cronTrigger *entity.Trigger) error {
	allowDay, err := m.isAllowDay(cronTrigger.AllowDays)
	if err != nil {
		return err
	}
	if !allowDay {
		return nil
	}

	if err := m.createCronTask(cronTrigger); err != nil {
		return err
	}

	// 获取全局上下文
	triggerCtx, err := m.getTriggerCtx(cronTrigger, []byte("{}"))
	realActionArgs, err := m.getRealActionArgs(triggerCtx, cronTrigger.Action.Args, cronTrigger.Action.ActionType)
	actor, ok := m.triggerActorMap[cronTrigger.Action.ActionType]
	if !ok {
		log.Warnf("Not found actor for actionType=%s", cronTrigger.Action.ActionType)
	}

	return actor(ctx, cronTrigger, realActionArgs)
}

// ConsumeTriggerEvent 消费触发器事件
func (m *WorkflowTriggerCommandService) ConsumeTriggerEvent(ctx context.Context, req *dto.TriggerEventDTO) error {
	triggerEvent := req.Message
	triggerEventKey := triggerEvent.Key()
	queryTriggerDTO := &dto.QueryTriggerDTO{Event: triggerEventKey, Status: entity.EnabledTrigger.IntValue()}
	triggers, err := m.triggerRepo.QueryByName(queryTriggerDTO)
	if err != nil {
		return err
	}

	// 没有eventName对应事件信息，直接返回
	if utils.IsZero(triggers) {
		return nil
	}

	for _, eventTrigger := range triggers {
		if err := m.fireEventTriggerIfMatchCondition(eventTrigger, triggerEvent); err != nil {
			log.Warnf("[%s]Failed to fire trigger, caused by %s",
				logs.GetFlowTraceID(eventTrigger.DefID, eventTrigger.InstID), err)
		}
	}

	return nil
}

func (m *WorkflowTriggerCommandService) fireEventTriggerIfMatchCondition(eventTrigger *entity.Trigger,
	triggerEvent pulsar.Message) error {
	// 获取全局上下文
	triggerCtx, err := m.getTriggerCtx(eventTrigger, triggerEvent.Payload())
	// 匹配条件表达式是否成立，不成立直接返回
	match, err := m.exprEvaluator.Match(triggerCtx, eventTrigger.Condition)
	if err != nil {
		return err
	}
	if !match {
		return nil
	}
	allowDay, err := m.isAllowDay(eventTrigger.AllowDays)
	if err != nil {
		return err
	}
	if !allowDay {
		return nil
	}

	realActionArgs, err := m.getRealActionArgs(triggerCtx, eventTrigger.Action.Args, eventTrigger.Action.ActionType)
	actor, ok := m.triggerActorMap[eventTrigger.Action.ActionType]
	if !ok {
		log.Warnf("Not found actor for actionType=%s", eventTrigger.Action.ActionType)
	}

	return actor(context.Background(), eventTrigger, realActionArgs)
}

// isAllowDay 判断当前时间是否可以执行
func (m *WorkflowTriggerCommandService) isAllowDay(allowDaysPolicy entity.AllowDaysPolicy) (bool, error) {
	allowDaysChecker, err := common.NewDefaultAllowDaysChecker()
	if err != nil {
		return false, err
	}

	check, err := allowDaysChecker.Check(time.Now(), allowDaysPolicy)
	if err != nil {
		return false, err
	}

	return check, nil
}

// createCronTask 创建定时任务
func (m *WorkflowTriggerCommandService) createCronTask(cronTrigger *entity.Trigger) error {
	nextTime, err := utils.GetNextTimeByExpr(cronTrigger.Expr, time.Now())
	if err != nil {
		return err
	}

	// 消息间隔时间>=10天
	if nextTime.Unix()-time.Now().Unix() >= constants.TdmqMaxCacheTime {
		addCronJobDTO := &remote.AddCronJobDTO{
			CronStr: cronTrigger.Expr,
			CronTriggerEvent: event.CronTriggerEvent{
				TriggerID:  cronTrigger.TriggerID,
				DefID:      cronTrigger.DefID,
				DefVersion: cronTrigger.DefVersion,
			},
		}

		return m.remoteRepo.AddCronJob(addCronJobDTO)
	}

	cronTriggerEvent := &event.CronTriggerEvent{TriggerID: cronTrigger.TriggerID, DefID: cronTrigger.DefID}
	return m.eventBusRepo.SendCronPresetEvent(context.Background(), nextTime, cronTriggerEvent)
}

// getTriggerCtx 获取触发器上下文
func (m *WorkflowTriggerCommandService) getTriggerCtx(triggerEntity *entity.Trigger, eventValue []byte) (
	map[string]interface{}, error) {
	// 事件上下文
	eventValueMap := map[string]interface{}{}
	if err := json.Unmarshal(eventValue, &eventValueMap); err != nil {
		return nil, err
	}
	eventCtxMap := map[string]interface{}{"event": eventValueMap}

	// 定义级别返回的是定义的上下文
	if triggerEntity.Level == entity.DefTrigger {
		return m.getDefTriggerCtx(triggerEntity.DefID, triggerEntity.DefVersion, eventCtxMap)
	}

	// 定义级别以下的返回的是定义级别以下的上下文
	return m.getInstTriggerCtx(triggerEntity.DefID, triggerEntity.InstID, eventCtxMap)
}

// getDefTriggerCtx 合并事件与流程定义上下文
func (m *WorkflowTriggerCommandService) getDefTriggerCtx(defID string, defVersion int,
	eventCtxMap map[string]interface{}) (map[string]interface{}, error) {
	workflowDefCtx, err := m.buildWorkflowDefCtx(defID, defVersion)
	if err != nil {
		return nil, err
	}

	return utils.MergeMap(workflowDefCtx, eventCtxMap)
}

// buildWorkflowDefCtx 构建流程定义上下文
func (m *WorkflowTriggerCommandService) buildWorkflowDefCtx(defID string,
	defVersion int) (map[string]interface{}, error) {
	getWorkflowDefDTO := &dto.GetWorkflowDefDTO{
		DefID:   defID,
		Version: defVersion,
		Status:  entity.Enabled.IntValue(),
	}

	workflowDef, err := m.workflowDefRepo.Get(getWorkflowDefDTO)
	if err != nil {
		return nil, err
	}

	return entity.DefConvertToCtx(workflowDef)
}

// getInstTriggerCtx 合并流程实例与外部事件上下文
func (m *WorkflowTriggerCommandService) getInstTriggerCtx(defID, instID string,
	eventCtxMap map[string]interface{}) (map[string]interface{}, error) {
	// 实例级别返回合并外部事件后的上下文
	getWorkflowInstDTO := &dto.GetWorkflowInstDTO{DefID: defID, InstID: instID}
	workflowInstCtx, err := m.workflowInstRepo.GetWorkflowInstCtx(getWorkflowInstDTO)
	if err != nil {
		return nil, err
	}

	return utils.MergeMap(workflowInstCtx, eventCtxMap)
}

// getRealActionArgs 解析actionMap
func (m *WorkflowTriggerCommandService) getRealActionArgs(ctxMap map[string]interface{},
	actionArgs interface{}, actionType entity.ActionType) (interface{}, error) {
	actionArgsMap, err := utils.StructToMap(actionArgs)
	if err != nil {
		return nil, err
	}
	parsedActionArgsMap, err := m.exprEvaluator.EvaluateMap(ctxMap, actionArgsMap)
	if err != nil {
		return nil, err
	}

	return entity.GetActionArgs(parsedActionArgsMap, actionType)
}

// CronCallBack 定时回调
func (m *WorkflowTriggerCommandService) CronCallBack(ctx context.Context, params string) error {
	// 解析参数
	cronTriggerEvent := event.CronTriggerEvent{}
	if err := json.Unmarshal([]byte(params), &cronTriggerEvent); err != nil {
		return err
	}

	getTriggerDTO := &dto.GetTriggerDTO{
		DefID:     cronTriggerEvent.DefID,
		TriggerID: cronTriggerEvent.TriggerID,
	}
	cronTrigger, err := m.triggerRepo.Get(getTriggerDTO)
	if err != nil {
		return err
	}

	if cronTrigger.Status != entity.EnabledTrigger {
		jobName := fmt.Sprintf("%s:%d", constants.PxCronJobNamePrefix, cronTriggerEvent.TriggerID)
		return m.remoteRepo.CancelCronJob(jobName)
	}

	return m.fireCronTriggerIfMatchCondition(ctx, err, cronTrigger)
}
