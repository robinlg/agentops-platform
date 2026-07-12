package runtime

import (
	"testing"

	"github.com/robinlg/agentops-platform/internal/model"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string { return &s }

func TestPromptBuilder_Build(t *testing.T) {
	builder := NewPromptBuilder()

	tests := []struct {
		name         string
		agent        *model.AgentM
		messages     []*model.MessageM
		wantRoles    []model.MessageRole
		wantContents []string
	}{
		{
			name:  "有 SystemPrompt 时应在最前面注入 system 消息",
			agent: &model.AgentM{SystemPrompt: strPtr("你是一个助手")},
			messages: []*model.MessageM{
				{Role: model.MessageRoleUser, Content: "你好"},
				{Role: model.MessageRoleAssistant, Content: "你好，有什么可以帮你"},
			},
			wantRoles:    []model.MessageRole{model.MessageRoleSystem, model.MessageRoleUser, model.MessageRoleAssistant},
			wantContents: []string{"你是一个助手", "你好", "你好，有什么可以帮你"},
		},
		{
			name:  "SystemPrompt 为 nil 时不注入 system 消息",
			agent: &model.AgentM{SystemPrompt: nil},
			messages: []*model.MessageM{
				{Role: model.MessageRoleUser, Content: "在吗"},
			},
			wantRoles:    []model.MessageRole{model.MessageRoleUser},
			wantContents: []string{"在吗"},
		},
		{
			name:  "SystemPrompt 为空字符串时不注入 system 消息",
			agent: &model.AgentM{SystemPrompt: strPtr("")},
			messages: []*model.MessageM{
				{Role: model.MessageRoleUser, Content: "测试"},
			},
			wantRoles:    []model.MessageRole{model.MessageRoleUser},
			wantContents: []string{"测试"},
		},
		{
			name:         "没有历史消息且无 SystemPrompt 时返回空列表",
			agent:        &model.AgentM{SystemPrompt: nil},
			messages:     nil,
			wantRoles:    []model.MessageRole{},
			wantContents: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := builder.Build(tt.agent, tt.messages)

			assert.Len(t, got, len(tt.wantRoles))
			for i := range got {
				assert.Equal(t, tt.wantRoles[i], got[i].Role, "role at index %d", i)
				assert.Equal(t, tt.wantContents[i], got[i].Content, "content at index %d", i)
			}
		})
	}
}

func TestPromptBuilder_Build_PreservesOrder(t *testing.T) {
	builder := NewPromptBuilder()
	agent := &model.AgentM{SystemPrompt: strPtr("system")}
	messages := []*model.MessageM{
		{Role: model.MessageRoleUser, Content: "1"},
		{Role: model.MessageRoleAssistant, Content: "2"},
		{Role: model.MessageRoleUser, Content: "3"},
	}

	got := builder.Build(agent, messages)

	// 期望：system 在最前，后续消息保持原始顺序
	assert.Equal(t, []string{"system", "1", "2", "3"}, []string{
		got[0].Content, got[1].Content, got[2].Content, got[3].Content,
	})
}
