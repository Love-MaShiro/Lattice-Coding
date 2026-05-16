package rerank

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type Reranker interface {
	Name() string
	Rerank(ctx context.Context, req domain.RerankRequest) (*domain.RerankResult, error)
}
