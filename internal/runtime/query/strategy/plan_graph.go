package strategy

import "lattice-coding/internal/runtime/query"

func NewPlanGraphStrategy() NotImplementedStrategy {
	return NewNotImplementedStrategy(query.ExecutionModePlanGraph)
}
