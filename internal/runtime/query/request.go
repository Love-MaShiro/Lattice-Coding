package query

import "time"

type QueryRequest struct {
	RunID         string
	AgentID       string
	UserID        string
	SessionID     string
	ProjectID     string
	NodeID        string
	Input         string
	Mode          ExecutionMode
	Provider      string
	Model         string
	ModelConfigID uint64
	AllowedTools  []string
	WorkingDir    string
	Budget        QueryBudget
	Timeout       time.Duration
	Metadata      map[string]interface{}
}

type Request = QueryRequest
