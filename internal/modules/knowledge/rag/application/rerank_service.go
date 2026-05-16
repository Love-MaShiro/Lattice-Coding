package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type RerankService interface {
	Rerank(ctx context.Context, req domain.RerankRequest) (*domain.RerankResult, error)
}
