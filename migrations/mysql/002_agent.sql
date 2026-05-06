-- 002_agent.sql
-- Agent 配置管理模块

DROP TABLE IF EXISTS t_agent;

CREATE TABLE agents (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name             VARCHAR(100) NOT NULL,
    description      TEXT,
    provider_id      BIGINT UNSIGNED NOT NULL,
    model_config_id  BIGINT UNSIGNED NOT NULL,
    system_prompt    TEXT,
    tools            JSON,
    max_steps        INT NOT NULL DEFAULT 20,
    timeout          INT NOT NULL DEFAULT 1200,
    enabled          TINYINT(1) NOT NULL DEFAULT 1,
    created_at       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at       DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_name (name),
    INDEX idx_provider_id (provider_id),
    INDEX idx_model_config_id (model_config_id),
    INDEX idx_enabled (enabled),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
