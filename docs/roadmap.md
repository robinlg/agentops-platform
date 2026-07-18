# Roadmap

本文档规划 `agentops-platform` 的分阶段演进目标。整体参考 [Kubeagi easyai](https://github.com/kubeagi) 的资源模型，但按自身节奏做简化，第一版聚焦最小闭环。

版本映射一览表：

| 本项目资源          | easyai 参考           | 第一版 | 第二版 | 第三版 |
|--------------------|----------------------|:-----:|:-----:|:-----:|
| ModelProvider      | LLM                  | ✅ 做  |   —   |   —   |
| Agent              | Agent                | ✅ 简化版 | 🔼 增强 |   —   |
| Conversation       | Conversation         | ✅ 做  |   —   |   —   |
| Message            | Message              | ✅ 做  |   —   |   —   |
| AgentRun           | 会话运行记录 / Trace  | ✅ 简化版 | 🔼 完整 Trace |   —   |
| KnowledgeBase      | KnowledgeBase        |   —   | ✅ 做 |   —   |
| Document           | Dataset/Datasource/VersionedDataset |   —   | ✅ 简化版 |   —   |
| VectorStore        | VectorStore          |   —   | ✅ 简化版 |   —   |
| Tool               | Agent tool           |   —   |   —   | ✅ 做 |
| AgentService       | Worker / Model       |   —   |   —   | 后期做 |
| eai-apiserver      | 自定义 API Server    |  ❌ 不做  |   —   |   —   |
| eai-controller-manager | Controller Manager |  —   |   —   | 后期只做小 Operator |

图例：✅ 完成 / 🔼 增强 / — 未涉及 / ❌ 不做

---

## V1 · 最小闭环（当前版本）

**目标**：跑通 `模型接入 → 定义智能体 → 会话对话 → 落库 Trace` 的完整链路，让平台"真的能对话"。

### 已完成

- **ModelProvider**：OpenAI 兼容协议 CRUD，可挂接任意兼容网关。
- **Agent**：CRUD + `system_prompt / temperature / max_tokens / model` 覆盖。
- **Conversation / Message**：会话与消息落库；chat 接口在会话 ID 为空时自动新建会话，Title 自动截取首条消息。
- **Chat**：`POST /v1/agents/:id/chat` 一站式：加载上下文 → 落库用户消息 → 创建 AgentRun → 调 LLM → 落库 assistant → 回填 Trace → 返回。
- **AgentRun（简化版 Trace）**：`pending / running / success / failed` 状态机，记录 input / output / model / tokens / latency / error_message / started_at / finished_at。
- **LLM 抽象层**：`llm.Client` 接口 + OpenAI-Compatible 实现，未来可按 `provider_type` 分支扩展。
- **数据一致性**：删除 Conversation 时用事务级联清理 `messages` 与 `agent_runs`。
- **基础设施**：Gin HTTP Server / GORM / Wire DI / 结构化日志 / 单元测试（biz + runtime + llm）+ REST Client 联调脚本。

### 明确不做

- 用户 / 租户 / 鉴权（内部工具，先跳过）。
- 流式（SSE / WebSocket）响应。
- 自定义 API Server 抽象（`eai-apiserver`）。
- K8s CRD / Operator（`eai-controller-manager`）。

---

## V2 · 引入 RAG

**目标**：让 Agent 具备"读私有资料"的能力。围绕 `KnowledgeBase / Document / VectorStore` 三类资源做**简化版**闭环。

### 计划新增

- **KnowledgeBase**（做）
  - CRUD；一个 KB 归属一个 owner / 若干 Document。
  - 与 Agent 建立多对多绑定（Agent 检索时选择一个或多个 KB）。
- **Document**（简化版）
  - 上传 → 解析（PDF / MD / TXT）→ 分块（chunk）→ 入库。
  - 参考 easyai 的 `Dataset / Datasource / VersionedDataset` 三级抽象，V2 合并为单一 `Document` 资源，简化 API。
- **VectorStore**（简化版）
  - 抽象 `VectorStore.Client` 接口 + 一个默认实现（例如 pgvector 或 Chroma）。
  - 提供 `Upsert(chunks) / Query(query, topK)`。
- **Retriever（新的 runtime 组件）**
  - 在 `internal/runtime` 下新增 `Retriever`，在 `PromptBuilder.Build` 前完成向量召回并把结果拼进 prompt。
  - 支持简单的 "top-k + 相似度阈值" 策略，暂不做 Rerank。
- **AgentRun 增强**
  - 记录检索到的 Document / Chunk 引用（用于回答溯源）。
  - Trace 视图初步成型（列表 + 详情 + 检索片段回显）。

### 主要新增 API（示意）

- `POST /v1/knowledge-bases`、`GET /v1/knowledge-bases`、…
- `POST /v1/knowledge-bases/:id/documents`（上传）、`GET .../documents`、`DELETE .../documents/:did`
- `POST /v1/agents/:id/knowledge-bases`（绑定 KB）
- Chat 无需新接口，`AgentRun.retrieved_chunks` 自动附加。

### 里程碑

- [ ] KB / Document CRUD & 上传解析
- [ ] Embedding 抽象接入（同样通过 `ModelProvider` 扩展 `provider_type=embedding`）
- [ ] VectorStore 接入 pgvector（默认）
- [ ] Retriever 接入 chat 主流程
- [ ] AgentRun 追加检索 Trace 字段
- [ ] E2E：上传一本《员工手册》→ 提问命中原文

---

## V3 · Tool Calling

**目标**：让 Agent 具备"调用外部工具"的能力，从"聊天"进入"办事"。

### 计划新增

- **Tool 资源**
  - CRUD：`name / description / schema(JSON Schema) / endpoint / auth`。
  - 与 Agent 多对多绑定。
- **Runtime 层增强**
  - `ToolExecutor`：根据 LLM 返回的 `tool_calls` 反射调用 Tool 并把结果作为 `role=tool` 消息重新喂给 LLM。
  - 支持多轮 tool-use 循环（带最大轮数保护）。
- **LLM 层增强**
  - 扩展 `llm.Message` / `llm.ChatResult`，支持 `tool_calls / tool_call_id`。
  - 目前 OpenAI-Compatible 实现补齐 tools 参数。
- **AgentRun 增强**
  - 每一步 tool 调用记录一条子事件（`tool_name / arguments / result / latency_ms`），构成完整调用链。

### 主要新增 API（示意）

- `POST /v1/tools`、`GET /v1/tools`、…
- `POST /v1/agents/:id/tools`（绑定 Tool）
- Chat 保持不变，多轮 tool-use 由平台内部完成。

### 里程碑

- [ ] Tool CRUD + JSON Schema 校验
- [ ] Agent-Tool 绑定
- [ ] OpenAI Chat Completions `tools` 打通
- [ ] `ToolExecutor` 多轮循环
- [ ] AgentRun 记录多步 Trace

---

## 后期规划（未定版本号）

以下能力属于"平台化"演进，V3 之后视场景启动，暂不承诺时间点。

- **AgentService（对应 easyai Worker / Model）**：把"一次 Agent"抽象成可独立部署的 Service，具备扩缩容、隔离、版本管理能力。可能会以 K8s Deployment / CRD 形式落地。
- **小 Operator（对应 eai-controller-manager 的一个子集）**：不做完整 controller-manager，只针对 AgentService 做一个薄 Operator，负责生命周期管理。
- **多租户 / 鉴权 / 审计**：随着资源类型增多，安全能力自然长出来。
- **可观测性**：指标（Prometheus）+ 分布式 Trace（OpenTelemetry）+ 面板。

**明确不做**：`eai-apiserver` 这类"自定义 API Server 抽象"—— 直接沿用当前的 Gin + Wire 已足够轻量，不引入 K8s API Server 风格的抽象层。

---

## 版本节奏与验收标准

| 版本 | 主题 | 验收标准 |
|---|---|---|
| V1 | 最小闭环 | 能在 5 分钟内配好 provider + agent，通过 REST 完成多轮对话并在 `AgentRun` 中看到 Trace |
| V2 | RAG    | 上传一份 PDF，Agent 能引用其中内容回答问题，Trace 能看到命中的 chunk |
| V3 | Tool   | 定义一个天气查询 Tool，Agent 能自主调用并返回结构化结果，Trace 能看到多步调用链 |

