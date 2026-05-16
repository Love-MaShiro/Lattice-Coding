package embedding

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type Embedder interface {
	Model() string
	Embed(ctx context.Context, chunks []domain.Chunk) ([]domain.Embedding, error)
}
