package python

import (
	"bytes"
	"os/exec"
)

// Executor python 执行器
type Executor struct {
}

// NewExecutor 创建 python 执行器
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute 执行 python 脚本
func (e *Executor) Execute(args ...string) (string, error) {
	cmd := exec.Command("python", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
