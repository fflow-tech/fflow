package config

import (
	"context"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/provider"
	"golang.org/x/oauth2"
)

var (
	oauthGroupKey = config.NewGroupKey("auth", "OAUTH2") // WebCors 跨域配置
)

// OAuth2Config 应用配置
type OAuth2Config struct {
	*oauth2.Config
}

// GetOauth2Config 获取默认配置
func GetOauth2Config() *OAuth2Config {
	conf := OAuth2Config{Config: &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: "https://www.fflow.link/auth/api/v1/oauth2/callback",
	}}

	provider.GetConfigProvider().GetAny(context.Background(), corsGroupKey, &conf)
	return &conf
}
