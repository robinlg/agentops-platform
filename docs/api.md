# API 参考（REST）

本文档基于 [`internal/httpserver.go`](../internal/httpserver.go) 与 [`pkg/api/v1/*.proto`](../pkg/api/v1) 生成。所有接口默认监听 `HTTPOptions.Addr`（示例 `127.0.0.1:38443`）。

- 前缀：`/v1`
- Content-Type：`application/json`
- 响应包装：由 `onexlib/pkg/core.WriteResponse` 统一处理；错误使用 `internal/pkg/errno` 定义的错误码。
- 时间字段：`google.protobuf.Timestamp`（JSON 中为 RFC3339 字符串）。

联调脚本见 [`test/http/`](../test/http)。

---

## 目录

- [0. 健康检查](#0-健康检查)
- [1. ModelProvider](#1-modelprovider)
- [2. Agent](#2-agent)
- [3. Chat](#3-chat)
- [4. Conversation](#4-conversation)
- [5. AgentRun](#5-agentrun)
- [附录：公共对象](#附录公共对象)

---

## 0. 健康检查

### 0.1 存活探针

**接口路径**

`GET /healthz`

**输入参数**

无。

**输出参数**

HTTP 200 表示进程存活，无响应体。

---

## 1. ModelProvider

模型提供商配置，例如 OpenAI / Azure / 兼容协议网关的 `base_url` + `api_key`。Proto 契约：[`model_provider.proto`](../pkg/api/v1/model_provider.proto)。

### 1.1 新建模型提供商

**接口路径**

`POST /v1/model-providers`

**输入参数**（Body，`CreateModelProviderRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| name          | body | string | 是 | 名称 |
| provider_type | body | string | 是 | 类型（openai / anthropic / azure …） |
| base_url      | body | string | 是 | 访问 URL |
| api_key       | body | string | 是 | API Key |
| default_model | body | string | 是 | 默认模型 |

**输出参数**（`CreateModelProviderResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | int64 | 新创建模型提供商的主键 |

---

### 1.2 查询模型提供商列表

**接口路径**

`GET /v1/model-providers`

**输入参数**（Query，`ListModelProviderRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| offset | query | int64 | 否 | 偏移量，用于分页，默认 0 |
| limit  | query | int64 | 否 | 每页数量限制 |

**输出参数**（`ListModelProviderResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| total_count     | int64 | 满足条件的总条数 |
| model_providers | ModelProvider[] | 当前页的模型提供商列表（结构见[附录](#附录公共对象)） |

---

### 1.3 查询单个模型提供商

**接口路径**

`GET /v1/model-providers/:id`

**输入参数**（`GetModelProviderRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id | uri | int64 | 是 | 模型提供商 ID |

**输出参数**（`GetModelProviderResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| model_provider | ModelProvider | 模型提供商详情（结构见[附录](#附录公共对象)） |

---

### 1.4 更新模型提供商（部分更新）

**接口路径**

`PUT /v1/model-providers/:id`

**输入参数**（`UpdateModelProviderRequest`；所有业务字段均为 `optional`，未设置的字段不参与更新）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id            | uri  | int64  | 是 | 待更新的模型提供商 ID |
| name          | body | string | 否 | 名称 |
| provider_type | body | string | 否 | 类型 |
| base_url      | body | string | 否 | 访问 URL |
| api_key       | body | string | 否 | API Key |
| default_model | body | string | 否 | 默认模型 |

**输出参数**（`UpdateModelProviderResponse`）

无字段（空对象 `{}`）。

---

### 1.5 删除模型提供商

**接口路径**

`DELETE /v1/model-providers/:id`

**输入参数**（`DeleteModelProviderRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id | uri | int64 | 是 | 待删除的模型提供商 ID |

**输出参数**（`DeleteModelProviderResponse`）

无字段（空对象 `{}`）。

---

## 2. Agent

一个智能体 = 一份"人格 / 提示词 / 模型参数"配置，绑定到某个 ModelProvider。Proto 契约：[`agent.proto`](../pkg/api/v1/agent.proto)。

### 2.1 新建智能体

**接口路径**

`POST /v1/agents`

**输入参数**（Body，`CreateAgentRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| name              | body | string | 是 | 智能体名称 |
| description       | body | string | 否 | 描述 |
| system_prompt     | body | string | 否 | 系统提示词 |
| model_provider_id | body | int64  | 是 | 绑定的模型提供商 ID |
| model             | body | string | 否 | 覆盖 provider 的默认模型 |
| temperature       | body | double | 否 | 采样温度（默认 0.7） |
| max_tokens        | body | int32  | 否 | 最大 token（默认 2048） |

**输出参数**（`CreateAgentResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| id | int64 | 新创建智能体的主键 |

---

### 2.2 查询智能体列表

**接口路径**

`GET /v1/agents`

**输入参数**（Query，`ListAgentRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| offset | query | int64 | 否 | 偏移量，用于分页 |
| limit  | query | int64 | 否 | 每页数量限制 |

**输出参数**（`ListAgentResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| total_count | int64 | 总条数 |
| agents      | Agent[] | 当前页的智能体列表（结构见[附录](#附录公共对象)） |

---

### 2.3 查询单个智能体

**接口路径**

`GET /v1/agents/:id`

**输入参数**（`GetAgentRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id | uri | int64 | 是 | 智能体 ID |

**输出参数**（`GetAgentResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| agent | Agent | 智能体详情（结构见[附录](#附录公共对象)） |

---

### 2.4 更新智能体（部分更新）

**接口路径**

`PUT /v1/agents/:id`

**输入参数**（`UpdateAgentRequest`；所有业务字段均为 `optional`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id                | uri  | int64  | 是 | 待更新的智能体 ID |
| name              | body | string | 否 | 名称 |
| description       | body | string | 否 | 描述 |
| system_prompt     | body | string | 否 | 系统提示词 |
| model_provider_id | body | int64  | 否 | 模型提供商 ID |
| model             | body | string | 否 | 模型 |
| temperature       | body | double | 否 | 采样温度 |
| max_tokens        | body | int32  | 否 | 最大 token |

**输出参数**（`UpdateAgentResponse`）

无字段（空对象 `{}`）。

---

### 2.5 删除智能体

**接口路径**

`DELETE /v1/agents/:id`

**输入参数**（`DeleteAgentRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id | uri | int64 | 是 | 待删除的智能体 ID |

**输出参数**（`DeleteAgentResponse`）

无字段（空对象 `{}`）。

---

## 3. Chat

对话入口 —— **这是第一版的核心接口**。触发一次 LLM 调用，同时落库：用户消息、assistant 回复、`AgentRun` Trace。Proto 契约：[`chat.proto`](../pkg/api/v1/chat.proto)。

### 3.1 发起一次对话

**接口路径**

`POST /v1/agents/:id/chat`

**输入参数**（`CreateChatRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| agent_id        | uri  | int64  | 是 | 从 URI 绑定（`:id`） |
| conversation_id | body | int64  | 否 | 会话 ID；**不传则自动新建会话**，Title 取用户首条消息前 30 字 |
| message         | body | string | 是 | 本轮用户输入 |

**输出参数**（`CreateChatResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| conversation_id | int64  | 使用/新建的会话 ID |
| message_id      | int64  | assistant 回复消息 ID |
| run_id          | int64  | 本次 `AgentRun` 记录 ID |
| answer          | string | assistant 回复内容 |
| usage           | Usage  | Token 使用量（结构见下） |
| latency_ms      | int64  | 本次 LLM 调用耗时（毫秒） |

`Usage` 结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| prompt_tokens     | int32 | 输入 token 数 |
| completion_tokens | int32 | 输出 token 数 |
| total_tokens      | int32 | 总 token 数 |

关键行为：
- 会话不存在或不属于该 agent → 报错。
- LLM 调用失败 → `AgentRun.Status = failed` 落库，接口返回错误。
- 成功路径见 [architecture.md#5-一次-chat-调用的数据流](./architecture.md#5-一次-chat-调用的数据流)。

联调示例：[`test/http/chat.http`](../test/http/chat.http)。

---

## 4. Conversation

会话及其消息列表。会话不由 handler 显式创建 —— 它作为 chat 接口的副产物产生。Proto 契约：[`conversation.proto`](../pkg/api/v1/conversation.proto)。

### 4.1 查询会话列表

**接口路径**

`GET /v1/conversations`

**输入参数**（Query，`ListConversationRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| offset | query | int64 | 否 | 偏移量，用于分页 |
| limit  | query | int64 | 否 | 每页数量限制 |

**输出参数**（`ListConversationResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| total_count   | int64 | 总条数 |
| conversations | Conversation[] | 当前页的会话列表（结构见[附录](#附录公共对象)） |

---

### 4.2 查询会话的消息列表

**接口路径**

`GET /v1/conversations/:id/messages`

**输入参数**（`ListMessageRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| conversation_id | uri   | int64 | 是 | 会话 ID（从 URI 绑定，`:id`） |
| offset          | query | int64 | 否 | 偏移量，用于分页 |
| limit           | query | int64 | 否 | 每页数量限制 |

**输出参数**（`ListMessageResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| total_count | int64 | 总条数 |
| messages    | Message[] | 当前页的消息列表（结构见[附录](#附录公共对象)） |

---

### 4.3 删除会话

**接口路径**

`DELETE /v1/conversations/:id`

在事务中依次删除该会话下的 `messages`、`agent_runs`，再删除会话本身，避免孤儿数据。

**输入参数**（`DeleteConversationRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id | uri | int64 | 是 | 待删除的会话 ID |

**输出参数**（`DeleteConversationResponse`）

无字段（空对象 `{}`）。

---

## 5. AgentRun

一次 LLM 调用的运行 Trace：输入 / 输出 / 使用模型 / token / 耗时 / 状态。**这是第一版的 Trace 载体**。Proto 契约：[`agent_run.proto`](../pkg/api/v1/agent_run.proto)。

### 5.1 查询运行记录列表

**接口路径**

`GET /v1/agent-runs`

**输入参数**（Query，`ListAgentRunRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| offset | query | int64 | 否 | 偏移量，用于分页 |
| limit  | query | int64 | 否 | 每页数量限制 |

**输出参数**（`ListAgentRunResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| total_count | int64 | 总条数 |
| agent_runs  | AgentRun[] | 当前页的运行记录列表（结构见[附录](#附录公共对象)） |

---

### 5.2 查询单条运行记录

**接口路径**

`GET /v1/agent-runs/:id`

**输入参数**（`GetAgentRunRequest`）

| 字段 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| id | uri | int64 | 是 | 运行记录 ID |

**输出参数**（`GetAgentRunResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| agent_run | AgentRun | 运行记录详情（结构见[附录](#附录公共对象)） |

---

## 附录：公共对象

以下资源在多个接口的响应中复用。

### ModelProvider

| 字段 | 类型 | 说明 |
|---|---|---|
| id            | int64     | 主键 |
| name          | string    | 名称 |
| provider_type | string    | 类型 |
| base_url      | string    | 访问 URL |
| api_key       | string    | API Key（建议由业务层脱敏） |
| default_model | string    | 默认模型 |
| created_at    | Timestamp | 创建时间 |
| updated_at    | Timestamp | 更新时间 |

### Agent

| 字段 | 类型 | 说明 |
|---|---|---|
| id                | int64     | 主键 |
| name              | string    | 名称 |
| description       | string    | 描述 |
| system_prompt     | string    | 系统提示词 |
| model_provider_id | int64     | 绑定的模型提供商 ID |
| model             | string    | 模型 |
| temperature       | double    | 采样温度 |
| max_tokens        | int32     | 最大 token |
| created_at        | Timestamp | 创建时间 |
| updated_at        | Timestamp | 更新时间 |

### Conversation

| 字段 | 类型 | 说明 |
|---|---|---|
| id         | int64     | 主键 |
| agent_id   | int64     | 智能体 ID |
| title      | string    | 会话标题 |
| created_at | Timestamp | 创建时间 |
| updated_at | Timestamp | 更新时间 |

### Message

| 字段 | 类型 | 说明 |
|---|---|---|
| id              | int64     | 主键 |
| conversation_id | int64     | 会话 ID |
| role            | string    | 角色：`system / user / assistant / tool`（见 [`message_role.go`](../internal/model/message_role.go)） |
| content         | string    | 消息内容 |
| created_at      | Timestamp | 创建时间 |

### AgentRun

| 字段 | 类型 | 说明 |
|---|---|---|
| id                | int64     | 主键 |
| agent_id          | int64     | 关联的智能体 |
| conversation_id   | int64     | 关联的会话 |
| status            | string    | `pending / running / success / failed` |
| input             | string    | 用户本轮输入 |
| output            | string    | assistant 回复 |
| model             | string    | LLM 实际使用的模型 |
| prompt_tokens     | int32     | 输入 token 数 |
| completion_tokens | int32     | 输出 token 数 |
| total_tokens      | int32     | 总 token 数 |
| latency_ms        | int64     | 耗时（毫秒） |
| error_message     | string    | 失败时的错误信息 |
| started_at        | Timestamp | 开始时间 |
| finished_at       | Timestamp | 结束时间 |

---

## 附录：常用状态码 & 错误

- `errno.ErrPageNotFound`：命中 404 兜底路由。
- 其余业务错误码定义于 [`internal/pkg/errno/code.go`](../internal/pkg/errno/code.go)。

## 附录：本地联调

REST Client 脚本目录：[`test/http/`](../test/http)。VS Code / GoLand 均支持直接点击运行。
