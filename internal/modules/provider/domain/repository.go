package domain

import "context"

type ProviderRepository interface {
	Create(ctx context.Context, provider *Provider) error
	Update(ctx context.Context, provider *Provider) error
	FindByID(ctx context.Context, id uint64) (*Provider, error)
	FindPage(ctx context.Context, req *PageRequest) (*PageResult[*Provider], error)
	DeleteByID(ctx context.Context, id uint64) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	UpdateEnabled(ctx context.Context, id uint64, enabled bool) error
	UpdateHealthStatus(ctx context.Context, id uint64, status HealthStatus, lastError string) error
}

type ModelConfigRepository interface {
	Create(ctx context.Context, modelConfig *ModelConfig) error
	Update(ctx context.Context, modelConfig *ModelConfig) error
	FindByID(ctx context.Context, id uint64) (*ModelConfig, error)
	FindByProviderID(ctx context.Context, providerID uint64) ([]*ModelConfig, error)
	FindPage(ctx context.Context, req *PageRequest) (*PageResult[*ModelConfig], error)
	DeleteByID(ctx context.Context, id uint64) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	UpdateEnabled(ctx context.Context, id uint64, enabled bool) error
}

type ProviderHealthRepository interface {
	Create(ctx context.Context, health *ProviderHealth) error
	FindLatestByProviderID(ctx context.Context, providerID uint64) (*ProviderHealth, error)
	FindLatestByModelConfigID(ctx context.Context, modelConfigID uint64) (*ProviderHealth, error)
}
