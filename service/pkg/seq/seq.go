// Package seq 根据 ip 取模作为 nodeid 来生成分布式 ID，如果集群比较大且 ip 比较分散的情况有概率可能会重复，后续可以改成从 MySQL/Redis 获取
package seq

import (
	"fmt"
	"math/big"

	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	nodeID, err := ipToInt64()
	if err != nil {
		log.Errorf("ip to int64 failed, err: %v", err)
		return
	}

	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Errorf("snowflake new node failed, err: %v", err)
		return
	}
}

// NewUint64 生成 uint64 类型的分布式唯一 ID
func NewUint64() (uint64, error) {
	if node == nil {
		return 0, fmt.Errorf("init snowflake node failed")
	}
	return uint64(node.Generate().Int64()), nil
}

// NewUint 生成 uint 类型的分布式唯一 ID
func NewUint() (uint, error) {
	if node == nil {
		return 0, fmt.Errorf("init snowflake node failed")
	}
	return uint(node.Generate().Int64()), nil
}

// NewString 生成 string 类型的分布式唯一 ID
func NewString() (string, error) {
	if node == nil {
		return "", fmt.Errorf("init snowflake node failed")
	}
	return node.Generate().String(), nil
}

func ipToInt64() (int64, error) {
	ip, err := utils.GetOutboundIP()
	if err != nil {
		return 0, err
	}

	i := big.NewInt(0)
	i.SetBytes([]byte(ip))
	return i.Int64() % 1024, nil
}
