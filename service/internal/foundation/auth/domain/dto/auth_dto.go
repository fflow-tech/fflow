package dto

import (
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/login"
	"time"
)

// LoginReqDTO 登录请求
type LoginReqDTO struct {
	Type          string `json:"type,omitempty"`
	AutoLogin     bool   `json:"autoLogin,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	EmailReceiver string `json:"emailReceiver,omitempty"`
	Captcha       string `json:"captcha,omitempty"`
}

// LoginRspDTO 登录结果
type LoginRspDTO struct {
	Status           string `json:"status,omitempty"`
	Type             string `json:"type,omitempty"`
	CurrentAuthority string `json:"currentAuthority,omitempty"`
	Message          string `json:"message,omitempty"`
}

// NewSucceedLoginRsp 登录成功返回
func NewSucceedLoginRsp(t string) *LoginRspDTO {
	return &LoginRspDTO{Status: "ok", Type: t, CurrentAuthority: "admin"}
}

// NewFailedLoginRsp 登录失败返回
func NewFailedLoginRsp(message string) *LoginRspDTO {
	return &LoginRspDTO{Status: "fail", Message: message}
}

// NewFailedLoginRspDTO 新建
func NewFailedLoginRspDTO() *LoginRspDTO {
	return &LoginRspDTO{
		Status:           "error",
		CurrentAuthority: "guest",
	}
}

// CurrentUserReqDTO 当前用户请求
type CurrentUserReqDTO struct {
}

// CurrentUserData 当前用户信息
type CurrentUserData = login.CurrentUserData

// CurrentUserRspDTO 当前用户结果
type CurrentUserRspDTO struct {
	Success bool             `json:"success"`
	Data    *CurrentUserData `json:"data"`
}

// OutLoginReqDTO 取消登录请求
type OutLoginReqDTO struct {
}

// OutLoginRspDTO 取消登录结果
type OutLoginRspDTO struct {
	Success bool `json:"success,omitempty"`
}

// Oauth2CallbackReqDTO 登录回调请求
type Oauth2CallbackReqDTO struct {
	Code string `json:"code,omitempty"`
}

// Oauth2CallbackRspDTO 回调结果
type Oauth2CallbackRspDTO struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	NickName string `json:"nickName,omitempty"`
	AuthType string `json:"authType,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// GetCaptchaReqDTO 发送请求
type GetCaptchaReqDTO struct {
	EmailReceiver string `json:"emailReceiver,omitempty"`
}

// GetCaptchaRspDTO 发送结果
type GetCaptchaRspDTO struct {
	Success bool `json:"success,omitempty"`
}

// VerifyCaptchaReqDTO 验证请求
type VerifyCaptchaReqDTO struct {
	EmailReceiver string `json:"emailReceiver,omitempty"`
	Captcha       string `json:"captcha,omitempty"`
}

// VerifyCaptchaRspDTO 验证请求返回
type VerifyCaptchaRspDTO struct {
	IsValidCaptcha bool   `json:"isValidCode,omitempty"`
	Username       string `json:"username,omitempty"`
	NickName       string `json:"nickName,omitempty"`
	AuthType       string `json:"authType,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
	Status         int    `json:"status,omitempty"`
}

// CreateUserDTO 创建用户请求
type CreateUserDTO struct {
	Username string `json:"username,omitempty"`
	NickName string `json:"nickName,omitempty"`
	AuthType string `json:"authType,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Status   int    `json:"status,omitempty"`
}

// UpdateUserDTO 更新用户请求
type UpdateUserDTO struct {
	ID     int `form:"id,omitempty" json:"id,omitempty"`
	Status int `form:"status,omitempty" json:"status,omitempty"`
}

// GetUserDTO 查询用户请求
type GetUserDTO struct {
	ID       int    `form:"id,omitempty" json:"id,omitempty"`
	Username string `form:"username,omitempty" json:"username,omitempty"`
	Email    string `form:"email,omitempty" json:"email,omitempty"`
	Phone    string `form:"phone,omitempty" json:"phone,omitempty"`
}

// UserDTO 查询用户请求
type UserDTO struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	NickName string `json:"nickName,omitempty"`
	AuthType string `json:"authType,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Status   int    `json:"status,omitempty"`
}

// DeleteUserDTO 删除用户请求
type DeleteUserDTO struct {
	ID int `form:"id,omitempty" json:"id,omitempty"`
}

// PageQueryUserDTO 分页查询用户请求
type PageQueryUserDTO struct {
	*constants.PageQuery
	*constants.Order
}

// CreateNamespaceDTO 创建Namespace请求
type CreateNamespaceDTO struct {
	Namespace   string `form:"namespace,omitempty" json:"namespace,omitempty"`
	Creator     string `form:"creator,omitempty" json:"creator,omitempty"`
	Description string `form:"description,omitempty" json:"description,omitempty"`
}

// UpdateNamespaceDTO 更新Namespace请求
type UpdateNamespaceDTO struct {
	ID string `form:"id,omitempty" json:"id,omitempty"`
}

// GetNamespaceDTO 查询Namespace请求
type GetNamespaceDTO struct {
	ID        string `form:"id,omitempty" json:"id,omitempty"`
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"`
}

