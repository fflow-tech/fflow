package storage

import (
	"github.com/fflow-tech/fflow/service/pkg/mysql"
)

// Transaction 执行事务
type Transaction interface {
	Transaction(fun func(*mysql.Client) error) error
}
