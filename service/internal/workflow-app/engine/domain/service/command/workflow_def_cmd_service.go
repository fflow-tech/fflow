package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/execution"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/trigger"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/service/command/validator"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

const (
	defaultDefInitVersion  = 1 // 默认定义初始版本号
	versionUpgradeStepSize = 1 // 版本更新步长
)

// WorkflowDefCommandService 写服务
type WorkflowDefCommandService struct {
	workflowDefRepo ports.WorkflowDefRepository
	workflowUpdater execution.WorkflowUpdater
	triggerRegistry trigger.Registry
	cacheRepo       ports.CacheRepository
}

// NewWorkflowDefCommandService 新建服务
func NewWorkflowDefCommandService(repoProviderSet *ports.RepoProviderSet,
	workflowUpdater *execution.DefaultWorkflowUpdater,
	triggerRegistry *trigger.DefaultRegistry) *WorkflowDefCommandService {
	return &WorkflowDefCommandService{
		workflowDefRepo: repoProviderSet.WorkflowDefRepo(),
		cacheRepo:       repoProviderSet.CacheRepo(),
		workflowUpdater: workflowUpdater,
		triggerRegistry: triggerRegistry,
	}
}

// CreateWorkflowDef 创建工作流定义
func (m *WorkflowDefCommandService) CreateWorkflowDef(ctx context.Context,
	req *dto.CreateWorkflowDefDTO) (string, error) {
	// 检查传入content格式
	if err := validator.ValidateDefJson(req.DefJson); err != nil {
		return "", err
	}
	// DAO层创建工作流定义
	req.Version = defaultDefInitVersion
	if err := m.initCreateWorkflowDefDTO(req); err != nil {
		return "", err
	}
	defID, err := m.workflowDefRepo.Create(req)
	if err != nil {
		log.Errorf("Failed to create workflow def, caused by %s", err)
		return "", err
	}
	// 创建子流程定义，如果子流程定义创建失败，则直接返回失败，不需要注册和发送事件
	if err := m.CreateSubworkflowDef(ctx, &dto.CreateSubworkflowDefDTO{
		ParentDefID:      req.DefID,
		ParentDefVersion: req.Version,
		DefJson:          req.DefJson}); err != nil {
		log.Errorf("[%d]Failed to create subworkflow def, caused by %s", defID, err)
		return "", err
	}

	// 发送创建事件
	if err := m.workflowUpdater.SendWorkflowDefExternalEvent(req.Namespace, req.DefID, event.DefCreate); err != nil {
		return "", err
	}
	return defID, nil
}

func (m *WorkflowDefCommandService) initCreateWorkflowDefDTO(req *dto.CreateWorkflowDefDTO) error {
	defJsonWorkflowEntity, err := m.getDefJsonWorkflowEntity(req.DefJson)
	if err != nil {
		return err
	}
	if req.Namespace == "" {
		req.Namespace = defJsonWorkflowEntity.Namespace
	}
	m.initWorkflowDefName(req, defJsonWorkflowEntity)
	m.initWorkflowDesc(req, defJsonWorkflowEntity)
	return nil
}

func (m *WorkflowDefCommandService) initWorkflowDefName(req *dto.CreateWorkflowDefDTO,
	workflowEntity *entity.WorkflowDef) {
	if req.Name != "" {
		// 当req中带有name时则直接使用
		return
	}
	req.Name = workflowEntity.Name
	return
}

func (m *WorkflowDefCommandService) initWorkflowDesc(req *dto.CreateWorkflowDefDTO,
	workflowEntity *entity.WorkflowDef) {
	if req.Description != "" {
		// 当req中带有Description时则直接使用
		return
	}
	req.Description = workflowEntity.Desc
	return
}

func (m *WorkflowDefCommandService) getDefJsonWorkflowEntity(defJson string) (*entity.WorkflowDef, error) {
	workflowDefEntity := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(defJson), workflowDefEntity); err != nil {
		return nil, err
	}
	return workflowDefEntity, nil
}

