package nodeexecutor

import (
	"context"
	"fmt"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestServiceOpenAINodeExecutor_SuccessExecute(t *testing.T) {
	tests := []struct {
		name        string
		args        *entity.OpenAIArgs
		expectError bool
	}{
		{
			name: "成功执行 - 使用 Prompt",
			args: &entity.OpenAIArgs{
				BaseURL:     "https://openrouter.ai/api/v1",
				APIKey:      "your-api-key",
				Model:       "deepseek/deepseek-r1-distill-qwen-32b:free",
				Prompt:      "Hello, how are you?",
				Temperature: 0.7,
			},
			expectError: false,
		},
		{
			name: "成功执行 - 使用 Messages",
			args: &entity.OpenAIArgs{
				BaseURL: "https://openrouter.ai/api/v1",
				APIKey:  "your-api-key",
				Model:   "deepseek/deepseek-r1-distill-qwen-32b:free",
				Messages: []map[string]string{
					{
						"role":    "user",
						"content": "Hello, how are you?",
					},
				},
				Temperature: 0.7,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备测试环境
			executor := &ServiceOpenAINodeExecutor{}

			// 执行测试
			nodeInst := &entity.NodeInst{}
			err := executor.Execute(context.Background(), nodeInst, tt.args)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nodeInst.Input)
				if len(tt.args.Messages) > 0 {
					fmt.Println(nodeInst.Output)
					assert.Contains(t, nodeInst.Input, "prompt")
				}
			}
		})
	}
}

func TestServiceOpenAINodeExecutor_FailExecute(t *testing.T) {
	tests := []struct {
		name        string
		args        *entity.OpenAIArgs
		expectError bool
	}{
		{
			name: "成功执行 - 使用 Prompt",
			args: &entity.OpenAIArgs{
				APIKey:      "test-api-key",
				Model:       "gpt-3.5-turbo",
				Prompt:      "Hello",
				Temperature: 0.7,
			},
			expectError: false,
		},
		{
			name: "失败 - 缺少 API Key",
			args: &entity.OpenAIArgs{
				Model:  "gpt-3.5-turbo",
				Prompt: "Hello",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备测试环境
			executor := &ServiceOpenAINodeExecutor{}

			// 执行测试
			nodeInst := &entity.NodeInst{}
			err := executor.Execute(context.Background(), nodeInst, tt.args)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nodeInst.Input)
				if len(tt.args.Messages) > 0 {
					assert.Contains(t, nodeInst.Input, "prompt")
				}
			}
		})
	}
}

func TestServiceOpenAINodeExecutor_call(t *testing.T) {
	tests := []struct {
		name        string
		args        *entity.OpenAIArgs
		expectError bool
	}{
		{
			name: "成功调用 API",
			args: &entity.OpenAIArgs{
				APIKey: "test-api-key",
				Model:  "gpt-3.5-turbo",
				Messages: []map[string]string{
					{
						"role":    "user",
						"content": "Hello",
					},
				},
				Temperature: 0.7,
			},
			expectError: false,
		},
		{
			name: "失败 - 空 Messages 和 Prompt",
			args: &entity.OpenAIArgs{
				APIKey: "test-api-key",
				Model:  "gpt-3.5-turbo",
			},
			expectError: true,
		},
		{
			name: "失败 - 无效的 Model",
			args: &entity.OpenAIArgs{
				APIKey: "test-api-key",
				Model:  "",
				Prompt: "Hello",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &ServiceOpenAINodeExecutor{}
			result, err := executor.call(context.Background(), tt.args)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestServiceOpenAINodeExecutor_Polling(t *testing.T) {
	executor := &ServiceOpenAINodeExecutor{}
	err := executor.Polling(context.Background(), &entity.NodeInst{}, &entity.OpenAIArgs{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestServiceOpenAINodeExecutor_Cancel(t *testing.T) {
	executor := &ServiceOpenAINodeExecutor{}
	err := executor.Cancel(context.Background(), &entity.NodeInst{}, &entity.OpenAIArgs{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}
