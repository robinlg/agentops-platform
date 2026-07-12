package store

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/apiserver/model"
	genericstore "github.com/robinlg/onexlib/pkg/store"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// AgentStore 定义了 agent 模块在 store 层所实现的方法
type AgentStore interface {
	Create(ctx context.Context, obj *model.AgentM) error
	Update(ctx context.Context, obj *model.AgentM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.AgentM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.AgentM, error)

	AgentExpansion
}

// AgentExpansion 定义了智能体操作的附加方法
type AgentExpansion interface{}

// agentStore 是 AgentStore 接口的实现.
type agentStore struct {
	*genericstore.Store[model.AgentM]
}

// 确保 agentStore 实现了 AgentStore 接口.
var _ AgentStore = (*agentStore)(nil)

// newAgentStore 创建 agentStore 的实例.
func newAgentStore(store *datastore) *agentStore {
	return &agentStore{
		Store: genericstore.NewStore[model.AgentM](store, NewLogger()),
	}
}
