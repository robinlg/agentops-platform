package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robinlg/onexlib/pkg/core"
)

// ListConversations 列出会话
func (h *Handler) ListConversations(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.ConversationV1().List)
}

// ListConversationMessages 列出会话的消息
func (h *Handler) ListConversationMessages(c *gin.Context) {
	core.HandleUriQueryRequest(c, h.biz.ConversationV1().ListMessages)
}

// DeleteConversation 删除会话
func (h *Handler) DeleteConversation(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.ConversationV1().Delete)
}
