package api

import "lattice-coding/internal/modules/chat/application"

type CreateSessionRequest struct {
	Title   string `json:"title"`
	AgentID uint64 `json:"agent_id" binding:"required"`
}

type CreateMessageRequest struct {
	SessionID uint64 `json:"session_id" binding:"required"`
	Role      string `json:"role"`
	Content   string `json:"content" binding:"required"`
	Meta      string `json:"meta"`
}

type CompletionRequest struct {
	AgentID   uint64 `json:"agent_id"`
	SessionID uint64 `json:"session_id"`
	Message   string `json:"message" binding:"required"`
}

type SessionPageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

type MessageListQuery struct {
	Limit int `form:"limit" json:"limit"`
}

func ToCreateSessionCommand(req *CreateSessionRequest) *application.CreateSessionCommand {
	return &application.CreateSessionCommand{
		Title:   req.Title,
		AgentID: req.AgentID,
	}
}

func ToCreateMessageCommand(req *CreateMessageRequest) *application.CreateMessageCommand {
	return &application.CreateMessageCommand{
		SessionID: req.SessionID,
		Role:      req.Role,
		Content:   req.Content,
		Meta:      req.Meta,
	}
}

func ToCompletionCommand(req *CompletionRequest) *application.CompletionCommand {
	return &application.CompletionCommand{
		AgentID:   req.AgentID,
		SessionID: req.SessionID,
		Message:   req.Message,
	}
}