// CreateWorkflowDefs 批量创建工作流定义
// 版本都是默认的版本，目前创建子流程时会调用这个方法
func (m *WorkflowDefCommandService) CreateWorkflowDefs(ctx context.Context, reqs []*dto.CreateWorkflowDefDTO) error {
	// 检查传入的content格式，并初始化定义版本
	for _, req := range reqs {
		if err := validator.ValidateDefJson(req.DefJson); err != nil {
			return err
		}
		req.Version = defaultDefInitVersion
	}

	// DAO层批量创建工作流定义
	err := m.workflowDefRepo.BatchCreate(reqs)
	if err != nil {
		log.Errorf("Failed to batch create def, caused by %s", err)
		return err
	}

	for _, req := range reqs {
		// 注册定时触发事件
		if err := m.doRegisterTrigger(req.DefID, req.Version, req.DefJson); err != nil {
			return fmt.Errorf("failed to register trigger: %w", err)
		}
		// 发送创建事件
		if err := m.workflowUpdater.SendWorkflowDefExternalEvent(req.Namespace, req.DefID, event.DefCreate); err != nil {
			return err
		}
	}

	return nil
}

// CreateSubworkflowDef 创建子流程定义
func (m *WorkflowDefCommandService) CreateSubworkflowDef(ctx context.Context, req *dto.CreateSubworkflowDefDTO) error {
	subworkflows, subWorkflowSlice, err := getSubworkflowsEntityAndSlice(req.DefJson)
	if err != nil {
		return err
	}
	// 获取创建流程实例请求体数组
	var subReqs []*dto.CreateWorkflowDefDTO
	for idx, subworkflow := range subworkflows {
		for refName, workflowDef := range subworkflow {
			subworkflowDefJson, err := getSubworkflowJsonStr(subWorkflowSlice, idx, refName)
			if err != nil {
				return err
			}
			subReq := &dto.CreateWorkflowDefDTO{
				Namespace:   req.Namespace,
				Version:     defaultDefInitVersion,
				ParentDefID: req.ParentDefID,
				Attribute:   po.WorkflowDefAttr{RefName: refName, ParentDefVersion: req.ParentDefVersion},
				Name:        workflowDef.Name,
				Creator:     workflowDef.Creator,
				Description: workflowDef.Desc,
				DefJson:     subworkflowDefJson,
			}
			subReqs = append(subReqs, subReq)
		}
	}
	// 如果不存在子流程则直接返回
	if len(subReqs) == 0 {
		return nil
	}
	if err := m.CreateWorkflowDefs(ctx, subReqs); err != nil {
		return err
	}

	log.Infof("Create subworkflows of workflow [%d] success", req.ParentDefID)
	return nil
}

