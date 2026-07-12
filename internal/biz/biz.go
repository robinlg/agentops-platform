package biz

//go:generate mockgen -destination mock_biz.go -package biz github.com/robinlg/agentops-platform/internal/biz IBiz

import (
	"github.com/google/wire"
	agentv1 "github.com/robinlg/agentops-platform/internal/biz/v1/agent"
	agentrunv1 "github.com/robinlg/agentops-platform/internal/biz/v1/agent-run"
	chatv1 "github.com/robinlg/agentops-platform/internal/biz/v1/chat"
	conversationv1 "github.com/robinlg/agentops-platform/internal/biz/v1/conversation"
	modelproviderv1 "github.com/robinlg/agentops-platform/internal/biz/v1/model-provider"
	"github.com/robinlg/agentops-platform/internal/llm"
	"github.com/robinlg/agentops-platform/internal/runtime"
	"github.com/robinlg/agentops-platform/internal/store"
)

// ProviderSet 是一个 Wire 的 Provider 集合，用于声明依赖注入的规则.
// 包含 NewBiz 构造函数，用于生成 biz 实例.
// wire.Bind 用于将接口 IBiz 与具体实现 *biz 绑定，
// 这样依赖 IBiz 的地方会自动注入 *biz 实例.
var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

// IBiz 定义了业务层需要实现的方法
type IBiz interface {
	// ModelProviderV1 模型提供商业务接口
	ModelProviderV1() modelproviderv1.ModelProviderBiz
	// AgentV1 智能体业务接口
	AgentV1() agentv1.AgentBiz
	// ChatV1 对话业务接口
	ChatV1() chatv1.ChatBiz
	// ConversationV1 会话业务接口
	ConversationV1() conversationv1.ConversationBiz
	// AgentRunV1 智能体运行记录业务接口
	AgentRunV1() agentrunv1.AgentRunBiz
}

// biz 是 IBiz 的一个具体实现
type biz struct {
	store         store.IStore
	promptBuilder runtime.PromptBuilder
	llmClient     llm.Client
}

// 确保 biz 实现了 IBiz 接口
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例
func NewBiz(store store.IStore, promptBuilder runtime.PromptBuilder, llmClient llm.Client) *biz {
	return &biz{store: store, promptBuilder: promptBuilder, llmClient: llmClient}
}

// ModelProviderV1 返回一个实现了 ModelProviderBiz 接口的实例
func (b *biz) ModelProviderV1() modelproviderv1.ModelProviderBiz {
	return modelproviderv1.New(b.store)
}

// AgentV1 返回一个实现了 AgentBiz 接口的实例
func (b *biz) AgentV1() agentv1.AgentBiz {
	return agentv1.New(b.store)
}

// ChatV1 返回一个实现了 ChatBiz 接口的实例
func (b *biz) ChatV1() chatv1.ChatBiz {
	return chatv1.New(b.store, b.promptBuilder, b.llmClient)
}

// ConversationV1 返回一个实现了 ConversationBiz 接口的实例
func (b *biz) ConversationV1() conversationv1.ConversationBiz {
	return conversationv1.New(b.store)
}

// AgentRunV1 返回一个实现了 AgentRunBiz 接口的实例
func (b *biz) AgentRunV1() agentrunv1.AgentRunBiz {
	return agentrunv1.New(b.store)
}
