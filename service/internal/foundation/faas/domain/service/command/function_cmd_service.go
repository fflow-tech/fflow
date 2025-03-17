package command

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/internal/foundation/faas/pkg/runtimecontext"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/localcache"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

const defaultVersion = 1

// FunctionCommandService 写服务
type FunctionCommandService struct {
	functionRepo     ports.FunctionRepository
	functionExecutor *execution.CodeExecutor
	cache            localcache.Client
}

// NewFunctionCommandService 新建服务
func NewFunctionCommandService(repoProviderSet *ports.RepoProviderSet,
	functionExecutor *execution.CodeExecutor) (*FunctionCommandService, error) {
	cache, err := localcache.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &FunctionCommandService{
		functionRepo:     repoProviderSet.FunctionRepo(),
		functionExecutor: functionExecutor,
		cache:            cache,
	}, nil
}

func getCacheKey(req *dto.CallFunctionReqDTO) string {
	return fmt.Sprintf("func:%s:%s", req.Namespace, req.Function)
}

func (m *FunctionCommandService) getFunction(req *dto.CallFunctionReqDTO) (*entity.Function, error) {
	f := entity.Function{}
	// 从缓存中获取失败的错误直接忽略
	err := m.cache.Get(getCacheKey(req), &f)
	if err == nil {
		return &f, nil
	}

	if err != nil {
		log.Infof("Get [%s:%s] function info from cache failed: %v", req.Namespace, req.Function, err)
	}

	getFuncDTO, err := m.convertCallDTOToGetDTO(req)
	if err != nil {
		return nil, err
	}
	function, err := m.functionRepo.Get(getFuncDTO)
	if err != nil {
		return nil, err
	}
	_ = m.cache.Append(getCacheKey(req), function)
	return function, nil
}

// CallFunction 执行函数
func (m *FunctionCommandService) CallFunction(ctx context.Context, req *dto.CallFunctionReqDTO) (interface{}, error) {
	function, err := m.getFunction(req)
	if err != nil {
		return nil, err
	}

	// 执行函数，避免消耗太多资源，打日志即可
	runtimeCtx := runtimecontext.NewRuntimeContext(withDebugMode(ctx, false), function, req.Request)
	result, _, err := m.functionExecutor.Execute(runtimeCtx, req, function)

	log.Infof("Run function [%s:%s] req: %+v, result: %+v, err: %s",
		function.Namespace, function.Name, req, result, err)

	return result, err
}

// convertCallDTOToGetDTO 获取函数所属的服务,并返回 GetFunctionDTO
func (m *FunctionCommandService) convertCallDTOToGetDTO(c *dto.CallFunctionReqDTO) (*dto.GetFunctionReqDTO, error) {
	return &dto.GetFunctionReqDTO{
		Namespace: c.Namespace,
		Function:  c.Function,
	}, nil
}

// DebugFunction 调试函数
func (m *FunctionCommandService) DebugFunction(ctx context.Context, req *dto.DebugFunctionDTO) (
	*dto.DebugFunctionRspDTO, error) {
	function := &entity.Function{
		Namespace: req.Namespace,
		Code:      req.Code,
		Language:  entity.GetLanguageTypeByStrValue(req.Language),
	}

	// 执行函数
	runtimeCtx := runtimecontext.NewRuntimeContext(withDebugMode(ctx, true), function, nil)
	result, funcLogs, err := m.functionExecutor.Debug(runtimeCtx, function, req.Input)

	// 执行的错误从接口中返回便于调用方处理
	if err != nil {
		return &dto.DebugFunctionRspDTO{Result: result, Logs: funcLogs, Error: err.Error()}, nil
	}
	return &dto.DebugFunctionRspDTO{Result: result, Logs: funcLogs, Error: ""}, nil
}

// CreateFunction 创建函数方法
func (m *FunctionCommandService) CreateFunction(d *dto.CreateFunctionReqDTO) (uint, error) {
	function, _ := m.functionRepo.Get(&dto.GetFunctionReqDTO{
		Namespace: d.Namespace,
		Function:  d.Function,
	})
	if function != nil {
		// 创建函数时不允许重名
		return 0, fmt.Errorf("the function name[%s] existed in this namespace[%s]", d.Function, d.Namespace)
	}

	createFunctionDTO := &dto.CreateFunctionDTO{
		Namespace:    d.Namespace,
		Function:     d.Function,
		Language:     entity.GetLanguageTypeByStrValue(d.Language),
		Code:         d.Code,
		Creator:      d.Creator,
		Updater:      d.Creator,
		Description:  d.Description,
		InputSchema:  d.InputSchema,
		OutputSchema: d.OutputSchema,
		Version:      defaultVersion,
		Token:        utils.GenerateToken(),
	}
	return m.functionRepo.Create(createFunctionDTO)
}

