package command

import (
	"context"
	"github.com/fflow-tech/fflow/service/internal/demo-app/blank-demo/domain/dto"
)

// CrawlerCommandService 写服务
type CrawlerCommandService struct {
}

// NewCrawlerCommandService 新建服务
func NewCrawlerCommandService() (*CrawlerCommandService, error) {
	return &CrawlerCommandService{}, nil
}

// StartCollect 开始采集
func (m *CrawlerCommandService) StartCollect(ctx context.Context, req *dto.StartCollectReqDTO) error {
	return nil
}
