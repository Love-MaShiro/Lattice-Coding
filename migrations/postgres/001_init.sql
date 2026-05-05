CREATE EXTENSION IF NOT EXISTS vector;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS knowledge_chunk (
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL,
    chunk_index INT NOT NULL,
    content TEXT NOT NULL,
    metadata JSONB,
    embedding vector(1536) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted SMALLINT NOT NULL DEFAULT 0
);

CREATE TRIGGER trg_knowledge_chunk_updated_at
BEFORE UPDATE ON knowledge_chunk
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_knowledge_chunk_doc ON knowledge_chunk (document_id, deleted, chunk_index);
CREATE INDEX idx_knowledge_chunk_created_at ON knowledge_chunk (created_at);
CREATE INDEX idx_knowledge_chunk_metadata ON knowledge_chunk USING gin (metadata);
CREATE INDEX idx_knowledge_chunk_embedding_ivfflat ON knowledge_chunk
USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
