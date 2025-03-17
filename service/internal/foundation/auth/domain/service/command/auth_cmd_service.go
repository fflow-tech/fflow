package command

import (
	"context"
	"fmt"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math/rand"

	"github.com/google/go-github/github"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/config"
	ac "github.com/fflow-tech/fflow/service/internal/foundation/auth/pkg/constants"
	"github.com/fflow-tech/fflow/service/internal/foundation/auth/repository/repo"
	"github.com/fflow-tech/fflow/service/pkg/email"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/redis"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

// AuthCommandService 写服务
type AuthCommandService struct {
	redisClient        *redis.Client
	userRepo           ports.UserRepository
	namespaceRepo      ports.NamespaceRepository
	namespaceTokenRepo ports.NamespaceTokenRepository
}

// NewAuthCommandService 新建服务
func NewAuthCommandService(redisClient *redis.Client,
	userRepo *repo.UserRepo,
	namespaceRepo *repo.NamespaceRepo,
	namespaceTokenRepo *repo.NamespaceTokenRepo) (
	*AuthCommandService, error) {
	return &AuthCommandService{
		redisClient:        redisClient,
		userRepo:           userRepo,
		namespaceRepo:      namespaceRepo,
		namespaceTokenRepo: namespaceTokenRepo,
	}, nil
}

// Login 登录函数
func (m *AuthCommandService) Login(ctx context.Context, req *dto.LoginReqDTO) (*dto.LoginRspDTO, error) {
	if req.Username == config.GetAppConfig().AdminUsername && req.Password == config.GetAppConfig().AdminPassword {
		return &dto.LoginRspDTO{
			Status:           "ok",
			Type:             req.Type,
			CurrentAuthority: "admin",
		}, nil
	}

	return dto.NewFailedLoginRspDTO(), nil
}

// OutLogin 登出
func (m *AuthCommandService) OutLogin(ctx context.Context, req *dto.OutLoginReqDTO) (*dto.OutLoginRspDTO, error) {
	return &dto.OutLoginRspDTO{Success: true}, nil
}

// Oauth2Callback Github登录回调
func (m *AuthCommandService) Oauth2Callback(ctx context.Context, req *dto.Oauth2CallbackReqDTO) (
	*dto.Oauth2CallbackRspDTO, error) {
	log.Infof("The oauth2 callback req is %s", utils.StructToJsonStr(req))

	oauth2Config := config.GetOauth2Config()
	token, err := oauth2Config.Exchange(ctx, req.Code)
	if err != nil {
		log.Errorf("Failed to exchange token: %v\n", err)
		return nil, err
	}

	// 使用Token创建一个Github实例
	client := github.NewClient(oauth2Config.Client(ctx, token))

	// 获取用户信息
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("Failed to get user: %v\n", err)
		return nil, err
	}

	userEntity, _ := m.userRepo.CreateIfNotExists(&dto.CreateUserDTO{
		Username: *user.Email,
		NickName: *user.Name,
		Email:    *user.Email,
		Avatar:   *user.AvatarURL,
		AuthType: entity.Github.String(),
		Status:   entity.Enabled.IntValue(),
	})
	return convertor.AuthConvertor.ConvertEntityToOauth2CallbackRspDTO(userEntity)
}

// GetCaptcha 发送验证码
func (m *AuthCommandService) GetCaptcha(ctx context.Context, req *dto.GetCaptchaReqDTO) (
	*dto.GetCaptchaRspDTO, error) {
	if req == nil || !email.IsEmailValid(req.EmailReceiver) {
		return nil, fmt.Errorf("invalid email address")
	}

	code := generateCaptcha(4)

	emailConfig := config.GetEmailConfig()
	e := email.Email{
		To:      []string{req.EmailReceiver},
		Subject: "[FFlow] - 验证码邮件",
		Body:    fmt.Sprintf("验证码：%s", code),
	}
	if err := e.Send(email.SMTPServer{
		Host:     emailConfig.Host,
		Port:     emailConfig.Port,
		From:     emailConfig.From,
		Password: emailConfig.Password,
	}); err != nil {
		return nil, err
	}

	if err := m.redisClient.Set(ctx, getVerificationCodeCacheKey(req.EmailReceiver), code, 5*60); err != nil {
		return nil, err
	}

	return &dto.GetCaptchaRspDTO{Success: true}, nil
}

// VerifyCaptcha 验证验证码
func (m *AuthCommandService) VerifyCaptcha(ctx context.Context, req *dto.VerifyCaptchaReqDTO) (
	*dto.VerifyCaptchaRspDTO, error) {
	captcha, err := m.redisClient.Get(ctx, getVerificationCodeCacheKey(req.EmailReceiver))
	if err != nil {
		log.Errorf("Failed to get captcha from redis, caused by %v", err)
		return &dto.VerifyCaptchaRspDTO{IsValidCaptcha: false}, nil
	}
	if captcha == req.Captcha {
		m.redisClient.Del(ctx, getVerificationCodeCacheKey(req.EmailReceiver))
		user, _ := m.userRepo.CreateIfNotExists(&dto.CreateUserDTO{
			Username: req.EmailReceiver,
			NickName: req.EmailReceiver,
			Email:    req.EmailReceiver,
			Avatar:   ac.DefaultAvatar,
			AuthType: entity.Email.String(),
			Status:   entity.Enabled.IntValue(),
		})

		captchaDTO, err := convertor.AuthConvertor.ConvertEntityToVerifyCaptchaDTO(user)
		if err != nil {
			return nil, err
		}
		captchaDTO.IsValidCaptcha = true
		return captchaDTO, nil
	}
	return &dto.VerifyCaptchaRspDTO{IsValidCaptcha: false}, nil
}

// ValidateToken 校验 Token
func (m *AuthCommandService) ValidateToken(ctx context.Context, req *dto.ValidateTokenReqDTO) error {
	_, err := m.namespaceTokenRepo.Get(&dto.GetNamespaceTokenDTO{Namespace: req.Namespace, Token: req.AccessToken})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invalid token")
		}

		return err
	}
	return nil
}

func getVerificationCodeCacheKey(email string) string {
	return fmt.Sprintf("auth:verification_code:%s", email)
}

// generateCaptcha 生成验证码
func generateCaptcha(n int) string {
	var captcha string
	for i := 0; i < n; i++ {
		captcha += string(rune(rand.Intn(10) + '0'))
	}
	return captcha
}
