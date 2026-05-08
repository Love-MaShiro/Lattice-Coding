-- 003_chat.sql
-- Chat 会话与消息模块

DROP TABLE IF EXISTS chat_message;
DROP TABLE IF EXISTS chat_session;

CREATE TABLE chat_session (
    id                           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    title                        VARCHAR(200) NOT NULL DEFAULT '',
    agent_id                     BIGINT UNSIGNED NOT NULL,
    model_config_id              BIGINT UNSIGNED NOT NULL,
    status                       VARCHAR(20) NOT NULL DEFAULT 'active',
    summary                      LONGTEXT,
    summarized_until_message_id  BIGINT UNSIGNED NOT NULL DEFAULT 0,
    meta                         JSON,
    created_at                   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at                   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at                   DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_model_config_id (model_config_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_agent_deleted_at (agent_id, deleted_at),
    INDEX idx_status_deleted_at (status, deleted_at),
    INDEX idx_summarized_until_message_id (summarized_until_message_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE chat_message (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    session_id    BIGINT UNSIGNED NOT NULL,
    role          VARCHAR(20) NOT NULL,
    content       LONGTEXT NOT NULL,
    token_count   INT NOT NULL DEFAULT 0,
    meta          JSON,
    created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at    DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_session_id (session_id),
    INDEX idx_role (role),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_session_deleted_at (session_id, deleted_at),
    INDEX idx_session_id_deleted_at (session_id, id, deleted_at),
    INDEX idx_session_role_deleted_at (session_id, role, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
