package utils

import (
	"context"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/log"
)

// ShutdownGraceful 优雅关闭
// 这里主要是用来关闭自身的资源
func ShutdownGraceful(fs ...func(chan struct{}) error) {
	log.Infof("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f func(chan struct{}) error) {
			defer wg.Done()

			c := make(chan struct{}, 1)
			go f(c)
			select {
			case <-c:
			case <-ctx.Done():
			}
		}(f)
	}

	wg.Wait()
	log.Infof("Server exit")
}
