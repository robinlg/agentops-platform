package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robinlg/onexlib/pkg/core"
)

// CreateModelProvider 创建模型提供商
func (h *Handler) CreateModelProvider(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.ModelProviderV1().Create)
}

// UpdateModelProvider 更新模型提供商
func (h *Handler) UpdateModelProvider(c *gin.Context) {
	core.HandleUriJSONRequest(c, h.biz.ModelProviderV1().Update)
}

// DeleteModelProvider 删除模型提供商
func (h *Handler) DeleteModelProvider(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.ModelProviderV1().Delete)
}

// GetModelProvider 获取模型提供商
func (h *Handler) GetModelProvider(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.ModelProviderV1().Get)
}

// ListModelProviders 列出模型提供商
func (h *Handler) ListModelProviders(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.ModelProviderV1().List)
}
