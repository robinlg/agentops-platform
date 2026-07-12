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