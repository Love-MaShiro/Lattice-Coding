package api

import (
	"lattice-coding/internal/modules/agent/application"
)

type CreateAgentRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	ProviderID    uint64 `json:"provider_id" binding:"required"`
	ModelConfigID uint64 `json:"model_config_id" binding:"required"`
	SystemPrompt  string `json:"system_prompt"`
	Tools         string `json:"tools"`
	MaxSteps      int    `json:"max_steps"`
	Timeout       int    `json:"timeout"`
	Enabled       *bool  `json:"enabled"`
}

type UpdateAgentRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ProviderID    uint64 `json:"provider_id"`
	ModelConfigID uint64 `json:"model_config_id"`
	SystemPrompt  string `json:"system_prompt"`
	Tools         string `json:"tools"`
	MaxSteps      int    `json:"max_steps"`
	Timeout       int    `json:"timeout"`
	Enabled       *bool  `json:"enabled"`
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
	timeout := 1200
	if req.Timeout > 0 {
		timeout = req.Timeout
	}
	return &application.CreateAgentCommand{
		Name:          req.Name,
		Description:   req.Description,
		ProviderID:    req.ProviderID,
		ModelConfigID: req.ModelConfigID,
		SystemPrompt:  req.SystemPrompt,
		Tools:         req.Tools,
		MaxSteps:      maxSteps,
		Timeout:       timeout,
		Enabled:       enabled,
	}
}

func ToUpdateAgentCommand(req *UpdateAgentRequest) *application.UpdateAgentCommand {
	return &application.UpdateAgentCommand{
		Name:          req.Name,
		Description:   req.Description,
		ProviderID:    req.ProviderID,
		ModelConfigID: req.ModelConfigID,
		SystemPrompt:  req.SystemPrompt,
		Tools:         req.Tools,
		MaxSteps:      req.MaxSteps,
		Timeout:       req.Timeout,
		Enabled:       req.Enabled,
	}
}
