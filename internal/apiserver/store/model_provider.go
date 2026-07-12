package store

import (
	"context"

	"github.com/robinlg/agentops-platform/internal/apiserver/model"
	genericstore "github.com/robinlg/onexlib/pkg/store"
	"github.com/robinlg/onexlib/pkg/store/where"
)

// ModelProviderStore 定义了 model_provider 模块在 store 层所实现的方法
type ModelProviderStore interface {
	Create(ctx context.Context, obj *model.ModelProviderM) error
	Update(ctx context.Context, obj *model.ModelProviderM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.ModelProviderM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.ModelProviderM, error)

	ModelProviderExpansion
}

// ModelProviderExpansion 定义了模型提供商操作的附加方法
type ModelProviderExpansion interface{}

// modelProviderStore 是 ModelProviderStore 接口的实现.
type modelProviderStore struct {
	*genericstore.Store[model.ModelProviderM]
}

// 确保 modelProviderStore 实现了 ModelProviderStore 接口.
var _ ModelProviderStore = (*modelProviderStore)(nil)

// newModelProviderStore 创建 modelProviderStore 的实例.
func newModelProviderStore(store *datastore) *modelProviderStore {
	return &modelProviderStore{
		Store: genericstore.NewStore[model.ModelProviderM](store, NewLogger()),
	}
}
