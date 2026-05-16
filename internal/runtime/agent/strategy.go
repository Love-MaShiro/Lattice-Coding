package agent

import (
	"context"
)

type Runtime interface {
	Run(ctx context.Context, req Request) (*Result, error)
	Strategy(name string) (ExecutionStrategy, bool)
}

type ExecutionStrategy interface {
	Name() string
	Execute(ctx context.Context, req Request) (*Result, error)
}

type DirectStrategy interface {
	ExecutionStrategy
}

type FunctionCallingStrategy interface {
	ExecutionStrategy
}

type ReActStrategy interface {
	ExecutionStrategy
}
