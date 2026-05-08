package api

import (
	"time"

	"lattice-coding/internal/modules/chat/application"
	"lattice-coding/internal/modules/chat/domain"
)

type SessionResponse struct {
	ID                       uint64    `json:"id"`
	Title                    string    `json:"title"`
	AgentID                  uint64    `json:"agent_id"`
	ModelConfigID            uint64    `json:"model_config_id"`
	Status                   string    `json:"status"`
	Summary                  string    `json:"summary"`
	SummarizedUntilMessageID uint64    `json:"summarized_until_message_id"`
	Meta                     string    `json:"meta"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type MessageResponse struct {
	ID         uint64    `json:"id"`
	SessionID  uint64    `json:"session_id"`
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	TokenCount int       `json:"token_count"`
	Meta       string    `json:"meta"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CompletionResponse struct {
	SessionID uint64          `json:"session_id"`
	Message   MessageResponse `json:"message"`
	Content   string          `json:"content"`
}

type SessionPageResponse struct {
	Items    []SessionResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

func ToSessionResponse(dto *application.SessionDTO) SessionResponse {
	if dto == nil {
		return SessionResponse{}
	}
	return SessionResponse{
		ID:                       dto.ID,
		Title:                    dto.Title,
		AgentID:                  dto.AgentID,
		ModelConfigID:            dto.ModelConfigID,
		Status:                   dto.Status,
		Summary:                  dto.Summary,
		SummarizedUntilMessageID: dto.SummarizedUntilMessageID,
		Meta:                     dto.Meta,
		CreatedAt:                dto.CreatedAt,
		UpdatedAt:                dto.UpdatedAt,
	}
}

func ToMessageResponse(dto *application.MessageDTO) MessageResponse {
	if dto == nil {
		return MessageResponse{}
	}
	return MessageResponse{
		ID:         dto.ID,
		SessionID:  dto.SessionID,
		Role:       dto.Role,
		Content:    dto.Content,
		TokenCount: dto.TokenCount,
		Meta:       dto.Meta,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}

func ToMessageResponses(dtos []*application.MessageDTO) []MessageResponse {
	items := make([]MessageResponse, 0, len(dtos))
	for _, dto := range dtos {
		items = append(items, ToMessageResponse(dto))
	}
	return items
}

func ToCompletionResponse(dto *application.CompletionDTO) CompletionResponse {
	if dto == nil {
		return CompletionResponse{}
	}
	return CompletionResponse{
		SessionID: dto.SessionID,
		Message:   ToMessageResponse(dto.Message),
		Content:   dto.Content,
	}
}

func ToSessionPageResponse(result *domain.PageResult[*application.SessionDTO]) *SessionPageResponse {
	if result == nil {
		return nil
	}
	items := make([]SessionResponse, 0, len(result.Items))
	for _, dto := range result.Items {
		items = append(items, ToSessionResponse(dto))
	}
	return &SessionPageResponse{
		Items:    items,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}
}
