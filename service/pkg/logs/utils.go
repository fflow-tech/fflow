package logs

import (
	"fmt"
)

// GetFlowTraceID 生成流程的TRACE_ID
func GetFlowTraceID(defID, instID interface{}) string {
	switch defID.(type) {
	case uint64, int64, uint, int:
		return fmt.Sprintf("$$%d$$##%d##", defID, instID)
	default:
		return fmt.Sprintf("$$%s$$##%s##", defID, instID)
	}
}
