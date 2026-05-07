package application

import "time"

type AgentDTO struct {
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

type AgentToolDTO struct {
	ID        uint64    `json:"id"`
	ToolID    uint64    `json:"tool_id"`
	ToolType  string    `json:"tool_type"`
	CreatedAt time.Time `json:"created_at"`
}

type AgentDetailDTO struct {
	AgentDTO
	Tools []*AgentToolDTO `json:"tools"`
}
