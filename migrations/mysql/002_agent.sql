-- 002_agent.sql
-- Agent 配置管理模块

DROP TABLE IF EXISTS agent_tool;
DROP TABLE IF EXISTS agent;

CREATE TABLE agent (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name             VARCHAR(100) NOT NULL,
    description      TEXT,
    agent_type       VARCHAR(50) NOT NULL DEFAULT 'customer_service',
    model_config_id  BIGINT UNSIGNED NOT NULL,
    system_prompt    TEXT,
    temperature      DECIMAL(3,2) NOT NULL DEFAULT 0.70,
    top_p            DECIMAL(3,2) NOT NULL DEFAULT 1.00,
    max_tokens       INT NOT NULL DEFAULT 4096,
    max_context_turns INT NOT NULL DEFAULT 10,
    max_steps        INT NOT NULL DEFAULT 20,
    enabled          TINYINT(1) NOT NULL DEFAULT 1,
    created_at       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at       DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_name (name),
    INDEX idx_agent_type (agent_type),
    INDEX idx_model_config_id (model_config_id),
    INDEX idx_enabled (enabled),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE agent_tool (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    agent_id   BIGINT UNSIGNED NOT NULL,
    tool_id    BIGINT UNSIGNED NOT NULL,
    tool_type  VARCHAR(50) NOT NULL DEFAULT '',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_tool_id (tool_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
