package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type EmbedCommand struct {
	Chunks   []domain.Chunk
	Model    string
	Metadata string
}

type EmbeddingService interface {
	Embed(ctx context.Context, cmd EmbedCommand) ([]domain.Embedding, error)
}
