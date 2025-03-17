package command

import (
	"github.com/billcobbler/casbin-redis-watcher/v2"
	"github.com/casbin/casbin/v2"
	ga "github.com/casbin/gorm-adapter/v3"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
)

// RbacClient 客户端
type RbacClient struct {
	*casbin.Enforcer
}

// NewRbacClient 获取一个客户端
func NewRbacClient(mySQLClient *mysql.Client, redisConfig config.RedisConfig) *RbacClient {
	a, err := ga.NewAdapterByDBUseTableName(mySQLClient.DB, "fflow", "rbac_policy")
	if err != nil {
		log.Errorf("Failed to new casbin client, caused by %s", err.Error())
		panic(err)
	}
	e, err := casbin.NewEnforcer("rbac_model.conf", a)
	if err != nil {
		log.Errorf("Failed to new casbin client, caused by %s", err.Error())
		panic(err)
	}

	w, err := rediswatcher.NewWatcher(redisConfig.Address, rediswatcher.Password(redisConfig.Password))
	if err != nil {
		log.Errorf("Failed to new casbin client, caused by %s", err)
		panic(err)
	}
	if err := w.SetUpdateCallback(func(s string) {
		// 当有权限更新时，会执行这个回调，具体更新内存的权限逻辑可以自定义
		if err := e.LoadPolicy(); err != nil {
			log.Errorf("Failed to load policy, caused by %s", err)
		}
	}); err != nil {
		log.Errorf("Failed to new casbin client, caused by %s", err)
		panic(err)
	}
	if err := e.SetWatcher(w); err != nil {
		log.Errorf("Failed to new casbin client, caused by %s", err)
		panic(err)
	}

	client := &RbacClient{e}
	if err := e.SetWatcher(w); err != nil {
		panic(err)
	}

	return client
}
