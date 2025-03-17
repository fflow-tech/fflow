package execution

import "time"

var (
	instTerminalErrFormat             = "[%s]workflow inst is already terminal"                // 实例已经终止错误模板
	instRepeatOperationErrFormat      = "[%s]workflow inst is already [%s]"                    // 实例重复操作错误模板
	notGetLockErr                     = "[%d]failed to get lock:[%s], caused by %s: %w"        // 没有获取到锁错误模板
	instLockRetryDelay                = 50 * time.Millisecond                                  // 重试拿实例锁的时间间隔
	instLockRetry                     = 20                                                     // 重试拿流程实例锁的次数
	instLockExpireTime                = 20 * time.Second                                       // 流程实例锁的过期时间
	defLockRetryDelay                 = 50 * time.Millisecond                                  //  重试拿流程定义锁的时间间隔
	defLockRetry                      = 20                                                     // 重试拿定义锁的次数
	defLockExpireTime                 = 20 * time.Second                                       // 流程定义锁的过期时间
	sendAlertKeyTtl                   = 100 * 365 * 24 * 60 * 60 * time.Second                 // 避免重复发送超时消息的缓存 key 超时时间
	successConditionNotMatchErrFormat = "node `successCondition` [%s] is not match, output=%s" // 成功条件不满足的错误模板
)
