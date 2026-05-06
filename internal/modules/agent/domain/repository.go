package domain

import "context"

type AgentRepository interface {
	Create(ctx context.Context, agent *Agent) error
	Update(ctx context.Context, agent *Agent) error
	FindByID(ctx context.Context, id uint64) (*Agent, error)
	FindPage(ctx context.Context, req *PageRequest) (*PageResult[*Agent], error)
	DeleteByID(ctx context.Context, id uint64) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	UpdateEnabled(ctx context.Context, id uint64, enabled bool) error
}

type AgentReferenceChecker interface {
	HasModelConfigReferences(ctx context.Context, modelConfigID uint64) (bool, error)
}
