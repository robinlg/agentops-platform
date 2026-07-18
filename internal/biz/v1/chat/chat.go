package chat

import (
	"context"
	"time"

	"github.com/robinlg/agentops-platform/internal/llm"
	"github.com/robinlg/agentops-platform/internal/model"
	"github.com/robinlg/agentops-platform/internal/pkg/log"
	"github.com/robinlg/agentops-platform/internal/runtime"
	"github.com/robinlg/agentops-platform/internal/store"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/robinlg/onexlib/pkg/store/where"
	"gorm.io/gorm/clause"
)

const (
	// 查询会话最近的历史消息的数量
	historyMessageNumber = 50
)

// ChatBiz 定义处理智能体请求所需的方法
type ChatBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateChatRequest) (*apiv1.CreateChatResponse, error)

	ChatExpansion
}

// ChatExpansion 定义智能体的扩展方法
type ChatExpansion interface {
}

// chatBiz 是 ChatBiz 接口的实现.
type chatBiz struct {
	store         store.IStore
	promptBuilder runtime.PromptBuilder
	llmClient     llm.Client
}

// 确保 chatBiz 实现了 ChatBiz 接口
var _ ChatBiz = (*chatBiz)(nil)

// New 创建 chatBiz 实例.
func New(store store.IStore, promptBuilder runtime.PromptBuilder, llmClient llm.Client) *chatBiz {
	return &chatBiz{store: store, promptBuilder: promptBuilder, llmClient: llmClient}
}

// Create 处理智能体对话请求：加载智能体、模型提供商及会话信息，随后调用模型完成一次对话.
func (b *chatBiz) Create(ctx context.Context, rq *apiv1.CreateChatRequest) (*apiv1.CreateChatResponse, error) {
	// 1. 加载智能体、模型提供商及会话信息
	cc, err := b.loadChatContext(ctx, rq)
	if err != nil {
		return nil, err
	}

	// 2. 保存用户消息
	userMsg := model.MessageM{
		ConversationID: cc.conversation.ID,
		Role:           model.MessageRoleUser,
		Content:        rq.GetMessage(),
	}
	if err = b.store.Message().Create(ctx, &userMsg); err != nil {
		log.Errorw("Failed to create user message", "conversation_id", cc.conversation.ID, "err", err)
		return nil, err
	}

	// 3. 创建智能体运行记录（状态为 running）
	input := rq.GetMessage()
	ag := model.AgentRunM{
		AgentID:        cc.agent.ID,
		ConversationID: cc.conversation.ID,
		Status:         model.AgentRunStatusRunning,
		Input:          &input,
		Model:          cc.agent.Model,
	}
	if err = b.store.AgentRun().Create(ctx, &ag); err != nil {
		log.Errorw("Failed to create agent run", "agent_id", cc.agent.ID, "conversation_id", cc.conversation.ID, "err", err)
		return nil, err
	}

	// 4. 查询历史消息并构建 prompt
	llmMessages, err := b.buildLLMMessages(ctx, cc.agent, cc.conversation.ID)
	if err != nil {
		return nil, err
	}

	// 5. 调用大模型完成一次对话
	startedAt := time.Now()
	result, err := b.llmClient.Chat(ctx, cc.modelProvider, llmMessages)
	if err != nil {
		log.Errorw("Failed to chat with LLM", "agent_run_id", ag.ID, "model_provider_id", cc.modelProvider.ID, "err", err)
		// 调用失败：将智能体运行记录标记为 failed
		b.markAgentRunFailed(ctx, &ag)
		return nil, err
	}

	// 6. 保存大模型响应消息
	assistantMsg := model.MessageM{
		ConversationID: cc.conversation.ID,
		Role:           model.MessageRoleAssistant,
		Content:        result.Content,
	}
	if err = b.store.Message().Create(ctx, &assistantMsg); err != nil {
		log.Errorw("Failed to create assistant message", "conversation_id", cc.conversation.ID, "err", err)
		return nil, err
	}

	// 7. 标记智能体运行成功：回填输出内容、实际使用的模型、token 使用量、耗时及结束时间
	latencyMs, err := b.markAgentRunSuccess(ctx, &ag, result, startedAt)
	if err != nil {
		return nil, err
	}

	return &apiv1.CreateChatResponse{
		ConversationId: cc.conversation.ID,
		MessageId:      assistantMsg.ID,
		RunId:          ag.ID,
		Answer:         result.Content,
		Usage: &apiv1.Usage{
			PromptTokens:     result.Usage.PromptTokens,
			CompletionTokens: result.Usage.CompletionTokens,
			TotalTokens:      result.Usage.TotalTokens,
		},
		LatencyMs: latencyMs,
	}, nil
}

// chatContext 汇总一次对话所需的上下文数据：智能体、模型提供商、会话.
type chatContext struct {
	agent         *model.AgentM
	modelProvider *model.ModelProviderM
	conversation  *model.ConversationM
}

