package strategy

import (
	"context"

	"lattice-coding/internal/runtime/query"
)

type QueryStrategy = query.QueryStrategy

type NotImplementedStrategy struct {
	mode query.ExecutionMode
}

func NewNotImplementedStrategy(mode query.ExecutionMode) NotImplementedStrategy {
	return NotImplementedStrategy{mode: mode}
}

func (s NotImplementedStrategy) Mode() query.ExecutionMode {
	return s.mode
}

func (s NotImplementedStrategy) Execute(context.Context, *query.QueryState) (*query.QueryResult, error) {
	return nil, query.ErrModeNotSupported.WithMessage("query strategy not implemented: " + s.mode.String())
}
