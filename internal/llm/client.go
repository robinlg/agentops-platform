package llm

//go:generate mockgen -destination mock_client.go -package llm github.com/robinlg/agentops-platform/internal/llm Client

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/model"
)

// Client 定义了与大模型交互的统一接口。
// 不同的模型提供商（OpenAI、Anthropic、Azure 等）通过实现该接口来适配。
type Client interface {
	// Chat 发起一次对话补全请求
	Chat(ctx context.Context, provider *model.ModelProviderM, messages []Message) (*ChatResult, error)
}

// NewClient 根据模型提供商类型创建对应的 Client 实现。
// 目前统一使用 OpenAI 兼容协议实现；后续可按 provider.ProviderType 分支返回不同实现。
func NewClient() Client {
	return newOpenAICompatibleClient()
}
