package application

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"lattice-coding/internal/modules/run/domain"
	runtimetool "lattice-coding/internal/runtime/tool"
)

type ToolInvocationRecorder struct {
	repo domain.ToolInvocationRepository
}

func NewToolInvocationRecorder(repo domain.ToolInvocationRepository) *ToolInvocationRecorder {
	return &ToolInvocationRecorder{repo: repo}
}

func (r *ToolInvocationRecorder) Start(ctx context.Context, record runtimetool.ToolInvocationStarted) (string, error) {
	inputJSON, err := marshalJSON(record.Input)
	if err != nil {
		return "", err
	}
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	invocation := &domain.ToolInvocation{
		ID:        id,
		RunID:     record.RunID,
		NodeID:    record.NodeID,
		ToolName:  record.ToolName,
		InputJSON: inputJSON,
		Status:    domain.ToolInvocationStatusStarted,
		StartedAt: record.StartedAt,
	}
	if invocation.StartedAt.IsZero() {
		invocation.StartedAt = time.Now()
	}
	if err := r.repo.Create(ctx, invocation); err != nil {
		return "", err
	}
	return id, nil
}

func (r *ToolInvocationRecorder) Finish(ctx context.Context, record runtimetool.ToolInvocationFinished) error {
	outputJSON, err := marshalJSON(record.Result)
	if err != nil {
		return err
	}
	completedAt := record.CompletedAt
	if completedAt.IsZero() {
		completedAt = time.Now()
	}
	status := record.Status
	if status == "" {
		status = domain.ToolInvocationStatusSuccess
		if record.Result.IsError {
			status = domain.ToolInvocationStatusFailed
		}
	}
	invocation := &domain.ToolInvocation{
		ID:            record.ID,
		RunID:         record.RunID,
		NodeID:        record.NodeID,
		ToolName:      record.ToolName,
		OutputJSON:    outputJSON,
		IsError:       record.Result.IsError,
		LatencyMs:     record.LatencyMs,
		Status:        status,
		FullResultRef: record.FullResultRef,
		CompletedAt:   &completedAt,
	}
	return r.repo.Update(ctx, invocation)
}

func marshalJSON(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
