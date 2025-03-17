package monitor

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// maxLabelBytes metrics 标签最大字节数
	maxLabelBytes = 512
	// cutLabelSuffix 被切的标签后缀
	cutLabelSuffix = "..."
)

// SlowSQLRecordTotal 慢查询SQL记录
var SlowSQLRecordTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "slow_sql_record_total",
	Help: "慢查询SQL记录",
}, []string{
	"costs", // 查询耗时(ms)
	"sql",   // sql语句
	"_name",
	"_type",
}).MustCurryWith(prometheus.Labels{"_name": "慢查询SQL记录", "_type": "counter"})

// ReportSlowSQLRecord 上报SQL慢查询记录 cost单位为ms
func ReportSlowSQLRecord(cost int, sql string) {
	if len(sql) > maxLabelBytes {
		sql = sql[:maxLabelBytes-len(cutLabelSuffix)] + cutLabelSuffix
	}
	SlowSQLRecordTotal.WithLabelValues(strconv.Itoa(cost), sql).Inc()
}
