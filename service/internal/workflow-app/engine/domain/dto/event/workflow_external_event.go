package event

// DefCreateEvent 流程定义创建事件
type DefCreateEvent struct {
	BasicEvent
	DefID string `json:"def_id"` // 流程定义ID
}

// DefUpdateEvent 流程定义更新事件
type DefUpdateEvent struct {
	BasicEvent
	DefID string `json:"def_id"` // 流程定义ID
}

// DefEnableEvent 流程定义激活事件
type DefEnableEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
}

// DefDisableEvent 流程定义取消激活事件
type DefDisableEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
}

// WorkflowStartEvent 流程实例启动事件
type WorkflowStartEvent struct {
	BasicEvent
	DefID     string `json:"def_id"`     // 流程定义ID
	InstID    string `json:"inst_id"`    // 流程实例ID
	InstName  string `json:"instName"`   // 流程实例名称
	StartNode string `json:"start_node"` // 开始的节点
}

// WorkflowSuccessEvent 流程实例成功事件
type WorkflowSuccessEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
}

// WorkflowCancelEvent 流程实例被取消事件
type WorkflowCancelEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
}

// WorkflowPauseEvent 流程实例被暂停事件
type WorkflowPauseEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
}

// WorkflowResumeEvent 流程实例被恢复事件
type WorkflowResumeEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
}

// WorkflowTimeoutEvent 流程实例超时事件
type WorkflowTimeoutEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
}

// WorkflowNearTimeoutEvent 流程实例接近超时事件
type WorkflowNearTimeoutEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
}

// WorkflowFailEvent 流程实例失败事件
type WorkflowFailEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
}

// NodeStartEvent 节点启动事件
type NodeStartEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}

// NodeSucceedEvent 节点成功事件
type NodeSucceedEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}

// NodeCancelEvent 节点被取消事件
type NodeCancelEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}

// NodeSkipEvent 节点被跳过事件
type NodeSkipEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
	Node       string `json:"node"`        // 节点引用名称
}

// NodeCancelSkipEvent 节点被取消跳过事件
type NodeCancelSkipEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`      // 流程定义ID
	DefVersion string `json:"def_version"` // 流程定义版本号
	InstID     string `json:"inst_id"`     // 流程实例ID
	Node       string `json:"node"`        // 节点引用名称
}

// NodeFailEvent 节点执行失败
type NodeFailEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}

// NodeWaitEvent 节点等待事件
type NodeWaitEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}

// NodeAsynWaitEvent 节点异步等待事件
type NodeAsynWaitEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}

// NodeTimeoutEvent 节点超时事件
type NodeTimeoutEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}

// NodeNearTimeoutEvent 节点接近超时事件
type NodeNearTimeoutEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`       // 流程定义ID
	DefVersion string `json:"def_version"`  // 流程定义版本号
	InstID     string `json:"inst_id"`      // 流程实例ID
	Node       string `json:"node"`         // 节点引用名称
	NodeInstID string `json:"node_inst_id"` // 节点实例ID
}
