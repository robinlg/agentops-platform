package apiserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robinlg/agentops-platform/internal/apiserver/biz"
	"github.com/robinlg/agentops-platform/internal/pkg/log"
	"github.com/robinlg/agentops-platform/internal/pkg/server"
	genericoptions "github.com/robinlg/onexlib/pkg/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置
type Config struct {
	ServerMode        string
	EnableMemoryStore bool
	TLSOptions        *genericoptions.TLSOptions
	HTTPOptions       *genericoptions.HTTPOptions
	PostgreSQLOptions *genericoptions.PostgreSQLOptions
}

type Server struct {
	srv server.Server
}

// ServerConfig 包含服务器的核心依赖和配置
type ServerConfig struct {
	cfg *Config
	biz biz.IBiz
}

// NewServer 根据配置创建服务器.
func (cfg *Config) NewServer() (*Server, error) {
	// 创建服务配置，这些配置可用来创建服务器
	srv, err := InitializeWebServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{srv: srv}, nil
}

// Run 运行应用.
func (s *Server) Run() error {
	go s.srv.RunOrDie()

	// 创建一个 os.Signal 类型的 channel，用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞程序，等待从 quit channel 中接收到信号
	<-quit

	log.Infow("Shutting down server...")

	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	s.srv.GracefulStop(ctx)

	log.Infow("Server exited")
	return nil
}

// NewDB 创建一个 *gorm.DB 实例.
func (cfg *Config) NewDB() (*gorm.DB, error) {
	if !cfg.EnableMemoryStore {
		log.Infow("Initializing database connection", "type", "postgresql", "addr", cfg.PostgreSQLOptions.Addr)
		return cfg.PostgreSQLOptions.NewDB()
	}

	log.Infow("Initializing database connection", "type", "memory", "engine", "SQLite")
	// 使用SQLite内存模式配置数据库
	// ?cache=shared 用于设置 SQLite 的缓存模式为 共享缓存模式 (shared)。
	// 默认情况下，SQLite 的每个数据库连接拥有自己的独立缓存，这种模式称为 专用缓存 (private)。
	// 使用 共享缓存模式 (shared) 后，不同连接可以共享同一个内存中的数据库和缓存。
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Errorw("Failed to create database connection", "err", err)
		return nil, err
	}

	return db, nil
}

// ProvideDB 根据配置提供一个数据库实例。
func ProvideDB(cfg *Config) (*gorm.DB, error) {
	return cfg.NewDB()
}

func NewWebServer(serverMode string, serverConfig *ServerConfig) (server.Server, error) {
	// 根据服务模式创建对应的服务实例
	switch serverMode {
	case GinServerMode:
		return serverConfig.NewGinServer(), nil
	default:
		return serverConfig.NewGinServer(), nil
	}
}
