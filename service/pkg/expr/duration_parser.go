package expr

import (
	"fmt"
	"regexp"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var durationTypeMap = map[string]time.Duration{
	"s": time.Second,        // 秒
	"m": time.Minute,        // 分钟
	"h": time.Hour,          // 小时
	"d": 24 * time.Hour,     // 天
	"w": 7 * 24 * time.Hour, // 周
}

var (
	durationFormatRE, _ = regexp.Compile("^([1-9][0-9]*)+[smhdw]$") // 校验表达式格式的正则
	durationTypeRE, _   = regexp.Compile("[smhdw]$")                // 获取表达式的单位的正则
	durationPrefixRE, _ = regexp.Compile("^[0-9]+")                 // 获取表达式的数字的正则
)

// ParseDuration 获取持续时间
func ParseDuration(duration string) (time.Duration, error) {
	// 校验配置格式
	if !durationFormatRE.MatchString(duration) {
		return 0, fmt.Errorf("duration format err %s", duration)
	}

	// 获取配置单位类型
	durationType := durationTypeRE.FindString(duration)

	// 获取配置前缀数量
	durationPrefix := durationPrefixRE.FindString(duration)
	durationNum, err := utils.StrToUInt64(durationPrefix)
	if err != nil {
		return 0, err
	}

	timeDuration, ok := durationTypeMap[durationType]
	if !ok {
		return 0, fmt.Errorf("illegal duration format:%s", duration)
	}

	return time.Duration(durationNum) * timeDuration, nil
}
