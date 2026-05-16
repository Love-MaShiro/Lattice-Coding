package domain

import "time"

const (
	RunStatusRunning   = "running"
	RunStatusCompleted = "completed"
	RunStatusFailed    = "failed"

	ToolInvocationStatusStarted = "started"
	ToolInvocationStatusSuccess = "success"
	ToolInvocationStatusFailed  = "failed"
)

type Run struct {
	ID          string
	AgentID     string
	SessionID   string
	WorkflowID  string
	Status      string
	Input       string
	Output      string
	Error       string
	StartedAt   time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ToolInvocation struct {
	ID            string
	RunID         string
	NodeID        string
	ToolName      string
	InputJSON     string
	OutputJSON    string
	IsError       bool
	LatencyMs     int64
	Status        string
	FullResultRef string
	StartedAt     time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PageRequest struct {
	Page     int
	PageSize int
}

type PageResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}
