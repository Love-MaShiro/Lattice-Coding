package api

import (
	"lattice-coding/internal/modules/agent/application"
)

type CreateAgentRequest struct {
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	AgentType       string  `json:"agent_type"`
	ModelConfigID   uint64  `json:"model_config_id" binding:"required"`
	SystemPrompt    string  `json:"system_prompt"`
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"top_p"`
	MaxTokens       int     `json:"max_tokens"`
	MaxContextTurns int     `json:"max_context_turns"`
	MaxSteps        int     `json:"max_steps"`
	Enabled         *bool   `json:"enabled"`
}

type UpdateAgentRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	AgentType       string  `json:"agent_type"`
	ModelConfigID   uint64  `json:"model_config_id"`
	SystemPrompt    string  `json:"system_prompt"`
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"top_p"`
	MaxTokens       int     `json:"max_tokens"`
	MaxContextTurns int     `json:"max_context_turns"`
	MaxSteps        int     `json:"max_steps"`
	Enabled         *bool   `json:"enabled"`
}

type AgentPageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

func ToCreateAgentCommand(req *CreateAgentRequest) *application.CreateAgentCommand {
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	maxSteps := 20
	if req.MaxSteps > 0 {
		maxSteps = req.MaxSteps
	}
	return &application.CreateAgentCommand{
		Name:            req.Name,
		Description:     req.Description,
		AgentType:       req.AgentType,
		ModelConfigID:   req.ModelConfigID,
		SystemPrompt:    req.SystemPrompt,
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		MaxTokens:       req.MaxTokens,
		MaxContextTurns: req.MaxContextTurns,
		MaxSteps:        maxSteps,
		Enabled:         enabled,
	}
}

func ToUpdateAgentCommand(req *UpdateAgentRequest) *application.UpdateAgentCommand {
	return &application.UpdateAgentCommand{
		Name:            req.Name,
		Description:     req.Description,
		AgentType:       req.AgentType,
		ModelConfigID:   req.ModelConfigID,
		SystemPrompt:    req.SystemPrompt,
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		MaxTokens:       req.MaxTokens,
		MaxContextTurns: req.MaxContextTurns,
		MaxSteps:        req.MaxSteps,
		Enabled:         req.Enabled,
	}
}
