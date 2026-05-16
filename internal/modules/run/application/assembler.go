package application

import "lattice-coding/internal/modules/run/domain"

func ToRunDTO(run *domain.Run) *RunDTO {
	if run == nil {
		return nil
	}
	return &RunDTO{
		ID:          run.ID,
		AgentID:     run.AgentID,
		SessionID:   run.SessionID,
		WorkflowID:  run.WorkflowID,
		Status:      run.Status,
		Input:       run.Input,
		Output:      run.Output,
		Error:       run.Error,
		StartedAt:   run.StartedAt,
		CompletedAt: run.CompletedAt,
		CreatedAt:   run.CreatedAt,
		UpdatedAt:   run.UpdatedAt,
	}
}

func ToRunDTOs(runs []*domain.Run) []*RunDTO {
	items := make([]*RunDTO, 0, len(runs))
	for _, run := range runs {
		items = append(items, ToRunDTO(run))
	}
	return items
}

func ToToolInvocationDTO(invocation *domain.ToolInvocation) *ToolInvocationDTO {
	if invocation == nil {
		return nil
	}
	return &ToolInvocationDTO{
		ID:            invocation.ID,
		RunID:         invocation.RunID,
		NodeID:        invocation.NodeID,
		ToolName:      invocation.ToolName,
		InputJSON:     invocation.InputJSON,
		OutputJSON:    invocation.OutputJSON,
		IsError:       invocation.IsError,
		LatencyMs:     invocation.LatencyMs,
		Status:        invocation.Status,
		FullResultRef: invocation.FullResultRef,
		StartedAt:     invocation.StartedAt,
		CompletedAt:   invocation.CompletedAt,
		CreatedAt:     invocation.CreatedAt,
		UpdatedAt:     invocation.UpdatedAt,
	}
}

func ToToolInvocationDTOs(invocations []*domain.ToolInvocation) []*ToolInvocationDTO {
	items := make([]*ToolInvocationDTO, 0, len(invocations))
	for _, invocation := range invocations {
		items = append(items, ToToolInvocationDTO(invocation))
	}
	return items
}
