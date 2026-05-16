package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type ChunkCommand struct {
	Document domain.Document
	Parsed   domain.ParsedDocument
	Strategy domain.ChunkStrategy
	Metadata string
}

type ChunkingService interface {
	Chunk(ctx context.Context, cmd ChunkCommand) ([]domain.Chunk, error)
}
