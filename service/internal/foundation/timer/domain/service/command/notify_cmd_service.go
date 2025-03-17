package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/ports"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/concurrency"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/monitor"
	"github.com/fflow-tech/fflow/service/pkg/limiter"
	"github.com/fflow-tech/fflow/service/pkg/remote"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/gorhill/cronexpr"
	"github.com/fflow-tech/fflow/service/pkg/log"
)

type reporter interface {
	ReportTriggerRecord(app string)
	ReportTimerCostRecord(app string, cost float64)
}

type trafficPool interface {
	Get(opts ...limiter.Option) error
}

// NotifyCommandService 通知服务
type NotifyCommandService struct {
	remoteRepo      ports.RemoteRepository
	timerTaskRepo   ports.TimerTaskRepository
	timerDefRepo    ports.TimerDefRepository
	pollingTaskRepo ports.PollingTaskRepository
	workerPool      concurrency.WorkerPool
	reporter        reporter
	trafficPool     trafficPool
}

// NewNotifyCommandService 通知服务构造函数
func NewNotifyCommandService(repoSet *ports.RepoProviderSet, workerPool *concurrency.GoWorkerPool,
	trafficPool *limiter.TrafficPool, reporter *monitor.Reporter) *NotifyCommandService {
	return &NotifyCommandService{
		remoteRepo:      repoSet.RemoteRepo(),
		timerTaskRepo:   repoSet.TimerTaskRepo(),
		timerDefRepo:    repoSet.TimerDefRepo(),
		pollingTaskRepo: repoSet.PollingTaskRepo(),
		workerPool:      workerPool,
		reporter:        reporter,
		trafficPool:     trafficPool,
	}
}

// DeleteRunHistories 删除过期的历史记录
func (n *NotifyCommandService) DeleteRunHistories() error {
	return n.timerTaskRepo.DeleteRunHistories()
}

// TimerListSendNotify 定时器批量触发 内部谨慎使用
func (n *NotifyCommandService) TimerListSendNotify(defIDs []string) error {
	for _, defID := range defIDs {
		if err := n.SendNotify(defID); err != nil {
			log.Errorf("Failed to TimerListSendNotify, defID %v caused by %v", defID, err)
			return err
		}
	}
	return nil
}

// SendNotify 发送通知.
func (n *NotifyCommandService) SendNotify(hashID string) error {
	log.Infof("SendNotify hashID %v", hashID)
	timerDef, err := n.timerDefRepo.GetTimerDef(&dto.GetTimerDefDTO{
		DefID: hashID,
	})
	if err != nil {
		log.Errorf("Failed to SendNotify GetTimerDef, defID: %s, err: %v", hashID, err)
		return err
	}

	// 对 pending 表进行清点，由于依赖于 save 表，需要顺序化执行.
	saveTask, err := n.timerTaskRepo.GetSaveTimerTask(hashID)
	if err != nil {
		log.Errorf("failed to get save timer task, defID: %s, err: %v", hashID, err)
	} else {
		if err := n.delPendingTimerTask(hashID, time.Unix(0, saveTask.UnixTime)); err != nil {
			log.Errorf("failed to delete pending timer task, defID: %s, err: %v", hashID, err)
		}
	}

	ok, err := n.checkCanNotify(timerDef)
	if err != nil {
		log.Errorf("Failed to SendNotify checkCanNotify, caused by %v", err)
		return err
	}

	if !ok {
		log.Warnf("Warn to SendNotify checkCanNotify, caused by not can notify")
		return nil
	}

	if err = n.workerPool.Submit(func() {
		// 执行通知
		if err := n.notify(timerDef, saveTask); err != nil {
			log.Errorf("notify failed, defID: %s, err: %v", hashID, err)
		}
	}); err != nil {
		log.Errorf("submit notify task to worker pool failed, defID: %s, err: %v", hashID, err)
	}
	return n.registerNext(timerDef)
}

// notifyPreprocess 通知前处理动作.
func (n *NotifyCommandService) notifyPreprocess(timer *entity.TimerDef,
	saveTask *dto.SaveTimerTaskDTO) (time.Time, error) {
	getTokenBegin := time.Now()

	// 执行限流操作.
	err := n.trafficPool.Get()

	startTime := time.Now()
	log.Infof("get token from traffic pool, cost: %d seconds, err: %v",
		startTime.Sub(getTokenBegin)/time.Second, err)

	// 定时器触发记录上报.
	n.workerPool.Submit(func() {
		n.reporter.ReportTriggerRecord(timer.App)
	})

	// 定时器触发时延数据上报.
	if saveTask != nil && saveTask.UnixTime > 0 {
		n.workerPool.Submit(func() {
			costTime := getCostTimeOfMillisecond(saveTask.UnixTime, startTime.UnixNano())
			n.reporter.ReportTimerCostRecord(timer.App, float64(costTime))
			log.Infof("create timer_trigger_cost record successfully, app: %s, defID: %s,"+
				" costTime: %d, taskInfo: %d, endTime: %d", timer.App, timer.DefID,
				costTime, saveTask.UnixTime, startTime.UnixNano())
		})
	}

	return startTime, err
}

