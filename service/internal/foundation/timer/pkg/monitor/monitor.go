package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// reportLabel 上报标签.
type reportLabel string

// monitorComponent 监控指标类型.
type monitorComponent string

const (
	// 定时器触发记录总数.
	timerTriggerRecord        = "timer_trigger_record_total"
	timerTriggerRecordSummary = "触发记录总数"

	// 定时器触发耗时.
	timerTriggerCost        = "timer_trigger_cost"
	timerTriggerCostSummary = "定时器触发耗时"

	// 处于激活态的定时器总数.
	timerEnabledTotal        = "timer_enabled_total"
	timerEnabledTotalSummary = "激活态定时器总数"

	// 未触发定时器数量.
	timerFailedTotal        = "timer_failed_total"
	timerFailedTotalSummary = "未触发定时器数量"

	// 上报标签: 天机阁默认要求的上报字段.
	reportName reportLabel = "_name"
	reportType reportLabel = "_type"
	// 上报标签：定时器所属应用.
	timerApp reportLabel = "timer_app"

	// 通用标签.
	label = "label"
	timer = "timer"

	// 计数器.
	counter monitorComponent = "counter"
	// 摘要.
	summary monitorComponent = "summary"
	// 仪表盘.
	gauge monitorComponent = "gauge"
	// 直方图.
	histogram monitorComponent = "histogram"
)

// Reporter 监控上报服务.
type Reporter struct {
	triggerRecorder      *prometheus.CounterVec
	timeCostRecorder     prometheus.ObserverVec
	timerEnabledRecorder *prometheus.GaugeVec
	failedTimerRecorder  *prometheus.GaugeVec
}

var reporter = newReporter()

// GetReporter 获取单例上报服务.
func GetReporter() *Reporter {
	return reporter
}

// newReporter 监控上报服务构造器.
func newReporter() *Reporter {
	return &Reporter{
		// 定时器触发记录.
		triggerRecorder: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: timerTriggerRecord,
			Help: timerTriggerRecordSummary,
		}, []string{
			string(timerApp),
			string(reportName),
			string(reportType),
		}).MustCurryWith(prometheus.Labels{string(reportName): timerTriggerRecordSummary,
			string(reportType): string(counter)}),

		// 定时器延时记录.
		timeCostRecorder: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name:       timerTriggerCost,
			Help:       timerTriggerCostSummary,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001, 0.9999: 0.00001},
		}, []string{
			string(timerApp),
			string(reportName),
			string(reportType),
		}).MustCurryWith(prometheus.Labels{string(reportName): timerTriggerCostSummary,
			string(reportType): string(summary)}),

		// 处于激活态的定时器总数.
		timerEnabledRecorder: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: timerEnabledTotal,
			Help: timerEnabledTotalSummary,
		}, []string{
			label,
			string(reportName),
			string(reportType),
		}).MustCurryWith(prometheus.Labels{string(reportName): timerEnabledTotalSummary,
			string(reportType): string(gauge)}),

		// 未触发定时器数量.
		failedTimerRecorder: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: timerFailedTotal,
			Help: timerFailedTotalSummary,
		}, []string{
			label,
			string(reportName),
			string(reportType),
		}).MustCurryWith(prometheus.Labels{string(reportName): timerFailedTotalSummary,
			string(reportType): string(gauge)}),
	}
}

// ReportTriggerRecord 上报触发记录.
func (r *Reporter) ReportTriggerRecord(app string) {
	r.triggerRecorder.WithLabelValues(app).Inc()
}

// ReportTimerCostRecord 上报定时器耗时记录.
func (r *Reporter) ReportTimerCostRecord(app string, cost float64) {
	r.timeCostRecorder.WithLabelValues(app).Observe(cost)
}

// ReportTimerEnabledTotalNum 上报激活态定时器总数.
func (r *Reporter) ReportTimerEnabledTotalNum(total float64) {
	r.timerEnabledRecorder.WithLabelValues(timer).Set(total)
}

// ReportTimerFailedNum 上报未触发定时器数量.
func (r *Reporter) ReportTimerFailedNum(total float64) {
	r.failedTimerRecorder.WithLabelValues(timer).Set(total)
}
