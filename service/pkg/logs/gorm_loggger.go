package logs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// ErrRecordNotFound not found error(查询结果为空时抛出)
var ErrRecordNotFound = fmt.Errorf("record not found")

// Config 日志配置
type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  logger.LogLevel
}

// NewGormLogger 构造方法
func NewGormLogger(config Config) *gormLogger {
	var (
		infoStr      = "%s [info] "
		warnStr      = "%s warn] "
		errStr       = "%s [error] "
		traceStr     = "%s [%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
	)

	return &gormLogger{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// gormLogger 日志结构体
type gormLogger struct {
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode 通过此方法设置日志级别，低于次level的日志信息不会打印
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		RecordSQLLog(SQLRecord{
			Level:  logger.Info,
			Format: l.infoStr + msg,
			Args:   append([]interface{}{utils.FileWithLineNum()}, data...),
		})
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		RecordSQLLog(SQLRecord{
			Level:  logger.Warn,
			Format: l.warnStr + msg,
			Args:   append([]interface{}{utils.FileWithLineNum()}, data...),
		})
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		RecordSQLLog(SQLRecord{
			Level:  logger.Error,
			Format: l.warnStr + msg,
			Args:   append([]interface{}{utils.FileWithLineNum()}, data...),
		})
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	// error 级别：SQL 执行出错时
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, ErrRecordNotFound) ||
		!l.IgnoreRecordNotFoundError):
		sql, rows := fc()

		RecordSQLLog(SQLRecord{
			Level:  logger.Error,
			Format: l.traceErrStr,
			Args:   append([]interface{}{utils.FileWithLineNum()}, err, float64(elapsed.Nanoseconds())/1e6, rows, sql),
		})
	// warn 级别：用于 SQL 执行时间超过预设的慢查询阈值时
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()

		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		RecordSQLLog(SQLRecord{
			Level:  logger.Warn,
			Format: l.traceWarnStr,
			Args: append([]interface{}{utils.FileWithLineNum()}, slowLog, float64(elapsed.Nanoseconds())/1e6, rows,
				sql),
		})
	case l.LogLevel == logger.Info:
		sql, rows := fc()

		RecordSQLLog(SQLRecord{
			Level:  logger.Info,
			Format: l.traceStr,
			Args:   append([]interface{}{utils.FileWithLineNum()}, float64(elapsed.Nanoseconds())/1e6, rows, sql),
		})
	// 默认不打印任何东西
	default:
	}
}
