package domain

import "context"

type PackRepository interface {
	CreateWithItems(ctx context.Context, pack *KnowledgePack) error
	FindByIDWithItems(ctx context.Context, id uint64) (*KnowledgePack, error)
	FindByKeyWithItems(ctx context.Context, packKey string) (*KnowledgePack, error)
	DeleteByID(ctx context.Context, id uint64) error
}
