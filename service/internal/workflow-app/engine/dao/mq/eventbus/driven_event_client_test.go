package eventbus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/dao/mq/tdmq"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/log"

	"github.com/apache/pulsar-client-go/pulsar"

	pkgtdmq "github.com/fflow-tech/fflow/service/pkg/mq/tdmq"
)

var (
	driveClient, _ = pkgtdmq.GetTDMQClient(config.GetTDMQConfig())
)

// TestDriveEventClient 测试事件总线 驱动事件Client 定时事件和延时事件如果是本地测试会有分钟级别的延时，本地时间和TDMQ上面的时间会有差异
func TestDriveEventClient(t *testing.T) {
	type args struct {
		handle func(pulsar.Message) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Normal Case", args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originMQClient := tdmq.NewClient(driveClient)
			mqClient := NewDriverEventClient(originMQClient)
			consumer, err := mqClient.NewConsumer(context.Background(),
				config.GetEventConfig().DriveEventGroup,
				func(ctx context.Context, message interface{}) error {
					msg, ok := (message).(pulsar.Message)
					if !ok {
						return fmt.Errorf("not pulsar.Message")
					}
					eventType, err := event.GetEventType(msg.Payload())
					if err != nil {
						return err
					}
					log.Infof("Consume message=%s type:%s", string(msg.Payload()), eventType)
					return nil
				})
			if err != nil {
				return
			}
			defer consumer.Close()
			// 因为延时消息和定时消息 需要消费者已在线的情况才能使用 所以这里需要等待一下
			time.Sleep(5 * time.Second)
		})
	}
}
