package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// defaultHomeDir 定义放置服务配置的默认目录.
	defaultHomeDir = ".agentops"

	// defaultConfigName 指定服务的默认配置文件名.
	defaultConfigName = "config.yaml"

	// defaultConfigDir 定义项目内的配置文件目录.
	defaultConfigDir = "./config"
)

// onInitialize 设置需要读取的配置文件名、环境变量，并将其内容读取到 viper 中
func onInitialize() {
	if configFile != "" {
		// 从命令行选项指定的配置文件中读取
		viper.SetConfigFile(configFile)
	} else {
		// 使用默认配置文件路径和名称
		for _, dir := range searchDirs() {
			// 将 dir 目录加入到配置文件的搜索路径
			viper.AddConfigPath(dir)
		}

		// 设置配置文件格式为 YAML
		viper.SetConfigType("yaml")

		// 配置文件名称（没有文件扩展名）
		viper.SetConfigName(defaultConfigName)
	}

	// 读取配置文件.如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read viper configuration file, err: %v", err)
	}

	// 打印当前使用的配置文件，方便调试
	log.Printf("Using config file: %s", viper.ConfigFileUsed())
}

// searchDirs 返回默认的配置文件搜索目录
func searchDirs() []string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	// 如果获取用户主目录失败，则打印错误信息并退出程序
	cobra.CheckErr(err)
	return []string{filepath.Join(homeDir, defaultHomeDir), ".", defaultConfigDir}
}

// filePath 获取默认配置文件的完整路径
func filePath() string {
	home, err := os.UserHomeDir()
	// 如果不能获取用户主目录，则记录错误并返回空路径
	cobra.CheckErr(err)
	return filepath.Join(home, defaultHomeDir, defaultConfigName)
}
