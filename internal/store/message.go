package store

//go:generate mockgen -destination mock_message.go -package store github.com/robinlg/agentops-platform/internal/store MessageStore

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/model"
	genericstore "github.com/robinlg/onexlib/pkg/store"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// MessageStore 定义了 messages 模块在 store 层所实现的方法
type MessageStore interface {
	Create(ctx context.Context, obj *model.MessageM) error
	Update(ctx context.Context, obj *model.MessageM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.MessageM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.MessageM, error)

	MessageExpansion
}

// MessageExpansion 定义了消息操作的附加方法
type MessageExpansion interface{}

// messageStore 是 MessageStore 接口的实现.
type messageStore struct {
	*genericstore.Store[model.MessageM]
}

// 确保 messageStore 实现了 MessageStore 接口.
var _ MessageStore = (*messageStore)(nil)

// newMessageStore 创建 messageStore 的实例.
func newMessageStore(store *datastore) *messageStore {
	return &messageStore{
		Store: genericstore.NewStore[model.MessageM](store, NewLogger()),
	}
}
