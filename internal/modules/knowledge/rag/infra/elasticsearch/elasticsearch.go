package elasticsearch

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type SearchRequest struct {
	Query     string
	IndexName string
	SourceIDs []uint64
	Limit     int
	Filters   string
}

type KeywordStore interface {
	Index(ctx context.Context, chunks []domain.Chunk, indexName string) error
	Search(ctx context.Context, req SearchRequest) ([]domain.Evidence, error)
	DeleteByChunkIDs(ctx context.Context, indexName string, chunkIDs []uint64) error
}
