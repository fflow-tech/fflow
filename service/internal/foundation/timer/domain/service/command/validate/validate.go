// Package validate web 参数校验
package validate

import (
	"fmt"
	"strings"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

const (
	// legalCronLength 合法定时表达式参数
	legalCronLength = 7
	// legalTimeFormat 合法的时间格式
	legalTimeFormat = "2006-01-02 15:04:05"
)

// CheckCreateDefParam  创建定时器定义参数校验
func CheckCreateDefParam(req interface{}) error {
	createTimerDefDTO := req.(*dto.CreateTimerDefDTO)
	// 校验定时配置参数
	if err := validateTimerParam(createTimerDefDTO); err != nil {
		return err
	}

	if err := validateEndTime(createTimerDefDTO); err != nil {
		return err
	}

	if err := validateExecuteTimeLimit(createTimerDefDTO); err != nil {
		return err
	}

	// 校验回调参数信息
	return validateNotifyParams(createTimerDefDTO)
}

// validateDelayTimerFiled 校验
func validateDelayTimerFiled(d *dto.CreateTimerDefDTO) error {
	if utils.IsZero(d.DelayTime) || (d.TriggerType != entity.TriggerOnce.ToInt()) {
		return fmt.Errorf("createDef failed `DelayTime` must not be zero, "+
			"`TriggerType` must equal 1, DelayTime:[%s]、TriggerType:[%d]", d.DelayTime, d.TriggerType)
	}

	return nil
}

// validateCronTimerFiled 校验定时器字段
func validateCronTimerFiled(d *dto.CreateTimerDefDTO) error {
	if utils.IsZero(d.Cron) || utils.IsZero(d.TriggerType) {
		return fmt.Errorf("createDef filed `Cron`、`TriggerType` must not be zero, "+
			"Cron:[%s]、TriggerType:[%d]", d.Cron, d.TriggerType)
	}

	return validateCronLength(d.Cron)
}

// validateCronLength 定时表达式长度校验
func validateCronLength(cronExpr string) error {
	split := strings.Split(cronExpr, " ")
	if len(split) != legalCronLength {
		return fmt.Errorf("CreateDef filed `Cron` expr must 7 bit, please check cron expr")
	}

	return nil
}

// validateNotifyParams 校验回调参数信息
func validateNotifyParams(d *dto.CreateTimerDefDTO) error {
	switch d.NotifyType {
	case entity.RPC.ToInt():
		return validateRPCParam(d)
	case entity.HTTP.ToInt():
		return validateHTTPParam(d)
	default:
		return fmt.Errorf("failed to validateNotifyParams , caused by NotifyType:%v error", d.NotifyType)
	}
}

func validateRPCParam(d *dto.CreateTimerDefDTO) error {
	rpcParam := d.NotifyRpcParam
	if utils.IsZero(rpcParam) {
		return fmt.Errorf("validateHTTPParam NotifyRpcParam equal `RPC`, `NotifyRpcParam` must not be zero")
	}

	if utils.IsZero(rpcParam.Service) || (utils.IsZero(rpcParam.RpcName) && utils.IsZero(rpcParam.Method)) {
		return fmt.Errorf("validateHTTPParam `NotifyRpcParam` filed `service`、`rpcName&method` must not "+
			"be zero, service:[%s]、rpcName:[%s]、method:[%s]", rpcParam.Service, rpcParam.RpcName, rpcParam.Method)
	}

	return nil
}

func validateHTTPParam(d *dto.CreateTimerDefDTO) error {
	httpParam := d.NotifyHttpParam
	if utils.IsZero(httpParam) {
		return fmt.Errorf("validateHTTPParam NotifyType equal `HTTP` ,`NotifyHttpParam` must not be zero")
	}

	if utils.IsZero(httpParam.Method) || utils.IsZero(httpParam.Url) {
		return fmt.Errorf("validateHTTPParam  `NotifyHttpParam` filed `Method`、`Url` must not be zero"+
			" Method:[%s]、Url:[%s]", httpParam.Method, httpParam.Url)
	}

	return nil
}

// validateTimerParam 校验定时参数信息
func validateTimerParam(d *dto.CreateTimerDefDTO) error {
	// 如果是定时器类型，校验 cron
	if d.TimerType == entity.CronTimer.ToInt() {
		return validateCronTimerFiled(d)
	}

	return validateDelayTimerFiled(d)
}

// validateEndTime 校验 EndTime
func validateEndTime(d *dto.CreateTimerDefDTO) error {
	if utils.IsZero(d.EndTime) {
		return nil
	}

	endTime, err := time.ParseInLocation(legalTimeFormat, d.EndTime, time.Local)
	if err != nil {
		return err
	}

	nowTime := time.Now()
	if !nowTime.Before(endTime) {
		return fmt.Errorf("validateEndTime `EndTime` must  after NowTime, EndTime:[%s]、NowTime:[%s]",
			d.EndTime, nowTime.Format(legalTimeFormat))
	}

	return nil
}

// validateExecuteTimeLimit 校验执行时长限制.
func validateExecuteTimeLimit(d *dto.CreateTimerDefDTO) error {
	if d.ExecuteTimeLimit < 0 || d.ExecuteTimeLimit > 15 {
		return fmt.Errorf("value of executeTimeLimit must between 0 and 15")
	}
	return nil
}
