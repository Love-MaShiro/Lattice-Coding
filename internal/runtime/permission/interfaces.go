package permission

import "context"

type Decision string

const (
	DecisionAllow Decision = "allow"
	DecisionDeny  Decision = "deny"
	DecisionAsk   Decision = "ask"
)

type Checker interface {
	Check(ctx context.Context, req Request) (*Result, error)
}

type Policy interface {
	Evaluate(ctx context.Context, req Request) (*Result, error)
}

type Request struct {
	RunID      string
	AgentID    string
	UserID     string
	ToolName   string
	Input      map[string]interface{}
	WorkingDir string
	Metadata   map[string]interface{}
}

type Result struct {
	Decision Decision
	Reason   string
	Metadata map[string]interface{}
}