// loadChatContext 根据请求加载对话所需的智能体、模型提供商与会话信息，
// 并校验会话确实属于该智能体。
func (b *chatBiz) loadChatContext(ctx context.Context, rq *apiv1.CreateChatRequest) (*chatContext, error) {
	// 根据请求中的智能体 ID 查询智能体信息
	agentM, err := b.store.Agent().Get(ctx, where.F("id", rq.GetAgentId()))
	if err != nil {
		log.Errorw("Failed to get agent", "agent_id", rq.GetAgentId(), "err", err)
		return nil, err
	}

	// 根据智能体关联的模型提供商 ID，查询对应的模型提供商配置（BaseURL、APIKey 等）
	modelProviderM, err := b.store.ModelProvider().Get(ctx, where.F("id", agentM.ModelProviderID))
	if err != nil {
		log.Errorw("Failed to get model provider", "model_provider_id", agentM.ModelProviderID, "err", err)
		return nil, err
	}

	// 根据智能体 ID 和会话 ID 查询会话信息，确保该会话属于当前智能体；
	// 若请求未指定 conversation_id，则为该智能体新建一个会话。
	conversationM, err := b.getOrCreateConversation(ctx, agentM, rq)
	if err != nil {
		return nil, err
	}

	return &chatContext{agent: agentM, modelProvider: modelProviderM, conversation: conversationM}, nil
}

// getOrCreateConversation 根据请求获取或新建会话：
//   - 若请求携带 conversation_id，则按 (agent_id, id) 查询已存在的会话；
//   - 若未携带，则以该智能体新建一个会话（Title 使用用户首条消息的前若干字符）。
func (b *chatBiz) getOrCreateConversation(ctx context.Context, agentM *model.AgentM, rq *apiv1.CreateChatRequest) (*model.ConversationM, error) {
	// 请求携带 conversation_id：按 (agent_id, id) 查询，确保会话属于该智能体
	if rq.ConversationId != nil {
		conversationM, err := b.store.Conversation().Get(ctx, where.F("agent_id", agentM.ID, "id", rq.GetConversationId()))
		if err != nil {
			log.Errorw("Failed to get conversation", "agent_id", agentM.ID, "conversation_id", rq.GetConversationId(), "err", err)
			return nil, err
		}
		return conversationM, nil
	}

	// 未携带 conversation_id：为该智能体新建一个会话，Title 取用户首条消息的前若干字符
	title := buildConversationTitle(rq.GetMessage())
	conversationM := &model.ConversationM{
		AgentID: agentM.ID,
		Title:   &title,
	}
	if err := b.store.Conversation().Create(ctx, conversationM); err != nil {
		log.Errorw("Failed to create conversation", "agent_id", agentM.ID, "err", err)
		return nil, err
	}
	return conversationM, nil
}

// buildLLMMessages 拉取指定会话最近的历史消息，并借助 PromptBuilder 构建
// 用于调用大模型的 prompt 消息列表。
func (b *chatBiz) buildLLMMessages(ctx context.Context, agentM *model.AgentM, conversationID int64) ([]llm.Message, error) {
	// 查询该会话最近的历史消息，用于构建上下文
	_, messages, err := b.store.Message().List(ctx,
		where.F("conversation_id", conversationID).
			C(clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Name: "id"}, Desc: true}}}).
			P(1, historyMessageNumber),
	)
	if err != nil {
		log.Errorw("Failed to list history messages", "conversation_id", conversationID, "err", err)
		return nil, err
	}

	// 构建 prompt 消息列表
	return b.promptBuilder.Build(agentM, messages), nil
}

// markAgentRunSuccess 将智能体运行记录标记为成功，并回填输出内容、实际使用的模型、
// token 使用量、耗时及结束时间。返回计算得到的耗时（毫秒）以供调用方复用。
func (b *chatBiz) markAgentRunSuccess(ctx context.Context, ag *model.AgentRunM, result *llm.ChatResult, startedAt time.Time) (int64, error) {
	finishedAt := time.Now()
	latencyMs := finishedAt.Sub(startedAt).Milliseconds()
	ag.Status = model.AgentRunStatusSuccess
	ag.Output = &result.Content
	ag.Model = &result.Model
	ag.PromptTokens = &result.Usage.PromptTokens
	ag.CompletionTokens = &result.Usage.CompletionTokens
	ag.TotalTokens = &result.Usage.TotalTokens
	ag.LatencyMs = &latencyMs
	ag.FinishedAt = &finishedAt
	if err := b.store.AgentRun().Update(ctx, ag); err != nil {
		log.Errorw("Failed to update agent run status to success", "agent_run_id", ag.ID, "err", err)
		return 0, err
	}
	return latencyMs, nil
}

// markAgentRunFailed 将智能体运行记录标记为失败。更新失败仅记录日志，不再向上抛出错误，
// 以免掩盖真正导致运行失败的原始错误。
func (b *chatBiz) markAgentRunFailed(ctx context.Context, ag *model.AgentRunM) {
	ag.Status = model.AgentRunStatusFailed
	if err := b.store.AgentRun().Update(ctx, ag); err != nil {
		log.Errorw("Failed to update agent run status to failed", "agent_run_id", ag.ID, "err", err)
	}
}

// conversationTitleMaxRunes 新建会话时自动生成标题的最大字符数（按 rune 计，避免多字节字符被截半）
const conversationTitleMaxRunes = 30

// buildConversationTitle 使用用户首条消息生成新建会话的默认标题：
// 截取前 conversationTitleMaxRunes 个字符，超出部分用 "..." 表示。
func buildConversationTitle(message string) string {
	runes := []rune(message)
	if len(runes) <= conversationTitleMaxRunes {
		return message
	}
	return string(runes[:conversationTitleMaxRunes]) + "..."
}
