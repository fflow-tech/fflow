package nodeexecutor

import (
	"context"
	"fmt"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/ports"
	"github.com/fflow-tech/fflow/service/pkg/utils"
	"github.com/sashabaranov/go-openai"
)

// ServiceOpenAINodeExecutor implements a node executor for OpenAI API calls
type ServiceOpenAINodeExecutor struct {
	remoteRepo ports.RemoteRepository
}

// NewServiceOpenAINodeExecutor creates a new instance of ServiceOpenAINodeExecutor
func NewServiceOpenAINodeExecutor(remoteRepo ports.RemoteRepository) *ServiceOpenAINodeExecutor {
	return &ServiceOpenAINodeExecutor{remoteRepo: remoteRepo}
}

// Execute 执行节点
func (d *ServiceOpenAINodeExecutor) Execute(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	args := originArgs.(*entity.OpenAIArgs)
	nodeInst.Input = map[string]interface{}{"prompt": args.Prompt}

	rsp, err := d.call(ctx, args)
	if err != nil {
		return err
	}

	nodeInst.Output = rsp
	return nil
}

// Polling 轮询节点
func (d *ServiceOpenAINodeExecutor) Polling(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	return fmt.Errorf("not implemented")
}

// Cancel 取消执行节点
func (d *ServiceOpenAINodeExecutor) Cancel(ctx context.Context,
	nodeInst *entity.NodeInst, originArgs interface{}) error {
	return fmt.Errorf("not implemented")
}

// call 组装 request 并发送 OpenAI API 请求
func (d *ServiceOpenAINodeExecutor) call(ctx context.Context, args *entity.OpenAIArgs) (map[string]interface{}, error) {
	if args.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API Key is required")
	}

	config := openai.DefaultConfig(args.APIKey)
	config.BaseURL = args.BaseURL

	// 为每个请求创建新的 client
	client := openai.NewClientWithConfig(config)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// 转换消息格式
	messages := make([]openai.ChatCompletionMessage, 0, len(args.Messages))
	for _, msg := range args.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg["role"],
			Content: msg["content"],
		})
	}

	// 如果消息列表为空，则添加用户消息
	if len(args.Messages) == 0 {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "user",
			Content: args.Prompt,
		})
	}

	// 构建请求
	req := openai.ChatCompletionRequest{
		Model:       args.Model,
		Messages:    messages,
		Temperature: float32(args.Temperature),
		MaxTokens:   args.MaxTokens,
		Stream:      args.Stream,
	}

	// 发送请求
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("failed to get response from OpenAI API")
	}

	return utils.StructToMap(resp.Choices[0].Message)
}
