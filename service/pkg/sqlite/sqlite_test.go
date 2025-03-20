package sqlite

import (
	"testing"

	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetMySQLClient(t *testing.T) {
	// 准备测试配置
	cfg := config.MySQLConfig{
		Dsn:                       "test_dsn",
		SkipDefaultTransaction:    true,
		SlowThreshold:             200,
		IgnoreRecordNotFoundError: true,
	}

	// 测试获取客户端
	client, err := GetMySQLClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestMySQLClient_CreateTables(t *testing.T) {
	// 准备测试客户端
	cfg := config.MySQLConfig{
		Dsn: "test_create_tables",
	}
	client, err := GetMySQLClient(cfg)
	assert.NoError(t, err)

	// 测试创建表
	mysqlClient := &MySQLClient{Client: client}
	err = mysqlClient.CreateTables()
	assert.NoError(t, err)

	// 验证表是否创建成功
	var tables []string
	err = client.DB.Raw("SELECT name FROM sqlite_master WHERE type='table'").Pluck("name", &tables).Error
	assert.NoError(t, err)

	// 验证必要的表是否存在
	expectedTables := []string{
		"history_node_inst",
		"history_workflow_inst",
		"node_inst",
		"trigger",
		"workflow_def",
		"workflow_inst",
	}
	for _, table := range expectedTables {
		assert.Contains(t, tables, table)
	}
}

func TestMySQLClient_Close(t *testing.T) {
	// 准备测试客户端
	cfg := config.MySQLConfig{
		Dsn: "test_close",
	}
	client, err := GetMySQLClient(cfg)
	assert.NoError(t, err)

	mysqlClient := &MySQLClient{Client: client}
	err = mysqlClient.Close()
	assert.NoError(t, err)
}
