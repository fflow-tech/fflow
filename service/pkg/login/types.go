package login

// CurrentUserData 当前用户信息
type CurrentUserData struct {
	Namespace string `json:"namespace,omitempty"`
	Username  string `json:"username,omitempty"`
	NickName  string `json:"nickName,omitempty"`
	AuthType  string `json:"authType,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Status    int    `json:"status,omitempty"`
}
