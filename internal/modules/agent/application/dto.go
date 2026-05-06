package application

import "time"

type AgentDTO struct {
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
