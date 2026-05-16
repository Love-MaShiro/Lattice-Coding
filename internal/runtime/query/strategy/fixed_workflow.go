package strategy

import "lattice-coding/internal/runtime/query"

func NewFixedWorkflowStrategy() NotImplementedStrategy {
	return NewNotImplementedStrategy(query.ExecutionModeFixedWorkflow)
}
