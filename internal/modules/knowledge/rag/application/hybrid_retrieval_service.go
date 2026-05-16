package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type HybridRetrievalService interface {
	RetrieveHybrid(ctx context.Context, req domain.RetrievalRequest) (*domain.RetrievalResult, error)
}
