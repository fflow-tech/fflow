package event

// WorkflowStartDriveEvent 流程启动事件
type WorkflowStartDriveEvent struct {
	BasicEvent
	DefID          string `json:"def_id"`
	DefVersion     int    `json:"def_version"`
	InstID         string `json:"inst_id"`
	FromResumeInst bool   `json:"from_resume_inst"`
}

// NodeScheduleDriveEvent 节点被调度事件
type NodeScheduleDriveEvent struct {
	BasicEvent
	DefID       string   `json:"def_id"`
	DefVersion  int      `json:"def_version"`
	InstID      string   `json:"inst_id"`
	NodeInstIDs []string `json:"node_inst_id"`
}

// NodeExecuteDriveEvent 节点执行调度事件
type NodeExecuteDriveEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`
	DefVersion int    `json:"def_version"`
	InstID     string `json:"inst_id"`
	NodeInstID string `json:"node_inst_id"`
}

// NodePollDriveEvent 节点轮询调度事件
type NodePollDriveEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`
	DefVersion int    `json:"def_version"`
	InstID     string `json:"inst_id"`
	NodeInstID string `json:"node_inst_id"`
}

// NodeCompleteDriveEvent 节点完成调度事件
type NodeCompleteDriveEvent struct {
	BasicEvent
	DefID          string `json:"def_id"`
	DefVersion     int    `json:"def_version"`
	InstID         string `json:"inst_id"`
	NodeInstID     string `json:"node_inst_id"`
	FromResumeInst bool   `json:"from_resume_inst"` // 来自恢复实例
}

// NodeRetryDriveEvent 节点重试事件
type NodeRetryDriveEvent struct {
	BasicEvent
	DefID      string `json:"def_id"`
	DefVersion int    `json:"def_version"`
	InstID     string `json:"inst_id"`
	NodeInstID string `json:"node_inst_id"`
}
