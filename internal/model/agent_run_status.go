package model

// AgentRunStatus 定义了智能体运行状态的类型，用于 AgentRunM.Status 字段
type AgentRunStatus = string

// 智能体运行状态枚举值.
const (
	// AgentRunStatusPending 表示运行任务已创建但尚未开始执行
	AgentRunStatusPending AgentRunStatus = "pending"
	// AgentRunStatusRunning 表示运行任务正在执行中
	AgentRunStatusRunning AgentRunStatus = "running"
	// AgentRunStatusSuccess 表示运行任务已成功完成
	AgentRunStatusSuccess AgentRunStatus = "success"
	// AgentRunStatusFailed 表示运行任务执行失败
	AgentRunStatusFailed AgentRunStatus = "failed"
)
