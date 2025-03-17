// Package monitor 系统相关监控组件
package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// requestRecordTotal 系统请求记录
var requestRecordTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "request_record_total",
	Help: "慢查询SQL记录",
}, []string{
	"url",    // 请求路径
	"costs",  // 耗时（ms）
	"source", // 请求来源
	"code",   // 返回code
	"_name",  // 天机阁默认要求的上报字段
	"_type",  // 天机阁默认要求的上报字段
}).MustCurryWith(prometheus.Labels{"_name": "系统请求日志", "_type": "counter"})

// ReportRequestRecord 上报系统请求记录
func ReportRequestRecord(url, costs, source, code string) {
	requestRecordTotal.WithLabelValues(url, costs, source, code).Inc()
}
