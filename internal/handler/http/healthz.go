package http

import (
	"github.com/gin-gonic/gin"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(c *gin.Context) {
	// 返回 JSON 响应
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
