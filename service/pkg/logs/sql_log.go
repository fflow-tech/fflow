package logs

import (
	"github.com/fflow-tech/fflow/service/pkg/log"
	"gorm.io/gorm/logger"
)

// SQLRecord SQL 日志记录
type SQLRecord struct {
	Level  logger.LogLevel // 日志打印level
	Format string          // 格式化str
	Args   []interface{}   // format参数
}

// RecordSQLLog 记录 SQL 日志
func RecordSQLLog(detail SQLRecord) {
	switch detail.Level {
	case logger.Info:
		sqlLog().Infof("[SQLRecord] "+detail.Format, detail.Args...)
	case logger.Warn:
		sqlLog().Warnf("[SQLRecord] "+detail.Format, detail.Args...)
	case logger.Error:
		sqlLog().Errorf("[SQLRecord] "+detail.Format, detail.Args...)
	default:
		log.Warnf("SQL log  only support level [info、warn、Error]")
	}
}

// GetSQLLogName SQL 日志插件的名称
func GetSQLLogName() string {
	return "sql_log"
}

func sqlLog() log.Logger {
	return log.GetDefaultLogger()
}
