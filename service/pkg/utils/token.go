package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateToken 生成token
func GenerateToken() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
