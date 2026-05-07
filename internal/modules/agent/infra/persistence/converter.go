package persistence

import (
	"lattice-coding/internal/modules/agent/domain"
)

func ConvertAgentToPO(src *domain.Agent, dst *AgentPO) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.AgentType = string(src.AgentType)
	dst.ModelConfigID = src.ModelConfigID
	dst.SystemPrompt = src.SystemPrompt
	dst.Temperature = src.Temperature
	dst.TopP = src.TopP
	dst.MaxTokens = src.MaxTokens
	dst.MaxContextTurns = src.MaxContextTurns
	dst.MaxSteps = src.MaxSteps
	dst.Enabled = src.Enabled
}

func ConvertPOToAgent(src *AgentPO) *domain.Agent {
	return &domain.Agent{
		ID:              src.ID,
		Name:            src.Name,
		Description:     src.Description,
		AgentType:       domain.AgentType(src.AgentType),
		ModelConfigID:   src.ModelConfigID,
		SystemPrompt:    src.SystemPrompt,
		Temperature:     src.Temperature,
		TopP:            src.TopP,
		MaxTokens:       src.MaxTokens,
		MaxContextTurns: src.MaxContextTurns,
		MaxSteps:        src.MaxSteps,
		Enabled:         src.Enabled,
		CreatedAt:       src.CreatedAt,
		UpdatedAt:       src.UpdatedAt,
	}
}

func ConvertAgentToolToPO(src *domain.AgentTool, dst *AgentToolPO) {
	dst.ID = src.ID
	dst.AgentID = src.AgentID
	dst.ToolID = src.ToolID
	dst.ToolType = src.ToolType
}

func ConvertPOToAgentTool(src *AgentToolPO) *domain.AgentTool {
	return &domain.AgentTool{
		ID:        src.ID,
		AgentID:   src.AgentID,
		ToolID:    src.ToolID,
		ToolType:  src.ToolType,
		CreatedAt: src.CreatedAt,
	}
}

func ConvertPOsToAgentTools(pos []AgentToolPO) []*domain.AgentTool {
	result := make([]*domain.AgentTool, len(pos))
	for i := range pos {
		result[i] = ConvertPOToAgentTool(&pos[i])
	}
	return result
}
