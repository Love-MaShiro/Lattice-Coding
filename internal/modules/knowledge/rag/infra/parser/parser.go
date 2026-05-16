package parser

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type Parser interface {
	Name() string
	Supports(documentType domain.DocumentType) bool
	Parse(ctx context.Context, document domain.Document, payload []byte) (*domain.ParsedDocument, error)
}
