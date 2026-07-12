package store

//go:generate mockgen -destination mock_agent_run.go -package store github.com/robinlg/agentops-platform/internal/store AgentRunStore

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/model"
	genericstore "github.com/robinlg/onexlib/pkg/store"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// AgentRunStore 定义了 agent_runs 模块在 store 层所实现的方法
type AgentRunStore interface {
	Create(ctx context.Context, obj *model.AgentRunM) error
	Update(ctx context.Context, obj *model.AgentRunM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.AgentRunM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.AgentRunM, error)

	AgentRunExpansion
}

// AgentRunExpansion 定义了智能体运行记录操作的附加方法
type AgentRunExpansion interface{}

// agentRunStore 是 AgentRunStore 接口的实现.
type agentRunStore struct {
	*genericstore.Store[model.AgentRunM]
}

// 确保 agentRunStore 实现了 AgentRunStore 接口.
var _ AgentRunStore = (*agentRunStore)(nil)

// newAgentRunStore 创建 agentRunStore 的实例.
func newAgentRunStore(store *datastore) *agentRunStore {
	return &agentRunStore{
		Store: genericstore.NewStore[model.AgentRunM](store, NewLogger()),
	}
}
