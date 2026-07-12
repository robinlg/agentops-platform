# API 手工测试脚本 (.http)

本目录存放 `.http` 格式的接口测试文件，可直接在 **GoLand / IntelliJ IDEA** 内置的 HTTP Client，
或 **VSCode 的 [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) 插件** 里点击运行。

> ⚠️ 这些文件不是 Go 单元测试，`go test` 不会执行它们。

## 目录结构

```
test/http/
├── README.md                # 本说明文件
├── http-client.env.json     # 环境变量（BASE_URL 等，GoLand 使用）
├── healthz.http             # 健康检查
└── model_provider.http      # 模型提供商 CRUD 全流程
```

## 前置条件

启动本地服务：

```bash
go run ./cmd/server --config ./config/config.yaml
```

## 使用方式

### GoLand / IntelliJ IDEA

1. 打开任意 `.http` 文件
2. 右上角选择环境 `dev`（对应 `http-client.env.json` 里的配置）
3. 每个请求上方会有绿色 ▶️ 按钮，点击即可发送

### VSCode

1. 安装 [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) 插件
2. 打开 `.http` 文件
3. 每个 `###` 分隔的请求上方会出现 `Send Request` 链接，点击即可发送
4. 文件里已经用 `@baseUrl` 等 in-file 变量声明，可直接工作

## 修改默认地址

编辑 [http-client.env.json](./http-client.env.json)，修改 `baseUrl` 值即可。
