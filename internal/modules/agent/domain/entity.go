package domain

import "time"

type Agent struct {
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
	Deleted       bool      `json:"deleted"`
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
