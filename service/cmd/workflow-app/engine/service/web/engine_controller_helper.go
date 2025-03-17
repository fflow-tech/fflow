package web

import (
	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/login"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

func bindReq(c *gin.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	secretKey := config.GetAppConfig().SecretKey
	currentUser, err := login.GetUserInfoFromCookie(c, secretKey)
	if err != nil {
		return err
	}

	log.Infof("Current user: %s", utils.StructToJsonStr(currentUser))
	if c.Request.Method == "GET" {
		login.SetNamespace(req, currentUser)
	} else {
		login.SetNamespace(req, currentUser)
		login.SetOperator(req, currentUser)
	}
	return nil
}
