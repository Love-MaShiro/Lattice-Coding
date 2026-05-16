package query

import (
	"time"

	"lattice-coding/internal/runtime/llm"
)

type QueryRequest struct {
	RunID         string
	AgentID       string
	UserID        string
	SessionID     string
	ProjectID     string
	NodeID        string
	Input         string
	Messages      []llm.Message
	SystemPrompt  string
	Summary       string
	Mode          ExecutionMode
	Provider      string
	Model         string
	ModelConfigID uint64
	Temperature   *float64
	TopP          *float64
	MaxTokens     int
	AllowedTools  []string
	WorkingDir    string
	Variables     map[string]interface{}
	Stream        bool
	Budget        QueryBudget
	Timeout       time.Duration
	Metadata      map[string]interface{}
}

type Request = QueryRequest
