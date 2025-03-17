package repo

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/cache/redis"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/dao/storage"
)

// PollingTaskRepo 轮询任务实体
type PollingTaskRepo struct {
	pollingTaskRepo storage.PollingTaskDAO
}

// NewPollingTaskRepo 实体构造函数
func NewPollingTaskRepo(d *redis.PollingTask) *PollingTaskRepo {
	return &PollingTaskRepo{pollingTaskRepo: d}
}

// GetTaskBucketID 获取hash所属桶ID
func (t *PollingTaskRepo) GetTaskBucketID(hashID string) (string, error) {
	return t.pollingTaskRepo.GetTaskBucketID(hashID)
}

// SetBucketNum 设置桶数量
func (t *PollingTaskRepo) SetBucketNum(num int) error {
	return t.pollingTaskRepo.SetBucketNum(num)
}

// SetTimeSlice 设置时间切片
func (t *PollingTaskRepo) SetTimeSlice(timeDuration string) error {
	return t.pollingTaskRepo.SetTimeSlice(timeDuration)
}

// GetTimeSlice 获取时间切片
func (t *PollingTaskRepo) GetTimeSlice(timeDuration string) error {
	return t.pollingTaskRepo.GetTimeSlice(timeDuration)
}

// SuccessTimeSlice 标记时间切片成功
func (t *PollingTaskRepo) SuccessTimeSlice(timeDuration string) error {
	return t.pollingTaskRepo.SuccessTimeSlice(timeDuration)
}

// GetBucketNum 获取桶数量
func (t *PollingTaskRepo) GetBucketNum() int {
	return t.pollingTaskRepo.GetBucketNum()
}
