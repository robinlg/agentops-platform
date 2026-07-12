package model

// MessageRole 定义了消息角色的类型，用于 MessageM.Role 字段
type MessageRole = string

// 消息角色枚举值，遵循 OpenAI Chat Completion API 的角色定义
const (
	// MessageRoleSystem 表示系统消息（system prompt），通常用于设定对话背景与规则
	MessageRoleSystem MessageRole = "system"
	// MessageRoleUser 表示由用户发送的消息
	MessageRoleUser MessageRole = "user"
	// MessageRoleAssistant 表示由模型（智能体）生成的回复消息
	MessageRoleAssistant MessageRole = "assistant"
	// MessageRoleTool 表示由工具/函数调用返回的消息
	MessageRoleTool MessageRole = "tool"
)
