package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/robinlg/agentops-platform/internal/model"
)

// defaultHTTPTimeout 定义调用大模型的默认超时时间
const defaultHTTPTimeout = 60 * time.Second

// chatCompletionRequest 表示发送给 /chat/completions 接口的请求体，遵循 OpenAI 协议
type chatCompletionRequest struct {
	// Model 指定使用的模型名称
	Model string `json:"model"`
	// Messages 对话上下文消息列表
	Messages []Message `json:"messages"`
}

// chatCompletionResponse 表示 /chat/completions 接口的响应体，遵循 OpenAI 协议
type chatCompletionResponse struct {
	// Model 实际使用的模型名称
	Model string `json:"model"`
	// Choices 模型生成的候选回复列表
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	// Usage token 使用量统计
	Usage struct {
		PromptTokens     int32 `json:"prompt_tokens"`
		CompletionTokens int32 `json:"completion_tokens"`
		TotalTokens      int32 `json:"total_tokens"`
	} `json:"usage"`
	// Error 错误信息（请求失败时返回）
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// openAICompatibleClient 是基于 OpenAI Chat Completion API 协议的 Client 实现，
// 适用于所有兼容 OpenAI 接口的模型提供商（OpenAI、DeepSeek、Moonshot、通义千问等）。
type openAICompatibleClient struct {
	httpClient *http.Client
}

// 确保 openAICompatibleClient 实现了 Client 接口。
var _ Client = (*openAICompatibleClient)(nil)

// newOpenAICompatibleClient 创建 openAICompatibleClient 实例。
func newOpenAICompatibleClient() *openAICompatibleClient {
	return &openAICompatibleClient{
		httpClient: &http.Client{Timeout: defaultHTTPTimeout},
	}
}

// Chat 实现 Client 接口的 Chat 方法。
// 使用 provider.BaseURL + provider.APIKey 组装 HTTP 请求，调用 /chat/completions 接口，
// 解析响应并填充 ChatResult。
func (c *openAICompatibleClient) Chat(ctx context.Context, provider *model.ModelProviderM, messages []Message) (*ChatResult, error) {
	// 组装请求体
	reqBody := chatCompletionRequest{
		Model:    provider.DefaultModel,
		Messages: messages,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal chat request failed: %w", err)
	}

	// 拼接请求地址：{base_url}/chat/completions（去除 base_url 末尾多余的斜杠）
	url := strings.TrimRight(provider.BaseURL, "/") + "/chat/completions"

	// 构建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create chat request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)

	// 发起请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send chat request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read chat response failed: %w", err)
	}

	// 解析响应
	var completionResp chatCompletionResponse
	if err = json.Unmarshal(respBytes, &completionResp); err != nil {
		return nil, fmt.Errorf("unmarshal chat response failed: %w, body: %s", err, string(respBytes))
	}

	// 处理非 2xx 状态码或业务错误
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if completionResp.Error != nil {
			return nil, fmt.Errorf("chat request failed with status %d: %s", resp.StatusCode, completionResp.Error.Message)
		}
		return nil, fmt.Errorf("chat request failed with status %d, body: %s", resp.StatusCode, string(respBytes))
	}

	// 校验是否有有效回复
	if len(completionResp.Choices) == 0 {
		return nil, fmt.Errorf("chat response contains no choices, body: %s", string(respBytes))
	}

	// 填充返回结果
	return &ChatResult{
		Content: completionResp.Choices[0].Message.Content,
		Model:   completionResp.Model,
		Usage: Usage{
			PromptTokens:     completionResp.Usage.PromptTokens,
			CompletionTokens: completionResp.Usage.CompletionTokens,
			TotalTokens:      completionResp.Usage.TotalTokens,
		},
	}, nil
}
