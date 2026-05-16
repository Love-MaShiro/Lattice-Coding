package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type KeywordIndexCommand struct {
	Chunks    []domain.Chunk
	IndexName string
	Metadata  string
}

type KeywordIndexService interface {
	Index(ctx context.Context, cmd KeywordIndexCommand) ([]domain.KeywordIndex, error)
}
