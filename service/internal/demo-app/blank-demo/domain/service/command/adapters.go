package command

// Adapters 适配器
type Adapters struct {
	*CrawlerCommandService
}

// NewCommandAdapters 初始化适配器
func NewCommandAdapters(crawlerCommandService *CrawlerCommandService) *Adapters {
	return &Adapters{CrawlerCommandService: crawlerCommandService}
}
