package domain

import "context"

type SourceRepository interface {
	Create(ctx context.Context, source *Source) error
	FindByID(ctx context.Context, id uint64) (*Source, error)
	FindByKey(ctx context.Context, sourceKey string) (*Source, error)
	Update(ctx context.Context, source *Source) error
	DeleteByID(ctx context.Context, id uint64) error
}

type DocumentRepository interface {
	Create(ctx context.Context, document *Document) error
	FindByID(ctx context.Context, id uint64) (*Document, error)
	FindByKey(ctx context.Context, documentKey string) (*Document, error)
	Update(ctx context.Context, document *Document) error
	DeleteByID(ctx context.Context, id uint64) error
}

type ParsedDocumentRepository interface {
	Create(ctx context.Context, parsed *ParsedDocument) error
	FindByDocumentID(ctx context.Context, documentID uint64) (*ParsedDocument, error)
}

type ChunkRepository interface {
	BatchCreate(ctx context.Context, chunks []Chunk) error
	FindByDocumentID(ctx context.Context, documentID uint64) ([]Chunk, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]Chunk, error)
	DeleteByDocumentID(ctx context.Context, documentID uint64) error
}

type EmbeddingRepository interface {
	BatchCreate(ctx context.Context, embeddings []Embedding) error
	FindByChunkIDs(ctx context.Context, chunkIDs []uint64) ([]Embedding, error)
	UpdateStatus(ctx context.Context, id uint64, status EmbeddingStatus, errorMessage string) error
}

type KeywordIndexRepository interface {
	BatchCreate(ctx context.Context, indexes []KeywordIndex) error
	FindByChunkIDs(ctx context.Context, chunkIDs []uint64) ([]KeywordIndex, error)
	UpdateStatus(ctx context.Context, id uint64, status KeywordIndexStatus, errorMessage string) error
}

type EvidenceRepository interface {
	BatchCreate(ctx context.Context, evidence []Evidence) error
	FindByQuery(ctx context.Context, query string, limit int) ([]Evidence, error)
}

type RetrievalTraceRepository interface {
	Create(ctx context.Context, trace *RetrievalTrace) error
}

type RerankResultRepository interface {
	Create(ctx context.Context, result *RerankResult) error
}

type EvidenceGateRepository interface {
	Create(ctx context.Context, result *EvidenceGateResult) error
}
