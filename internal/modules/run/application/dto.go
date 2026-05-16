package application

import "time"

type RunDTO struct {
	ID          string     `json:"id"`
	AgentID     string     `json:"agent_id,omitempty"`
	SessionID   string     `json:"session_id,omitempty"`
	WorkflowID  string     `json:"workflow_id,omitempty"`
	Status      string     `json:"status"`
	Input       string     `json:"input,omitempty"`
	Output      string     `json:"output,omitempty"`
	Error       string     `json:"error,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ToolInvocationDTO struct {
	ID            string     `json:"id"`
	RunID         string     `json:"run_id,omitempty"`
	NodeID        string     `json:"node_id,omitempty"`
	ToolName      string     `json:"tool_name"`
	InputJSON     string     `json:"input_json,omitempty"`
	OutputJSON    string     `json:"output_json,omitempty"`
	IsError       bool       `json:"is_error"`
	LatencyMs     int64      `json:"latency_ms"`
	Status        string     `json:"status"`
	FullResultRef string     `json:"full_result_ref,omitempty"`
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type RunPageQuery struct {
	Page     int
	PageSize int
}

type PageResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}
