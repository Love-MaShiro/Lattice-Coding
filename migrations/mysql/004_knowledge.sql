-- 004_knowledge.sql
-- Knowledge source metadata and retrieval trace tables.

CREATE TABLE IF NOT EXISTS knowledge_source (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name          VARCHAR(200) NOT NULL,
    source_type   VARCHAR(40) NOT NULL,
    description   TEXT,
    uri           VARCHAR(1000) NOT NULL DEFAULT '',
    owner         VARCHAR(100) NOT NULL DEFAULT '',
    status        VARCHAR(20) NOT NULL DEFAULT 'active',
    config        JSON,
    meta          JSON,
    created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at    DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_name (name),
    INDEX idx_source_type (source_type),
    INDEX idx_owner (owner),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_source_type_status_deleted_at (source_type, status, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS knowledge_document (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    source_id      BIGINT UNSIGNED NOT NULL,
    title          VARCHAR(500) NOT NULL,
    document_type  VARCHAR(40) NOT NULL DEFAULT '',
    uri            VARCHAR(1000) NOT NULL DEFAULT '',
    author         VARCHAR(200) NOT NULL DEFAULT '',
    version        VARCHAR(100) NOT NULL DEFAULT '',
    summary        TEXT,
    content_hash   VARCHAR(128) NOT NULL DEFAULT '',
    status         VARCHAR(20) NOT NULL DEFAULT 'pending',
    meta           JSON,
    created_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at     DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_source_id (source_id),
    INDEX idx_title (title),
    INDEX idx_document_type (document_type),
    INDEX idx_content_hash (content_hash),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_source_status_deleted_at (source_id, status, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE knowledge_document ADD COLUMN IF NOT EXISTS source_id BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE knowledge_document ADD COLUMN IF NOT EXISTS document_type VARCHAR(40) NOT NULL DEFAULT '';
ALTER TABLE knowledge_document ADD COLUMN IF NOT EXISTS uri VARCHAR(1000) NOT NULL DEFAULT '';
ALTER TABLE knowledge_document ADD COLUMN IF NOT EXISTS author VARCHAR(200) NOT NULL DEFAULT '';
ALTER TABLE knowledge_document ADD COLUMN IF NOT EXISTS version VARCHAR(100) NOT NULL DEFAULT '';
ALTER TABLE knowledge_document ADD COLUMN IF NOT EXISTS summary TEXT;
ALTER TABLE knowledge_document MODIFY COLUMN content_hash VARCHAR(128) NOT NULL DEFAULT '';
CREATE INDEX idx_knowledge_document_source_id ON knowledge_document (source_id);
CREATE INDEX idx_knowledge_document_document_type ON knowledge_document (document_type);

CREATE TABLE IF NOT EXISTS retrieval_trace (
    id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    query             TEXT NOT NULL,
    intent            VARCHAR(80) NOT NULL DEFAULT '',
    route             VARCHAR(40) NOT NULL,
    confidence        DOUBLE NOT NULL DEFAULT 0,
    reason            TEXT,
    evidence_count    INT NOT NULL DEFAULT 0,
    selected_sources  JSON,
    meta              JSON,
    created_at        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at        DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_intent (intent),
    INDEX idx_route (route),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_route_created_at (route, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS structured_query_tool (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name            VARCHAR(120) NOT NULL,
    description     TEXT,
    domain          VARCHAR(80) NOT NULL DEFAULT '',
    source_id       BIGINT UNSIGNED NOT NULL DEFAULT 0,
    db_type         VARCHAR(30) NOT NULL DEFAULT 'mysql',
    sql_template    TEXT NOT NULL,
    params_schema   JSON,
    result_schema   JSON,
    enabled         TINYINT(1) NOT NULL DEFAULT 1,
    read_only       TINYINT(1) NOT NULL DEFAULT 1,
    timeout_ms      INT NOT NULL DEFAULT 3000,
    max_rows        INT NOT NULL DEFAULT 100,
    meta            JSON,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_structured_tool_name (name),
    INDEX idx_domain (domain),
    INDEX idx_source_id (source_id),
    INDEX idx_enabled (enabled),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_domain_enabled_deleted_at (domain, enabled, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS context_source (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    source_key   VARCHAR(120) NOT NULL,
    kind         VARCHAR(50) NOT NULL,
    name         VARCHAR(200) NOT NULL DEFAULT '',
    uri          VARCHAR(1000) NOT NULL DEFAULT '',
    scope        VARCHAR(200) NOT NULL DEFAULT '',
    metadata     JSON,
    created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at   DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_context_source_key (source_key),
    INDEX idx_kind (kind),
    INDEX idx_name (name),
    INDEX idx_scope (scope),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_kind_scope_deleted_at (kind, scope, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS context_candidate (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    candidate_key   VARCHAR(120) NOT NULL,
    source_key      VARCHAR(120) NOT NULL,
    source_kind     VARCHAR(50) NOT NULL,
    title           VARCHAR(500) NOT NULL DEFAULT '',
    content         LONGTEXT,
    location        VARCHAR(500) NOT NULL DEFAULT '',
    score           DOUBLE NOT NULL DEFAULT 0,
    token_estimate  INT NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    metadata        JSON,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_context_candidate_key (candidate_key),
    INDEX idx_source_key (source_key),
    INDEX idx_source_kind (source_kind),
    INDEX idx_score (score),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_source_status_deleted_at (source_key, status, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS context_signal (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    candidate_id  BIGINT UNSIGNED NOT NULL,
    signal_key    VARCHAR(120) NOT NULL,
    kind          VARCHAR(50) NOT NULL,
    weight        DOUBLE NOT NULL DEFAULT 0,
    reason        TEXT,
    metadata      JSON,
    created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at    DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_candidate_id (candidate_id),
    INDEX idx_signal_key (signal_key),
    INDEX idx_kind (kind),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_candidate_kind_deleted_at (candidate_id, kind, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS context_policy (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    policy_key   VARCHAR(120) NOT NULL,
    name         VARCHAR(200) NOT NULL DEFAULT '',
    description  TEXT,
    max_tokens   INT NOT NULL DEFAULT 0,
    max_items    INT NOT NULL DEFAULT 0,
    rules        JSON,
    metadata     JSON,
    created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at   DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_context_policy_key (policy_key),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS knowledge_pack (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    pack_key        VARCHAR(120) NOT NULL,
    query           TEXT,
    intent          VARCHAR(80) NOT NULL DEFAULT '',
    route           VARCHAR(80) NOT NULL DEFAULT '',
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    token_estimate  INT NOT NULL DEFAULT 0,
    max_tokens      INT NOT NULL DEFAULT 0,
    prompt_context  LONGTEXT,
    warnings        JSON,
    options         JSON,
    meta            JSON,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_knowledge_pack_key (pack_key),
    INDEX idx_intent (intent),
    INDEX idx_route (route),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_route_status_deleted_at (route, status, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS knowledge_item (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    pack_id         BIGINT UNSIGNED NOT NULL,
    item_key        VARCHAR(120) NOT NULL,
    source_kind     VARCHAR(50) NOT NULL,
    source_id       VARCHAR(200) NOT NULL DEFAULT '',
    source_type     VARCHAR(80) NOT NULL DEFAULT '',
    title           VARCHAR(500) NOT NULL DEFAULT '',
    content         LONGTEXT,
    location        VARCHAR(500) NOT NULL DEFAULT '',
    score           DOUBLE NOT NULL DEFAULT 0,
    token_estimate  INT NOT NULL DEFAULT 0,
    citation_key    VARCHAR(120) NOT NULL DEFAULT '',
    metadata        JSON,
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at      DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_pack_id (pack_id),
    INDEX idx_item_key (item_key),
    INDEX idx_source_kind (source_kind),
    INDEX idx_source_id (source_id),
    INDEX idx_source_type (source_type),
    INDEX idx_citation_key (citation_key),
    INDEX idx_sort_order (sort_order),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_pack_sort_deleted_at (pack_id, sort_order, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS knowledge_citation (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    pack_id        BIGINT UNSIGNED NOT NULL,
    citation_key   VARCHAR(120) NOT NULL,
    source_kind    VARCHAR(50) NOT NULL,
    source_id      VARCHAR(200) NOT NULL DEFAULT '',
    title          VARCHAR(500) NOT NULL DEFAULT '',
    location       VARCHAR(500) NOT NULL DEFAULT '',
    uri            VARCHAR(1000) NOT NULL DEFAULT '',
    score          DOUBLE NOT NULL DEFAULT 0,
    metadata       JSON,
    sort_order     INT NOT NULL DEFAULT 0,
    created_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at     DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_pack_id (pack_id),
    INDEX idx_citation_key (citation_key),
    INDEX idx_source_kind (source_kind),
    INDEX idx_source_id (source_id),
    INDEX idx_sort_order (sort_order),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_pack_sort_deleted_at (pack_id, sort_order, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
