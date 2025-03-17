package execution

import (
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"

	"github.com/panjf2000/ants/v2"
)

var (
	archiveDefaultInitPageIndex = 1
	archiveDefaultPageSize      = 50
	archiveMaxQueryTimes        = 1000
)

// HistoryArchiver 历史归档器
type HistoryArchiver interface {
	Archive(*dto.ArchiveHistoryWorkflowInstsDTO) error // 历史归档
}

// DefaultHistoryArchiver 默认的历史归档器
type DefaultHistoryArchiver struct {
	workflowDefRepo     ports.WorkflowDefRepository
	workflowInstRepo    ports.WorkflowInstRepository
	nodeInstRepo        ports.NodeInstRepository
	workflowArchiveRepo ports.WorkflowArchiveRepository
	goroutinePool       *ants.Pool
}

// NewDefaultHistoryArchiver 新建历史归档器
func NewDefaultHistoryArchiver(repoProviderSet *ports.RepoProviderSet) (*DefaultHistoryArchiver, error) {
	goroutinePool, err := ants.NewPool(config.GetHistoryArchiverConfig().GoroutinePoolSize)
	if err != nil {
		return nil, err
	}

	return &DefaultHistoryArchiver{
		workflowDefRepo:     repoProviderSet.WorkflowDefRepo(),
		workflowInstRepo:    repoProviderSet.WorkflowInstRepo(),
		nodeInstRepo:        repoProviderSet.NodeInstRepo(),
		workflowArchiveRepo: repoProviderSet.WorkflowArchiveRepo(),
		goroutinePool:       goroutinePool,
	}, nil
}

// Archive 归档
func (m *DefaultHistoryArchiver) Archive(req *dto.ArchiveHistoryWorkflowInstsDTO) error {
	return m.doArchive(req)
}

func (m *DefaultHistoryArchiver) doArchive(req *dto.ArchiveHistoryWorkflowInstsDTO) error {
	query := &dto.PageQueryWorkflowDefDTO{
		PageQuery:     constants.NewPageQuery(defaultInitPageIndex, defaultPageSize),
		DefID:         req.DefID,
		ReadFromSlave: true,
	}

	for i := defaultInitPageIndex; i <= archiveMaxQueryTimes; i++ {
		query.PageIndex = i
		workflowDefList, err := m.workflowDefRepo.PageQueryLastVersion(query)
		if err != nil {
			archiveLog().Errorf("Failed to page query workflow def, caused by %s", err)
			break
		}
		if len(workflowDefList) == 0 {
			break
		}
		for _, workflowDef := range workflowDefList {
			workflowDefCopy := workflowDef
			m.goroutinePool.Submit(func() {
				if !config.GetHistoryArchiverConfig().EnableArchive {
					log.Infof("History archiving is disabled, skip it, defID=%d", workflowDef.DefID)
					return
				}

				if err := m.archiveOneDef(workflowDefCopy); err != nil {
					archiveLog().Warnf("Failed to archive workflow=%d, caused by %+v", workflowDefCopy.DefID, err) // 非关键异常，只记录警告
				}
			})
		}
	}
	return nil
}

func (m *DefaultHistoryArchiver) archiveOneDef(workflowDef *entity.WorkflowDef) error {
	startTime := time.Now()
	archiveLog().Infof("Start to archive def=%d", workflowDef.DefID)
	defer func() {
		archiveLog().Infof("Finish to archive def=%d, costs %dms",
			workflowDef.DefID, time.Since(startTime).Milliseconds())
	}()

	pageQueryInsts := &dto.GetWorkflowInstListDTO{
		DefID:         workflowDef.DefID,
		CreatedBefore: getCreatedBefore(),
		Status:        entity.InstSucceed,
		PageQuery:     constants.NewPageQuery(archiveDefaultInitPageIndex, archiveDefaultPageSize),
		ReadFromSlave: true,
	}

	for i := defaultInitPageIndex; i <= maxQueryTimes; i++ {
		insts, _, err := m.workflowInstRepo.PageQuery(pageQueryInsts)
		if err != nil {
			archiveLog().Errorf("Failed to page query insts, caused by %s", err)
			return err
		}
		if len(insts) == 0 {
			break
		}

		for _, inst := range insts {
			if err := m.archiveOneInst(inst); err != nil {
				archiveLog().Infof("Failed to archive def=%d inst=%d, caused by %s",
					workflowDef.DefID, inst.InstID, err)
				continue
			}
		}
	}

	return nil
}

func getCreatedBefore() time.Time {
	return time.Now().Add(time.Duration(config.GetHistoryArchiverConfig().KeepDataDuration) * -24 * time.Hour)
}

func (m *DefaultHistoryArchiver) archiveOneInst(inst *entity.WorkflowInst) (err error) {
	if inst == nil || inst.WorkflowDef == nil {
		return fmt.Errorf("failed to archive one inst, caused by inst or def is nil")
	}

	startTime := time.Now()
	archiveLog().Infof("Start to archive def=%d, inst=%d...", inst.WorkflowDef.DefID, inst.InstID)
	defer logs.DumpPanicStack("HistoryArchiver",
		fmt.Errorf("failed to archive inst=%d: %w", inst.InstID, err))
	defer func() {
		archiveLog().Infof("Finish to archive def=%d, inst=%d, costs %dms",
			inst.WorkflowDef.DefID, inst.InstID, time.Since(startTime).Milliseconds())
	}()

	if err := m.workflowArchiveRepo.ArchiveWorkflowInst(&dto.ArchiveWorkflowInstsDTO{
		DefID:   inst.WorkflowDef.DefID,
		InstIDs: []string{inst.InstID},
	}); err != nil {
		archiveLog().Infof("Failed to archive def=%d inst=%s, caused by %s",
			inst.WorkflowDef.DefID, inst.InstID, err)
		return err
	}

	return m.archiveInstNodeInsts(inst)
}

func (m *DefaultHistoryArchiver) archiveInstNodeInsts(inst *entity.WorkflowInst) error {
	pageQueryNodeInsts := &dto.PageQueryNodeInstDTO{
		DefID:         inst.WorkflowDef.DefID,
		InstID:        inst.InstID,
		PageQuery:     constants.NewPageQuery(defaultInitPageIndex, defaultPageSize),
		ReadFromSlave: true,
	}

	for i := defaultInitPageIndex; i <= maxQueryTimes; i++ {
		nodeInsts, err := m.nodeInstRepo.PageQuery(pageQueryNodeInsts)
		if err != nil {
			return err
		}
		if len(nodeInsts) == 0 {
			break
		}

		var nodeInstIDs []string
		for _, nodeInst := range nodeInsts {
			nodeInstIDs = append(nodeInstIDs, nodeInst.NodeInstID)
		}
		if err := m.workflowArchiveRepo.ArchiveNodeInst(&dto.ArchiveNodeInstsDTO{
			DefID:       inst.WorkflowDef.DefID,
			InstID:      inst.InstID,
			NodeInstIDs: nodeInstIDs,
		},
		); err != nil {
			archiveLog().Infof("Failed to archive inst=%d, caused by %s", inst.InstID, err)
			continue
		}
	}
	return nil
}

func archiveLog() log.Logger {
	return log.GetDefaultLogger()
}
