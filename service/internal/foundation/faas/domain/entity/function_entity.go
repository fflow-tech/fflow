package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LanguageType 语言的类型
type LanguageType struct {
	strValue   string
	fileSuffix string
}

var languages = map[string]LanguageType{
	Js.strValue:     Js,
	Golang.strValue: Golang,
}

// UnmarshalJSON 重写反序列化方法
func (l *LanguageType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	var ok bool
	*l, ok = languages[strings.ToLower(j)]
	if !ok {
		return fmt.Errorf("not found %s language type", j)
	}

	return nil
}

// MarshalJSON 重写序列化方法
func (l LanguageType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(l.strValue)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// NewLanguageType 初始化
func NewLanguageType(strValue string, fileSuffix string) LanguageType {
	return LanguageType{strValue: strValue, fileSuffix: fileSuffix}
}

var (
	Js     = NewLanguageType("javascript", ".js")
	Golang = NewLanguageType("golang", ".go")
)

// String 转换成字符串
func (l LanguageType) String() string {
	return l.strValue
}

// FileSuffix 文件后缀
func (l LanguageType) FileSuffix() string {
	return l.fileSuffix
}

var (
	strLanguageTypeMap = map[string]LanguageType{
		Js.String():     Js,
		Golang.String(): Golang,
	}
)

// GetLanguageTypeByStrValue 通过字符串值返回枚举
func GetLanguageTypeByStrValue(s string) LanguageType {
	return strLanguageTypeMap[s]
}

// RunStatus 运行状态
type RunStatus string

const (
	Running RunStatus = "running"
	Succeed RunStatus = "succeed"
	Failed  RunStatus = "failed"
)

// Function 函数实体
type Function struct {
	ID           int          `json:"id,omitempty"`
	Namespace    string       `json:"namespace,omitempty"`     // 命名空间
	Creator      string       `json:"creator,omitempty"`       // 创建人
	Language     LanguageType `json:"language,omitempty"`      // 所使用的语言
	Code         string       `json:"code,omitempty"`          // 代码
	Token        string       `json:"token,omitempty"`         // token
	InputSchema  string       `json:"input_schema,omitempty"`  // 函数入参格式
	OutputSchema string       `json:"output_schema,omitempty"` // 函数返回结果格式
	Description  string       `json:"description,omitempty"`   // 描述
	Name         string       `json:"function,omitempty"`      // 函数名
	Updater      string       `json:"updater,omitempty"`       // 更新人
	Version      int          `json:"version,omitempty"`       // 版本号
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`    // 更新时间
	CreatedAt    time.Time    `json:"created_at,omitempty"`    // 创建时间
}

// Metadata 函数元数据
type Metadata struct {
	*Function
}

func (f *Metadata) Namespace() string {
	return f.Function.Namespace
}
func (f *Metadata) ID() string {
	return strconv.Itoa(f.Function.ID)
}
func (f *Metadata) Name() string {
	return f.Function.Name

}
func (f *Metadata) Version() int {
	return f.Function.Version
}
func (f *Metadata) Attribute(key string) (any, error) {
	return nil, nil
}

// RunHistory 执行历史实体
type RunHistory struct {
	ID        uint      `json:"id,omitempty"`         // 执行记录 ID
	Namespace string    `json:"namespace,omitempty"`  // 命名空间
	Name      string    `json:"name,omitempty"`       // 函数名称
	Operator  string    `json:"operator,omitempty"`   // 执行者
	Input     string    `json:"input,omitempty"`      // 函数入参
	Output    string    `json:"output,omitempty"`     // 执行结果
	Log       string    `json:"log,omitempty"`        // 执行日志
	CostTime  int       `json:"cost_time,omitempty"`  // 执行耗时
	Version   int       `json:"version,omitempty"`    // 版本号
	Status    string    `json:"status,omitempty"`     // 当前状态
	UpdatedAt time.Time `json:"updated_at,omitempty"` // 更新时间
	CreatedAt time.Time `json:"created_at,omitempty"` // 创建时间
}
