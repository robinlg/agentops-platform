package conversion

import (
	"github.com/robinlg/agentops-platform/internal/model"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/robinlg/onexlib/pkg/core"
)

// ConversationMToConversationV1 将模型层的 ConversationM 转换为 API 层的 Conversation
func ConversationMToConversationV1(conversationModel *model.ConversationM) *apiv1.Conversation {
	var protoConversation apiv1.Conversation
	_ = core.CopyWithConverters(&protoConversation, conversationModel)
	return &protoConversation
}

// MessageMToMessageV1 将模型层的 MessageM 转换为 API 层的 Message
func MessageMToMessageV1(messageModel *model.MessageM) *apiv1.Message {
	var protoMessage apiv1.Message
	_ = core.CopyWithConverters(&protoMessage, messageModel)
	return &protoMessage
}
