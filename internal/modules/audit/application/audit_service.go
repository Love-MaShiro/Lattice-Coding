package application

import (
	"context"
	"encoding/json"

	"lattice-coding/internal/modules/audit/domain"
	runtimetool "lattice-coding/internal/runtime/tool"
)

type AuditRecorder struct {
	repo domain.AuditLogRepository
}

func NewAuditRecorder(repo domain.AuditLogRepository) *AuditRecorder {
	return &AuditRecorder{repo: repo}
}

func (r *AuditRecorder) Record(ctx context.Context, event runtimetool.AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	log := &domain.AuditLog{
		RunID:        event.Request.Context.RunID,
		TraceID:      event.Request.Context.TraceID,
		EventType:    event.EventType,
		ToolName:     event.Request.Name,
		ResourceType: "tool",
		ResourceID:   event.Request.ID,
		Message:      event.Result.Error,
		PayloadJSON:  string(payload),
	}
	if log.EventType == "" {
		log.EventType = "tool_event"
	}
	return r.repo.Create(ctx, log)
}
