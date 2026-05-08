package application

import "lattice-coding/internal/modules/chat/domain"

func ToSessionDTO(session *domain.ChatSession) *SessionDTO {
	if session == nil {
		return nil
	}
	return &SessionDTO{
		ID:                       session.ID,
		Title:                    session.Title,
		AgentID:                  session.AgentID,
		ModelConfigID:            session.ModelConfigID,
		Status:                   string(session.Status),
		Summary:                  session.Summary,
		SummarizedUntilMessageID: session.SummarizedUntilMessageID,
		Meta:                     session.Meta,
		CreatedAt:                session.CreatedAt,
		UpdatedAt:                session.UpdatedAt,
	}
}

func ToSessionDTOs(sessions []*domain.ChatSession) []*SessionDTO {
	dtos := make([]*SessionDTO, len(sessions))
	for i, session := range sessions {
		dtos[i] = ToSessionDTO(session)
	}
	return dtos
}

func ToMessageDTO(message *domain.ChatMessage) *MessageDTO {
	if message == nil {
		return nil
	}
	return &MessageDTO{
		ID:         message.ID,
		SessionID:  message.SessionID,
		Role:       string(message.Role),
		Content:    message.Content,
		TokenCount: message.TokenCount,
		Meta:       message.Meta,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
	}
}

func ToMessageDTOs(messages []*domain.ChatMessage) []*MessageDTO {
	dtos := make([]*MessageDTO, len(messages))
	for i, message := range messages {
		dtos[i] = ToMessageDTO(message)
	}
	return dtos
}
