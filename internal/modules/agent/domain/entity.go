package domain

import "time"

type Agent struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AgentType       AgentType `json:"agent_type"`
	ModelConfigID   uint64    `json:"model_config_id"`
	SystemPrompt    string    `json:"system_prompt"`
	Temperature     float64   `json:"temperature"`
	TopP            float64   `json:"top_p"`
	MaxTokens       int       `json:"max_tokens"`
	MaxContextTurns int       `json:"max_context_turns"`
	MaxSteps        int       `json:"max_steps"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AgentTool struct {
	ID        uint64    `json:"id"`
	AgentID   uint64    `json:"agent_id"`
	ToolID    uint64    `json:"tool_id"`
	ToolType  string    `json:"tool_type"`
	CreatedAt time.Time `json:"created_at"`
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
