package command

// Adapters 适配器
type Adapters struct {
	*TimerDefCommandService
	*TimerTaskCommandService
	*PollingTaskCommandService
	*NotifyCommandService
	*AppCommandService
}

// NewCommandAdapters 新建命名适配器
func NewCommandAdapters(defService *TimerDefCommandService,
	timerTaskService *TimerTaskCommandService,
	pollingTaskService *PollingTaskCommandService,
	notifyService *NotifyCommandService,
	appService *AppCommandService) *Adapters {
	return &Adapters{TimerDefCommandService: defService,
		TimerTaskCommandService:   timerTaskService,
		PollingTaskCommandService: pollingTaskService,
		NotifyCommandService:      notifyService,
		AppCommandService:         appService,
	}
}
