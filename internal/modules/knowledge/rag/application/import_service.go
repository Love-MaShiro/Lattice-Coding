package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/rag/domain"
)

type ImportCommand struct {
	Source   domain.Source
	Document domain.Document
	Payload  []byte
	Metadata string
}

type ImportResult struct {
	Source   domain.Source
	Document domain.Document
}

type ImportService interface {
	Import(ctx context.Context, cmd ImportCommand) (*ImportResult, error)
}
