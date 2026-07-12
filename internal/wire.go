//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/robinlg/agentops-platform/internal/biz"
	"github.com/robinlg/agentops-platform/internal/llm"
	"github.com/robinlg/agentops-platform/internal/pkg/server"
	"github.com/robinlg/agentops-platform/internal/runtime"
	"github.com/robinlg/agentops-platform/internal/store"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		// 提供 Web 服务器实例，并从 *Config 中提取 ServerMode 字段用于注入
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		// 构建 ServerConfig 结构体，* 表示注入其全部字段
		wire.Struct(new(ServerConfig), "*"),
		// 提供 store 数据访问层和 biz 业务逻辑层的 Provider 集合
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		// 提供 PromptBuilder 实例
		runtime.NewPromptBuilder,
		// 提供 LLM Client 实例
		llm.NewClient,
		// 提供数据库实例
		ProvideDB,
	)
	return nil, nil
}
