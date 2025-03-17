package query

// Adapters 查询适配器
type Adapters struct {
	*TimerDefQueryService
	*TimerTaskQueryService
	*AppQueryService
}

// NewQueryAdapters 新建查询服务
func NewQueryAdapters(defService *TimerDefQueryService,
	timerTaskService *TimerTaskQueryService,
	appQueryService *AppQueryService) *Adapters {
	return &Adapters{
		TimerDefQueryService:  defService,
		TimerTaskQueryService: timerTaskService,
		AppQueryService:       appQueryService,
	}
}
