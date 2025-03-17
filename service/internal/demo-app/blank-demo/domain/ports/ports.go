package ports

import (
	"context"

	"github.com/fflow-tech/fflow/service/internal/demo-app/blank-demo/domain/dto"
)

// CommandPorts 写入接口
type CommandPorts interface {
	CrawlerCommandPorts
}

// QueryPorts 读入接口
type QueryPorts interface {
	CrawlerQueryPorts
}

// CrawlerCommandPorts 函数相关接口
type CrawlerCommandPorts interface {
	StartCollect(ctx context.Context, req *dto.StartCollectReqDTO) error
}

// CrawlerQueryPorts 鉴权相关接口
type CrawlerQueryPorts interface {
}
