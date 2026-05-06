package api

import (
	"lattice-coding/internal/modules/agent/application"
	"lattice-coding/internal/modules/agent/domain"
	"time"
)

type AgentResponse struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ProviderID    uint64    `json:"provider_id"`
	ModelConfigID uint64    `json:"model_config_id"`
	SystemPrompt  string    `json:"system_prompt"`
	Tools         string    `json:"tools"`
	MaxSteps      int       `json:"max_steps"`
	Timeout       int       `json:"timeout"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToAgentResponse(dto *application.AgentDTO) AgentResponse {
	if dto == nil {
		return AgentResponse{}
	}
	return AgentResponse{
		ID:            dto.ID,
		Name:          dto.Name,
		Description:   dto.Description,
		ProviderID:    dto.ProviderID,
		ModelConfigID: dto.ModelConfigID,
		SystemPrompt:  dto.SystemPrompt,
		Tools:         dto.Tools,
		MaxSteps:      dto.MaxSteps,
		Timeout:       dto.Timeout,
		Enabled:       dto.Enabled,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

func ToAgentPageResponse(result *domain.PageResult[*application.AgentDTO]) *AgentPageResponse {
	if result == nil {
		return nil
	}
	items := make([]AgentResponse, 0, len(result.Items))
	for _, dto := range result.Items {
		items = append(items, ToAgentResponse(dto))
	}
	return &AgentPageResponse{
		Items:    items,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}
}

type AgentPageResponse struct {
	Items    []AgentResponse `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
