package domain

import "time"

type ChatSession struct {
	ID                       uint64        `json:"id"`
	Title                    string        `json:"title"`
	AgentID                  uint64        `json:"agent_id"`
	ModelConfigID            uint64        `json:"model_config_id"`
	Status                   SessionStatus `json:"status"`
	Summary                  string        `json:"summary"`
	SummarizedUntilMessageID uint64        `json:"summarized_until_message_id"`
	Meta                     string        `json:"meta"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}

type ChatMessage struct {
	ID         uint64      `json:"id"`
	SessionID  uint64      `json:"session_id"`
	Role       MessageRole `json:"role"`
	Content    string      `json:"content"`
	TokenCount int         `json:"token_count"`
	Meta       string      `json:"meta"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type PageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}
