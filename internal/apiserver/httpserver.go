package apiserver

import (
	"context"

	"github.com/gin-gonic/gin"
	handler "github.com/robinlg/agentops-platform/internal/apiserver/handler/http"
	"github.com/robinlg/agentops-platform/internal/pkg/errno"
	"github.com/robinlg/agentops-platform/internal/pkg/server"
	"github.com/robinlg/onexlib/pkg/core"
)

// ginServer 定义一个使用 Gin 框架开发的 HTTP 服务器
type ginServer struct {
	srv server.Server
}

// 确保 *ginServer 实现了 server.Server 接口.
var _ server.Server = (*ginServer)(nil)

// NewGinServer 初始化一个新的 Gin 服务器实例
func (c *ServerConfig) NewGinServer() server.Server {
	// 创建 Gin 引擎
	engine := gin.New()

	// 注册 REST API 路由
	c.InstallRESTAPI(engine)

	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, c.cfg.TLSOptions, engine)

	return &ginServer{srv: httpsrv}
}

// InstallRESTAPI 注册 API 路由。路由的路径和 HTTP 方法，严格遵循 REST 规范.
func (c *ServerConfig) InstallRESTAPI(engine *gin.Engine) {
	// 注册业务无关的 API 接口
	InstallGenericAPI(engine)

	// 创建核心业务处理器
	handler := handler.NewHandler(c.biz)

	// 注册健康检查接口
	engine.GET("/healthz", handler.Healthz)

	// 注册 v1 版本 API 路由分组
	v1 := engine.Group("/v1")
	{
		// 模型提供商相关路由
		modelProvider := v1.Group("/model-providers")
		{
			modelProvider.POST("", handler.CreateModelProvider)
			modelProvider.PUT("/:id", handler.UpdateModelProvider)
			modelProvider.DELETE("/:id", handler.DeleteModelProvider)
			modelProvider.GET("/:id", handler.GetModelProvider)
			modelProvider.GET("", handler.ListModelProviders)
		}
	}
}

// InstallGenericAPI 注册业务无关的路由，例如 pprof、404 处理等.
func InstallGenericAPI(engine *gin.Engine) {
	// 注册 404 路由处理
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})
}

// RunOrDie 启动 Gin 服务器，出错则程序崩溃退出
func (s *ginServer) RunOrDie() {
	s.srv.RunOrDie()
}

// GracefulStop 优雅停止服务器
func (s *ginServer) GracefulStop(ctx context.Context) {
	s.srv.GracefulStop(ctx)
}
