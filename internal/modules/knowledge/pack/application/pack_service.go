package application

import (
	"context"

	"lattice-coding/internal/modules/knowledge/pack/domain"
)

type PackService interface {
	CreatePack(ctx context.Context, pack *domain.KnowledgePack) (*domain.KnowledgePack, error)
	GetPack(ctx context.Context, id uint64) (*domain.KnowledgePack, error)
	GetPackByKey(ctx context.Context, packKey string) (*domain.KnowledgePack, error)
}
