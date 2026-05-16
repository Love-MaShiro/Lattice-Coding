package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type ParseCommand struct {
	Document domain.Document
	Payload  []byte
	Metadata string
}

type ParserService interface {
	Parse(ctx context.Context, cmd ParseCommand) (*domain.ParsedDocument, error)
}
