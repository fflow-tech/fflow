package execution

import (
	"fmt"
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/cache"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/expr"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// needDelay 判断执行前是否需要等待
func needDelay(nodeInst *entity.NodeInst) bool {
	waitConf := nodeInst.BasicNodeDef.Wait
	return waitConf.Duration != "" || waitConf.Expr != ""
}

// needRetry 判断节点是否重试
func needRetry(curNodeInst *entity.NodeInst) bool {
	return curNodeInst.Status == entity.NodeInstFailed && curNodeInst.BasicNodeDef.Retry.Count > curNodeInst.RetryCount
}

// getRetryDelay 获取重试的延迟时间
func getRetryDelay(curNodeInst *entity.NodeInst) (time.Duration, error) {
	// 获取重试延迟时间配置 这里的时间校验不在这里做 应该放在配置读取时就做
	retryDelay, err := expr.ParseDuration(curNodeInst.BasicNodeDef.Retry.Duration)
	if err != nil {
		return 0, err
	}
	// 获取真正重试等待时间
	if curNodeInst.BasicNodeDef.Retry.Policy == entity.ExponentialBackoff {
		// 当为指数退避时 则每次时间跟次数相乘 这里 +1 是因为 RetryCount 初始值为0
		retryDelay = time.Duration(curNodeInst.RetryCount+1) * retryDelay
	}
	return retryDelay, nil
}

// appendExecutePath 添加执行路径
func appendExecutePath(inst *entity.WorkflowInst, decideResult *entity.DecideResult) {
	if len(inst.ExecutePath) == 0 {
		inst.ExecutePath = [][]string{}
	}

	if len(decideResult.NodesToBeScheduled) > 0 {
		inst.ExecutePath = append(inst.ExecutePath, entity.GetNodeRefNames(decideResult.NodesToBeScheduled))
	}
}

// getMaxNodeInstID 获取最大的节点实例 ID
func getMaxNodeInstID(inst *entity.WorkflowInst) string {
	nodeInsts := inst.SchedNodeInsts
	if len(nodeInsts) == 0 {
		return ""
	}

	maxID := inst.BeforeLastRestartMaxNodeInstID
	for _, nodeInst := range nodeInsts {
		maxID = utils.MaxUint64Str(nodeInst.NodeInstID, maxID)
	}

	return maxID
}

// getWorkflowInstLockName 获取流程实例锁的名称
func getWorkflowInstLockName(instID string) string {
	return strings.Join([]string{utils.GetEnv(), "inst", instID}, ":")
}

// getWorkflowDefLockName 获取流程定义锁的名称
func getWorkflowDefLockName(defID string) string {
	return strings.Join([]string{utils.GetEnv(), "def", defID}, ":")
}

// matchCondition 判断条件是否匹配
func matchCondition(evaluator expr.Evaluator, inst *entity.WorkflowInst, condition string) (bool, error) {
	ctx, err := entity.ConvertToCtx(inst)
	if err != nil {
		return false, err
	}

	match, err := evaluator.Match(ctx, condition)
	if err != nil {
		return false, err
	}

	return match, nil
}

// matchConditionForCurNodeInst 判断当前节点条件是否匹配，这样做主要是可以通过 this 的方法引用当前节点的上下文
func matchConditionForCurNodeInst(evaluator expr.Evaluator, inst *entity.WorkflowInst,
	curNodeInst *entity.NodeInst, condition string) (bool, error) {
	ctx, err := entity.ConvertToCtx(inst)
	if err != nil {
		return false, err
	}

	if err := entity.AppendNodeInfoToCtx(ctx, curNodeInst); err != nil {
		return false, err
	}

	// 追加当前节点实例的信息到上下文
	if err := entity.AppendNodeInfoToCtxKey(ctx, curNodeInst, constants.ThisNode); err != nil {
		return false, err
	}

	match, err := evaluator.Match(ctx, condition)
	if err != nil {
		return false, err
	}

	return match, nil
}

// evaluateMapForCurNodeInst 根据当前节点实例计算值
func evaluateMapForCurNodeInst(evaluator expr.Evaluator, inst *entity.WorkflowInst,
	curNodeInst *entity.NodeInst, oldMap map[string]interface{}) (map[string]interface{}, error) {
	ctx, err := entity.ConvertToCtx(inst)
	if err != nil {
		return nil, err
	}

	if err := entity.AppendNodeInfoToCtx(ctx, curNodeInst); err != nil {
		return nil, err
	}

	// 追加当前节点实例的信息到上下文
	if err := entity.AppendNodeInfoToCtxKey(ctx, curNodeInst, constants.ThisNode); err != nil {
		return nil, err
	}

	return evaluator.EvaluateMap(ctx, oldMap)
}

// GetInstDistributeLock 获取流程实例操作的分布式锁
func GetInstDistributeLock(cacheRepo ports.CacheRepository, instID string) (cache.DistributeLock, error) {
	lockName := getWorkflowInstLockName(instID)
	lock := cacheRepo.GetDistributeLockWithRetry(lockName, instLockExpireTime, instLockRetry, instLockRetryDelay)
	if err := lock.Lock(); err != nil {
		return nil, fmt.Errorf(notGetLockErr, instID, lockName, err, errno.Unavailable)
	}
	return lock, nil
}

// GetDefDistributeLock 获取流程定义操作的分布式锁
func GetDefDistributeLock(cacheRepo ports.CacheRepository, defID string) (cache.DistributeLock, error) {
	lockName := getWorkflowDefLockName(defID)
	lock := cacheRepo.GetDistributeLockWithRetry(lockName, defLockExpireTime, defLockRetry, defLockRetryDelay)
	if err := lock.Lock(); err != nil {
		return nil, fmt.Errorf(notGetLockErr, defID, lockName, err, errno.Unavailable)
	}
	return lock, nil
}

// getInstTimeoutSendAlertKey 获取实例超时告警是否已已经发送的标记
func getInstTimeoutSendAlertKey(inst *entity.WorkflowInst, template common.MsgTemplate) string {
	return strings.Join([]string{utils.GetEnv(),
		string(template),
		inst.InstID},
		":")
}

// getNodeTimeoutSendAlertKey 获取节点超时告警是否已已经发送的标记
func getNodeTimeoutSendAlertKey(nodeInst *entity.NodeInst, template common.MsgTemplate) string {
	return strings.Join([]string{utils.GetEnv(),
		string(template),
		nodeInst.NodeInstID},
		":")
}
