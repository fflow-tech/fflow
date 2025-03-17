package ginfilter

import (
	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/login"
	"net/http"
)

const cookieMaxAge = 24 * 60 * 60

// Auth 权限校验
// 这里只是最基本的权限校验，后面需要统一通过类似 casbin 框架来实现
func Auth(authConfig *config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(login.SessionCookieName)
		if err != nil {
			log.Errorf("Failed to get token str, tokenStr=%s, caused by %s", tokenStr, err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := login.GetUserInfoFromCookie(c, authConfig.SecretKey)
		if err != nil {
			log.Errorf("Failed to get user info, caused by %s", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if len(user.Username) > 0 && len(user.Email) > 0 && len(user.Avatar) > 0 {
			c.SetCookie(login.SessionCookieName, tokenStr,
				cookieMaxAge, "/", authConfig.Domain, false, true)
			c.Next()
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
