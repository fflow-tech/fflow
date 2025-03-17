// Package convertor 提供对象与对象的转换功能
package convertor

import (
	"fmt"
	"strconv"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
)

var (
	EventConvertor = &eventConvertorImpl{} // 转换器
)

type eventConvertorImpl struct {
}

// ConvertEntityToNodeExternalEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeExternalEvent(nodeInst *entity.NodeInst,
	eventType event.ExternalEventType) (interface{}, error) {
	switch eventType {
	case event.NodeSuccess:
		return c.ConvertEntityToNodeSuccessEvent(nodeInst)
	case event.NodeFail:
		return c.ConvertEntityToNodeFailEvent(nodeInst)
	case event.NodeStart:
		return c.ConvertEntityToNodeStartEvent(nodeInst)
	case event.NodeCancel:
		return c.ConvertEntityToNodeCancelEvent(nodeInst)
	case event.NodeTimeout:
		return c.ConvertEntityToNodeTimeoutEvent(nodeInst)
	case event.NodeNearTimeout:
		return c.ConvertEntityToNodeNearTimeoutEvent(nodeInst)
	case event.NodeWait:
		return c.ConvertEntityToNodeWaitEvent(nodeInst)
	case event.NodeAsynWait:
		return c.ConvertEntityToNodeAsynWaitEvent(nodeInst)
	default:
		return nil, fmt.Errorf("unsupport eventType=%s", eventType)
	}
}

// ConvertEntityToWorkflowExternalEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowExternalEvent(inst *entity.WorkflowInst, nodeRefName string,
	eventType event.ExternalEventType) (interface{}, error) {
	switch eventType {
	case event.NodeSkip:
		return c.ConvertEntityToNodeSkipEvent(inst, nodeRefName)
	case event.NodeCancelSkip:
		return c.ConvertEntityToNodeCancelSkipEvent(inst, nodeRefName)
	case event.WorkflowPause:
		return c.ConvertEntityToWorkflowPauseEvent(inst)
	case event.WorkflowResume:
		return c.ConvertEntityToWorkflowResumeEvent(inst)
	case event.WorkflowSuccess:
		return c.ConvertEntityToWorkflowSuccessEvent(inst)
	case event.WorkflowStart:
		// 流程重启的时候才通过这里发送, 因为实例已经生成了, 所以走这个路径
		return c.ConvertEntityToWorkflowRestartEvent(inst)
	case event.WorkflowFail:
		return c.ConvertEntityToWorkflowFailEvent(inst)
	case event.WorkflowCancel:
		return c.ConvertEntityToWorkflowCancelEvent(inst)
	case event.WorkflowNearTimeout:
		return c.ConvertEntityToWorkflowNearTimeoutEvent(inst)
	case event.WorkflowTimeout:
		return c.ConvertEntityToWorkflowTimeoutEvent(inst)
	default:
		return nil, fmt.Errorf("Unsupport eventType=%s", eventType)
	}
}

// ConvertEntityToNodeDriveEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeDriveEvent(nodeInst *entity.NodeInst,
	eventType event.DriveEventType) (interface{}, error) {
	switch eventType {
	case event.NodeExecuteDrive:
		return c.ConvertEntityToNodeExecuteDriveEvent(nodeInst)
	case event.NodePollDrive:
		return c.ConvertEntityToNodePollDriveEvent(nodeInst)
	case event.NodeCompleteDrive:
		return c.ConvertEntityToNodeCompleteDriveEvent(nodeInst, false)
	default:
		return nil, fmt.Errorf("Unsupport eventType=%s", eventType)
	}
}

// ConvertEntityToNodeStartEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeStartEvent(nodeInst *entity.NodeInst) (event.NodeStartEvent, error) {
	externalEvent := event.NodeStartEvent{
		BasicEvent: c.newNodeInstExternalBasicEventWithReason(event.NodeStart, nodeInst,
			nodeInst.Operator.CancelledOperator, nodeInst.Reason.CancelledReason),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeSuccessEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeSuccessEvent(nodeInst *entity.NodeInst) (
	event.NodeSucceedEvent, error) {
	externalEvent := event.NodeSucceedEvent{
		BasicEvent: c.newNodeInstExternalBasicEvent(event.NodeSuccess, nodeInst),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeFailEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeFailEvent(nodeInst *entity.NodeInst) (event.NodeFailEvent, error) {
	externalEvent := event.NodeFailEvent{
		BasicEvent: c.newNodeInstExternalBasicEventWithReason(event.NodeFail, nodeInst,
			nodeInst.Operator.FailedOperator, nodeInst.Reason.FailedReason),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeCancelEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeCancelEvent(nodeInst *entity.NodeInst) (event.NodeCancelEvent, error) {
	externalEvent := event.NodeCancelEvent{
		BasicEvent: c.newNodeInstExternalBasicEventWithReason(event.NodeCancel, nodeInst,
			nodeInst.Operator.CancelledOperator, nodeInst.Reason.CancelledReason),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeNearTimeoutEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeNearTimeoutEvent(nodeInst *entity.NodeInst) (
	event.NodeNearTimeoutEvent, error) {
	externalEvent := event.NodeNearTimeoutEvent{
		BasicEvent: c.newNodeInstExternalBasicEvent(event.NodeNearTimeout, nodeInst),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeTimeoutEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeTimeoutEvent(nodeInst *entity.NodeInst) (
	event.NodeTimeoutEvent, error) {
	externalEvent := event.NodeTimeoutEvent{
		BasicEvent: c.newNodeInstExternalBasicEvent(event.NodeTimeout, nodeInst),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeSkipEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeSkipEvent(inst *entity.WorkflowInst, nodeRefName string) (
	event.NodeSkipEvent, error) {
	externalEvent := event.NodeSkipEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEvent(event.NodeSkip, inst),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
		Node:       nodeRefName,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeCancelSkipEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeCancelSkipEvent(inst *entity.WorkflowInst, nodeRefName string) (
	event.NodeCancelSkipEvent, error) {
	externalEvent := event.NodeCancelSkipEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEvent(event.NodeCancelSkip, inst),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
		Node:       nodeRefName,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeWaitEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeWaitEvent(nodeInst *entity.NodeInst) (event.NodeWaitEvent, error) {
	externalEvent := event.NodeWaitEvent{
		BasicEvent: c.newNodeInstExternalBasicEvent(event.NodeWait, nodeInst),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeAsynWaitEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeAsynWaitEvent(
	nodeInst *entity.NodeInst) (event.NodeAsynWaitEvent, error) {
	externalEvent := event.NodeAsynWaitEvent{
		BasicEvent: c.newNodeInstExternalBasicEvent(event.NodeAsynWait, nodeInst),
		DefID:      nodeInst.DefID,
		DefVersion: strconv.Itoa(nodeInst.DefVersion),
		InstID:     nodeInst.InstID,
		Node:       nodeInst.BasicNodeDef.RefName,
		NodeInstID: nodeInst.NodeInstID,
	}
	return externalEvent, nil
}

// ConvertEntityToNodeExecuteDriveEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeExecuteDriveEvent(inst *entity.NodeInst) (interface{}, error) {
	driveEvent := event.NodeExecuteDriveEvent{
		BasicEvent: event.NewDriveBasicEvent(event.NodeExecuteDrive, inst.Namespace, inst.DefID, inst.InstID),
		DefID:      inst.DefID,
		DefVersion: inst.DefVersion,
		InstID:     inst.InstID,
		NodeInstID: inst.NodeInstID,
	}

	return driveEvent, nil
}

// ConvertEntityToNodePollDriveEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodePollDriveEvent(nodeInst *entity.NodeInst) (interface{}, error) {
	driveEvent := event.NodePollDriveEvent{
		BasicEvent: event.NewDriveBasicEvent(event.NodePollDrive, nodeInst.Namespace, nodeInst.DefID, nodeInst.InstID),
		DefID:      nodeInst.DefID,
		DefVersion: nodeInst.DefVersion,
		InstID:     nodeInst.InstID,
		NodeInstID: nodeInst.NodeInstID,
	}

	return driveEvent, nil
}

// ConvertEntityToNodeScheduleDriveEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeScheduleDriveEvent(inst *entity.WorkflowInst,
	nodesToBeScheduled []string) (interface{}, error) {
	driveEvent := event.NodeScheduleDriveEvent{
		BasicEvent: event.NewDriveBasicEvent(event.NodeScheduleDrive,
			inst.WorkflowDef.Namespace, inst.WorkflowDef.DefID, inst.InstID),
		DefID:       inst.WorkflowDef.DefID,
		DefVersion:  inst.WorkflowDef.Version,
		InstID:      inst.InstID,
		NodeInstIDs: nodesToBeScheduled,
	}

	return driveEvent, nil
}

// ConvertEntityToNodeCompleteDriveEvent 转换
func (c *eventConvertorImpl) ConvertEntityToNodeCompleteDriveEvent(nodeInst *entity.NodeInst,
	fromResumeInst bool) (interface{}, error) {
	driveEvent := event.NodeCompleteDriveEvent{
		BasicEvent: event.NewDriveBasicEvent(event.NodeCompleteDrive, nodeInst.Namespace,
			nodeInst.DefID, nodeInst.InstID),
		DefID:          nodeInst.DefID,
		DefVersion:     nodeInst.DefVersion,
		InstID:         nodeInst.InstID,
		NodeInstID:     nodeInst.NodeInstID,
		FromResumeInst: fromResumeInst,
	}

	return driveEvent, nil
}

// ConvertEntityToWorkflowStartEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowStartEvent(inst *entity.WorkflowInst) (
	event.WorkflowStartEvent, error) {
	externalEvent := event.WorkflowStartEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowStart, inst, inst.Creator,
			inst.Reason.StartReason),
		DefID:    inst.WorkflowDef.DefID,
		InstID:   inst.InstID,
		InstName: inst.Name,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowPauseEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowPauseEvent(inst *entity.WorkflowInst) (
	event.WorkflowPauseEvent, error) {
	externalEvent := event.WorkflowPauseEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowPause, inst,
			inst.Operator.PauseOperator, inst.Reason.StartReason),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowResumeEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowResumeEvent(inst *entity.WorkflowInst) (
	event.WorkflowResumeEvent, error) {
	externalEvent := event.WorkflowResumeEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowPause, inst,
			inst.Operator.PauseOperator, inst.Reason.StartReason),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowSuccessEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowSuccessEvent(inst *entity.WorkflowInst) (
	event.WorkflowSuccessEvent, error) {
	externalEvent := event.WorkflowSuccessEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowSuccess, inst,
			inst.Operator.SucceedOperator, inst.Reason.SucceedReason),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowCancelEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowCancelEvent(inst *entity.WorkflowInst) (
	event.WorkflowCancelEvent, error) {
	externalEvent := event.WorkflowCancelEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowCancel, inst,
			inst.Operator.CancelledOperator, inst.Reason.CancelledReason),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowTimeoutEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowTimeoutEvent(inst *entity.WorkflowInst) (
	event.WorkflowTimeoutEvent, error) {
	externalEvent := event.WorkflowTimeoutEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowTimeout, inst,
			"", inst.Reason.StartReason),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowNearTimeoutEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowNearTimeoutEvent(inst *entity.WorkflowInst) (
	event.WorkflowNearTimeoutEvent, error) {
	externalEvent := event.WorkflowNearTimeoutEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowNearTimeout, inst,
			"", inst.Reason.StartReason),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowFailEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowFailEvent(inst *entity.WorkflowInst) (
	event.WorkflowFailEvent, error) {
	externalEvent := event.WorkflowFailEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowFail, inst,
			inst.Operator.FailedOperator, inst.Reason.FailedRootCause.FailedReason),
		DefID:      inst.WorkflowDef.DefID,
		DefVersion: strconv.Itoa(inst.WorkflowDef.Version),
		InstID:     inst.InstID,
	}

	return externalEvent, nil
}

// ConvertEntityToWorkflowRestartEvent 转换
func (c *eventConvertorImpl) ConvertEntityToWorkflowRestartEvent(inst *entity.WorkflowInst) (
	event.WorkflowStartEvent, error) {
	externalEvent := event.WorkflowStartEvent{
		BasicEvent: c.newWorkflowInstExternalBasicEventWithReason(event.WorkflowStart, inst,
			inst.Creator, inst.Reason.RestartReason),
		DefID:     inst.WorkflowDef.DefID,
		InstID:    inst.InstID,
		InstName:  inst.Name,
		StartNode: inst.LastRestartNode,
	}

	return externalEvent, nil
}

// newNodeInstExternalBasicEvent
func (c *eventConvertorImpl) newNodeInstExternalBasicEvent(eventType event.ExternalEventType,
	nodeInst *entity.NodeInst) event.BasicEvent {
	return event.NewExternalBasicEvent(eventType, nodeInst.Namespace, nodeInst.DefID, nodeInst.InstID, nodeInst.Biz)
}

// newWorkflowInstExternalBasicEvent
func (c *eventConvertorImpl) newWorkflowInstExternalBasicEvent(eventType event.ExternalEventType,
	inst *entity.WorkflowInst) event.BasicEvent {
	return event.NewExternalBasicEvent(eventType, inst.WorkflowDef.Namespace,
		inst.WorkflowDef.DefID, inst.InstID, inst.Biz)
}

// newNodeInstExternalBasicEventWithReason
func (c *eventConvertorImpl) newNodeInstExternalBasicEventWithReason(eventType event.ExternalEventType,
	nodeInst *entity.NodeInst, operator, reason string) event.BasicEvent {
	return event.NewExternalBasicEventWithReason(eventType, nodeInst.Namespace,
		nodeInst.DefID, nodeInst.InstID, operator, reason, nodeInst.Biz)
}

// newWorkflowInstExternalBasicEventWithReason
func (c *eventConvertorImpl) newWorkflowInstExternalBasicEventWithReason(eventType event.ExternalEventType,
	inst *entity.WorkflowInst, operator, reason string) event.BasicEvent {
	return event.NewExternalBasicEventWithReason(eventType, inst.WorkflowDef.Namespace,
		inst.WorkflowDef.DefID, inst.InstID, operator, reason, inst.Biz)
}
