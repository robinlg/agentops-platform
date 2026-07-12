package http

import (
	"github.com/robinlg/agentops-platform/internal/biz"
)

// Handler 处理请求
type Handler struct {
	biz biz.IBiz
}

// NewHandler 创建新的 Handler 实例
func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}