// notify 通知.
func (n *NotifyCommandService) notify(timer *entity.TimerDef, saveTask *dto.SaveTimerTaskDTO) error {
	var resp map[string]interface{}

	startTime, err := n.notifyPreprocess(timer, saveTask)
	if err == nil {
		seconds := 15 * time.Second
		if timer.ExecuteTimeLimit > 0 {
			seconds = time.Duration(timer.ExecuteTimeLimit) * time.Second
		}
		switch timer.NotifyType {
		case entity.HTTP:
			resp, err = n.notifyHttp(timer.NotifyHttpParam, seconds)
		default:
			log.Errorf("Failed to notify , caused by notify type:%d", timer.NotifyType)
			resp, err = nil, nil
		}
	}

	n.notifyPostProcess(timer, startTime, resp, err)
	return err
}

// notifyPostProcess 通知后处理.
func (n *NotifyCommandService) notifyPostProcess(timer *entity.TimerDef,
	startTime time.Time, notifyResp map[string]interface{}, notifyErr error) {
	if notifyErr != nil {
		n.workerPool.Submit(func() {
			n.notifyFailureHandler(timer, startTime, notifyErr)
		})
		return
	}

	n.workerPool.Submit(func() {
		if err := n.createTimerHistory(timer, startTime, time.Now(),
			entity.Succeed.String(), utils.MapToStr(notifyResp)); err != nil {
			log.Errorf("failed to createTimerHistory succeed, caused by %v", err)
		}
	})
}

