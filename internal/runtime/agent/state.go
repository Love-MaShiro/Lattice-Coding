package agent

import (
	"time"
)

const (
	StrategyReAct = "react"
)

type Request struct {
	RunID         string
	AgentID       string
	UserID        string
	SessionID     string
	ProjectID     string
	NodeID        string
	Input         string
	Strategy      string
	Provider      string
	Model         string
	ModelConfigID uint64
	AllowedTools  []string
	WorkingDir    string
	MaxIterations int
	MaxToolCalls  int
	Timeout       time.Duration
	Metadata      map[string]interface{}
	Query         map[string]interface{}
}

type Result struct {
	RunID    string
	Content  string
	Messages []Message
	Metadata map[string]interface{}
}

type Message struct {
	Role    string
	Content string
}

type ReActStep struct {
	Iteration   int                    `json:"iteration"`
	Reason      string                 `json:"reason,omitempty"`
	Action      string                 `json:"action,omitempty"`
	ActionInput map[string]interface{} `json:"action_input,omitempty"`
	Observation string                 `json:"observation,omitempty"`
	IsError     bool                   `json:"is_error"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt time.Time              `json:"completed_at"`
}

type state struct {
	request   Request
	steps     []ReActStep
	toolCalls int
}

func (s *state) appendStep(step ReActStep) {
	s.steps = append(s.steps, step)
	if step.Action != "" && step.Action != ReActActionFinal && step.Action != "parse_react_action" {
		s.toolCalls++
	}
}
