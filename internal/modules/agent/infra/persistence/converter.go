package persistence

import (
	"lattice-coding/internal/modules/agent/domain"
)

func ConvertAgentToPO(src *domain.Agent, dst *AgentPO) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.ProviderID = src.ProviderID
	dst.ModelConfigID = src.ModelConfigID
	dst.SystemPrompt = src.SystemPrompt
	dst.Tools = src.Tools
	if dst.Tools == "" {
		dst.Tools = "[]"
	}
	dst.MaxSteps = src.MaxSteps
	dst.Timeout = src.Timeout
	dst.Enabled = src.Enabled
}

func ConvertPOToAgent(src *AgentPO) *domain.Agent {
	return &domain.Agent{
		ID:            src.ID,
		Name:          src.Name,
		Description:   src.Description,
		ProviderID:    src.ProviderID,
		ModelConfigID: src.ModelConfigID,
		SystemPrompt:  src.SystemPrompt,
		Tools:         src.Tools,
		MaxSteps:      src.MaxSteps,
		Timeout:       src.Timeout,
		Enabled:       src.Enabled,
		CreatedAt:     src.CreatedAt,
		UpdatedAt:     src.UpdatedAt,
	}
}
