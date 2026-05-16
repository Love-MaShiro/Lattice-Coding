CREATE TABLE IF NOT EXISTS providers (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    provider_type VARCHAR(50) NOT NULL,
    base_url VARCHAR(255),
    auth_type VARCHAR(50) NOT NULL DEFAULT 'api_key',
    api_key_ciphertext TEXT,
    auth_config_ciphertext TEXT,
    config TEXT,
    enabled TINYINT NOT NULL DEFAULT 1,
    health_status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    last_checked_at DATETIME(3),
    last_error TEXT,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    KEY idx_provider_name (name),
    KEY idx_provider_type (provider_type),
    KEY idx_provider_enabled (enabled),
    KEY idx_provider_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS model_configs (
    id BIGINT NOT NULL AUTO_INCREMENT,
    provider_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    model_type VARCHAR(50) NOT NULL DEFAULT 'chat',
    params TEXT,
    capabilities TEXT,
    is_default TINYINT NOT NULL DEFAULT 0,
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    KEY idx_model_config_provider (provider_id),
    KEY idx_model_config_model (model),
    KEY idx_model_config_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS provider_healths (
    id BIGINT NOT NULL AUTO_INCREMENT,
    provider_id BIGINT NOT NULL,
    model_config_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    latency_ms BIGINT NOT NULL DEFAULT 0,
    error_code VARCHAR(50),
    error_message TEXT,
    checked_at DATETIME(3) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_provider_health_provider (provider_id),
    KEY idx_provider_health_model_config (model_config_id),
    KEY idx_provider_health_checked_at (checked_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS agents (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    agent_type VARCHAR(50) NOT NULL DEFAULT 'customer_service',
    model_config_id BIGINT NOT NULL,
    system_prompt TEXT,
    temperature DECIMAL(3,2) NOT NULL DEFAULT 0.70,
    top_p DECIMAL(3,2) NOT NULL DEFAULT 1.00,
    max_tokens INT NOT NULL DEFAULT 4096,
    max_context_turns INT NOT NULL DEFAULT 10,
    max_steps INT NOT NULL DEFAULT 20,
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    KEY idx_agent_name (name),
    KEY idx_agent_type (agent_type),
    KEY idx_agent_model_config (model_config_id),
    KEY idx_agent_enabled (enabled),
    KEY idx_agent_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS agent_tool (
    id BIGINT NOT NULL AUTO_INCREMENT,
    agent_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    tool_type VARCHAR(50) NOT NULL,
    config JSON,
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uk_agent_tool (agent_id, name, deleted),
    KEY idx_agent_tool_agent (agent_id, deleted),
    KEY idx_agent_tool_type (tool_type, deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mcp_server (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    base_url VARCHAR(500) NOT NULL,
    transport VARCHAR(20) NOT NULL DEFAULT 'http',
    auth JSON,
    config JSON,
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uk_mcp_server_name (name, deleted),
    KEY idx_mcp_server_enabled (enabled, deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_session (
    id BIGINT NOT NULL AUTO_INCREMENT,
    title VARCHAR(200),
    agent_id BIGINT,
    model_config_id BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    meta JSON,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_chat_session_agent (agent_id, deleted),
    KEY idx_chat_session_model_config (model_config_id, deleted),
    KEY idx_chat_session_status (status, deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_message (
    id BIGINT NOT NULL AUTO_INCREMENT,
    session_id BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL,
    content LONGTEXT NOT NULL,
    token_count INT NOT NULL DEFAULT 0,
    meta JSON,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_chat_message_session (session_id, deleted),
    KEY idx_chat_message_role (role, deleted),
    KEY idx_chat_message_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS run (
    id BIGINT NOT NULL AUTO_INCREMENT,
    run_id VARCHAR(64) NOT NULL,
    agent_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    input LONGTEXT,
    output LONGTEXT,
    error LONGTEXT,
    started_at DATETIME(3),
    completed_at DATETIME(3),
    duration_ms INT,
    meta JSON,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uk_run_run_id (run_id),
    KEY idx_run_agent (agent_id, deleted),
    KEY idx_run_status (status, deleted),
    KEY idx_run_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS knowledge_document (
    id BIGINT NOT NULL AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    source_uri VARCHAR(500),
    content_hash VARCHAR(64),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    chunk_count INT NOT NULL DEFAULT 0,
    meta JSON,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uk_knowledge_doc_hash (content_hash, deleted),
    KEY idx_knowledge_doc_status (status, deleted),
    KEY idx_knowledge_doc_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS safety_check (
    id BIGINT NOT NULL AUTO_INCREMENT,
    check_type VARCHAR(30) NOT NULL,
    input TEXT NOT NULL,
    allowed TINYINT NOT NULL DEFAULT 0,
    reason VARCHAR(500),
    meta JSON,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_safety_check_type (check_type, deleted),
    KEY idx_safety_check_allowed (allowed, deleted),
    KEY idx_safety_check_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS audit_log (
    id BIGINT NOT NULL AUTO_INCREMENT,
    trace_id VARCHAR(80),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id BIGINT,
    request_method VARCHAR(10),
    request_path VARCHAR(255),
    status_code INT,
    message VARCHAR(500),
    extra JSON,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_audit_log_trace (trace_id),
    KEY idx_audit_log_action (action, created_at),
    KEY idx_audit_log_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
