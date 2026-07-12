CREATE TABLE model_providers
(
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(128) NOT NULL,
    provider_type VARCHAR(64)  NOT NULL,
    base_url      TEXT         NOT NULL,
    api_key       TEXT         NOT NULL,
    default_model VARCHAR(128) NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT
    ON TABLE model_providers IS '模型提供商表';
COMMENT
    ON COLUMN model_providers.name IS '模型提供商名称';
COMMENT
    ON COLUMN model_providers.provider_type IS '模型提供商类型';
COMMENT
    ON COLUMN model_providers.base_url IS '模型提供商URL';
COMMENT
    ON COLUMN model_providers.api_key IS '模型提供商API KEY';
COMMENT
    ON COLUMN model_providers.default_model IS '默认模型';
COMMENT
    ON COLUMN model_providers.created_at IS '创建时间';
COMMENT
    ON COLUMN model_providers.updated_at IS '更新时间';

CREATE TABLE agents
(
    id                BIGSERIAL PRIMARY KEY,
    name              VARCHAR(128) NOT NULL,
    description       TEXT,
    system_prompt     TEXT,
    model_provider_id BIGINT       NOT NULL,
    model             VARCHAR(128),
    temperature       DOUBLE PRECISION      DEFAULT 0.7,
    max_tokens        INT                   DEFAULT 2048,
    created_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT
    ON TABLE agents IS '智能体表';
COMMENT
    ON COLUMN agents.name IS '智能体名称';
COMMENT
    ON COLUMN agents.description IS '智能体描述';
COMMENT
    ON COLUMN agents.system_prompt IS '系统提示';
COMMENT
    ON COLUMN agents.model_provider_id IS '模型提供商ID';
COMMENT
    ON COLUMN agents.model IS '模型';
COMMENT
    ON COLUMN agents.temperature IS '热度';
COMMENT
    ON COLUMN agents.max_tokens IS '最大token';
COMMENT
    ON COLUMN agents.created_at IS '创建时间';
COMMENT
    ON COLUMN agents.updated_at IS '更新时间';

CREATE TABLE conversations
(
    id         BIGSERIAL PRIMARY KEY,
    agent_id   BIGINT       NOT NULL,
    title      VARCHAR(256),
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT
    ON TABLE conversations IS '会话表';
COMMENT
    ON COLUMN conversations.agent_id IS '智能体ID';
COMMENT
    ON COLUMN conversations.title IS '会话标题';
COMMENT
    ON COLUMN conversations.created_at IS '创建时间';
COMMENT
    ON COLUMN conversations.updated_at IS '更新时间';

CREATE TABLE messages
(
    id              BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT      NOT NULL,
    role            VARCHAR(32) NOT NULL,
    content         TEXT        NOT NULL,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT
    ON TABLE messages IS '消息表';
COMMENT
    ON COLUMN messages.conversation_id IS '会话ID';
COMMENT
    ON COLUMN messages.role IS '角色（user/assistant/system/tool）';
COMMENT
    ON COLUMN messages.content IS '消息内容';
COMMENT
    ON COLUMN messages.created_at IS '创建时间';

CREATE TABLE agent_runs
(
    id                BIGSERIAL PRIMARY KEY,
    agent_id          BIGINT       NOT NULL,
    conversation_id   BIGINT       NOT NULL,
    status            VARCHAR(32)  NOT NULL,
    input             TEXT,
    output            TEXT,
    model             VARCHAR(128),
    prompt_tokens     INT                   DEFAULT 0,
    completion_tokens INT                   DEFAULT 0,
    total_tokens      INT                   DEFAULT 0,
    latency_ms        BIGINT                DEFAULT 0,
    error_message     TEXT,
    started_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at       TIMESTAMP
);

COMMENT
    ON TABLE agent_runs IS '智能体运行记录表';
COMMENT
    ON COLUMN agent_runs.agent_id IS '智能体ID';
COMMENT
    ON COLUMN agent_runs.conversation_id IS '会话ID';
COMMENT
    ON COLUMN agent_runs.status IS '运行状态（pending/running/success/failed）';
COMMENT
    ON COLUMN agent_runs.input IS '输入内容';
COMMENT
    ON COLUMN agent_runs.output IS '输出内容';
COMMENT
    ON COLUMN agent_runs.model IS '使用的模型';
COMMENT
    ON COLUMN agent_runs.prompt_tokens IS '输入token数';
COMMENT
    ON COLUMN agent_runs.completion_tokens IS '输出token数';
COMMENT
    ON COLUMN agent_runs.total_tokens IS '总token数';
COMMENT
    ON COLUMN agent_runs.latency_ms IS '耗时（毫秒）';
COMMENT
    ON COLUMN agent_runs.error_message IS '错误信息';
COMMENT
    ON COLUMN agent_runs.started_at IS '开始时间';
COMMENT
    ON COLUMN agent_runs.finished_at IS '结束时间';