// NamespaceDTO 查询Namespace请求
type NamespaceDTO struct {
	ID string `form:"id,omitempty" json:"id,omitempty"`
}

// DeleteNamespaceDTO 删除Namespace请求
type DeleteNamespaceDTO struct {
	ID string `form:"id,omitempty" json:"id,omitempty"`
}

// PageQueryNamespaceDTO 分页查询Namespace请求
type PageQueryNamespaceDTO struct {
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"`
	*constants.PageQuery
	*constants.Order
}

// CreateNamespaceTokenDTO 创建NamespaceToken请求
type CreateNamespaceTokenDTO struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Token     string `json:"token,omitempty"`
	Creator   string `json:"creator,omitempty"`
}

// UpdateNamespaceTokenDTO 更新NamespaceToken请求
type UpdateNamespaceTokenDTO struct {
	ID string `json:"id,omitempty"`
}

// GetNamespaceTokenDTO 查询NamespaceToken请求
type GetNamespaceTokenDTO struct {
	Namespace string `json:"namespace,omitempty"`
	Token     string `json:"token,omitempty"`
}

// NamespaceTokenDTO 查询NamespaceToken请求
type NamespaceTokenDTO struct {
	Name      string    `json:"name,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	Token     string    `json:"token,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// DeleteNamespaceTokenDTO 删除NamespaceToken请求
type DeleteNamespaceTokenDTO struct {
	ID string `json:"id,omitempty"`
}

// PageQueryNamespaceTokenDTO 分页查询NamespaceToken请求
type PageQueryNamespaceTokenDTO struct {
	Namespace string    `json:"namespace,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	*constants.PageQuery
	*constants.Order
}

// ValidateTokenReqDTO 校验 Token 请求
type ValidateTokenReqDTO struct {
	Namespace   string `json:"namespace,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

// GetUserNamespacesReqDTO 获取用户的命名空间
type GetUserNamespacesReqDTO struct {
	Username string `json:"username,omitempty"`
}

// GetUserNamespacesRspDTO 获取用户的命名空间
type GetUserNamespacesRspDTO struct {
	CurrentNamespace string   `json:"currentNamespace,omitempty"`
	Namespaces       []string `json:"namespaces,omitempty"`
}

// SetUserCurrentNamespaceReqDTO 设置用户当前的命名空间
type SetUserCurrentNamespaceReqDTO struct {
	CurrentNamespace string `form:"namespace,omitempty" json:"namespace,omitempty"`
}

// GetNamespacesReqDTO 获取所有命名空间
type GetNamespacesReqDTO struct {
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"`
	Username  string `form:"username,omitempty" json:"username,omitempty"`
	*constants.PageQuery
	*constants.Order
}

// NamespacePermissionDTO 命名空间权限
type NamespacePermissionDTO struct {
	Namespace   string    `json:"namespace"`
	Permissions []string  `json:"permissions"`
	Creator     string    `json:"creator,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

// RegisterNamespaceAPITokenReqDTO 注册命名空间 Token 请求
type RegisterNamespaceAPITokenReqDTO struct {
	Name      string `form:"name,omitempty"  json:"name,omitempty"`
	Namespace string `form:"namespace,omitempty" json:"namespace,omitempty"`
	Creator   string `form:"creator,omitempty" json:"creator,omitempty"`
}
