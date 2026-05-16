package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type RetrievalService interface {
	Retrieve(ctx context.Context, req domain.RetrievalRequest) (*domain.RetrievalResult, error)
}
