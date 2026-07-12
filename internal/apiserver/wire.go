//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"
	"github.com/robinlg/agentops-platform/internal/apiserver/biz"
	"github.com/robinlg/agentops-platform/internal/apiserver/store"
	"github.com/robinlg/agentops-platform/internal/pkg/server"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB, // 提供数据库实例
	)
	return nil, nil
}
