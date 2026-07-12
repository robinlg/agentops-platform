package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robinlg/onexlib/pkg/core"
)

// CreateAgent 创建智能体
func (h *Handler) CreateAgent(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.AgentV1().Create)
}

// UpdateAgent 更新智能体
func (h *Handler) UpdateAgent(c *gin.Context) {
	core.HandleUriJSONRequest(c, h.biz.AgentV1().Update)
}

// DeleteAgent 删除智能体
func (h *Handler) DeleteAgent(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.AgentV1().Delete)
}

// GetAgent 获取智能体
func (h *Handler) GetAgent(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.AgentV1().Get)
}

// ListAgents 列出智能体
func (h *Handler) ListAgents(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.AgentV1().List)
}

// CreateChat 创建对话
func (h *Handler) CreateChat(c *gin.Context) {
	core.HandleUriJSONRequest(c, h.biz.ChatV1().Create)
}
