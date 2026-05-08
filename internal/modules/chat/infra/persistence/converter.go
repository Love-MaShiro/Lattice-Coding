package persistence

import "lattice-coding/internal/modules/chat/domain"

func ConvertSessionToPO(src *domain.ChatSession, dst *ChatSessionPO) {
	dst.ID = src.ID
	dst.Title = src.Title
	dst.AgentID = src.AgentID
	dst.ModelConfigID = src.ModelConfigID
	dst.Status = string(src.Status)
	dst.Summary = src.Summary
	dst.SummarizedUntilMessageID = src.SummarizedUntilMessageID
	dst.Meta = src.Meta
	if dst.Status == "" {
		dst.Status = string(domain.SessionStatusActive)
	}
}

func ConvertPOToSession(src *ChatSessionPO) *domain.ChatSession {
	return &domain.ChatSession{
		ID:                       src.ID,
		Title:                    src.Title,
		AgentID:                  src.AgentID,
		ModelConfigID:            src.ModelConfigID,
		Status:                   domain.SessionStatus(src.Status),
		Summary:                  src.Summary,
		SummarizedUntilMessageID: src.SummarizedUntilMessageID,
		Meta:                     src.Meta,
		CreatedAt:                src.CreatedAt,
		UpdatedAt:                src.UpdatedAt,
	}
}

func ConvertMessageToPO(src *domain.ChatMessage, dst *ChatMessagePO) {
	dst.ID = src.ID
	dst.SessionID = src.SessionID
	dst.Role = string(src.Role)
	dst.Content = src.Content
	dst.TokenCount = src.TokenCount
	dst.Meta = src.Meta
}

func ConvertPOToMessage(src *ChatMessagePO) *domain.ChatMessage {
	return &domain.ChatMessage{
		ID:         src.ID,
		SessionID:  src.SessionID,
		Role:       domain.MessageRole(src.Role),
		Content:    src.Content,
		TokenCount: src.TokenCount,
		Meta:       src.Meta,
		CreatedAt:  src.CreatedAt,
		UpdatedAt:  src.UpdatedAt,
	}
}

func ConvertPOsToMessages(pos []ChatMessagePO) []*domain.ChatMessage {
	items := make([]*domain.ChatMessage, len(pos))
	for i := range pos {
		items[i] = ConvertPOToMessage(&pos[i])
	}
	return items
}
