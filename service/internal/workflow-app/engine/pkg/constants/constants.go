// Package constants 提供通用常量的定义
package constants

var (
	ServiceName                   = "engine"              // 引擎服务名称
	TdmqMaxCacheTime        int64 = 1 * 60 * 60 * 24 * 10 // tdmq最大缓存时间，单位s
	MinCronIntervalTime     int64 = 1 * 60                // 最小定时间隔时间
	MaxIntervalTimeCronExpr       = "0 0 0 1 */1  * *"    // 最大间隔时间表达式(一个月)
	PxCronJobNamePrefix           = "cronTrigger"         // 分布式任务name前缀
	CronExprByteLength            = 6                     // 定时表达式位数
	ThisNode                      = "this"                // 当前节点
)
