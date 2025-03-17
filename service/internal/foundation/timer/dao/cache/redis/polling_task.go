package redis

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/redis"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/fflow-tech/fflow/service/pkg/log"
)

const (
	bucketNum        = "bucketNum"
	minBucketNum     = 3
	maxBucketNum     = 10000
	defaultBucketNum = 3

	timeSliceTimeout  = 3 * 24 * 60 * 60 // 时间分片的超时时间 现在先设置3天
	timeSliceMaxValue = "success"        // 时间分片最大数值
	timeSliceMinValue = 40               // 时间分片最小数值

	distributeLockTime = 2 * time.Second // 分布式锁的时间
)

// PollingTask 轮询任务
type PollingTask struct {
	config      config.PollingTaskConfig
	redisClient *redis.Client
}

// NewPollingTaskClient 新建轮询器的客户端
func NewPollingTaskClient(redisClient *redis.Client) *PollingTask {
	return &PollingTask{
		redisClient: redisClient,
		// TODO(@weixxxu): 遵循依赖注入原则.
		config: config.GetPollingTaskConfig(),
	}
}

// GetBucketNum 获取桶数量
func (p *PollingTask) GetBucketNum() int {
	return p.config.BucketNum
}

// GetTaskBucketID 获取任务所属桶名
func (p *PollingTask) GetTaskBucketID(defID string) (string, error) {
	defIDNum, err := strconv.Atoi(defID[len(defID)-1:]) // 这里截取最后一位
	if err != nil {
		return "", err
	}
	bucketID := p.getRandBucketNum(defIDNum, p.GetBucketNum())
	log.Infof("GetTaskBucketID defID %v bucketID %v", defID, bucketID)
	return strconv.Itoa(bucketID), nil
}

func (p *PollingTask) getRandBucketNum(defID, bucketNum int) int {
	rand.Seed(time.Now().UnixNano() + int64(defID))
	return rand.Intn(bucketNum)
}

// SetBucketNum 设置桶的数量 桶的数量写到redis中 这里只能增大 不能变小
func (p *PollingTask) SetBucketNum(num int) error {
	if num < minBucketNum || num > maxBucketNum {
		return fmt.Errorf("failed to SetBucketNum, caused by num err %d", num)
	}
	return p.redisClient.SetIntKey(context.Background(), bucketNum, num)
}

// SetTimeSlice 设置时间片值
func (p *PollingTask) SetTimeSlice(timeDuration string) error {
	lockName := p.getDistributeLockName(timeDuration)
	lock := p.redisClient.GetDistributeLock(lockName, distributeLockTime)
	// 使用个分布式锁来设值
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()
	v, err := p.redisClient.Exists(context.Background(), timeDuration)
	if err != nil {
		return err
	}

	if v == 1 { // 当值存在时 因为锁的问题这里需要再判断一次
		if err := p.seizeTimeDuration(timeDuration); err != nil {
			return err
		}
	}

	timeout := time.Now().Add(timeSliceMinValue * time.Second) // 设置超时时间值
	return p.setTimeDurationValue(timeDuration, timeout.Format(dto.TimerTriggerTimeFormat))
}

// GetTimeSlice 获取时间片
func (p *PollingTask) GetTimeSlice(timeDuration string) error {
	v, err := p.redisClient.Exists(context.Background(), timeDuration)
	if err != nil {
		return err
	}
	if v == 0 {
		// 当key不存在时 则获得这个时间分片的权限
		return p.SetTimeSlice(timeDuration)
	}

	// 判断能否抢占任务
	if err := p.seizeTimeDuration(timeDuration); err != nil {
		return err
	}
	// 当执行任务超时时 则直接抢占任务
	return p.SetTimeSlice(timeDuration)
}

// deleteTimeSlice 删除时间切片
func (p *PollingTask) deleteTimeSlice(timeSliceString string) error {
	return p.redisClient.Del(context.Background(), timeSliceString)
}

// checkTimeSliceSuccess 检查时间片是否完成
func (p *PollingTask) checkTimeSliceSuccess(timeSliceString string) bool {
	return timeSliceString == timeSliceMaxValue
}

// checkTimeSliceTimeout 检查时间片是否超时
func (p *PollingTask) checkTimeSliceTimeout(timeSliceString string) bool {
	duration, err := time.ParseInLocation(dto.TimerTriggerTimeFormat, timeSliceString, time.Local)
	if err != nil {
		log.Errorf("Failed to checkTimeSliceTimeout ParseInLocation, caused by %v", err)
		return false
	}

	log.Infof("checkTimeSliceTimeout duration is:%v nowTime is:%v", duration, time.Now())
	if time.Now().Before(duration) {
		return false
	}
	return true
}

// SuccessTimeSlice 标记时间片任务完成
func (p *PollingTask) SuccessTimeSlice(timeDuration string) error {
	lockName := p.getDistributeLockName(timeDuration)
	lock := p.redisClient.GetDistributeLock(lockName, distributeLockTime)
	// 使用个分布式锁来设值
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	timeSliceString, err := p.redisClient.Get(context.Background(), timeDuration)
	if err != nil {
		return err
	}

	if p.checkTimeSliceSuccess(timeSliceString) { // 当任务已经完成时 直接返回err
		return fmt.Errorf(" Failed to checkTimeSliceSuccess GoroutineID:%d, caused by timeSlice is success",
			utils.GetCurrentGoroutineID())
	}
	return p.setTimeDurationValue(timeDuration, timeSliceMaxValue)
}

func (p *PollingTask) getDistributeLockName(timeDuration string) string {
	return fmt.Sprintf("LOCK_%s", timeDuration)
}

func (p *PollingTask) setTimeDurationValue(timeDuration, timeSliceValue string) error {
	return p.redisClient.SetEx(context.Background(), timeDuration, timeSliceValue, timeSliceTimeout)
}

// seizeTimeDuration 抢占时间片
func (p *PollingTask) seizeTimeDuration(timeDuration string) error {
	timeSliceString, err := p.redisClient.Get(context.Background(), timeDuration)
	if err != nil {
		return err
	}

	// 当任务已经完成时 直接返回err
	if p.checkTimeSliceSuccess(timeSliceString) {
		return fmt.Errorf("failed to checkTimeDurationValueCanSet, caused by timeSlice is done")
	}

	// 当已经被别的任务抢占时 需要判断是否已超时
	if !p.checkTimeSliceTimeout(timeSliceString) {
		return fmt.Errorf("failed to checkTimeDurationValueCanSet, caused by timeSliceString is runnig")
	}
	return nil
}
