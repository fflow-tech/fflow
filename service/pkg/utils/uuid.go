package utils

import (
	"github.com/google/uuid"
)

// GetUUID 获取 uuid
func GetUUID() string {
	id := uuid.New()
	return id.String()
}
