package store

//go:generate mockgen -destination mock_conversation.go -package store github.com/robinlg/agentops-platform/internal/store ConversationStore

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/model"
	genericstore "github.com/robinlg/onexlib/pkg/store"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// ConversationStore 定义了 conversations 模块在 store 层所实现的方法
type ConversationStore interface {
	Create(ctx context.Context, obj *model.ConversationM) error
	Update(ctx context.Context, obj *model.ConversationM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.ConversationM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.ConversationM, error)

	ConversationExpansion
}

// ConversationExpansion 定义了会话操作的附加方法
type ConversationExpansion interface{}

// conversationStore 是 ConversationStore 接口的实现.
type conversationStore struct {
	*genericstore.Store[model.ConversationM]
}

// 确保 conversationStore 实现了 ConversationStore 接口.
var _ ConversationStore = (*conversationStore)(nil)

// newConversationStore 创建 conversationStore 的实例.
func newConversationStore(store *datastore) *conversationStore {
	return &conversationStore{
		Store: genericstore.NewStore[model.ConversationM](store, NewLogger()),
	}
}
