// Package convertor 仓储层结构体到领域实体的转换
package convertor

import (
	"encoding/json"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/jinzhu/copier"
)

var (
	DefConvertor      = &defConvertorImpl{}
	InstConvertor     = &instConvertorImpl{}
	NodeInstConvertor = &nodeInstConvertorImpl{}
	TriggerConvertor  = &triggerConvertorImpl{}
)

type defConvertorImpl struct {
}

type instConvertorImpl struct {
}

type nodeInstConvertorImpl struct {
}

type triggerConvertorImpl struct {
}

// ConvertPOToEntity 转换成实体
func (*defConvertorImpl) ConvertPOToEntity(p *po.WorkflowDefPO) (*entity.WorkflowDef, error) {
	def := &entity.WorkflowDef{}
	if err := json.Unmarshal([]byte(p.DefJson), &def); err != nil {
		return nil, err
	}
	if err := copier.Copy(def, p); err != nil {
		return nil, err
	}
	def.ID = utils.UintToStr(p.ID)
	def.DefID = utils.Uint64ToStr(p.DefID)
	def.ParentDefID = utils.Uint64ToStr(p.ParentDefID)

	def.Status = entity.GetDefStatusByIntValue(p.Status)
	def.RefName = p.Attribute.RefName
	def.ParentDefVersion = p.Attribute.ParentDefVersion
	def.Desc = p.Description
	return def, nil
}

// ConvertPOToEntity 转换
func (*instConvertorImpl) ConvertPOToEntity(p *po.WorkflowInstPO,
	schedNodeInsts []*entity.NodeInst, curNodeInstID string) (*entity.WorkflowInst, error) {
	e := &entity.WorkflowInst{}
	if err := json.Unmarshal([]byte(p.Context), e); err != nil {
		return nil, err
	}
	err := copier.Copy(e, p)
	e.Status = entity.GetInstStatusByIntValue(p.Status)
	e.PreStatus = entity.GetInstStatusByIntValue(p.Status)
	e.NodeInstsCount = len(schedNodeInsts)
	e.InstID = utils.UintToStr(p.ID)
	e.SchedNodeInsts = filterSchedNodeInsts(e, schedNodeInsts)
	e.CurNodeInst = getNodeInstByID(schedNodeInsts, curNodeInstID)
	if e.Operator == nil {
		e.Operator = &entity.InstOperator{}
	}
	if e.Reason == nil {
		e.Reason = &entity.InstReason{}
	}

	return e, err
}

func filterSchedNodeInsts(inst *entity.WorkflowInst, nodeInsts []*entity.NodeInst) []*entity.NodeInst {
	// 0. 没有重启过, 只需要把所有最新的节点拿出来
	if inst.BeforeLastRestartMaxNodeInstID == "" {
		return entity.GetLastNodeInsts(nodeInsts)
	}

	lastRestartNodeInst := entity.GetNodeInstByRefName(entity.GetOldestNodeInsts(nodeInsts), inst.LastRestartNode)

	// 1. 重启的节点没执行过, 把所有最新的节点拿出来
	if lastRestartNodeInst == nil {
		return entity.GetLastNodeInsts(nodeInsts)
	}

	// 2. 重启的节点在当前节点前, 要筛选出重启节点前的节点实例 + 重启后的节点实例
	return filterSchedNodeInstsAfterLastRestart(inst, nodeInsts, lastRestartNodeInst)
}

func filterSchedNodeInstsAfterLastRestart(inst *entity.WorkflowInst, nodeInsts []*entity.NodeInst,
	lastRestartNodeInst *entity.NodeInst) []*entity.NodeInst {
	// 以第一次走的路径为准, 后面重启后的路径可能和之前的不同, 这里不考虑那种场景
	scheduledBeforeRestartNode := getScheduledBeforeRestartNode(nodeInsts, lastRestartNodeInst)

	r := []*entity.NodeInst{}
	lastNodeInsts := entity.GetLastNodeInsts(nodeInsts)
	for _, nodeInst := range lastNodeInsts {
		// 0. 如果在重启的节点前被调度的, 加入到被调度的队列里面
		if _, exists := scheduledBeforeRestartNode[nodeInst.BasicNodeDef.RefName]; exists {
			r = append(r, nodeInst)
			continue
		}
		// 1. 如果在重启的时间点后被调度的, 加入到被调度的队列里面
		if nodeInst.NodeInstID > inst.BeforeLastRestartMaxNodeInstID {
			r = append(r, nodeInst)
		}
	}

	return r
}

func getScheduledBeforeRestartNode(nodeInsts []*entity.NodeInst,
	lastRestartNodeInst *entity.NodeInst) map[string]struct{} {
	oldestNodeInsts := entity.GetOldestNodeInsts(nodeInsts)
	scheduledBeforeRestartNode := map[string]struct{}{}
	for _, nodeInst := range oldestNodeInsts {
		if nodeInst.NodeInstID < lastRestartNodeInst.NodeInstID {
			scheduledBeforeRestartNode[nodeInst.BasicNodeDef.RefName] = struct{}{}
		}
	}
	return scheduledBeforeRestartNode
}

func getNodeInstByID(curNodeInsts []*entity.NodeInst, nodeInstID string) *entity.NodeInst {
	for _, nodeInst := range curNodeInsts {
		if nodeInst.NodeInstID == nodeInstID {
			return nodeInst
		}
	}
	return nil
}

// ConvertPOToEntity 转换
func (*nodeInstConvertorImpl) ConvertPOToEntity(p *po.NodeInstPO) (*entity.NodeInst, error) {
	e := &entity.NodeInst{}
	if err := json.Unmarshal([]byte(p.Context), e); err != nil {
		return nil, err
	}
	if err := copier.Copy(e, p); err != nil {
		return nil, err
	}
	e.InstID = utils.Uint64ToStr(p.InstID)
	e.DefID = utils.Uint64ToStr(p.DefID)

	e.PreStatus = entity.GetNodeInstStatus(p.Status)
	e.Status = entity.GetNodeInstStatus(p.Status)
	e.NodeInstID = utils.UintToStr(p.ID)
	return e, nil
}

// ConvertPOToEntity 转换
func (*triggerConvertorImpl) ConvertPOToEntity(p *po.TriggerPO) (*entity.Trigger, error) {
	t := &entity.Trigger{}
	if err := json.Unmarshal([]byte(p.Attribute), t); err != nil {
		return nil, err
	}
	t.TriggerID = utils.UintToStr(p.ID)
	t.Event = p.Event
	t.DefID = utils.Uint64ToStr(p.DefID)
	t.DefVersion = p.DefVersion
	t.InstID = utils.Uint64ToStr(p.InstID)
	t.Status = entity.TriggerStatus(p.Status)
	t.Level = entity.TriggerLevel(p.Level)
	return t, nil
}
