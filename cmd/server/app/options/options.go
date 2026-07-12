package options

import (
	"github.com/robinlg/agentops-platform/internal/apiserver"
	genericoptions "github.com/robinlg/onexlib/pkg/options"
)

// ServerOptions 包含服务器配置选项(命令行选项)
type ServerOptions struct {
	// ServerMode 定义服务器模式：Gin HTTP
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`
	// EnableMemoryStore 指示是否启用内存数据库（用于测试或开发环境）
	EnableMemoryStore bool `json:"enable-memory-store" mapstructure:"enable-memory-store"`
	// PostgreSQLOptions 包含 PostgreSQL 配置选项
	PostgreSQLOptions *genericoptions.PostgreSQLOptions `json:"postgresql" mapstructure:"postgresql"`
	// HTTPOptions 包含 HTTP 配置选项
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
}

// NewServerOptions 创建带有默认值的 ServerOptions 实例
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:        apiserver.GinServerMode,
		EnableMemoryStore: false,
		PostgreSQLOptions: genericoptions.NewPostgreSQLOptions(),
		HTTPOptions:       genericoptions.NewHTTPOptions(),
	}

	return opts
}

// Config 基于 ServerOptions 构建 apiserver.Config.
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:        o.ServerMode,
		EnableMemoryStore: o.EnableMemoryStore,
		PostgreSQLOptions: o.PostgreSQLOptions,
		HTTPOptions:       o.HTTPOptions,
	}, nil
}
