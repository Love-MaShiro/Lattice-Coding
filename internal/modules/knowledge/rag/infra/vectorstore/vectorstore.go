package vectorstore

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type SearchRequest struct {
	QueryVector []float64
	SourceIDs   []uint64
	Limit       int
	Filters     string
}

type VectorStore interface {
	Upsert(ctx context.Context, embeddings []domain.Embedding) error
	Search(ctx context.Context, req SearchRequest) ([]domain.Evidence, error)
	DeleteByChunkIDs(ctx context.Context, chunkIDs []uint64) error
}
