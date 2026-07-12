package llm

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robinlg/agentops-platform/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAICompatibleClient_Chat_Success(t *testing.T) {
	var (
		gotPath   string
		gotAuth   string
		gotMethod string
		gotBody   chatCompletionRequest
	)

	// 启动一个模拟大模型服务
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotMethod = r.Method

		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"model": "deepseek-chat",
			"choices": [
				{"message": {"role": "assistant", "content": "Kubernetes Controller 是控制循环"}, "finish_reason": "stop"}
			],
			"usage": {"prompt_tokens": 10, "completion_tokens": 20, "total_tokens": 30}
		}`))
	}))
	defer srv.Close()

	client := newOpenAICompatibleClient()
	provider := &model.ModelProviderM{
		BaseURL:      srv.URL,
		APIKey:       "sk-test-key",
		DefaultModel: "deepseek-chat",
	}
	messages := []Message{
		{Role: model.MessageRoleSystem, Content: "你是助手"},
		{Role: model.MessageRoleUser, Content: "解释一下 Controller"},
	}

	result, err := client.Chat(context.Background(), provider, messages)
	require.NoError(t, err)

	// 校验请求格式
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/chat/completions", gotPath)
	assert.Equal(t, "Bearer sk-test-key", gotAuth)
	assert.Equal(t, "deepseek-chat", gotBody.Model)
	assert.Equal(t, messages, gotBody.Messages)

	// 校验响应解析
	assert.Equal(t, "Kubernetes Controller 是控制循环", result.Content)
	assert.Equal(t, "deepseek-chat", result.Model)
	assert.Equal(t, int32(10), result.Usage.PromptTokens)
	assert.Equal(t, int32(20), result.Usage.CompletionTokens)
	assert.Equal(t, int32(30), result.Usage.TotalTokens)
}

func TestOpenAICompatibleClient_Chat_TrimsTrailingSlash(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"model":"m","choices":[{"message":{"content":"ok"}}],"usage":{}}`))
	}))
	defer srv.Close()

	client := newOpenAICompatibleClient()
	// BaseURL 末尾带斜杠，应被正确处理，不出现双斜杠
	provider := &model.ModelProviderM{BaseURL: srv.URL + "/", APIKey: "k", DefaultModel: "m"}

	_, err := client.Chat(context.Background(), provider, []Message{{Role: model.MessageRoleUser, Content: "hi"}})
	require.NoError(t, err)
	assert.Equal(t, "/chat/completions", gotPath)
}

func TestOpenAICompatibleClient_Chat_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": {"message": "invalid api key", "type": "auth_error"}}`))
	}))
	defer srv.Close()

	client := newOpenAICompatibleClient()
	provider := &model.ModelProviderM{BaseURL: srv.URL, APIKey: "bad", DefaultModel: "m"}

	_, err := client.Chat(context.Background(), provider, []Message{{Role: model.MessageRoleUser, Content: "hi"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid api key")
}

func TestOpenAICompatibleClient_Chat_EmptyChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"model": "m", "choices": [], "usage": {}}`))
	}))
	defer srv.Close()

	client := newOpenAICompatibleClient()
	provider := &model.ModelProviderM{BaseURL: srv.URL, APIKey: "k", DefaultModel: "m"}

	_, err := client.Chat(context.Background(), provider, []Message{{Role: model.MessageRoleUser, Content: "hi"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no choices")
}

func TestOpenAICompatibleClient_Chat_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not a json`))
	}))
	defer srv.Close()

	client := newOpenAICompatibleClient()
	provider := &model.ModelProviderM{BaseURL: srv.URL, APIKey: "k", DefaultModel: "m"}

	_, err := client.Chat(context.Background(), provider, []Message{{Role: model.MessageRoleUser, Content: "hi"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}
