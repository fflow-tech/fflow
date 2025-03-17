package execution

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution/common"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"

	"github.com/panjf2000/ants/v2"
)

var (
	defaultInitPageIndex = 1
	defaultPageSize      = 20
	maxQueryTimes        = 1000
)

// TimeoutChecker 超时检查器
type TimeoutChecker interface {
	CheckAll() error
}

// DefaultTimeoutChecker 默认的超时检查器
type DefaultTimeoutChecker struct {
	workflowDefRepo    ports.WorkflowDefRepository
	workflowInstRepo   ports.WorkflowInstRepository
	nodeInstRepo       ports.NodeInstRepository
	instTimeoutChecker common.InstTimeoutChecker
	workflowRunner     WorkflowRunner
	nodeRunner         NodeRunner
	goroutinePool      *ants.Pool
}

func checkerLog() log.Logger {
	return log.GetDefaultLogger()
}

// NewDefaultTimeoutChecker 新建超时检查器
func NewDefaultTimeoutChecker(repoProviderSet *ports.RepoProviderSet,
	workflowProviderSet *WorkflowProviderSet,
	workflowRunner *DefaultWorkflowRunner,
	nodeRunner *DefaultNodeRunner) (*DefaultTimeoutChecker, error) {
	goroutinePool, err := ants.NewPool(config.GetTimeoutCheckerConfig().GoroutinePoolSize)
	if err != nil {
		return nil, err
	}

	return &DefaultTimeoutChecker{workflowDefRepo: repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo:   repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:       repoProviderSet.NodeInstRepo(),
		instTimeoutChecker: workflowProviderSet.InstTimeoutChecker(),
		workflowRunner:     workflowRunner, nodeRunner: nodeRunner,
		goroutinePool: goroutinePool,
	}, nil
}

// CheckAll 检查超时
func (m *DefaultTimeoutChecker) CheckAll() error {
	return m.doCheckAll()
}

func (m *DefaultTimeoutChecker) doCheckAll() error {
	query := &dto.PageQueryWorkflowDefDTO{
		PageQuery:     constants.NewPageQuery(defaultInitPageIndex, defaultPageSize),
		ReadFromSlave: true,
	}

	for i := defaultInitPageIndex; i <= maxQueryTimes; i++ {
		query.PageIndex = i
		workflowDefList, err := m.workflowDefRepo.PageQueryLastVersion(query)
		if err != nil {
			checkerLog().Errorf("Failed to page query workflow def, caused by %s", err)
			break
		}
		if len(workflowDefList) == 0 {
			break
		}
		for _, workflowDef := range workflowDefList {
			workflowDefCopy := workflowDef
			m.goroutinePool.Submit(func() {
				if err := m.checkOneDef(workflowDefCopy); err != nil {
					checkerLog().Warnf("Failed to check workflow=%d, caused by %+v", workflowDef.DefID, err) // 非关键异常，只记录警告
				}
			})
		}
	}
	return nil
}

func (m *DefaultTimeoutChecker) checkOneDef(workflowDef *entity.WorkflowDef) error {
	startTime := time.Now()
	checkerLog().Infof("Start to check def=%d", workflowDef.DefID)
	defer func() {
		checkerLog().Infof("Finish to check def=%d, costs %dms", workflowDef.DefID, time.Since(startTime).Milliseconds())
	}()

	pageQueryInsts := &dto.GetWorkflowInstListDTO{
		DefID:         workflowDef.DefID,
		PageQuery:     constants.NewPageQuery(defaultInitPageIndex, defaultPageSize),
		ReadFromSlave: true,
	}

	for i := defaultInitPageIndex; i <= maxQueryTimes; i++ {
		pageQueryInsts.PageIndex = i
		insts, _, err := m.workflowInstRepo.PageQuery(pageQueryInsts)
		if err != nil {
			checkerLog().Errorf("Failed to page query insts, caused by %s", err)
			return err
		}
		if len(insts) == 0 {
			break
		}

		for _, inst := range insts {
			if err := m.checkOneInst(inst); err != nil {
				checkerLog().Infof("Failed to check inst=%d, caused by %s", inst.InstID, err)
				continue
			}
		}
	}

	return nil
}

