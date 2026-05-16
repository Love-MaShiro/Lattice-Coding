package domain

import "context"

type RunRepository interface {
	Create(ctx context.Context, run *Run) error
	Update(ctx context.Context, run *Run) error
	FindByID(ctx context.Context, id string) (*Run, error)
	FindPage(ctx context.Context, req PageRequest) (*PageResult[*Run], error)
}

type ToolInvocationRepository interface {
	Create(ctx context.Context, invocation *ToolInvocation) error
	Update(ctx context.Context, invocation *ToolInvocation) error
	FindByID(ctx context.Context, id string) (*ToolInvocation, error)
	FindByRunID(ctx context.Context, runID string) ([]*ToolInvocation, error)
}
