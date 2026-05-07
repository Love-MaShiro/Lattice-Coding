package api

import (
	"lattice-coding/internal/modules/agent/application"
	"lattice-coding/internal/modules/agent/domain"
	"time"
)

type AgentResponse struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AgentType       string    `json:"agent_type"`
	ModelConfigID   uint64    `json:"model_config_id"`
	SystemPrompt    string    `json:"system_prompt"`
	Temperature     float64   `json:"temperature"`
	TopP            float64   `json:"top_p"`
	MaxTokens       int       `json:"max_tokens"`
	MaxContextTurns int       `json:"max_context_turns"`
	MaxSteps        int       `json:"max_steps"`
	Enabled         bool      `json:"enabled"`
	ToolCount       int64     `json:"tool_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToAgentResponse(dto *application.AgentDTO) AgentResponse {
	if dto == nil {
		return AgentResponse{}
	}
	return AgentResponse{
		ID:              dto.ID,
		Name:            dto.Name,
		Description:     dto.Description,
		AgentType:       dto.AgentType,
		ModelConfigID:   dto.ModelConfigID,
		SystemPrompt:    dto.SystemPrompt,
		Temperature:     dto.Temperature,
		TopP:            dto.TopP,
		MaxTokens:       dto.MaxTokens,
		MaxContextTurns: dto.MaxContextTurns,
		MaxSteps:        dto.MaxSteps,
		Enabled:         dto.Enabled,
		ToolCount:       dto.ToolCount,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
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

type AgentToolResponse struct {
	ID        uint64 `json:"id"`
	ToolID    uint64 `json:"tool_id"`
	ToolType  string `json:"tool_type"`
	CreatedAt string `json:"created_at"`
}

type AgentDetailResponse struct {
	AgentResponse
	Tools []AgentToolResponse `json:"tools"`
}

func ToAgentToolResponse(dto *application.AgentToolDTO) AgentToolResponse {
	if dto == nil {
		return AgentToolResponse{}
	}
	return AgentToolResponse{
		ID:        dto.ID,
		ToolID:    dto.ToolID,
		ToolType:  dto.ToolType,
		CreatedAt: dto.CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func ToAgentDetailResponse(dto *application.AgentDetailDTO) AgentDetailResponse {
	if dto == nil {
		return AgentDetailResponse{}
	}
	tools := make([]AgentToolResponse, len(dto.Tools))
	for i, t := range dto.Tools {
		tools[i] = ToAgentToolResponse(t)
	}
	return AgentDetailResponse{
		AgentResponse: ToAgentResponse(&dto.AgentDTO),
		Tools:         tools,
	}
}
