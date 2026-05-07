package application

import (
	"lattice-coding/internal/modules/agent/domain"
)

func ToAgentDTO(a *domain.Agent) *AgentDTO {
	if a == nil {
		return nil
	}
	return &AgentDTO{
		ID:              a.ID,
		Name:            a.Name,
		Description:     a.Description,
		AgentType:       string(a.AgentType),
		ModelConfigID:   a.ModelConfigID,
		SystemPrompt:    a.SystemPrompt,
		Temperature:     a.Temperature,
		TopP:            a.TopP,
		MaxTokens:       a.MaxTokens,
		MaxContextTurns: a.MaxContextTurns,
		MaxSteps:        a.MaxSteps,
		Enabled:         a.Enabled,
		ToolCount:       0,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func ToAgentDTOWithToolCount(a *domain.Agent, toolCount int64) *AgentDTO {
	if a == nil {
		return nil
	}
	return &AgentDTO{
		ID:              a.ID,
		Name:            a.Name,
		Description:     a.Description,
		AgentType:       string(a.AgentType),
		ModelConfigID:   a.ModelConfigID,
		SystemPrompt:    a.SystemPrompt,
		Temperature:     a.Temperature,
		TopP:            a.TopP,
		MaxTokens:       a.MaxTokens,
		MaxContextTurns: a.MaxContextTurns,
		MaxSteps:        a.MaxSteps,
		Enabled:         a.Enabled,
		ToolCount:       toolCount,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func ToAgentDTOs(agents []*domain.Agent) []*AgentDTO {
	dtos := make([]*AgentDTO, len(agents))
	for i, a := range agents {
		dtos[i] = ToAgentDTO(a)
	}
	return dtos
}

func ToAgentToolDTO(t *domain.AgentTool) *AgentToolDTO {
	if t == nil {
		return nil
	}
	return &AgentToolDTO{
		ID:        t.ID,
		ToolID:    t.ToolID,
		ToolType:  t.ToolType,
		CreatedAt: t.CreatedAt,
	}
}

func ToAgentToolDTOs(tools []*domain.AgentTool) []*AgentToolDTO {
	dtos := make([]*AgentToolDTO, len(tools))
	for i, t := range tools {
		dtos[i] = ToAgentToolDTO(t)
	}
	return dtos
}

func ToAgentDetailDTO(agent *domain.Agent, tools []*domain.AgentTool) *AgentDetailDTO {
	if agent == nil {
		return nil
	}
	return &AgentDetailDTO{
		AgentDTO: *ToAgentDTO(agent),
		Tools:    ToAgentToolDTOs(tools),
	}
}
