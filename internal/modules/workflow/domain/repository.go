package domain

import "context"

type WorkflowRepository interface {
	CreateWithGraph(ctx context.Context, workflow *WorkflowDefinition) error
	UpdateWithGraph(ctx context.Context, workflow *WorkflowDefinition) error
	FindByIDWithGraph(ctx context.Context, id uint64) (*WorkflowDefinition, error)
	FindPage(ctx context.Context, req *PageRequest) (*PageResult[*WorkflowDefinition], error)
	DeleteWithGraph(ctx context.Context, id uint64) error
}
