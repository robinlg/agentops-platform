package conversation

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/pkg/conversion"
	"github.com/robinlg/agentops-platform/internal/store"
	apiv1 "github.com/robinlg/agentops-platform/pkg/api/v1"
	"github.com/robinlg/onexlib/pkg/store/where"
	"gorm.io/gorm/clause"
)

// ConversationBiz 定义处理会话请求所需的方法
type ConversationBiz interface {
	List(ctx context.Context, rq *apiv1.ListConversationRequest) (*apiv1.ListConversationResponse, error)
	ListMessages(ctx context.Context, rq *apiv1.ListMessageRequest) (*apiv1.ListMessageResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteConversationRequest) (*apiv1.DeleteConversationResponse, error)

	ConversationExpansion
}

// ConversationExpansion 定义会话的扩展方法
type ConversationExpansion interface {
}

// conversationBiz 是 ConversationBiz 接口的实现.
type conversationBiz struct {
	store store.IStore
}

// 确保 conversationBiz 实现了 ConversationBiz 接口
var _ ConversationBiz = (*conversationBiz)(nil)

func New(store store.IStore) *conversationBiz {
	return &conversationBiz{store: store}
}

// List 实现 ConversationBiz 接口中的 List 方法
func (b *conversationBiz) List(ctx context.Context, rq *apiv1.ListConversationRequest) (*apiv1.ListConversationResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))

	count, conversationList, err := b.store.Conversation().List(ctx, whr)
	if err != nil {
		return nil, err
	}
	conversations := make([]*apiv1.Conversation, 0, len(conversationList))
	for _, item := range conversationList {
		conversations = append(conversations, conversion.ConversationMToConversationV1(item))
	}

	return &apiv1.ListConversationResponse{TotalCount: count, Conversations: conversations}, nil
}

// ListMessages 实现 ConversationBiz 接口中的 ListMessages 方法
func (b *conversationBiz) ListMessages(ctx context.Context, rq *apiv1.ListMessageRequest) (*apiv1.ListMessageResponse, error) {
	whr := where.F("conversation_id", rq.GetConversationId()).
		C(clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Name: "id"}, Desc: false}}}).
		P(int(rq.GetOffset()), int(rq.GetLimit()))

	count, messageList, err := b.store.Message().List(ctx, whr)
	if err != nil {
		return nil, err
	}
	messages := make([]*apiv1.Message, 0, len(messageList))
	for _, item := range messageList {
		messages = append(messages, conversion.MessageMToMessageV1(item))
	}

	return &apiv1.ListMessageResponse{TotalCount: count, Messages: messages}, nil
}

// Delete 实现 ConversationBiz 接口中的 Delete 方法
func (b *conversationBiz) Delete(ctx context.Context, rq *apiv1.DeleteConversationRequest) (*apiv1.DeleteConversationResponse, error) {
	if err := b.store.Conversation().Delete(ctx, where.F("id", rq.GetId())); err != nil {
		return nil, err
	}

	return &apiv1.DeleteConversationResponse{}, nil
}
