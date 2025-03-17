// Package mysql 数据库的相关功能包
package mysql

import (
	"fmt"
	"sync"
	"time"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/logs"
	"gorm.io/plugin/dbresolver"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var (
	clientMap sync.Map
	mutex     sync.Mutex
)

// Client 客户端
type Client struct {
	*gorm.DB
	Config config.MySQLConfig
}

// NewClient 新建客户端
func NewClient(db *gorm.DB) *Client {
	return &Client{
		DB: db,
	}
}

// ReadFromSlave 从备节点读取数据
func (c *Client) ReadFromSlave(readFromSlave bool) *gorm.DB {
	if readFromSlave {
		return c.DB.Clauses(dbresolver.Read)
	}
	return c.DB.Clauses(dbresolver.Write)
}

// GetClient 获取一个数据库客户端
func GetClient(config config.MySQLConfig) (*Client, error) {
	if client, ok := clientMap.Load(config.Dsn); ok {
		return client.(*Client), nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		// 慢查询阈值为100ms, 日志级别为info即打印所有级别日志
		Logger: logs.NewGormLogger(logs.Config{
			SlowThreshold:             time.Duration(config.SlowThreshold) * time.Millisecond,
			IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
			LogLevel:                  glogger.Warn,
		})})

	client := &Client{DB: db, Config: config}
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
		return nil, err
	}

	clientMap.Store(config.Dsn, client)

	if config.SlaveDsn == "" {
		return client, nil
	}
	if err := db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(config.Dsn)},
		Replicas: []gorm.Dialector{mysql.Open(config.SlaveDsn)},
		Policy:   dbresolver.RandomPolicy{},
	})); err != nil {
		panic(fmt.Errorf("failed to register slave database, %w", err))
		return nil, err
	}

	return client, nil
}
