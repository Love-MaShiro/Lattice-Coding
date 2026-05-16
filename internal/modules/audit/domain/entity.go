package domain

import (
	"context"
	"time"
)

type AuditLog struct {
	ID           uint64
	RunID        string
	TraceID      string
	EventType    string
	ToolName     string
	ResourceType string
	ResourceID   string
	Message      string
	PayloadJSON  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
}
