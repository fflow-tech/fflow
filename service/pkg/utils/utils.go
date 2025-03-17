// Package utils 通用工具包
package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"mime/multipart"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorhill/cronexpr"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/imdario/mergo"
)

const (
	defaultTimeFormat = "2006-01-02 15:04:05"
	defaultEnv        = "dev"
)

// GetCurrentTimestamp 获取 yyyyMMddHHmmss 格式的当前时间戳
func GetCurrentTimestamp() string {
	return time.Now().Format("20060102150405")
}

// GetCurrentLogTimestamp 获取 2006-01-02 15:04:05.000 格式的当前时间戳
func GetCurrentLogTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// GetDefaultFormatTime 获取默认格式时间戳
func GetDefaultFormatTime(t time.Time) string {
	return t.Format(defaultTimeFormat)
}

// IsJSONString 判断一个字符串是否为 json 格式
func IsJSONString(s string) bool {
	return json.Valid([]byte(s))
}

// IsString 判断类型是否为string
func IsString(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case string:
		return true
	default:
		return false
	}
}

// JsonStrToMap json 字符串转换成map
func JsonStrToMap(j string) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	// 对于空字符串直接空的 map
	if j == "" {
		return m, nil
	}
	err := json.Unmarshal([]byte(j), &m)
	return m, err
}

// UintToStr 数字转换成字符串
func UintToStr(i uint) string {
	if i == 0 {
		return ""
	}

	return strconv.FormatUint(uint64(i), 10)
}

// Uint64ToStr 数字转换成字符串
func Uint64ToStr(i uint64) string {
	if i == 0 {
		return ""
	}

	return strconv.FormatUint(i, 10)
}

// CopyMap 从一个 map 赋值所有的 key 到另外一个 map
func CopyMap(m1 map[string]interface{}, m2 map[string]interface{}) {
	for k, v := range m1 {
		m2[k] = v
	}
}

// StringMapToInterfaceMap 讲一个 map[string]string 转为 map[string]interface{}
func StringMapToInterfaceMap(m1 map[string]string) map[string]interface{} {
	m2 := map[string]interface{}{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

// StructToMap 将 struct 转换成 map
func StructToMap(s interface{}) (map[string]interface{}, error) {
	if s == nil || s == "" {
		return map[string]interface{}{}, nil
	}

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return JsonStrToMap(string(jsonBytes))
}

// StrToInt64 字符串转整数
func StrToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	r, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string to int64: %w", err)
	}

	return r, nil
}

// StrToInt 字符串转整数
func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StrToUInt64 字符串转整数
func StrToUInt64(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}

	r, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string to uint64: %w", err)
	}

	return r, nil
}

// StrToUInt 字符串转整数
func StrToUInt(s string) (uint, error) {
	r, err := StrToUInt64(s)
	return uint(r), err
}

// StrContains 一个数组里面是否包含某个元素
func StrContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// MinUint64 取较小的值
func MinUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// MaxUint64 取较大的值
func MaxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// MaxUint64Str 取较大的值
func MaxUint64Str(a, b string) string {
	ai, _ := StrToUInt64(a)
	bi, _ := StrToUInt64(b)

	if ai > bi {
		return a
	}
	return b
}

// MinUint 取较小的值
func MinUint(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

// GetEnv 获取环境标识
func GetEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		log.Debugf("Cannot get `ENV` environment variable，use `dev` as default")
		return defaultEnv
	}

	return env
}

// BytesToJsonStr bytes转换成json字符串, 如果转换错误返回空字符串
func BytesToJsonStr(b []byte) string {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		return ""
	}

	return buffer.String()
}

// StructToJsonStr 将 struct 转换成 json
// 如果转换失败, 会返回 {} 作为默认值
func StructToJsonStr(s interface{}) string {
	if s == nil || s == "" {
		return "{}"
	}

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		log.Warnf("Failed to struct to json str, caused by %s, use `{}` as return value", err)
		return "{}"
	}

	return string(jsonBytes)
}

// MapToStr 将 map 转换成字符串
func MapToStr(m map[string]interface{}) string {
	if m == nil || len(m) == 0 {
		return ""
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Failed to struct to json str, caused by %s", err)
		return ""
	}

	return string(jsonBytes)
}

// ToOtherInterfaceValue 通过json的方式将一个结构体转换成另一个结构体
func ToOtherInterfaceValue(toValue interface{}, fromValue interface{}) error {
	fromBytes, err := json.Marshal(fromValue)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(fromBytes, toValue); err != nil {
		return err
	}
	return nil
}

// ParseIPFromAddr 从地址中解析出IP
func ParseIPFromAddr(remoteAddr net.Addr) string {
	switch addr := remoteAddr.(type) {
	case *net.UDPAddr:
		if addr.IP == nil {
			return ""
		}
		return addr.IP.String()
	case *net.TCPAddr:
		if addr.IP == nil {
			return ""
		}
		return addr.IP.String()
	default:
		return ""
	}
}

// HasAnyPrefix 判断列表里面的字符串有一个是str的前缀
func HasAnyPrefix(str string, prefixs []string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}

	return false
}

// ParseTemplate 解析模板
func ParseTemplate(text string, m interface{}) string {
	tmpl, _ := template.New("t").Parse(text)
	var b bytes.Buffer
	_ = tmpl.Option().Execute(&b, m)
	return b.String()
}

// InterfaceToSimpleJson interface 转为 simplejson 可以解析的 Json
func InterfaceToSimpleJson(i interface{}) (*simplejson.Json, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	j, err := simplejson.NewJson(b)
	if err != nil {
		return nil, err
	}

	return j, nil
}

