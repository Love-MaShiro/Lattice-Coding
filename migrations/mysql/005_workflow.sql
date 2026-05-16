-- 005_workflow.sql
-- Workflow definitions, nodes, and edges.

CREATE TABLE IF NOT EXISTS workflow (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name         VARCHAR(200) NOT NULL,
    description  TEXT,
    status       VARCHAR(20) NOT NULL DEFAULT 'draft',
    version      INT NOT NULL DEFAULT 1,
    meta         JSON,
    created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at   DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_name (name),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_status_deleted_at (status, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS workflow_node (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    workflow_id  BIGINT UNSIGNED NOT NULL,
    node_key     VARCHAR(120) NOT NULL,
    name         VARCHAR(200) NOT NULL DEFAULT '',
    type         VARCHAR(50) NOT NULL,
    config       JSON,
    position     JSON,
    sort_order   INT NOT NULL DEFAULT 0,
    meta         JSON,
    created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at   DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_workflow_id (workflow_id),
    INDEX idx_node_key (node_key),
    INDEX idx_type (type),
    INDEX idx_sort_order (sort_order),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_workflow_deleted_at (workflow_id, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS workflow_edge (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    workflow_id     BIGINT UNSIGNED NOT NULL,
    edge_key        VARCHAR(160) NOT NULL,
    source_key      VARCHAR(120) NOT NULL,
    target_key      VARCHAR(120) NOT NULL,
    condition_expr  TEXT,
    sort_order      INT NOT NULL DEFAULT 0,
    meta            JSON,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_workflow_id (workflow_id),
    INDEX idx_edge_key (edge_key),
    INDEX idx_source_key (source_key),
    INDEX idx_target_key (target_key),
    INDEX idx_sort_order (sort_order),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_workflow_deleted_at (workflow_id, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
