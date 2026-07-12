package llm

import "github.com/robinlg/agentops-platform/internal/model"

// Message 是发送给大模型的单条对话消息，遵循 OpenAI Chat Completion API 协议。
// 该类型仅用于与 LLM 交互，不会持久化到数据库；数据库消息请使用 model.MessageM。
type Message struct {
	// Role 表示消息角色，取值为 model.MessageRoleSystem/User/Assistant/Tool 之一
	Role model.MessageRole `json:"role"`
	// Content 表示消息文本内容
	Content string `json:"content"`
}

// Usage 表示一次对话的 token 使用量统计
type Usage struct {
	// PromptTokens 输入 token 数
	PromptTokens int32 `json:"prompt_tokens"`
	// CompletionTokens 输出 token 数
	CompletionTokens int32 `json:"completion_tokens"`
	// TotalTokens 总 token 数
	TotalTokens int32 `json:"total_tokens"`
}

// ChatResult 表示一次对话调用的返回结果
type ChatResult struct {
	// Content 大模型返回的回复内容
	Content string `json:"content"`
	// Model 实际使用的模型名称
	Model string `json:"model"`
	// Usage token 使用量统计
	Usage Usage `json:"usage"`
}