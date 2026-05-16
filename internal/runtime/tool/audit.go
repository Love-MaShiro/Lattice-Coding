package tool

import (
	"context"
	"time"
)

type AuditEvent struct {
	EventType  string                 `json:"event_type"`
	Request    ToolRequest            `json:"request"`
	Descriptor ToolDescriptor         `json:"descriptor"`
	Result     ToolResult             `json:"result"`
	StartedAt  time.Time              `json:"started_at"`
	FinishedAt time.Time              `json:"finished_at"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type AuditRecorder interface {
	Record(ctx context.Context, event AuditEvent) error
}

type AuditRecorderFunc func(ctx context.Context, event AuditEvent) error

func (f AuditRecorderFunc) Record(ctx context.Context, event AuditEvent) error {
	return f(ctx, event)
}

type NoopAuditRecorder struct{}

func (NoopAuditRecorder) Record(context.Context, AuditEvent) error {
	return nil
}
