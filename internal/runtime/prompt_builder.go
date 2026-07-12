package runtime

//go:generate mockgen -destination mock_prompt_builder.go -package runtime github.com/robinlg/agentops-platform/internal/runtime PromptBuilder

import (
	"github.com/robinlg/agentops-platform/internal/llm"
	"github.com/robinlg/agentops-platform/internal/model"
)

// PromptBuilder 定义 prompt 构建接口
type PromptBuilder interface {
	Build(agent *model.AgentM, messages []*model.MessageM) []llm.Message
}

// promptBuilder 是 PromptBuilder 接口的默认实现
type promptBuilder struct{}

// NewPromptBuilder 创建 PromptBuilder 实例
func NewPromptBuilder() PromptBuilder {
	return &promptBuilder{}
}

func (p *promptBuilder) Build(agent *model.AgentM, messages []*model.MessageM) []llm.Message {
	// SystemPrompt 为空时不注入 system 消息
	result := make([]llm.Message, 0, len(messages)+1)
	if agent.SystemPrompt != nil && len(*agent.SystemPrompt) > 0 {
		result = append(result, llm.Message{
			Role:    model.MessageRoleSystem,
			Content: *agent.SystemPrompt,
		})
	}

	for _, msg := range messages {
		result = append(result, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return result
}
