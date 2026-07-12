package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/robinlg/onexlib/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// 帮助信息文本.
const helpText = `Usage: main [flags] arg [arg...]

This is a pflag example.

Flags:
`

// GenerateConfig 保存代码生成的配置
type GenerateConfig struct {
	ModelPackagePath string
	GenerateFunc     func(g *gen.Generator)
}

// 预定义的生成配置
var generateConfigs = map[string]GenerateConfig{
		"mb": {ModelPackagePath: "../../internal/model", GenerateFunc: GenerateModels},
}

// 命令行参数.
var (
	addr       = pflag.StringP("addr", "a", "127.0.0.1:5432", "PostgreSQL host address.")
	username   = pflag.StringP("username", "u", "postgres", "Username to connect to the database.")
	password   = pflag.StringP("password", "p", "postgres123456", "Password to use when connecting to the database.")
	database   = pflag.StringP("db", "d", "agentops", "Database name to connect to.")
	modelPath  = pflag.String("model-pkg-path", "", "Generated model code's package name.")
	components = pflag.StringSlice("component", []string{"mb"}, "Generated model code's for specified component.")
	help       = pflag.BoolP("help", "h", false, "Show this help message.")
)

func main() {
	// 设置自定义的使用说明函数
	pflag.Usage = func() {
		fmt.Printf("%s", helpText)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	// 如果设置了帮助标志，则显示帮助信息并退出
	if *help {
		pflag.Usage()
		return
	}

	// 初始化数据库连接
	dbInstance, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 处理组件并生成代码
	for _, component := range *components {
		processComponent(component, dbInstance)
	}
}

// initializeDatabase 创建并返回一个数据库连接
func initializeDatabase() (*gorm.DB, error) {
	dbOptions := &db.PostgreSQLOptions{
		Addr:     *addr,
		Username: *username,
		Password: *password,
		Database: *database,
	}

	// 创建并返回数据库连接
	return db.NewPostgreSQL(dbOptions)
}

// processComponent 处理单个组件以生成代码
func processComponent(component string, dbInstance *gorm.DB) {
	config, ok := generateConfigs[component]
	if !ok {
		log.Printf("Component '%s' not found in configuration. Skipping.", component)
		return
	}

	// 解析模型包路径
	modelPkgPath := resolveModelPackagePath(config.ModelPackagePath)

	// 创建生成器实例
	generator := createGenerator(modelPkgPath)
	generator.UseDB(dbInstance)

	// 应用自定义生成器选项
	applyGeneratorOptions(generator)

	// 使用指定的函数生成模型
	config.GenerateFunc(generator)

	// 执行代码生成
	generator.Execute()
}

// resolveModelPackagePath 确定模型生成的包路径
func resolveModelPackagePath(defaultPath string) string {
	if *modelPath != "" {
		return *modelPath
	}
	absPath, err := filepath.Abs(defaultPath)
	if err != nil {
		log.Printf("Error resolving path: %v", err)
		return defaultPath
	}
	return absPath
}

// createGenerator 初始化并返回一个新的生成器实例
func createGenerator(packagePath string) *gen.Generator {
	return gen.NewGenerator(gen.Config{
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		ModelPkgPath:      packagePath,
		WithUnitTest:      true,
		FieldNullable:     true,  // 对于数据库中可空的字段，使用指针类型
		FieldSignable:     false, // 禁用无符号属性以提高兼容性
		FieldWithIndexTag: false, // 不包含 GORM 的索引标签
		FieldWithTypeTag:  false, // 不包含 GORM 的类型标签
	})
}

// applyGeneratorOptions 设置自定义生成器选项
func applyGeneratorOptions(g *gen.Generator) {
	// 为特定字段自定义 GORM 标签
	g.WithOpts(
		gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
		gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
	)
}

// GenerateModels 生成模型
func GenerateModels(g *gen.Generator) {
	g.GenerateModelAs(
		"model_providers",
		"ModelProviderM",
	)
	g.GenerateModelAs(
		"agents",
		"AgentM",
	)
	g.GenerateModelAs(
		"conversations",
		"ConversationM",
	)
	g.GenerateModelAs(
		"messages",
		"MessageM",
	)
	g.GenerateModelAs(
		"agent_runs",
		"AgentRunM",
	)
}
