package ginfilter

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

// Cors 跨域中间件
func Cors(corsConfig *config.CorsConfig) gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	log.Infof("Cors config: %+v", cfg)

	cfg.AllowMethods = corsConfig.AllowMethods

	cfg.AllowHeaders = corsConfig.AllowHeaders
	cfg.AllowCredentials = true

	allowOrigins := corsConfig.AllowOrigins
	cfg.AllowOrigins = allowOrigins

	// 防止因为读取配置失败导致 panic, 兜底允许所有的跨域请求
	if len(allowOrigins) == 0 {
		log.Warnf("Allow origins is empty, allow all origins")
		cfg.AllowAllOrigins = true
	}

	return cors.New(cfg)
}
