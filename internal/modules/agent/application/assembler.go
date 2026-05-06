package application

import (
	"lattice-coding/internal/modules/agent/domain"
)

func ToAgentDTO(a *domain.Agent) *AgentDTO {
	if a == nil {
		return nil
	}
	return &AgentDTO{
		ID:            a.ID,
		Name:          a.Name,
		Description:   a.Description,
		ProviderID:    a.ProviderID,
		ModelConfigID: a.ModelConfigID,
		SystemPrompt:  a.SystemPrompt,
		Tools:         a.Tools,
		MaxSteps:      a.MaxSteps,
		Timeout:       a.Timeout,
		Enabled:       a.Enabled,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func ToAgentDTOs(agents []*domain.Agent) []*AgentDTO {
	dtos := make([]*AgentDTO, len(agents))
	for i, a := range agents {
		dtos[i] = ToAgentDTO(a)
	}
	return dtos
}