// UpdateFunction 更新函数方法
func (m *FunctionCommandService) UpdateFunction(d *dto.UpdateFunctionDTO) (uint, error) {
	function, err := m.functionRepo.Get(&dto.GetFunctionReqDTO{
		Namespace: d.Namespace,
		Function:  d.Function,
	})
	if err != nil {
		return 0, err
	}
	// 创建一个新的版本
	newFunction := &dto.CreateFunctionDTO{
		Namespace:    d.Namespace,
		Function:     d.Function,
		Code:         d.Code,
		Description:  d.Description,
		Version:      function.Version + 1,
		Token:        function.Token,
		Updater:      d.Updater,
		Creator:      function.Creator,
		InputSchema:  d.InputSchema,
		OutputSchema: d.OutputSchema,
		Language:     function.Language,
	}
	return m.functionRepo.Create(newFunction)
}

// DeleteFunction 删除函数方法
func (m *FunctionCommandService) DeleteFunction(d *dto.DeleteFunctionDTO) error {
	return m.functionRepo.Delete(d)
}

// addHistory 添加执行历史记录
func (m *FunctionCommandService) addHistory(req *dto.CallFunctionReqDTO, function *entity.Function) (uint, error) {
	input, err := json.Marshal(req.Input)
	if err != nil {
		return 0, err
	}
	history, err := m.functionRepo.CreateRunHistory(&dto.CreateRunHistoryDTO{
		Namespace:    function.Namespace,
		FunctionName: function.Name,
		Version:      uint(function.Version),
		Operator:     req.Operator,
		Input:        string(input),
		Status:       string(entity.Running),
	})
	if err != nil {
		return 0, err
	}
	return history.ID, nil
}

// updateHistory 更新执行历史
func (m *FunctionCommandService) updateHistory(history *dto.UpdateRunHistoryDTO) {
	// 更新执行结果,报错不影响函数执行结果的返回
	historyErr := m.functionRepo.UpdateRunHistory(history)
	if historyErr != nil {
		log.Infof("update run history error, caused by %s", historyErr.Error())
	}
}

const (
	defaultInitPageIndex = 1
	defaultPageSize      = 100
	maxQueryTimes        = 10000
	minKeepDays          = 1
)

// BatchDeleteExpiredRunHistory 批量删除过期的执行历史
func (m *FunctionCommandService) BatchDeleteExpiredRunHistory(req *dto.BatchDeleteExpiredRunHistoryDTO) error {
	if req.KeepDays < minKeepDays {
		req.KeepDays = minKeepDays
	}

	expireTime := time.Now().Add(time.Duration(-req.KeepDays) * 24 * time.Hour)

	query := &dto.PageQueryRunHistoryDTO{
		PageQuery: constants.NewPageQuery(defaultInitPageIndex, defaultPageSize),
		IDs:       req.IDs,
		CreatedAt: expireTime,
	}

	curMaxID := uint(math.MaxUint32)
	for i := defaultInitPageIndex; i <= maxQueryTimes; i++ {
		query.MaxID = curMaxID
		histories, _, err := m.functionRepo.PageQueryRunHistory(query)
		if err != nil {
			return err
		}
		if len(histories) == 0 {
			break
		}

		ids := []uint{}
		for _, history := range histories {
			ids = append(ids, history.ID)
			curMaxID = utils.MinUint(curMaxID, history.ID)
		}

		startTime := time.Now()
		log.Infof("Start to batch delete run history ids: %v", ids)
		if err := m.functionRepo.BatchDeleteRunHistory(&dto.BatchDeleteRunHistoryDTO{IDs: ids}); err != nil {
			archiveLog().Infof("Batch delete run history=%d err: %v", ids, err)
			continue
		}
		log.Infof("Finish to batch delete run history ids: %v, costs %dms", ids,
			time.Since(startTime).Milliseconds())
	}
	return nil
}

func archiveLog() log.Logger {
	return log.GetDefaultLogger()
}

// withDebugMode 获取一个携带 debugMode 标记的 Context
// 注意 context 只能自上而下携带值，若后续需要扩展，
func withDebugMode(ctx context.Context, debugMode bool) context.Context {
	return context.WithValue(ctx, "debugMode", debugMode)
}
