package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robinlg/onexlib/pkg/core"
)

// ListAgentRuns 列出智能体运行记录
func (h *Handler) ListAgentRuns(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.AgentRunV1().List)
}

// GetAgentRun 获取单个智能体运行记录
func (h *Handler) GetAgentRun(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.AgentRunV1().Get)
}
