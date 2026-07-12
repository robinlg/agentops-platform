package app

import (
	"github.com/robinlg/agentops-platform/cmd/server/app/options"
	"github.com/robinlg/agentops-platform/internal/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string // 配置文件路径

func NewAgentOpsCommand() *cobra.Command {
	// 创建默认的应用命令行选项
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		// 指定命令的名字，该名字会出现在帮助信息中
		Use: "agentops",
		// 命令的简短描述
		Short: "这是一个智能体运营平台",
		// 命令的详细描述
		Long: `这是一个智能体运营平台`,
		// 命令出错时，不打印帮助信息。设置为 true 可以确保命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
		// 设置命令运行时的参数检查，不需要指定命令行参数。例如：./onexblog param1 param2
		Args: cobra.NoArgs,
	}

	// 初始化配置函数，在每个命令运行时调用
	cobra.OnInitialize(onInitialize)

	// cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	// 推荐使用配置文件来配置应用，便于管理配置项
	//cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the agentops configuration file.")

	return cmd
}

// run 是主运行逻辑，解析配置等其他逻辑，并启动服务器
func run(opts *options.ServerOptions) error {
	// 初始化日志
	log.Init(logOptions())
	defer log.Sync() // 确保日志在退出时被刷新到磁盘

	// 将 viper 中的配置解析到 opts
	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	// 获取应用配置
	// 将命令行选项和应用配置分开，可以更加灵活的处理 2 种不同类型的配置
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	// 创建服务器实例
	// 注意这里是联合服务器，因为可能同时启动多个不同类型的服务器.
	server, err := cfg.NewServer()
	if err != nil {
		return err
	}

	// 启动服务器
	return server.Run()
}

// logOptions 从 viper 中读取日志配置，构建 *log.Options 并返回
// 注意：viper.Get<Type>() 中 key 的名字需要使用 . 分割，以跟 YAML 中保持相同的缩进
func logOptions() *log.Options {
	opts := log.NewOptions()
	if viper.IsSet("log.disable-caller") {
		opts.DisableCaller = viper.GetBool("log.disable-caller")
	}
	if viper.IsSet("log.disable-stacktrace") {
		opts.DisableStacktrace = viper.GetBool("log.disable-stacktrace")
	}
	if viper.IsSet("log.level") {
		opts.Level = viper.GetString("log.level")
	}
	if viper.IsSet("log.format") {
		opts.Format = viper.GetString("log.format")
	}
	if viper.IsSet("log.output-paths") {
		opts.OutputPaths = viper.GetStringSlice("log.output-paths")
	}
	return opts
}
