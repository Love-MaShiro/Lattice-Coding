package application

import "time"

type SessionDTO struct {
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

type MessageDTO struct {
	ID         uint64    `json:"id"`
	SessionID  uint64    `json:"session_id"`
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	TokenCount int       `json:"token_count"`
	Meta       string    `json:"meta"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CompletionDTO struct {
	SessionID uint64      `json:"session_id"`
	Message   *MessageDTO `json:"message"`
	Content   string      `json:"content"`
}

type AgentRuntimeDTO struct {
	ID              uint64  `json:"id"`
	Name            string  `json:"name"`
	ModelConfigID   uint64  `json:"model_config_id"`
	SystemPrompt    string  `json:"system_prompt"`
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"top_p"`
	MaxTokens       int     `json:"max_tokens"`
	MaxContextTurns int     `json:"max_context_turns"`
	Enabled         bool    `json:"enabled"`
}
