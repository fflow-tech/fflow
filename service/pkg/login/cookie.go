package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"reflect"
	"time"
)

const SessionCookieName = "_fflow_session"

// GetUserInfoFromCookie 从 Cookie 中获取用户信息
func GetUserInfoFromCookie(c *gin.Context, secretKey string) (*CurrentUserData, error) {
	tokenStr, err := c.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return &CurrentUserData{
		Namespace: claims["namespace"].(string),
		Username:  claims["username"].(string),
		Email:     claims["email"].(string),
		Avatar:    claims["avatar"].(string),
	}, nil
}

// SetUserInfoToCookie 设置用户信息
func SetUserInfoToCookie(c *gin.Context, user *CurrentUserData, secretKey string) error {
	if user == nil {
		return fmt.Errorf("invalid callback user info")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["namespace"] = user.Namespace
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["avatar"] = user.Avatar
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Errorf("Failed to generate JWT token, caused by %s", err)
		return err
	}

	c.SetCookie(SessionCookieName, signedToken, 24*3600, "/", config.GetAppConfig().Domain, false, false)
	return nil
}

// SetOperator 设置操作人
func SetOperator(req interface{}, currentUser *CurrentUserData) {
	// 管理员默认看到所有内容
	if currentUser.Username == config.GetAppConfig().AdminUsername {
		return
	}

	v := reflect.ValueOf(req).Elem()
	if !v.IsValid() {
		return
	}

	fields := []string{"Creator", "Operator", "Updater"}
	for _, field := range fields {
		f := v.FieldByName(field)
		if f.IsValid() {
			f.Set(reflect.ValueOf(currentUser.Username))
		}
	}
}

// SetNamespace 设置命名空间
func SetNamespace(req interface{}, currentUser *CurrentUserData) {
	// 管理员默认看到所有内容
	if currentUser.Username == config.GetAppConfig().AdminUsername {
		return
	}

	v := reflect.ValueOf(req).Elem()
	if !v.IsValid() {
		return
	}

	namespace := v.FieldByName("Namespace")
	if namespace.IsValid() {
		namespace.Set(reflect.ValueOf(currentUser.Namespace))
	}
}