// notifyHttp http通知
func (n *NotifyCommandService) notifyHttp(notifyStr string, timeout time.Duration) (map[string]interface{}, error) {
	notifyHttpParam := &dto.NotifyHttpParam{}
	if err := json.Unmarshal([]byte(notifyStr), &notifyHttpParam); err != nil {
		return nil, err
	}
	body, err := utils.JsonStrToMap(notifyHttpParam.Body)
	if err != nil {
		return nil, err
	}
	header, err := n.getHttpHeader(notifyHttpParam)
	if err != nil {
		return nil, err
	}
	req := &remote.CallHTTPReqDTO{
		// 非 mock 模式
		MockMode: false,
		Method:   notifyHttpParam.Method,
		URL:      notifyHttpParam.Url,
		Header:   header,
		Body:     body,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return n.remoteRepo.CallHTTP(ctx, req)
}

// getHttpHeader 获取 http 请求头
func (n *NotifyCommandService) getHttpHeader(notifyHttpParam *dto.NotifyHttpParam) (map[string]string, error) {
	header := make(map[string]string)
	if utils.IsZero(notifyHttpParam.Header) {
		return header, nil
	}
	if err := json.Unmarshal([]byte(notifyHttpParam.Header), &header); err != nil {
		return nil, err
	}

	return header, nil
}

// notifyFailureHandler 通知失败处理
func (n *NotifyCommandService) notifyFailureHandler(timerDef *entity.TimerDef, startTime time.Time, err error) {
	log.Errorf("failed to notify, caused by %v", err)
	// 发送消息通知
	n.sendAlertMsg(timerDef, err)
	status := entity.Failed.String()
	if errors.Is(err, context.DeadlineExceeded) {
		status = entity.Timeout.String()
	}
	// 更新执行历史记录
	if err := n.createTimerHistory(timerDef, startTime, time.Now(), status, err.Error()); err != nil {
		log.Errorf("failed to createTimerHistory failed, caused by %v", err)
	}
}

// sendAlertMsg 发送告警信息
func (n *NotifyCommandService) sendAlertMsg(timerDef *entity.TimerDef, err error) {
	contentTemplate := config.GetMsgContentTemplate()
	msgInfo := map[string]interface{}{
		"DefID":     timerDef.DefID,
		"TimerName": timerDef.Name,
		"Fail":      err,
	}
	content, err := config.ParseTemplate(contentTemplate, msgInfo)
	if err != nil {
		return
	}
	msg := fmt.Sprintf("%s\n%s", config.GetMsgTitle(), content)
	if err := n.remoteRepo.SendMsgToUser(timerDef.Creator, msg); err != nil {
		log.Errorf("failed to SendMsgToUser, caused by %v", err)
	}
}

// registerNext 注册下一个
func (n *NotifyCommandService) registerNext(tickerDef *entity.TimerDef) error {
	// 注册前先清除 save 表记录.
	if err := n.timerTaskRepo.DeleteSaveTimerTask(tickerDef.DefID); err != nil {
		log.Errorf("failed to delete timer task timeSlice, defID: %s, err: %v", tickerDef.DefID, err)
	}
	canNext, err := n.checkTimerCanContinue(tickerDef)
	if err != nil {
		return err
	}
	if tickerDef.TimerType == entity.DelayTimer {
		return n.deleteTimerDef(tickerDef)
	}
	if !canNext {
		return nil
	}
	// 找到最近的下次触发时间并丢入对应的 work_table
	// 找到符合条件的就发送到期定时器
	expr, err := cronexpr.Parse(tickerDef.Cron)
	if err != nil {
		log.Errorf("failed to registerNext cronexpr parse failed, caused by %v%v", err)
		return err
	}
	nowTime := time.Now()
	nextTimeout := expr.Next(nowTime)
	if nextTimeout.UnixNano() < 0 {
		log.Errorf("failed to registerNext, invalid next time: %+v, defID: %s", nextTimeout, tickerDef.DefID)
		return nil
	}
	return n.addTimerTask(tickerDef.DefID, nextTimeout)
}
func (n *NotifyCommandService) checkTimerCanContinue(timerDef *entity.TimerDef) (bool, error) {
	if timerDef.TimerType == entity.DelayTimer {
		// 延时类型的定时器 只触发一次
		return false, nil
	}
	if timerDef.TriggerType == entity.TriggerOnce { //只触发一次 则不再触发
		return false, nil
	}
	return n.checkCanNotify(timerDef)
}
func (n *NotifyCommandService) checkCanNotify(timerDef *entity.TimerDef) (bool, error) {
	// 检查对应的timer的状态
	if timerDef.Status != entity.Enabled {
		log.Warnf("failed to SendNotify, caused by ticker status is disable")
		return false, nil
	}
	if timerDef.EndTime == "" {
		return true, nil
	}
	endTime, err := time.ParseInLocation(entity.DelayTimeFormat, timerDef.EndTime, time.Local)
	if err != nil {
		return false, err
	}
	nowTime := time.Now()
	if endTime.Before(nowTime) {
		log.Infof("nowTimeout timer defID %v is end nowTime is:%v endTime is:%v",
			timerDef.DefID, nowTime, endTime)
		return false, nil
	}
	return true, nil
}
func (n *NotifyCommandService) addTimerTask(defID string, nextTime time.Time) error {
	bucketNum, err := n.pollingTaskRepo.GetTaskBucketID(defID)
	if err != nil {
		return err
	}
	addDTO := &dto.AddTimerTaskDTO{
		HashID:    defID,
		TimerTime: nextTime,
		BucketID:  bucketNum,
	}
	return n.timerTaskRepo.AddTimerTask(addDTO)
}
func (n *NotifyCommandService) delPendingTimerTask(defID string, curTime time.Time) error {
	return n.timerTaskRepo.DelPendingTimerTask(defID, curTime)
}
func (n *NotifyCommandService) createTimerHistory(
	def *entity.TimerDef, startTime time.Time, endTime time.Time, status string, rsp string) error {
	create := &dto.CreateRunHistoryDTO{
		DefID:    def.DefID,
		Name:     def.Name,
		RunTimer: startTime.Format(dto.TimerTriggerTimeFormat),
		Status:   status,
		Output:   rsp,
	}
	if endTime != startTime {
		create.CostTime = getCostTimeOfMillisecond(startTime.UnixNano(), endTime.UnixNano())
	}
	_, err := n.timerTaskRepo.CreateHistory(create)
	return err
}

// deleteTimerDef 删除定时器定义
func (n *NotifyCommandService) deleteTimerDef(def *entity.TimerDef) error {
	if def.DeleteType != entity.TriggerDelete {
		return nil
	}
	log.Infof("deleteTimerDef defID: %d", def.DefID)
	return n.timerDefRepo.DeleteTimerDef(&dto.DeleteTimerDefDTO{
		DefID: def.DefID,
	})
}

// getCostTimeOfMillisecond 获取时间差，单位为 ms
func getCostTimeOfMillisecond(start int64, complete int64) int64 {
	costTime := complete - start
	if costTime < 0 {
		costTime = time.Now().UnixNano() - start
	}

	return costTime / 1e6
}

// ManualTriggerSend 手动触发发送
func (n *NotifyCommandService) ManualTriggerSend(hashID string) error {
	log.Infof("SendNotify hashID %v", hashID)
	timerDef, err := n.timerDefRepo.GetTimerDef(&dto.GetTimerDefDTO{
		DefID: hashID,
	})
	if err != nil {
		log.Errorf("Failed to ManualTriggerSend GetTimerDef, defID: %s, err: %v", hashID, err)
		return err
	}

	seconds := 15 * time.Second
	if timerDef.ExecuteTimeLimit > 0 {
		seconds = time.Duration(timerDef.ExecuteTimeLimit) * time.Second
	}

	var resp map[string]interface{}
	switch timerDef.NotifyType {
	case entity.HTTP:
		resp, err = n.notifyHttp(timerDef.NotifyHttpParam, seconds)
	default:
		log.Errorf("Failed to ManualTriggerSend , caused by notify type:%d", timerDef.NotifyType)
		resp, err = nil, nil
	}

	n.notifyPostProcess(timerDef, time.Now(), resp, err)
	return err
}

// ManualTriggerSendList 手动触发发送批量 内部谨慎使用
func (n *NotifyCommandService) ManualTriggerSendList(defIDs []string) error {
	for _, defID := range defIDs {
		if err := n.ManualTriggerSend(defID); err != nil {
			log.Errorf("Failed to ManualTriggerSendList, defID %v caused by %v", defID, err)
			return err
		}
	}
	return nil
}
