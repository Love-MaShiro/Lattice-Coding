package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type EvidenceGateService interface {
	Evaluate(ctx context.Context, req domain.EvidenceGateRequest) (*domain.EvidenceGateResult, error)
}
