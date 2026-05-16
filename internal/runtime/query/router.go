package query

type ExecutionModeRouter interface {
	Route(req QueryRequest) ExecutionMode
}

type DefaultExecutionModeRouter struct{}

func NewExecutionModeRouter() DefaultExecutionModeRouter {
	return DefaultExecutionModeRouter{}
}

func (DefaultExecutionModeRouter) Route(req QueryRequest) ExecutionMode {
	if req.Mode != "" {
		return req.Mode
	}
	if len(req.AllowedTools) > 0 {
		return ExecutionModePureReAct
	}
	return ExecutionModeDirectChat
}
