package entity

import (
	"bytes"
	"encoding/json"
)

// User 用户实体
type User struct {
	ID       string     `json:"id,omitempty"`
	NickName string     `json:"nickName,omitempty"`
	AuthType AuthType   `json:"authType,omitempty"`
	Username string     `json:"username,omitempty"`
	Password string     `json:"password,omitempty"`
	Email    string     `json:"email,omitempty"`
	Phone    string     `json:"phone,omitempty"`
	Avatar   string     `json:"avatar,omitempty"`
	Status   UserStatus `json:"status"`
}

// UserStatus 用户状态状态枚举
type UserStatus struct {
	intValue int
	strValue string
}

// 流程状态枚举类型
var (
	Disabled = UserStatus{1, "disabled"} // 未激活
	Enabled  = UserStatus{2, "enabled"}  // 已激活
)

// IntValue 整数值
func (s UserStatus) IntValue() int {
	return s.intValue
}

// String 整数值
func (s UserStatus) String() string {
	return s.strValue
}

var (
	intDefStatusMap = map[int]UserStatus{
		Disabled.IntValue(): Disabled,
		Enabled.IntValue():  Enabled,
	}
	strDefStatusMap = map[string]UserStatus{
		Disabled.String(): Disabled,
		Enabled.String():  Enabled,
	}
)

// GetUserStatusByIntValue 通过整数值返回状态枚举
func GetUserStatusByIntValue(i int) UserStatus {
	return intDefStatusMap[i]
}

// GetUserStatusByStrValue 通过字符串值返回状态枚举
func GetUserStatusByStrValue(s string) UserStatus {
	return strDefStatusMap[s]
}

// MarshalJSON 重写序列化方法
func (s UserStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON 重写反序列化方法
func (s *UserStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = strDefStatusMap[j]
	return nil
}

// AuthType 认证类型枚举
type AuthType string

// 认证类型枚举
var (
	Password AuthType = "password"
	Email    AuthType = "email"
	Github   AuthType = "github"
)

// String 字符串
func (s AuthType) String() string {
	return string(s)
}
