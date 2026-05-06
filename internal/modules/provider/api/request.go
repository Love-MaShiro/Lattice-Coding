package api

import (
	"lattice-coding/internal/modules/provider/application"
)

type CreateProviderRequest struct {
	Name         string `json:"name" binding:"required"`
	ProviderType string `json:"provider_type" binding:"required"`
	BaseURL      string `json:"base_url"`
	AuthType     string `json:"auth_type" binding:"required"`
	APIKey       string `json:"api_key"`
	AuthConfig   string `json:"auth_config"`
	Config       string `json:"config"`
	Enabled      *bool  `json:"enabled"`
}

type UpdateProviderRequest struct {
	Name         string `json:"name"`
	ProviderType string `json:"provider_type"`
	BaseURL      string `json:"base_url"`
	AuthType     string `json:"auth_type"`
	APIKey       string `json:"api_key"`
	AuthConfig   string `json:"auth_config"`
	Config       string `json:"config"`
	Enabled      *bool  `json:"enabled"`
}

type ProviderPageQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"`
}

type CreateModelConfigRequest struct {
	ProviderID   uint64 `json:"provider_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Model        string `json:"model" binding:"required"`
	ModelType    string `json:"model_type"`
	Params       string `json:"params"`
	Capabilities string `json:"capabilities"`
	IsDefault    bool   `json:"is_default"`
	Enabled      *bool  `json:"enabled"`
}

type UpdateModelConfigRequest struct {
	Name         string `json:"name"`
	Model        string `json:"model"`
	ModelType    string `json:"model_type"`
	Params       string `json:"params"`
	Capabilities string `json:"capabilities"`
	IsDefault    *bool  `json:"is_default"`
	Enabled      *bool  `json:"enabled"`
}

type ModelConfigPageQuery struct {
	Page       int    `form:"page" json:"page"`
	PageSize   int    `form:"page_size" json:"page_size"`
	ProviderID uint64 `form:"provider_id" json:"provider_id"`
	Keyword    string `form:"keyword" json:"keyword"`
}

func ToCreateProviderCommand(req *CreateProviderRequest) *application.CreateProviderCommand {
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	return &application.CreateProviderCommand{
		Name:         req.Name,
		ProviderType: req.ProviderType,
		BaseURL:      req.BaseURL,
		AuthType:     req.AuthType,
		APIKey:       req.APIKey,
		AuthConfig:   req.AuthConfig,
		Config:       req.Config,
		Enabled:      enabled,
	}
}

func ToUpdateProviderCommand(req *UpdateProviderRequest) *application.UpdateProviderCommand {
	return &application.UpdateProviderCommand{
		Name:         req.Name,
		ProviderType: req.ProviderType,
		BaseURL:      req.BaseURL,
		AuthType:     req.AuthType,
		APIKey:       req.APIKey,
		AuthConfig:   req.AuthConfig,
		Config:       req.Config,
		Enabled:      req.Enabled,
	}
}

func ToCreateModelConfigCommand(req *CreateModelConfigRequest) *application.CreateModelConfigCommand {
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	return &application.CreateModelConfigCommand{
		ProviderID:   req.ProviderID,
		Name:         req.Name,
		Model:        req.Model,
		ModelType:    req.ModelType,
		Params:       req.Params,
		Capabilities: req.Capabilities,
		IsDefault:    req.IsDefault,
		Enabled:      enabled,
	}
}

func ToUpdateModelConfigCommand(req *UpdateModelConfigRequest) *application.UpdateModelConfigCommand {
	return &application.UpdateModelConfigCommand{
		Name:         req.Name,
		Model:        req.Model,
		ModelType:    req.ModelType,
		Params:       req.Params,
		Capabilities: req.Capabilities,
		IsDefault:    req.IsDefault,
		Enabled:      req.Enabled,
	}
}
