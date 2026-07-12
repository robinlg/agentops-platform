package chat

import (
	"context"
	"errors"
	"testing"

	"github.com/robinlg/agentops-platform/internal/llm"
	"github.com/robinlg/agentops-platform/internal/model"
	"github.com/robinlg/agentops-platform/internal/runtime"
	"github.com/robinlg/agentops-platform/internal/store"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// chatMocks 汇总一次测试所需的全部 mock，简化各用例的初始化。
type chatMocks struct {
	store             *store.MockIStore
	agentStore        *store.MockAgentStore
	modelProvider     *store.MockModelProviderStore
	conversationStore *store.MockConversationStore
	messageStore      *store.MockMessageStore
	agentRunStore     *store.MockAgentRunStore
	promptBuilder     *runtime.MockPromptBuilder
	llmClient         *llm.MockClient
}

func newChatMocks(ctrl *gomock.Controller) *chatMocks {
	return &chatMocks{
		store:             store.NewMockIStore(ctrl),
		agentStore:        store.NewMockAgentStore(ctrl),
		modelProvider:     store.NewMockModelProviderStore(ctrl),
		conversationStore: store.NewMockConversationStore(ctrl),
		messageStore:      store.NewMockMessageStore(ctrl),
		agentRunStore:     store.NewMockAgentRunStore(ctrl),
		promptBuilder:     runtime.NewMockPromptBuilder(ctrl),
		llmClient:         llm.NewMockClient(ctrl),
	}
}

func (m *chatMocks) biz() *chatBiz {
	return New(m.store, m.promptBuilder, m.llmClient)
}

// expectLoadChatContext 配置 loadChatContext 中三次 store 查询的成功返回。
func (m *chatMocks) expectLoadChatContext() {
	m.store.EXPECT().Agent().Return(m.agentStore)
	m.agentStore.EXPECT().Get(gomock.Any(), gomock.Any()).
		Return(&model.AgentM{ID: 10, ModelProviderID: 20}, nil)

	m.store.EXPECT().ModelProvider().Return(m.modelProvider)
	m.modelProvider.EXPECT().Get(gomock.Any(), gomock.Any()).
		Return(&model.ModelProviderM{ID: 20, DefaultModel: "deepseek-chat"}, nil)

	m.store.EXPECT().Conversation().Return(m.conversationStore)
	m.conversationStore.EXPECT().Get(gomock.Any(), gomock.Any()).
		Return(&model.ConversationM{ID: 100, AgentID: 10}, nil)
}

func TestChatBiz_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newChatMocks(ctrl)
	m.expectLoadChatContext()

	// 保存用户消息
	m.store.EXPECT().Message().Return(m.messageStore)
	m.messageStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	// 创建 agent run（running）
	m.store.EXPECT().AgentRun().Return(m.agentRunStore)
	m.agentRunStore.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ag *model.AgentRunM) error {
			ag.ID = 555
			return nil
		})

	// 查询历史消息 + 构建 prompt
	m.store.EXPECT().Message().Return(m.messageStore)
	m.messageStore.EXPECT().List(gomock.Any(), gomock.Any()).
		Return(int64(1), []*model.MessageM{{Role: model.MessageRoleUser, Content: "hi"}}, nil)
	m.promptBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).
		Return([]llm.Message{{Role: model.MessageRoleUser, Content: "hi"}})

	// 调用大模型
	m.llmClient.EXPECT().Chat(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&llm.ChatResult{
			Content: "Kubernetes Controller 是控制循环",
			Model:   "deepseek-chat",
			Usage:   llm.Usage{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30},
		}, nil)

	// 保存助手消息
	m.store.EXPECT().Message().Return(m.messageStore)
	m.messageStore.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, msg *model.MessageM) error {
			msg.ID = 999
			return nil
		})

	// 标记 agent run 成功
	m.store.EXPECT().AgentRun().Return(m.agentRunStore)
	m.agentRunStore.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := m.biz().Create(context.Background(), &apiv1.CreateChatRequest{
		AgentId:        10,
		ConversationId: ptrInt64(100),
		Message:        "解释一下 Controller",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(100), resp.ConversationId)
	assert.Equal(t, int64(999), resp.MessageId)
	assert.Equal(t, int64(555), resp.RunId)
	assert.Equal(t, "Kubernetes Controller 是控制循环", resp.Answer)
	require.NotNil(t, resp.Usage)
	assert.Equal(t, int32(30), resp.Usage.TotalTokens)
}

func TestChatBiz_Create_LLMError_MarksRunFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newChatMocks(ctrl)
	m.expectLoadChatContext()

	m.store.EXPECT().Message().Return(m.messageStore)
	m.messageStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	m.store.EXPECT().AgentRun().Return(m.agentRunStore)
	m.agentRunStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	m.store.EXPECT().Message().Return(m.messageStore)
	m.messageStore.EXPECT().List(gomock.Any(), gomock.Any()).
		Return(int64(0), []*model.MessageM{}, nil)
	m.promptBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return([]llm.Message{})

	// 大模型调用失败
	m.llmClient.EXPECT().Chat(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("llm timeout"))

	// 失败后应将 agent run 标记为 failed（一次额外的 Update）
	m.store.EXPECT().AgentRun().Return(m.agentRunStore)
	m.agentRunStore.EXPECT().Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ag *model.AgentRunM) error {
			assert.Equal(t, model.AgentRunStatusFailed, ag.Status)
			return nil
		})

	_, err := m.biz().Create(context.Background(), &apiv1.CreateChatRequest{
		AgentId:        10,
		ConversationId: ptrInt64(100),
		Message:        "hi",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "llm timeout")
}

func TestChatBiz_Create_AgentNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newChatMocks(ctrl)

	// 加载智能体即失败，后续流程不应发生
	m.store.EXPECT().Agent().Return(m.agentStore)
	m.agentStore.EXPECT().Get(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("agent not found"))

	_, err := m.biz().Create(context.Background(), &apiv1.CreateChatRequest{AgentId: 999})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent not found")
}

func ptrInt64(v int64) *int64 { return &v }
