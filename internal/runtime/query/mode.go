package query

type ExecutionMode string

const (
	ExecutionModeDirectChat    ExecutionMode = "direct_chat"
	ExecutionModeFixedWorkflow ExecutionMode = "fixed_workflow"
	ExecutionModePlanGraph     ExecutionMode = "plan_graph"
	ExecutionModePureReAct     ExecutionMode = "pure_react"
)

func (m ExecutionMode) String() string {
	return string(m)
}