func (m *DefaultTimeoutChecker) checkOneInst(inst *entity.WorkflowInst) (err error) {
	if inst == nil || inst.WorkflowDef == nil {
		return fmt.Errorf("failed to check one inst, caused by inst or def is nil")
	}

	startTime := time.Now()
	checkerLog().Infof("Start to check inst=%d...", inst.InstID)
	defer logs.DumpPanicStack("TimeoutChecker",
		fmt.Errorf("failed to check inst=%d: %w", inst.InstID, err))
	defer func() {
		checkerLog().Infof("Finish to check inst=%d, costs %dms", inst.InstID, time.Since(startTime).Milliseconds())
	}()

	pageQueryNodeInsts := &dto.PageQueryNodeInstDTO{
		DefID:  inst.WorkflowDef.DefID,
		InstID: inst.InstID,
		Statuses: []entity.NodeInstStatus{entity.NodeInstRunning,
			entity.NodeInstScheduled,
			entity.NodeInstPaused,
			entity.NodeInstWaiting},
		PageQuery:     constants.NewPageQuery(defaultInitPageIndex, defaultPageSize),
		ReadFromSlave: true,
	}

	for i := defaultInitPageIndex; i <= maxQueryTimes; i++ {
		pageQueryNodeInsts.PageIndex = i
		nodeInsts, err := m.nodeInstRepo.PageQuery(pageQueryNodeInsts)
		if err != nil {
			return err
		}
		if len(nodeInsts) == 0 {
			break
		}

		for _, nodeInst := range nodeInsts {
			if err := m.checkOneNodeInst(nodeInst); err != nil {
				checkerLog().Warnf("Failed to check inst=%d nodeInst=%d, caused by %s",
					nodeInst.InstID, nodeInst.NodeInstID, err)
				continue
			}
		}
	}

	return m.checkInstTimeout(inst)
}

func (m *DefaultTimeoutChecker) checkInstTimeout(inst *entity.WorkflowInst) error {
	instTimeout, err := m.instTimeoutChecker.CheckWorkflowInst(inst)
	if err != nil {
		return err
	}
	if !instTimeout {
		return nil
	}

	setWorkflowInstTimeoutReq := &dto.SetWorkflowInstTimeoutDTO{
		DefID:  inst.WorkflowDef.DefID,
		InstID: inst.InstID,
		Reason: "workflow inst execute timeout",
	}
	checkerLog().Infof("The inst=%d is timeout", inst.InstID)
	return m.workflowRunner.SetTimeout(context.Background(), setWorkflowInstTimeoutReq)
}

func (m *DefaultTimeoutChecker) checkOneNodeInst(nodeInst *entity.NodeInst) error {
	err := m.checkNodeInstTimeout(nodeInst)
	if err != nil {
		return err
	}

	return m.checkNodeInstNearTimeout(nodeInst, err)
}

func (m *DefaultTimeoutChecker) checkNodeInstTimeout(nodeInst *entity.NodeInst) error {
	timeout, err := m.instTimeoutChecker.CheckNodeInst(nodeInst)
	if err != nil {
		return err
	}
	if !timeout {
		return nil
	}

	setTimeoutReq := &dto.SetNodeTimeoutDTO{
		DefID:       nodeInst.DefID,
		InstID:      nodeInst.InstID,
		NodeInstID:  nodeInst.NodeInstID,
		NodeRefName: nodeInst.BasicNodeDef.RefName,
		Operator:    "system",
		Reason:      "node execute timeout",
	}
	checkerLog().Infof("The inst=%d nodeInst=%d is timeout", nodeInst.InstID, nodeInst.NodeInstID)
	return m.nodeRunner.SetTimeout(context.Background(), setTimeoutReq)
}

func (m *DefaultTimeoutChecker) checkNodeInstNearTimeout(nodeInst *entity.NodeInst, err error) error {
	nearTimeout, err := m.instTimeoutChecker.CheckNodeInstNearTimeout(nodeInst)
	if err != nil {
		return err
	}
	if !nearTimeout {
		return nil
	}

	setNearTimeoutReq := &dto.SetNodeNearTimeoutDTO{
		DefID:       nodeInst.DefID,
		InstID:      nodeInst.InstID,
		NodeInstID:  nodeInst.NodeInstID,
		NodeRefName: nodeInst.BasicNodeDef.RefName,
		Operator:    "system",
		Reason:      "node execute near timeout",
	}

	checkerLog().Infof("The inst=%d nodeInst=%d is near timeout", nodeInst.InstID, nodeInst.NodeInstID)
	return m.nodeRunner.SetNearTimeout(context.Background(), setNearTimeoutReq)
}