// GetStrFromJson 从JSON中取出某个值
func GetStrFromJson(jsonBytes []byte, branch ...string) (string, error) {
	valueJson, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		return "", err
	}
	strValue := valueJson.GetPath(branch...).MustString()
	if strValue != "" {
		return strValue, nil
	}

	intValue := valueJson.GetPath(branch...).MustInt()
	if intValue != 0 {
		return strconv.Itoa(intValue), nil
	}

	floatValue := valueJson.GetPath(branch...).MustFloat64()
	if floatValue != 0 {
		return fmt.Sprintf("%.0f", floatValue), nil
	}

	return "", nil
}

// GetIntFromJson 从JSON中取出某个值
func GetIntFromJson(jsonBytes []byte, branch ...string) (int, error) {
	valueJson, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		return 0, err
	}
	intValue := valueJson.GetPath(branch...).MustInt()
	if intValue != 0 {
		return intValue, nil
	}

	strValue := valueJson.GetPath(branch...).MustString()
	return StrToInt(strValue)
}

// GetUInt64FromJson 从JSON中取出某个值
func GetUInt64FromJson(jsonBytes []byte, branch ...string) (uint64, error) {
	valueJson, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		return 0, err
	}
	uint64Value := valueJson.GetPath(branch...).MustUint64()
	if uint64Value != 0 {
		return uint64Value, nil
	}

	strValue := valueJson.GetPath(branch...).MustString()
	return StrToUInt64(strValue)
}

// IntToBytes 将 int 转换为 byte slice
func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuffer, binary.BigEndian, data)
	return bytebuffer.Bytes()
}

// BytesToInt 将  byte slice 转换为 int
func BytesToInt(bys []byte) (int, error) {
	if len(bys) != 8 {
		return 0, fmt.Errorf("BytesToInt byte slice len must be 8bytes")
	}
	bytebuffer := bytes.NewBuffer(bys)
	var data int64
	err := binary.Read(bytebuffer, binary.BigEndian, &data)
	return int(data), err
}

// StructToBytes 将结构体转换为字节数组
func StructToBytes(value interface{}) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("StructToBytes value can't be nil")
	}
	// 注册数据类型，否则无法反序列化
	gob.Register(value)
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(&value)
	return buf.Bytes(), err
}

// BytesToStruct 将 byte slice 转换为 Struct
func BytesToStruct(b []byte) (interface{}, error) {
	var value interface{}
	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(&value)
	return value, err
}

// MergeMap 合并map, 将b的所有key追加到a上
func MergeMap(a, b map[string]interface{}) (map[string]interface{}, error) {
	if err := mergo.Merge(&b, a); err != nil {
		return nil, err
	}

	return b, nil
}

// IsZero 判断当前值是否为空值， 如果为 nil 也认为当前值空值
func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}

	if reflect.ValueOf(v).IsZero() {
		return true
	}

	return false
}

// GetNextTimeByExpr 根据表达式获取
func GetNextTimeByExpr(expr string, nowTime time.Time) (time.Time, error) {
	if len(strings.Split(expr, " ")) == constants.CronExprByteLength {
		expr = fmt.Sprintf("%s *", expr)
	}
	cronExpr, err := cronexpr.Parse(expr)
	if err != nil {
		log.Errorf("Failed to parse cron expr, caused by %s", err)
		return time.Time{}, err
	}

	return cronExpr.Next(nowTime), nil
}

// GetCurrentGoroutineID 获取当前的协程ID
func GetCurrentGoroutineID() int {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	goIDStr := strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine")[1])
	goID, err := strconv.Atoi(goIDStr)
	if err != nil {
		return 0
	}
	return goID
}

// ReadFile 读取文件
func ReadFile(file *multipart.FileHeader) (string, error) {
	openFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openFile.Close()

	buffer := make([]byte, file.Size)
	_, err = openFile.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

// CopyFile 拷贝文件
func CopyFile(file *multipart.FileHeader, dir, fileName string) error {
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	buffer := make([]byte, file.Size)
	_, err = srcFile.Read(buffer)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dir + "/" + fileName)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = dstFile.Write(buffer)
	return err
}

// GetIntervalTimeByCronExpr 根据定时表达式获取定时间隔时间
// TODO(ankangan):目前只支持递增型(每隔x秒/x分钟/x小时/x日/x月)、定点型:每天早上8点等等
func GetIntervalTimeByCronExpr(cronExpr string) (int64, error) {
	firstTime, err := GetNextTimeByExpr(cronExpr, time.Now())
	if err != nil {
		return 0, err
	}

	secondTime, err := GetNextTimeByExpr(cronExpr, firstTime)
	if err != nil {
		return 0, err
	}

	return secondTime.Unix() - firstTime.Unix(), nil
}

// GetCostTime 获取执行时间，对于未到终态的情况，用当前时间计算耗费时间
func GetCostTime(start int64, complete int64) int64 {
	costTime := complete - start
	if costTime < 0 {
		costTime = time.Now().Unix() - start
	}

	return costTime
}

// AddElementsToSliceIfNotExists 如果元素不存在，往数组里面添加元素
func AddElementsToSliceIfNotExists(slice []string, elements ...string) []string {
	allKeys := map[string]bool{}
	result := []string{}
	for _, v := range slice {
		if _, ok := allKeys[v]; ok {
			continue
		}

		allKeys[v] = true
		result = append(result, v)
	}

	for _, v := range elements {
		if _, ok := allKeys[v]; ok {
			continue
		}

		allKeys[v] = true
		result = append(result, v)
	}
	return result
}

// DeleteElementsFromSlice 删除一个切片中另一个切片的数据
func DeleteElementsFromSlice(slice []string, elements ...string) []string {
	for _, element := range elements {
		for index, value := range slice {
			if value == element {
				slice = append(slice[:index], slice[index+1:]...)
				break
			}
		}
	}
	return slice
}
