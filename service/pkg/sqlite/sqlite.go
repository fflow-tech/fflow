// Package memory 提供内存版本的外部依赖实现，用于本地测试
package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/storage/po"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"github.com/fflow-tech/fflow/service/pkg/mysql"

	// 使用纯 Go 实现的 SQLite

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var (
	clientMap sync.Map
	mutex     sync.Mutex
)

// NewClient 新建内存客户端
func NewClient(db *gorm.DB) *mysql.Client {
	return mysql.NewClient(db)
}

// GetMySQLClient 获取一个内存数据库客户端
func GetMySQLClient(config config.MySQLConfig) (*mysql.Client, error) {
	// 使用DSN作为缓存键，确保相同配置复用同一个客户端
	if client, ok := clientMap.Load(config.Dsn); ok {
		return client.(*mysql.Client), nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	// 二次检查
	if client, ok := clientMap.Load(config.Dsn); ok {
		return client.(*mysql.Client), nil
	}

	// 创建 SQLite 共享内存数据库作为 MySQL 替代
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_pragma=foreign_keys(1)&mode=memory&_txlock=immediate"), &gorm.Config{
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		Logger: logs.NewGormLogger(logs.Config{
			SlowThreshold:             time.Duration(config.SlowThreshold) * time.Millisecond,
			IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
			LogLevel:                  glogger.Warn,
		}),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect in-memory database: %w", err)
	}

	// 启用WAL模式以提高性能
	db.Exec("PRAGMA journal_mode=WAL")

	// 启用外键约束
	db.Exec("PRAGMA foreign_keys=ON")

	client := &mysql.Client{DB: db, Config: config}
	clientMap.Store(config.Dsn, client)

	return mysql.NewClient(client.DB), nil
}

// MySQLClient 是一个包装了 mysql.Client 的本地类型
type MySQLClient struct {
	*mysql.Client
}

// CreateTables 创建必要的表，可根据服务中的表结构调用
func (c *MySQLClient) CreateTables() error {
	// 使用 GORM 的 AutoMigrate 方法创建表
	err := c.DB.AutoMigrate(
		&po.HistoryNodeInstPO{},
		&po.HistoryWorkflowInstPO{},
		&po.NodeInstPO{},
		&po.TriggerPO{},
		&po.WorkflowDefPO{},
		&po.WorkflowInstPO{},
	)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	// 查询已经创建成功的表
	tables := []string{"history_node_inst", "history_workflow_inst", "node_inst", "trigger", "workflow_def", "workflow_inst"}
	for _, table := range tables {
		if !c.DB.Migrator().HasTable(table) {
			return fmt.Errorf("table %s not created successfully", table)
		}
	}
	log.Infof("Successfully created tables: %v", tables)
	return nil
}

// InitTestData 初始化测试数据
func (c *MySQLClient) InitTestData() error {
	return nil
}

// Close 关闭数据库连接
func (c *MySQLClient) Close() error {
	// 获取原始的SQL连接
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
