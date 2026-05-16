package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/pack/domain"
)

type BuildPackCommand struct {
	PackKey    string
	Query      string
	Intent     string
	Route      string
	Candidates []domain.PackCandidate
	Options    string
	Metadata   string
}

type PackBuilder interface {
	Build(ctx context.Context, cmd BuildPackCommand) (*domain.KnowledgePack, error)
}