// UpdateWorkflowDef 更新工作流定义
func (m *WorkflowDefCommandService) UpdateWorkflowDef(ctx context.Context, req *dto.CreateWorkflowDefDTO) error {
	// 检查传入content格式
	if err := validator.ValidateDefJson(req.DefJson); err != nil {
		return err
	}

	lock, err := execution.GetDefDistributeLock(m.cacheRepo, req.DefID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 更新流程定义信息
	if err := m.updateWorkflowDefInfo(req); err != nil {
		return err
	}
	// 更新子流程定义信息
	if err := m.CreateSubworkflowDef(ctx, &dto.CreateSubworkflowDefDTO{
		ParentDefID:      req.DefID,
		ParentDefVersion: req.Version,
		DefJson:          req.DefJson}); err != nil {
		log.Errorf("[%d]Failed to create subworkflow def for update parent def, caused by %s", req.DefID, err)
		return err
	}
	// 发送更新事件
	return m.workflowUpdater.SendWorkflowDefExternalEvent(req.Namespace, req.DefID, event.DefUpdate)
}

// updateWorkflowDefInfo 更新流程定义数据层
func (m *WorkflowDefCommandService) updateWorkflowDefInfo(req *dto.CreateWorkflowDefDTO) error {
	workflowDef, err := m.getWorkflowDef(req.DefID)
	if err != nil {
		return err
	}
	if workflowDef.Status == entity.Enabled {
		return fmt.Errorf("def status is enabled")
	}
	req.Version = workflowDef.Version + versionUpgradeStepSize
	req.DefID = workflowDef.DefID
	if err := m.initCreateWorkflowDefDTO(req); err != nil {
		return err
	}
	if _, err := m.workflowDefRepo.Create(req); err != nil {
		return err
	}
	return nil
}

// EnableWorkflowDef 激活工作流定义
func (m *WorkflowDefCommandService) EnableWorkflowDef(ctx context.Context, req *dto.EnableWorkflowDefDTO) error {
	lock, err := execution.GetDefDistributeLock(m.cacheRepo, req.DefID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 更新流程定义信息
	workflowDef, err := m.updateWorkflowDefStatus(req.DefID, entity.Enabled)

	// 如果已经处于激活的状态，直接返回
	if errors.Is(err, errno.RepeatOperation) {
		return nil
	}

	if err != nil {
		return err
	}
	// 注册触发事件
	if err := m.registerTrigger(workflowDef); err != nil {
		return fmt.Errorf("failed to register trigger %s", err)
	}
	// 激活子流程
	if err := m.EnableSubworkflowDef(ctx, workflowDef); err != nil {
		return fmt.Errorf("enable subworkflow %d failed, err: %w", workflowDef.DefID, err)
	}
	// 发送激活流程定义事件
	return m.workflowUpdater.SendWorkflowDefExternalEvent(req.Namespace, req.DefID, event.DefEnable)
}

// EnableSubworkflowDef 激活子工作流定义
func (m *WorkflowDefCommandService) EnableSubworkflowDef(ctx context.Context, def *entity.WorkflowDef) error {
	for _, subWorkflow := range def.Subworkflows {
		for _, workflowDef := range subWorkflow {
			if err := m.EnableWorkflowDef(ctx, &dto.EnableWorkflowDefDTO{DefID: workflowDef.DefID}); err != nil {
				return err
			}
		}
	}
	return nil
}

// DisableWorkflowDef 去激活工作流定义
func (m *WorkflowDefCommandService) DisableWorkflowDef(ctx context.Context, req *dto.DisableWorkflowDefDTO) error {
	lock, err := execution.GetDefDistributeLock(m.cacheRepo, req.DefID)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 更新流程定义信息
	workflowDef, err := m.updateWorkflowDefStatus(req.DefID, entity.Disabled)

	// 如果已经处于未激活的状态，直接返回
	if errors.Is(err, errno.RepeatOperation) {
		return nil
	}

	if err != nil {
		return err
	}
	// 反注册定时触发事件
	if err := m.unregisterTrigger(workflowDef); err != nil {
		return fmt.Errorf("unregister trigger failed, err: %w", err)
	}
	// 去激活子流程
	if err := m.DisableSubWorkflowDef(ctx, workflowDef); err != nil {
		return fmt.Errorf("disable subworkflow %d failed, err: %w", workflowDef.DefID, err)
	}
	// 发送去激活流程定义事件
	return m.workflowUpdater.SendWorkflowDefExternalEvent(req.Namespace, req.DefID, event.DefDisable)
}

// DisableSubWorkflowDef 去激活子工作流定义
func (m *WorkflowDefCommandService) DisableSubWorkflowDef(ctx context.Context, def *entity.WorkflowDef) error {
	for _, subWorkflow := range def.Subworkflows {
		for _, workflowDef := range subWorkflow {
			if err := m.DisableWorkflowDef(ctx, &dto.DisableWorkflowDefDTO{DefID: workflowDef.DefID}); err != nil {
				return err
			}
		}
	}
	return nil
}

// getWorkflowDef 获取流程定义DTO数据
func (m *WorkflowDefCommandService) getWorkflowDef(defID string) (*entity.WorkflowDef, error) {
	getDefInfo := &dto.GetWorkflowDefDTO{DefID: defID}
	return m.workflowDefRepo.GetLastVersion(getDefInfo)
}

// updateWorkflowDefStatus 更新流程定义状态
func (m *WorkflowDefCommandService) updateWorkflowDefStatus(defID string,
	newStatus entity.DefStatus) (*entity.WorkflowDef, error) {
	workflowDef, err := m.getWorkflowDef(defID)
	if err != nil {
		return nil, err
	}
	// 如果修改前状态就是本身就直接返回
	if workflowDef.Status == newStatus {
		return nil, errno.RepeatOperation
	}
	updateDefDTO := &dto.UpdateWorkflowDefDTO{
		DefID:   workflowDef.DefID,
		Version: workflowDef.Version,
		Status:  newStatus,
	}
	if err := m.workflowDefRepo.UpdateStatus(updateDefDTO); err != nil {
		return nil, err
	}

	workflowDef.Status = newStatus
	return workflowDef, nil
}

// doRegisterTrigger DTO格式注册触发器
func (m *WorkflowDefCommandService) doRegisterTrigger(defID string, defVersion int, defJson string) error {
	def := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(defJson), def); err != nil {
		return err
	}
	def.DefID = defID
	def.Version = defVersion
	return m.registerTrigger(def)
}

// 注册触发器 这里不将错误吐出 跟流程定义设定有关
func (m *WorkflowDefCommandService) registerTrigger(def *entity.WorkflowDef) error {
	// 未激活不注册
	if def.Status == entity.Disabled {
		return nil
	}
	registerDTO := &dto.RegisterTriggersDTO{
		DefID:      def.DefID,
		DefVersion: def.Version,
		Operator:   def.Creator,
		Level:      entity.DefTrigger,
		Triggers:   def.Triggers,
	}
	return m.triggerRegistry.Register(registerDTO)
}

// unregisterTrigger 反注册触发器
func (m *WorkflowDefCommandService) unregisterTrigger(def *entity.WorkflowDef) error {
	unRegisterDTO := &dto.UnRegisterTriggersDTO{
		DefID:      def.DefID,
		DefVersion: def.Version,
		Level:      entity.DefTrigger,
	}
	return m.triggerRegistry.UnRegister(unRegisterDTO)
}

// UploadWorkflowDef 上传工作流定义
func (m *WorkflowDefCommandService) UploadWorkflowDef(ctx context.Context,
	d *dto.UploadWorkflowDefDTO) (string, error) {
	createDefDTO := &dto.CreateWorkflowDefDTO{}
	if err := copier.Copy(createDefDTO, d); err != nil {
		return "", err
	}
	// 当未带defID时即为新建
	if d.DefID == "" {
		return m.CreateWorkflowDef(ctx, createDefDTO)
	}
	return d.DefID, m.UpdateWorkflowDef(ctx, createDefDTO)
}

func getSubworkflowsEntityAndSlice(defJson string) ([]map[string]entity.WorkflowDef, []interface{}, error) {
	// 解析拿到当前流程中的所有子流程
	workflowDefEntity := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(defJson), workflowDefEntity); err != nil {
		return nil, nil, err
	}
	subWorkflows := workflowDefEntity.Subworkflows

	// 获得子流程数组
	subWorkflowSlice, err := getSubworkflowSlice([]byte(defJson))
	if err != nil {
		return nil, nil, err
	}

	return subWorkflows, subWorkflowSlice, nil
}

// getSubworkflowSlice 获取子流程配置对应的 slice
func getSubworkflowSlice(jsonBytes []byte) ([]interface{}, error) {
	valueJson, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		return nil, err
	}
	subWorkflows := valueJson.GetPath("subworkflows").MustArray()

	return subWorkflows, nil
}

// getSubworkflowJsonStr 获取子流程的 json 字符串
// subWorkflows: 子流程数组; idx: 子流程下标; refName: 子流程对应的 refName
func getSubworkflowJsonStr(subWorkflows []interface{}, idx int, refName string) (string, error) {
	if idx >= len(subWorkflows) {
		return "", fmt.Errorf("get subworkflow out of range with index %d", idx)
	}
	// 根据下标获取到当前子流程
	subDefJson, ok := subWorkflows[idx].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("subworkflow in the define is not a valid map: %v", subWorkflows[idx])
	}
	// 根据 refName 获取到子流程的配置
	defJsonInterface, ok := subDefJson[refName]
	if !ok {
		return "", fmt.Errorf("get subworkflow by refName %s failed", refName)
	}
	// 获取子流程的 DefJson
	defJsonStr, err := json.Marshal(defJsonInterface)
	if err != nil {
		return "", err
	}
	return string(defJsonStr), nil
}
