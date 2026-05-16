package tool

import (
	"context"
	"time"
)

type ToolInvocationStarted struct {
	RunID     string                 `json:"run_id,omitempty"`
	NodeID    string                 `json:"node_id,omitempty"`
	ToolName  string                 `json:"tool_name"`
	Input     map[string]interface{} `json:"input,omitempty"`
	StartedAt time.Time              `json:"started_at"`
}

type ToolInvocationFinished struct {
	ID            string     `json:"id,omitempty"`
	RunID         string     `json:"run_id,omitempty"`
	NodeID        string     `json:"node_id,omitempty"`
	ToolName      string     `json:"tool_name"`
	Result        ToolResult `json:"result"`
	Status        string     `json:"status"`
	CompletedAt   time.Time  `json:"completed_at"`
	LatencyMs     int64      `json:"latency_ms"`
	FullResultRef string     `json:"full_result_ref,omitempty"`
}

type ToolInvocationRecorder interface {
	Start(ctx context.Context, record ToolInvocationStarted) (string, error)
	Finish(ctx context.Context, record ToolInvocationFinished) error
}

type NoopToolInvocationRecorder struct{}

func (NoopToolInvocationRecorder) Start(context.Context, ToolInvocationStarted) (string, error) {
	return "", nil
}

func (NoopToolInvocationRecorder) Finish(context.Context, ToolInvocationFinished) error {
	return nil
}

type ToolInvocationRecorderFunc struct {
	StartFunc  func(ctx context.Context, record ToolInvocationStarted) (string, error)
	FinishFunc func(ctx context.Context, record ToolInvocationFinished) error
}

func (f ToolInvocationRecorderFunc) Start(ctx context.Context, record ToolInvocationStarted) (string, error) {
	if f.StartFunc == nil {
		return "", nil
	}
	return f.StartFunc(ctx, record)
}

func (f ToolInvocationRecorderFunc) Finish(ctx context.Context, record ToolInvocationFinished) error {
	if f.FinishFunc == nil {
		return nil
	}
	return f.FinishFunc(ctx, record)
}
