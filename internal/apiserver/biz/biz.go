package biz

import (
	"github.com/google/wire"
	modelproviderv1 "github.com/robinlg/agentops-platform/internal/apiserver/biz/v1/model-provider"
	"github.com/robinlg/agentops-platform/internal/apiserver/store"
)

//go:generate mockgen -destination mock_biz.go -package biz github.com/robinlg/agentops-platform/internal/apiserver/biz IBiz

// ProviderSet 是一个 Wire 的 Provider 集合，用于声明依赖注入的规则.
// 包含 NewBiz 构造函数，用于生成 biz 实例.
// wire.Bind 用于将接口 IBiz 与具体实现 *biz 绑定，
// 这样依赖 IBiz 的地方会自动注入 *biz 实例.
var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

// IBiz 定义了业务层需要实现的方法
type IBiz interface {
	// ModelProviderV1 模型提供商业务接口
	ModelProviderV1() modelproviderv1.ModelProviderBiz
}

// biz 是 IBiz 的一个具体实现
type biz struct {
	store store.IStore
}

// 确保 biz 实现了 IBiz 接口
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例
func NewBiz(store store.IStore) *biz {
	return &biz{store: store}
}

// ModelProviderV1 返回一个实现了 ModelProviderBiz 接口的实例
func (b *biz) ModelProviderV1() modelproviderv1.ModelProviderBiz {
	return modelproviderv1.New(b.store)
}
