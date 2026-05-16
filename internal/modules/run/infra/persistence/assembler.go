package persistence

import "lattice-coding/internal/modules/run/domain"

func ToRunPO(run *domain.Run) *RunPO {
	if run == nil {
		return nil
	}
	return &RunPO{
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

func ToRunDomain(po *RunPO) *domain.Run {
	if po == nil {
		return nil
	}
	return &domain.Run{
		ID:          po.ID,
		AgentID:     po.AgentID,
		SessionID:   po.SessionID,
		WorkflowID:  po.WorkflowID,
		Status:      po.Status,
		Input:       po.Input,
		Output:      po.Output,
		Error:       po.Error,
		StartedAt:   po.StartedAt,
		CompletedAt: po.CompletedAt,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

func ToInvocationPO(invocation *domain.ToolInvocation) *ToolInvocationPO {
	if invocation == nil {
		return nil
	}
	return &ToolInvocationPO{
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

func ToInvocationDomain(po *ToolInvocationPO) *domain.ToolInvocation {
	if po == nil {
		return nil
	}
	return &domain.ToolInvocation{
		ID:            po.ID,
		RunID:         po.RunID,
		NodeID:        po.NodeID,
		ToolName:      po.ToolName,
		InputJSON:     po.InputJSON,
		OutputJSON:    po.OutputJSON,
		IsError:       po.IsError,
		LatencyMs:     po.LatencyMs,
		Status:        po.Status,
		FullResultRef: po.FullResultRef,
		StartedAt:     po.StartedAt,
		CompletedAt:   po.CompletedAt,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}
